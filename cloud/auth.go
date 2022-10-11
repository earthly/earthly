package cloud

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	secretsapi "github.com/earthly/cloud-api/secrets"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/agent"
)

// Authenticate fetches a new auth token from the server and saves it to the client.
// The user should have credentials store on disk within the ~/.earthly directory.
// Credentials may be either email/password, ssh-based, or a custom token.
// Upon successful authenticate, the JWT provided by the server is stored in
// ~/.earthly/auth.jwt, and can be refreshed any time via another call to Authenticate().
func (c *client) Authenticate(ctx context.Context) error {
	var err error
	switch {
	case c.email != "" && c.password != "":
		err = c.loginWithPassowrd(ctx)
	case c.authCredToken != "":
		err = c.loginWithToken(ctx)
	default:
		err = c.loginWithSSH(ctx)
	}
	if err != nil {
		if errors.Is(err, ErrNoAuthorizedPublicKeys) || errors.Is(err, ErrNoSSHAgent) {
			return ErrUnauthorized
		}
		return err
	}
	return c.saveToken()
}

func (c *client) IsLoggedIn(ctx context.Context) bool {
	return c.authToken != "" || c.authCredToken != ""
}

func (c *client) FindSSHCredentials(ctx context.Context, emailToFind string) error {
	keys, err := c.GetPublicKeys(ctx)
	if err != nil {
		return err
	}
	challenge, err := c.getChallenge(ctx)
	if err != nil {
		return err
	}
	for _, key := range keys {
		credentials, err := c.getSSHCredentials(challenge, key)
		if err != nil {
			return err
		}
		c.authToken, c.authTokenExpiry, err = c.login(ctx, credentials)
		if errors.Is(err, ErrUnauthorized) {
			continue // try next key
		} else if err != nil {
			return err
		}
		email, _, err := c.ping(ctx)
		if err != nil {
			return err
		}
		if email == emailToFind {
			if err := c.SetSSHCredentials(ctx, email, key.String()); err != nil {
				return err
			}
			return nil
		}
	}
	return ErrNoAuthorizedPublicKeys
}

func (c *client) GetAuthToken(ctx context.Context) (string, error) {
	err := c.Authenticate(ctx) // Ensure the current token is valid
	if err != nil {
		return "", errors.Wrap(err, "could not authenticate")
	}
	return c.authToken, nil
}

func (c *client) SetSSHCredentials(ctx context.Context, email, sshKey string) error {
	sshKeyType, sshKeyBlob, _, err := parseSSHKey(sshKey)
	if err != nil {
		return err
	}

	c.password = ""
	c.authCredToken = ""
	c.email = email

	c.sshKeyBlob, err = base64.StdEncoding.DecodeString(sshKeyBlob)
	if err != nil {
		return errors.Wrap(err, "base64 decode failed")
	}

	authedEmail, _, _, err := c.WhoAmI(ctx)
	if err != nil {
		return err
	}
	if authedEmail != email {
		return errors.Errorf("failed to set correct email") // shouldn't happen
	}
	return c.saveCredentials(ctx, email, sshKeyType, sshKeyBlob)
}

func (c *client) DeleteCachedToken(ctx context.Context) error {
	var zero time.Time
	c.authToken = ""
	c.authTokenExpiry = zero
	tokenPath, err := c.getTokenPath(false)
	if err != nil {
		return err
	}
	if err = os.Remove(tokenPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return errors.Wrapf(err, "failed to delete cached token %s", tokenPath)
	}
	return nil
}

func (c *client) DeleteAuthCache(ctx context.Context) error {
	err := c.DeleteCachedToken(ctx)
	if err != nil {
		return err
	}
	err = c.deleteCachedCredentials()
	if err != nil {
		return err
	}
	return nil
}

func (c *client) SetTokenCredentials(ctx context.Context, token string) (string, error) {
	c.email = ""
	c.password = ""
	c.authCredToken = token
	email, _, _, err := c.WhoAmI(ctx)
	if err != nil {
		return "", err
	}
	err = c.saveCredentials(ctx, email, "token", token)
	if err != nil {
		return "", err
	}
	return email, nil
}

func (c *client) DisableSSHKeyGuessing(ctx context.Context) {
	c.disableSSHKeyGuessing = true
}

func (c *client) SetAuthTokenDir(ctx context.Context, path string) {
	c.authDir = path
}

func (c *client) SetPasswordCredentials(ctx context.Context, email, password string) error {
	c.authCredToken = ""
	c.email = email
	c.password = password
	_, _, _, err := c.WhoAmI(ctx)
	if err != nil {
		return err
	}
	return c.savePasswordCredentials(ctx, email, password)
}

// IsValidEmail returns true if email is valid
func IsValidEmail(ctx context.Context, email string) bool {
	if strings.Contains(email, " ") {
		return false
	}
	parts := strings.Split(email, "@")
	return len(parts) == 2
}

func (c *client) loadCredentials() error {
	credPath, err := c.getCredentialsPath(false)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(credPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return errors.Wrap(err, "failed to read file")
	}
	parts := strings.SplitN(string(data), " ", 3)
	if len(parts) != 3 {
		return nil
	}
	c.email = parts[0]
	credType := parts[1]
	credData := parts[2]
	switch credType {
	case "password":
		passwordBytes, err := base64.StdEncoding.DecodeString(credData)
		if err != nil {
			return errors.Wrap(err, "base64 decode failed")
		}
		c.password = string(passwordBytes)
	case "ssh-rsa":
		c.sshKeyBlob, err = base64.StdEncoding.DecodeString(credData)
		if err != nil {
			return errors.Wrap(err, "base64 decode failed")
		}
	case "token":
		c.authCredToken = credData
	default:
		c.warnFunc("unable to handle cached auth type %s", credType)
	}
	return nil
}

func (c *client) loadToken() error {
	tokenPath, err := c.getTokenPath(false)
	if err != nil {
		return err
	}
	if exists, _ := fileutil.FileExists(tokenPath); !exists {
		return nil
	}
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return errors.Wrap(err, "failed to read file")
	}
	parts := strings.SplitN(string(data), " ", 2)
	if len(parts) != 2 {
		// trigger re-authenticate and save a new token
		return nil
	}
	c.authToken = parts[0]
	c.authTokenExpiry, err = time.Parse(tokenExpiryLayout, parts[1])
	if err != nil {
		// trigger re-authenticate and save a new token
		return nil
	}
	return nil
}

func (c *client) getCredentialsPath(create bool) (string, error) {
	confDirPath := c.authDir
	if confDirPath == "" {
		if create {
			var err error
			confDirPath, err = cliutil.GetOrCreateEarthlyDir()
			if err != nil {
				return "", errors.Wrap(err, "cannot get .earthly dir")
			}
		} else {
			confDirPath = cliutil.GetEarthlyDir()
		}
	}
	credPath := filepath.Join(confDirPath, "auth.credentials")
	return credPath, nil
}

func (c *client) migrateOldToken() error {
	confDirPath := c.authDir
	if confDirPath == "" {
		confDirPath = cliutil.GetEarthlyDir()
	}
	tokenPath := filepath.Join(confDirPath, "auth.token")
	newPath := filepath.Join(confDirPath, "auth.credentials")
	if ok, _ := fileutil.FileExists(tokenPath); ok {
		if err := os.Rename(tokenPath, newPath); err != nil {
			return errors.Wrapf(err, "failed to migrate credentials from '%s' to '%s'", tokenPath, newPath)
		}
	}
	return nil
}

func (c *client) getTokenPath(create bool) (string, error) {
	confDirPath := c.authDir
	if confDirPath == "" {
		if create {
			var err error
			confDirPath, err = cliutil.GetOrCreateEarthlyDir()
			if err != nil {
				return "", errors.Wrap(err, "cannot get .earthly dir")
			}
		} else {
			confDirPath = cliutil.GetEarthlyDir()
		}
	}
	tokenPath := filepath.Join(confDirPath, "auth.jwt")
	return tokenPath, nil
}

func (c *client) saveCredentials(ctx context.Context, email, tokenType, tokenValue string) error {
	tokenPath, err := c.getCredentialsPath(true)
	if err != nil {
		return err
	}

	if !IsValidEmail(ctx, email) {
		return errors.Errorf("invalid email: %q", email)
	}
	if strings.Contains(tokenType, " ") {
		return errors.Errorf("invalid token type: %q", tokenType)
	}
	if strings.Contains(tokenValue, " ") {
		return errors.Errorf("invalid token value: %q", tokenValue)
	}

	data := []byte(email + " " + tokenType + " " + tokenValue)
	err = os.WriteFile(tokenPath, []byte(data), 0600)
	if err != nil {
		return errors.Wrapf(err, "failed to store auth credentials")
	}
	return nil
}

func (c *client) saveSSHCredentials(ctx context.Context, email, sshKey string) error {
	sshKeyType, sshKeyBlob, _, err := parseSSHKey(sshKey)
	if err != nil {
		return err
	}
	return c.saveCredentials(ctx, email, sshKeyType, sshKeyBlob)
}

func (c *client) savePasswordCredentials(ctx context.Context, email, password string) error {
	password64 := base64.StdEncoding.EncodeToString([]byte(password))
	return c.saveCredentials(ctx, email, "password", password64)
}

func (c *client) deleteCachedCredentials() error {
	c.email = ""
	c.password = ""
	c.authCredToken = ""
	credsPath, err := c.getCredentialsPath(false)
	if err != nil {
		return err
	}
	if err = os.Remove(credsPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return errors.Wrapf(err, "failed to delete cached credentials %s", credsPath)
	}
	return nil
}

// loads the following files:
//   - ~/.earthly/auth.credentials
//   - ~/.earthly/auth.jwt
//
// If a an old-style auth.token file exists, it is automatically migrated and removed.
func (c *client) loadAuthStorage() error {
	if err := c.migrateOldToken(); err != nil {
		return err
	}
	if err := c.loadToken(); err != nil {
		return err
	}
	if err := c.loadCredentials(); err != nil {
		return err
	}
	return nil
}

func (c *client) saveToken() error {
	path, err := c.getTokenPath(true)
	if err != nil {
		return err
	}
	data := []byte(fmt.Sprintf(
		"%s %s",
		c.authToken,
		c.authTokenExpiry.Format(tokenExpiryLayout)))
	if err = os.WriteFile(path, data, 0600); err != nil {
		return errors.Wrap(err, "failed writing auth token to disk")
	}
	return nil
}

func (c *client) loginWithPassowrd(ctx context.Context) error {
	var err error
	c.authCredToken = getPasswordAuthToken(c.email, c.password)
	c.authToken, c.authTokenExpiry, err = c.login(ctx, c.authCredToken)
	return err
}

func (c *client) loginWithToken(ctx context.Context) error {
	var err error
	c.authToken, c.authTokenExpiry, err = c.login(ctx, "token "+c.authCredToken)
	return err
}

func (c *client) loginWithSSH(ctx context.Context) error {
	if c.disableSSHKeyGuessing {
		return ErrNoAuthorizedPublicKeys
	}
	challenge, err := c.getChallenge(ctx)
	if err != nil {
		return err
	}
	keys, err := c.GetPublicKeys(ctx)
	if err != nil {
		return err
	}
	for _, key := range keys {
		credentials, err := c.getSSHCredentials(challenge, key)
		if err != nil {
			return err
		}
		c.authToken, c.authTokenExpiry, err = c.login(ctx, credentials)
		if errors.Is(err, ErrUnauthorized) {
			continue // try next key
		} else if err != nil {
			return err
		}
		email, _, err := c.ping(ctx)
		if err != nil {
			return err
		}
		return c.saveSSHCredentials(ctx, email, key.String())
	}
	return ErrNoAuthorizedPublicKeys
}

// login calls the login endpoint on the cloud server, passing the provided credentials.
// If auth succeeds, a new jwt token is returned with it's expiry date.
// ErrUnauthroized is returned if the credentials are not valid.
func (c *client) login(ctx context.Context, credentials string) (token string, expiry time.Time, err error) {
	var zero time.Time
	status, body, err := c.doCall(ctx, "POST", "/api/v0/account/login",
		withHeader("Authorization", credentials))
	if err != nil {
		return "", zero, errors.Wrap(err, "failed to execute login request")
	}
	if status == http.StatusUnauthorized {
		return "", zero, ErrUnauthorized
	}
	if status != http.StatusOK {
		return "", zero, errors.Errorf("unexpected status code from login: %d", status)
	}
	var resp secretsapi.LoginResponse
	err = c.jum.Unmarshal(body, &resp)
	if err != nil {
		return "", zero, errors.Wrap(err, "failed to unmarshal login response")
	}
	return resp.Token, resp.Expiry.AsTime().UTC(), nil
}

func (c *client) getChallenge(ctx context.Context) (string, error) {
	status, body, err := c.doCall(ctx, "GET", "/api/v0/account/auth-challenge")
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return "", errors.Errorf("failed to get auth challenge: %s", msg)
	}

	var challengeResponse secretsapi.AuthChallengeResponse
	err = c.jum.Unmarshal(body, &challengeResponse)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal challenge response")
	}

	return challengeResponse.Challenge, nil
}

func (c *client) signChallenge(challenge string, key *agent.Key) (string, error) {
	sig, err := c.sshAgent.Sign(key, []byte(challenge))
	if err != nil {
		return "", err
	}
	s := base64.StdEncoding.EncodeToString(sig.Blob)
	return s, nil
}

func (c *client) getSSHCredentials(challenge string, key *agent.Key) (credentials string, err error) {
	sig, err := c.signChallenge(challenge, key)
	if err != nil {
		return credentials, err
	}
	blob := base64.StdEncoding.EncodeToString(key.Blob)
	credentials = fmt.Sprintf("ssh-rsa %s %s", blob, sig)
	return credentials, nil
}

func getPasswordAuthToken(email, password string) string {
	email64 := base64.StdEncoding.EncodeToString([]byte(email))
	password64 := base64.StdEncoding.EncodeToString([]byte(password))
	return fmt.Sprintf("password %s %s", email64, password64)
}

func parseSSHKey(sshKey string) (string, string, string, error) {
	parts := strings.SplitN(sshKey, " ", 3)
	if len(parts) < 2 {
		return "", "", "", errors.Errorf("invalid sshKey")
	}
	sshKeyType := parts[0]
	sshKeyBlob := parts[1]
	sshKeyComment := ""
	if len(parts) == 3 {
		sshKeyComment = parts[2]
	}
	return sshKeyType, sshKeyBlob, sshKeyComment, nil
}
