package cloud

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	secretsapi "github.com/earthly/cloud-api/secrets"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/agent"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/earthly/earthly/util/hint"
)

type AuthMethod string

const (
	AuthMethodSSH       AuthMethod = "ssh"
	AuthMethodPassword  AuthMethod = "password"
	AuthMethodToken     AuthMethod = "token"
	AuthMethodCachedJWT AuthMethod = "cached jwt"
)

// TokenDetail contains token information
type TokenDetail struct {
	Name           string
	Write          bool
	Expiry         time.Time
	Indefinite     bool
	LastAccessedAt time.Time
}

func (c *Client) ListPublicKeys(ctx context.Context) ([]string, error) {
	status, body, err := c.doCall(ctx, "GET", "/api/v0/account/keys", withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to list public keys: %s", msg)
	}

	keys := []string{}
	for _, k := range strings.Split(string(body), "\n") {
		if k != "" {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func (c *Client) AddPublicKey(ctx context.Context, key string) error {
	key = strings.TrimSpace(key) + "\n"
	status, body, err := c.doCall(ctx, "PUT", "/api/v0/account/keys", withAuth(), withBody([]byte(key)))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to add public keys: %s", msg)
	}
	return nil
}

func (c *Client) RemovePublicKey(ctx context.Context, key string) error {
	key = strings.TrimSpace(key) + "\n"
	status, body, err := c.doCall(ctx, "DELETE", "/api/v0/account/keys", withAuth(), withBody([]byte(key)))
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to remove public keys: %s", msg)
	}
	return nil
}

func (c *Client) CreateToken(ctx context.Context, name string, write bool, expiry *time.Time, overWrite bool) (string, error) {
	name = url.QueryEscape(name)
	authToken := secretsapi.AuthToken{
		Write:        write,
		KeepExisting: !overWrite,
	}
	if expiry != nil {
		authToken.Expiry = timestamppb.New(expiry.UTC())
	}
	status, body, err := c.doCall(ctx, "PUT", "/api/v0/account/token/"+name, withAuth(), withJSONBody(&authToken))
	if err != nil {
		return "", err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		if status == http.StatusConflict {
			return "", hint.Wrap(errors.New(msg), "To overwrite the existing token, use the --overwrite flag")
		}
		return "", errors.Errorf("failed to create new token: %s", msg)
	}
	return string(body), nil
}

func (c *Client) ListTokens(ctx context.Context) ([]*TokenDetail, error) {
	status, body, err := c.doCall(ctx, "GET", "/api/v0/account/tokens", withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to list tokens: %s", msg)
	}

	var listTokensResponse secretsapi.ListAuthTokensResponse
	err = c.jum.Unmarshal(body, &listTokensResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal listTokens response")
	}

	tokenDetails := []*TokenDetail{}
	for _, token := range listTokensResponse.Tokens {
		var lastAccessedAt time.Time
		if token.LastAccessedAt != nil {
			lastAccessedAt = token.LastAccessedAt.AsTime()
		}
		tokenDetails = append(tokenDetails, &TokenDetail{
			Name:           token.Name,
			Write:          token.Write,
			Expiry:         token.Expiry.AsTime(),
			Indefinite:     token.Indefinite,
			LastAccessedAt: lastAccessedAt,
		})
	}
	return tokenDetails, nil
}

func (c *Client) RemoveToken(ctx context.Context, name string) error {
	name = url.QueryEscape(name)
	status, body, err := c.doCall(ctx, "DELETE", "/api/v0/account/token/"+name, withAuth())
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to delete token: %s", msg)
	}
	return nil
}

func (c *Client) WhoAmI(ctx context.Context) (string, AuthMethod, bool, error) {
	email, writeAccess, err := c.ping(ctx)
	if err != nil {
		return "", "", false, err
	}
	authMethod := c.lastAuthMethod
	if authMethod == "" {
		authMethod = AuthMethodCachedJWT
	}
	return email, authMethod, writeAccess, nil
}

func (c *Client) GetPublicKeys(ctx context.Context) ([]*agent.Key, error) {
	keys, err := c.sshAgent.List()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list ssh keys")
	}

	if c.forceSSHKey {
		// only return the matching SSH key
		for _, k := range keys {
			if bytes.Equal(k.Blob, c.sshKeyBlob) {
				return []*agent.Key{k}, nil
			}
		}
		return []*agent.Key{}, nil
	}

	// move most recently used key to the front
	sort.Slice(keys, func(i, j int) bool { return bytes.Equal(keys[i].Blob, c.sshKeyBlob) })
	return keys, nil
}

func (c *Client) RegisterEmail(ctx context.Context, email string) error {
	status, body, err := c.doCall(ctx, "PUT", fmt.Sprintf("/api/v0/account/create/%s", url.QueryEscape(email)))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to create account registration request: %s", msg)
	}
	return nil
}

func (c *Client) CreateAccount(ctx context.Context, email, verificationToken, password, publicKey string, termsConditionsPrivacy bool) error {
	if !IsValidEmail(ctx, email) {
		return errors.Errorf("invalid email: %q", email)
	}
	if publicKey != "" {
		var err error
		_, _, _, err = parseSSHKey(publicKey)
		if err != nil {
			return err
		}
	}
	createAccountRequest := secretsapi.CreateAccountRequest{
		Email:                 email,
		VerificationToken:     verificationToken,
		PublicKey:             publicKey,
		Password:              password,
		AcceptTermsConditions: termsConditionsPrivacy,
		AcceptPrivacyPolicy:   termsConditionsPrivacy,
	}
	status, body, err := c.doCall(ctx, "PUT", "/api/v0/account/create", withJSONBody(&createAccountRequest))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to create account: %s", msg)
	}

	// cache login preferences for future command runs
	if publicKey != "" {
		err = c.saveSSHCredentials(ctx, email, publicKey)
		if err != nil {
			c.warnFunc("failed to cache public ssh key: %s", err.Error())
		}
	} else {
		err = c.savePasswordCredentials(ctx, email, password)
		if err != nil {
			c.warnFunc("failed to cache password token: %s", err.Error())
		}
	}

	return nil
}

// ping calls the ping endpoint on the server,
// which is used to both test an auth token and retrieve the associated email address.
func (c *Client) ping(ctx context.Context) (email string, writeAccess bool, err error) {
	status, body, err := c.doCall(ctx, "GET", "/api/v0/account/ping", withAuth())
	if err != nil {
		return "", false, errors.Wrap(err, "failed executing ping request")
	}
	if status == http.StatusUnauthorized {
		return "", false, ErrUnauthorized
	}
	if status != http.StatusOK {
		return "", false, errors.Errorf("unexpected status code from ping: %d", status)
	}
	var resp secretsapi.PingResponse
	err = c.jum.Unmarshal(body, &resp)
	if err != nil {
		return "", false, errors.Wrap(err, "failed to unmarshal challenge response")
	}
	return resp.Email, resp.WriteAccess, nil
}

func (c *Client) AccountResetRequestToken(ctx context.Context, email string) error {
	email = url.QueryEscape(email)
	status, _, err := c.doCall(ctx, "PUT", "/api/v0/account/reset/"+email)
	if err != nil {
		return errors.Wrap(err, "failed executing account reset token request")
	}
	if status != http.StatusCreated {
		return errors.Errorf("unexpected status code from account reset token request: %d", status)
	}
	return nil
}

func (c *Client) AccountReset(ctx context.Context, email, token, password string) error {
	createAccountRequest := secretsapi.ResetPasswordRequest{
		Email:             email,
		VerificationToken: token,
		Password:          password,
	}
	status, _, err := c.doCall(ctx, "PUT", "/api/v0/account/reset", withJSONBody(&createAccountRequest))
	if err != nil {
		return errors.Wrap(err, "failed executing account reset request")
	}
	if status != http.StatusOK {
		return errors.Errorf("unexpected status code from account reset request: %d", status)
	}
	return nil
}
