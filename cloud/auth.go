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
func (c *Client) Authenticate(ctx context.Context) (AuthMethod, error) {
	authMethod, err := c.doLogin(ctx)
	if err != nil {
		if errors.Is(err, ErrNoAuthorizedPublicKeys) || errors.Is(err, ErrNoSSHAgent) {
			return "", ErrUnauthorized
		}
		return "", err
	}
	err = c.saveToken()
	if err != nil {
		return "", err
	}
	c.lastAuthMethod = authMethod
	return authMethod, nil
}

func (c *Client) doLogin(ctx context.Context) (AuthMethod, error) {
	if c.email != "" && c.password != "" {
		return AuthMethodPassword, c.loginWithPassword(ctx)
	}
	if c.authCredToken != "" {
		return AuthMethodToken, c.loginWithToken(ctx)
	}
	return AuthMethodSSH, c.loginWithSSH(ctx)
}

func (c *Client) IsLoggedIn(ctx context.Context) bool {
	return c.authToken != "" || c.authCredToken != ""
}

func (c *Client) FindSSHCredentials(ctx context.Context, emailToFind string) error {
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

func (c *Client) GetAuthToken(ctx context.Context) (string, error) {
	_, err := c.Authenticate(ctx) // Ensure the current token is valid
	if err != nil {
		return "", errors.Wrap(err, "could not authenticate")
	}
	return c.authToken, nil
}

func (c *Client) SetSSHCredentials(ctx context.Context, email, sshKey string) error {
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

func (c *Client) DeleteCachedToken(ctx context.Context) error {
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

func (c *Client) DeleteAuthCache(ctx context.Context) error {
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

func (c *Client) DisableAutoLogin(ctx context.Context) error {
	path, err := c.getDisableAutoLoginPath(true)
	if err != nil {
		return errors.Wrapf(err, "failed to get disable auto login path")
	}
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s", path)
	}
	err = f.Close()
	if err != nil {
		return errors.Wrapf(err, "failed to close %s", path)
	}
	return nil
}

func (c *Client) IsAutoLoginPermitted(ctx context.Context) (bool, error) {
	path, err := c.getDisableAutoLoginPath(true)
	if err != nil {
		return false, errors.Wrapf(err, "failed to get disable auto login path")
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil // auto login is allowed when the file does not exist
		}
		return false, err
	}
	if info.IsDir() {
		return false, fmt.Errorf("expected %s to be a file (not a directory)", path)
	}
	return false, nil
}

func (c *Client) EnableAutoLogin(ctx context.Context) error {
	canAutoLogin, err := c.IsAutoLoginPermitted(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to determine if auto login is permitted")
	}
	if canAutoLogin {
		return nil // already allowed
	}
	path, err := c.getDisableAutoLoginPath(true)
	if err != nil {
		return errors.Wrapf(err, "failed to get disable auto login path")
	}
	err = os.Remove(path)
	if err != nil {
		return errors.Wrapf(err, "failed to remove %s", path)
	}
	return nil
}

func (c *Client) SetTokenCredentials(ctx context.Context, token string) (string, error) {
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

func (c *Client) DisableSSHKeyGuessing(ctx context.Context) {
	c.disableSSHKeyGuessing = true
}

func (c *Client) SetAuthTokenDir(ctx context.Context, path string) {
	c.authDir = path
}

func (c *Client) SetPasswordCredentials(ctx context.Context, email, password string) error {
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

func (c *Client) loadCredentials() error {
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
	if credType == "password" {
		passwordBytes, err := base64.StdEncoding.DecodeString(credData)
		if err != nil {
			return errors.Wrap(err, "base64 decode failed")
		}
		c.password = string(passwordBytes)
	} else if credType == "token" {
		c.authCredToken = credData
	} else if strings.HasPrefix(credType, "ssh-") {
		c.sshKeyBlob, err = base64.StdEncoding.DecodeString(credData)
		if err != nil {
			return errors.Wrap(err, "base64 decode failed")
		}
	} else {
		c.warnFunc("unable to handle cached auth type %s", credType)
	}
	return nil
}

func (c *Client) loadToken() error {
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

func (c *Client) migrateOldToken() error {
	confDirPath := c.authDir
	if confDirPath == "" {
		confDirPath = cliutil.GetEarthlyDir(c.installationName)
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

func (c *Client) getCredentialsPath(create bool) (string, error) {
	return c.getAuthPath("auth.credentials", create)
}

func (c *Client) getTokenPath(create bool) (string, error) {
	return c.getAuthPath("auth.jwt", create)
}

func (c *Client) getDisableAutoLoginPath(create bool) (string, error) {
	return c.getAuthPath("do-not-login-automatically", create)
}

func (c *Client) getAuthPath(filename string, createEarthlyDir bool) (string, error) {
	confDirPath := c.authDir
	if confDirPath == "" {
		if createEarthlyDir {
			var err error
			confDirPath, err = cliutil.GetOrCreateEarthlyDir(c.installationName)
			if err != nil {
				return "", errors.Wrap(err, "cannot get .earthly dir")
			}
		} else {
			confDirPath = cliutil.GetEarthlyDir(c.installationName)
		}
	}
	tokenPath := filepath.Join(confDirPath, filename)
	return tokenPath, nil
}

func (c *Client) saveCredentials(ctx context.Context, email, tokenType, tokenValue string) error {
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

func (c *Client) saveSSHCredentials(ctx context.Context, email, sshKey string) error {
	sshKeyType, sshKeyBlob, _, err := parseSSHKey(sshKey)
	if err != nil {
		return err
	}
	return c.saveCredentials(ctx, email, sshKeyType, sshKeyBlob)
}

func (c *Client) savePasswordCredentials(ctx context.Context, email, password string) error {
	password64 := base64.StdEncoding.EncodeToString([]byte(password))
	return c.saveCredentials(ctx, email, "password", password64)
}

func (c *Client) deleteCachedCredentials() error {
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
func (c *Client) loadAuthStorage() error {
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

func (c *Client) saveToken() error {
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

func (c *Client) loginWithPassword(ctx context.Context) error {
	var err error
	c.authCredToken = getPasswordAuthToken(c.email, c.password)
	c.authToken, c.authTokenExpiry, err = c.login(ctx, c.authCredToken)
	return err
}

func (c *Client) loginWithToken(ctx context.Context) error {
	var err error
	c.authToken, c.authTokenExpiry, err = c.login(ctx, "token "+c.authCredToken)
	return err
}

func (c *Client) loginWithSSH(ctx context.Context) error {
	allowAutoLogin, err := c.IsAutoLoginPermitted(ctx)
	if err != nil {
		return err
	}
	if c.disableSSHKeyGuessing || !allowAutoLogin {
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
// ErrUnauthorized is returned if the credentials are not valid.
func (c *Client) login(ctx context.Context, credentials string) (token string, expiry time.Time, err error) {
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

func (c *Client) getChallenge(ctx context.Context) (string, error) {
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

func (c *Client) signChallenge(challenge string, key *agent.Key) (string, string, error) {
	sig, err := c.sshAgent.SignWithFlags(key, []byte(challenge), agent.SignatureFlagRsaSha512)
	if err != nil {
		return "", "", err
	}
	s := base64.StdEncoding.EncodeToString(sig.Blob)
	return sig.Format, s, nil
}

func (c *Client) getSSHCredentials(challenge string, key *agent.Key) (credentials string, err error) {
	sigFormat, sig, err := c.signChallenge(challenge, key)
	if err != nil {
		return credentials, err
	}
	blob := base64.StdEncoding.EncodeToString(key.Blob)
	credentials = fmt.Sprintf("%s %s %s", sigFormat, blob, sig)
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
