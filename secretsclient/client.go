package secretsclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/earthly/earthly/secretsclient/api"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/fileutil"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/agent"
)

// ErrAccountExists occurs account creation when an account already exists
var ErrAccountExists = errors.Errorf("account already exists")

// ErrUnauthorized occurs when a user is unauthorized to access a resource
var ErrUnauthorized = errors.Errorf("unauthorized")

// ErrNoAuthorizedPublicKeys occurs when no authorized public keys are found
var ErrNoAuthorizedPublicKeys = errors.Errorf("no authorized public keys found")

// OrgDetail contains an organization and details
type OrgDetail struct {
	Name  string
	Admin bool
}

// OrgPermissions contains permission details within an org
type OrgPermissions struct {
	User  string
	Path  string
	Write bool
}

// TokenDetail contains token information
type TokenDetail struct {
	Name   string
	Write  bool
	Expiry time.Time
}

// Client provides a client to the shared secrets service
type Client interface {
	RegisterEmail(email string) error
	CreateAccount(email, verificationToken, password, publicKey string, termsConditionsPrivacy bool) error
	Get(path string) ([]byte, error)
	Remove(path string) error
	Set(path string, data []byte) error
	List(path string) ([]string, error)
	GetPublicKeys() ([]*agent.Key, error)
	CreateOrg(org string) error
	Invite(org, user string, write bool) error
	ListOrgs() ([]*OrgDetail, error)
	ListOrgPermissions(path string) ([]*OrgPermissions, error)
	RevokePermission(path, user string) error
	ListPublicKeys() ([]string, error)
	AddPublickKey(string) error
	RemovePublickKey(string) error
	CreateToken(string, bool, *time.Time) (string, error)
	ListTokens() ([]*TokenDetail, error)
	RemoveToken(string) error
	WhoAmI() (string, string, bool, error)
	FindSSHAuth() (map[string][]string, error)
	SetLoginCredentials(string, string) error
	SetLoginToken(token string) (string, error)
	SetLoginSSH(email, sshKey string) error
	DeleteCachedCredentials() error
	DisableSSHKeyGuessing()
	SetAuthTokenDir(path string)
	RedeemOAuthToken(oauthToken string) (string, error)
}

type request struct {
	hasBody bool
	body    []byte

	hasAuth bool
	retry   bool
}
type requestOpt func(*request) error

func withAuth() requestOpt {
	return func(r *request) error {
		r.hasAuth = true
		return nil
	}
}

func withRetry() requestOpt {
	return func(r *request) error {
		r.retry = true
		return nil
	}
}

func withJSONBody(body proto.Message) requestOpt {
	return func(r *request) error {
		marshaler := jsonpb.Marshaler{}
		encodedBody, err := marshaler.MarshalToString(body)
		if err != nil {
			return err
		}

		r.hasBody = true
		r.body = []byte(encodedBody)
		return nil
	}
}

func withBody(body string) requestOpt {
	return func(r *request) error {
		r.hasBody = true
		r.body = []byte(body)
		return nil
	}
}

const maxAttempt = 10
const maxSleepBeforeRetry = time.Second * 3

func (c *client) doCall(method, url string, opts ...requestOpt) (int, string, error) {
	var r request
	for _, opt := range opts {
		err := opt(&r)
		if err != nil {
			return 0, "", err
		}
	}

	var status int
	var body string
	var err error
	duration := time.Millisecond * 100
	for attempt := 0; attempt < maxAttempt; attempt++ {
		status, body, err = c.doCallImp(r, method, url, opts...)
		if (err == nil && status < 500) || errors.Cause(err) == ErrNoAuthorizedPublicKeys || errors.Cause(err) == ErrNoSSHAgent ||
			(err != nil && strings.Contains(err.Error(), "failed to connect to ssh-agent")) {
			return status, body, err
		}
		if err != nil {
			c.warnFunc("retrying http request due to %v", err)
		} else {
			msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
			if err == nil {
				c.warnFunc("retrying http request due to unexpected status code %v: %v", status, msg)
			} else {
				c.warnFunc("retrying http request due to unexpected status code %v", status)
			}
		}
		if duration > maxSleepBeforeRetry {
			duration = maxSleepBeforeRetry
		}
		time.Sleep(duration)
		duration *= 2
	}
	return status, body, err
}

func (c *client) doCallImp(r request, method, url string, opts ...requestOpt) (int, string, error) {
	var bodyReader io.Reader
	var bodyLen int64
	if r.hasBody {
		bodyReader = bytes.NewReader(r.body)
		bodyLen = int64(len(r.body))
	}

	req, err := http.NewRequest(method, c.secretServer+url, bodyReader)
	if err != nil {
		return 0, "", err
	}
	if bodyReader != nil {
		req.ContentLength = bodyLen
	}
	if r.hasAuth {
		authToken, err := c.getAuthToken()
		if err != nil {
			return 0, "", err
		}
		req.Header.Add("Authorization", authToken)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}
	return resp.StatusCode, string(respBody), nil
}

type client struct {
	secretServer          string
	sshKeyBlob            []byte // sshKey to use
	forceSSHKey           bool   // if true only use the above ssh key, don't attempt to guess others
	sshAgent              agent.ExtendedAgent
	warnFunc              func(string, ...interface{})
	email                 string
	password              string
	authToken             string
	authTokenDir          string
	disableSSHKeyGuessing bool
	jm                    *jsonpb.Unmarshaler
}

// NewClient provides a new client
func NewClient(secretServer, agentSockPath, authTokenOverride string, warnFunc func(string, ...interface{})) (Client, error) {
	c := &client{
		secretServer: secretServer,
		sshAgent: &lazySSHAgent{
			sockPath: agentSockPath,
		},
		warnFunc: warnFunc,
		jm: &jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		},
	}
	if authTokenOverride != "" {
		c.authToken = authTokenOverride
	} else {
		err := c.loadAuthToken()
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *client) filterKeys(keys []*agent.Key) []*agent.Key {
	if len(c.sshKeyBlob) == 0 && !c.forceSSHKey {
		return keys
	}

	keys2 := []*agent.Key{}
	for _, k := range keys {
		if bytes.Equal(k.Blob, c.sshKeyBlob) {
			keys2 = append(keys2, k)
		}
	}
	return keys2
}

func (c *client) GetPublicKeys() ([]*agent.Key, error) {
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

func (c *client) RegisterEmail(email string) error {
	status, body, err := c.doCall("PUT", fmt.Sprintf("/api/v0/account/create/%s", url.QueryEscape(email)))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to create account registration request: %s", msg)
	}
	return nil
}

func getMessageFromJSON(r io.Reader) (string, error) {
	decoder := json.NewDecoder(r)
	msg := struct {
		Message string `json:"message"`
	}{}
	err := decoder.Decode(&msg)
	if err != nil {
		return "", err
	}
	return msg.Message, nil
}

func (c *client) getLastUsedPublicKey() (string, error) {
	data, err := os.ReadFile(path.Join(os.TempDir(), "last-used-public-key"))
	if err != nil {
		return "", errors.Wrap(err, "failed to read file")
	}
	return string(data), nil
}

func (c *client) savePublicKey(publicKey string) error {
	f, err := os.Create(path.Join(os.TempDir(), "last-used-public-key"))
	if err != nil {
		return errors.Wrap(err, "failed to create path")
	}
	defer f.Close()
	_, err = f.WriteString(publicKey)
	if err != nil {
		return errors.Wrap(err, "failed to write public key")
	}
	return nil
}

func (c *client) CreateAccount(email, verificationToken, password, publicKey string, termsConditionsPrivacy bool) error {
	if !IsValidEmail(email) {
		return errors.Errorf("invalid email: %q", email)
	}
	if publicKey != "" {
		var err error
		_, _, _, err = parseSSHKey(publicKey)
		if err != nil {
			return err
		}
	}
	createAccountRequest := api.CreateAccountRequest{
		Email:                 email,
		VerificationToken:     verificationToken,
		PublicKey:             publicKey,
		Password:              password,
		AcceptTermsConditions: termsConditionsPrivacy,
		AcceptPrivacyPolicy:   termsConditionsPrivacy,
	}
	status, body, err := c.doCall("PUT", "/api/v0/account/create", withJSONBody(&createAccountRequest))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to create account: %s", msg)
	}

	// cache login preferences for future command runs
	if publicKey != "" {
		err = c.saveSSHToken(email, publicKey)
		if err != nil {
			c.warnFunc("failed to cache public ssh key: %s", err.Error())
		}
	} else {
		err = c.savePasswordToken(email, password)
		if err != nil {
			c.warnFunc("failed to cache password token: %s", err.Error())
		}
	}

	return nil
}

func (c *client) getChallenge() (string, error) {
	status, body, err := c.doCall("GET", "/api/v0/account/auth-challenge")
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return "", errors.Errorf("failed to get auth challenge: %s", msg)
	}

	var challengeResponse api.AuthChallengeResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &challengeResponse)
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

func (c *client) tryAuth(challenge string, key *agent.Key) (string, string, error) {
	client := &http.Client{}

	sig, err := c.signChallenge(challenge, key)
	if err != nil {
		return "", "", err
	}

	blob := base64.StdEncoding.EncodeToString(key.Blob)
	authToken := fmt.Sprintf("ssh-rsa %s %s", blob, sig)

	url := fmt.Sprintf("%s/api/v0/account/ping", c.secretServer)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Authorization", authToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	var pingResponse api.PingResponse
	err = c.jm.Unmarshal(resp.Body, &pingResponse)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to unmarshal challenge response")
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return "", "", ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.Errorf("failed to get secret: %s", pingResponse.Message)
	}

	return pingResponse.Email, authToken, nil
}

func getPasswordAuthToken(email, password string) string {
	email64 := base64.StdEncoding.EncodeToString([]byte(email))
	password64 := base64.StdEncoding.EncodeToString([]byte(password))
	return fmt.Sprintf("password %s %s", email64, password64)
}

func (c *client) getAuthToken() (string, error) {
	if c.email != "" && c.password != "" {
		return getPasswordAuthToken(c.email, c.password), nil
	}
	if c.authToken != "" {
		return "token " + c.authToken, nil
	}

	if c.disableSSHKeyGuessing {
		return "", ErrNoAuthorizedPublicKeys
	}

	challenge, err := c.getChallenge()
	if err != nil {
		return "", err
	}

	keys, err := c.GetPublicKeys()
	if err != nil {
		return "", err
	}
	for _, key := range keys {
		email, authToken, err := c.tryAuth(challenge, key)
		if err == ErrUnauthorized {
			continue // try next key
		} else if err != nil {
			return "", err
		}
		c.saveSSHToken(email, key.String())
		return authToken, nil
	}
	return "", ErrNoAuthorizedPublicKeys
}

func (c *client) CreateOrg(org string) error {
	status, body, err := c.doCall("PUT", fmt.Sprintf("/api/v0/admin/organizations/%s", url.QueryEscape(org)), withAuth())
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to create org: %s", msg)
	}
	return nil
}

func (c *client) Remove(path string) error {
	if path == "" || path[0] != '/' {
		return errors.Errorf("invalid path")
	}
	status, body, err := c.doCall("DELETE", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to remove secret: %s", msg)
	}
	return nil
}

func (c *client) List(path string) ([]string, error) {
	if path != "" && !strings.HasSuffix(path, "/") {
		return nil, errors.Errorf("invalid path")
	}
	status, body, err := c.doCall("GET", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to list secrets: %s", msg)
	}
	if body == "" {
		return []string{}, nil
	}
	return strings.Split(body, "\n"), nil
}

func (c *client) Get(path string) ([]byte, error) {
	if path == "" || path[0] != '/' || strings.HasSuffix(path, "/") {
		return nil, errors.Errorf("invalid path")
	}
	status, body, err := c.doCall("GET", fmt.Sprintf("/api/v0/secrets%s", path), withAuth(), withRetry())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to get secret: %s", msg)
	}
	return []byte(body), nil
}

func (c *client) Set(path string, data []byte) error {
	if path == "" || path[0] != '/' {
		return errors.Errorf("invalid path")
	}
	status, body, err := c.doCall("PUT", fmt.Sprintf("/api/v0/secrets%s", path), withAuth(), withBody(string(data)))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to set secret: %s", msg)
	}
	return nil
}

func getOrgFromPath(path string) (string, bool) {
	if path == "" || path[0] != '/' {
		return "", false
	}

	parts := strings.SplitN(path, "/", 3)
	if len(parts) < 2 {
		return "", false
	}
	return parts[1], true
}

func (c *client) Invite(path, user string, write bool) error {
	orgName, ok := getOrgFromPath(path)
	if !ok {
		return errors.Errorf("invalid path")
	}

	permission := api.OrgPermissions{
		Path:  path,
		Email: user,
		Write: write,
	}

	status, body, err := c.doCall("PUT", fmt.Sprintf("/api/v0/admin/organizations/%s/permissions", url.QueryEscape(orgName)), withAuth(), withJSONBody(&permission))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to invite user into org: %s", msg)
	}
	return nil
}

func (c *client) RevokePermission(path, user string) error {
	orgName, ok := getOrgFromPath(path)
	if !ok {
		return errors.Errorf("invalid path")
	}

	permission := api.OrgPermissions{
		Path:  path,
		Email: user,
	}

	status, body, err := c.doCall("DELETE", fmt.Sprintf("/api/v0/admin/organizations/%s/permissions", url.QueryEscape(orgName)), withAuth(), withJSONBody(&permission))
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to revoke user from org: %s", msg)
	}
	return nil
}

func (c *client) ListOrgPermissions(path string) ([]*OrgPermissions, error) {
	orgName, ok := getOrgFromPath(path)
	if !ok {
		return nil, errors.Errorf("invalid path")
	}

	status, body, err := c.doCall("GET", fmt.Sprintf("/api/v0/admin/organizations/%s/permissions", url.QueryEscape(orgName)), withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to list org permissions: %s", msg)
	}

	var listOrgPermissionsResponse api.ListOrgPermissionsResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &listOrgPermissionsResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal list org permissions response")
	}

	res := []*OrgPermissions{}
	for _, perm := range listOrgPermissionsResponse.Permissions {
		if strings.Contains(perm.Path, path) {
			res = append(res, &OrgPermissions{
				Path:  perm.Path,
				User:  perm.Email,
				Write: perm.Write,
			})
		}
	}

	return res, nil
}

func (c *client) ListOrgs() ([]*OrgDetail, error) {
	status, body, err := c.doCall("GET", "/api/v0/admin/organizations", withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to list orgs: %s", msg)
	}

	var listOrgsResponse api.ListOrgsResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &listOrgsResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal list orgs response")
	}

	res := []*OrgDetail{}
	for _, org := range listOrgsResponse.Details {
		res = append(res, &OrgDetail{
			Name:  org.Name,
			Admin: org.Admin,
		})
	}

	return res, nil
}

func (c *client) ListPublicKeys() ([]string, error) {
	status, body, err := c.doCall("GET", "/api/v0/account/keys", withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to list public keys: %s", msg)
	}

	keys := []string{}
	for _, k := range strings.Split(body, "\n") {
		if k != "" {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func (c *client) AddPublickKey(key string) error {
	key = strings.TrimSpace(key) + "\n"
	status, body, err := c.doCall("PUT", "/api/v0/account/keys", withAuth(), withBody(key))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to add public keys: %s", msg)
	}
	return nil
}

func (c *client) RemovePublickKey(key string) error {
	key = strings.TrimSpace(key) + "\n"
	status, body, err := c.doCall("DELETE", "/api/v0/account/keys", withAuth(), withBody(key))
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to remove public keys: %s", msg)
	}
	return nil
}

func (c *client) CreateToken(name string, write bool, expiry *time.Time) (string, error) {
	name = url.QueryEscape(name)

	expiryPB, err := ptypes.TimestampProto(expiry.UTC())
	if err != nil {
		return "", errors.Wrap(err, "TimestampProto failed")
	}

	authToken := api.AuthToken{
		Write:  write,
		Expiry: expiryPB,
	}
	status, body, err := c.doCall("PUT", "/api/v0/account/token/"+name, withAuth(), withJSONBody(&authToken))
	if err != nil {
		return "", err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return "", errors.Errorf("failed to create new token: %s", msg)
	}
	return body, nil
}

func (c *client) ListTokens() ([]*TokenDetail, error) {
	status, body, err := c.doCall("GET", "/api/v0/account/tokens", withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to list tokens: %s", msg)
	}

	var listTokensResponse api.ListAuthTokensResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &listTokensResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal listTokens response")
	}

	tokenDetails := []*TokenDetail{}
	for _, token := range listTokensResponse.Tokens {
		expiry, err := ptypes.Timestamp(token.Expiry)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode expiry proto timestamp")
		}
		tokenDetails = append(tokenDetails, &TokenDetail{
			Name:   token.Name,
			Write:  token.Write,
			Expiry: expiry,
		})
	}
	return tokenDetails, nil
}

func (c *client) RemoveToken(name string) error {
	name = url.QueryEscape(name)
	status, body, err := c.doCall("DELETE", "/api/v0/account/token/"+name, withAuth())
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to delete token: %s", msg)
	}
	return nil
}

func (c *client) WhoAmI() (string, string, bool, error) {
	status, body, err := c.doCall("GET", "/api/v0/account/ping", withAuth())
	if err != nil {
		return "", "", false, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return "", "", false, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return "", "", false, errors.Errorf("failed to authenticate: %s", msg)
	}

	var pingResponse api.PingResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &pingResponse)
	if err != nil {
		return "", "", false, errors.Wrap(err, "failed to unmarshal ping response")
	}

	authType := "ssh"
	if c.password != "" {
		authType = "password"
	} else if c.authToken != "" {
		authType = "token"
	}

	return pingResponse.Email, authType, pingResponse.WriteAccess, nil
}

func (c *client) getAuthTokenPath(create bool) (string, error) {
	confDirPath := c.authTokenDir
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
	tokenPath := filepath.Join(confDirPath, "auth.token")
	return tokenPath, nil
}

// loads ~/.earthly/auth.token
// which is formatted as
// <email> <type> ...
func (c *client) loadAuthToken() error {
	tokenPath, err := c.getAuthTokenPath(false)
	if err != nil {
		return err
	}
	if !fileutil.FileExists(tokenPath) {
		return nil
	}
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}
	parts := strings.SplitN(string(data), " ", 3)
	if len(parts) != 3 {
		return nil
	}
	c.email = parts[0]
	authType := parts[1]
	authData := parts[2]
	switch authType {
	case "password":
		passwordBytes, err := base64.StdEncoding.DecodeString(authData)
		if err != nil {
			return errors.Wrap(err, "base64 decode failed")
		}
		c.password = string(passwordBytes)
	case "ssh-rsa":
		c.sshKeyBlob, err = base64.StdEncoding.DecodeString(authData)
		if err != nil {
			return errors.Wrap(err, "base64 decode failed")
		}
	case "token":
		c.authToken = authData
	default:
		c.warnFunc("unable to handle cached auth type %s", authType)
	}
	return nil
}

// IsValidEmail returns true if email is valid
func IsValidEmail(email string) bool {
	if strings.Contains(email, " ") {
		return false
	}
	parts := strings.Split(email, "@")
	return len(parts) == 2
}

func (c *client) saveToken(email, tokenType, tokenValue string) error {
	tokenPath, err := c.getAuthTokenPath(true)
	if err != nil {
		return err
	}

	if !IsValidEmail(email) {
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
		return errors.Wrapf(err, "failed to store auth token")
	}
	return nil
}

func (c *client) saveSSHToken(email, sshKey string) error {
	sshKeyType, sshKeyBlob, _, err := parseSSHKey(sshKey)
	if err != nil {
		return err
	}
	return c.saveToken(email, sshKeyType, sshKeyBlob)
}

func (c *client) savePasswordToken(email, password string) error {
	password64 := base64.StdEncoding.EncodeToString([]byte(password))
	return c.saveToken(email, "password", password64)
}

func (c *client) SetLoginCredentials(email, password string) error {
	c.authToken = ""
	c.email = email
	c.password = password
	_, _, _, err := c.WhoAmI()
	if err != nil {
		return err
	}
	return c.savePasswordToken(email, password)
}

func (c *client) SetLoginToken(token string) (string, error) {
	c.email = ""
	c.password = ""
	c.authToken = token
	email, _, _, err := c.WhoAmI()
	if err != nil {
		return "", err
	}
	err = c.saveToken(email, "token", token)
	if err != nil {
		return "", err
	}
	return email, nil
}

func (c *client) SetLoginPublicKey(email, key string) (string, error) {
	c.password = ""
	c.authToken = ""
	c.email = email

	sshKeyType, sshKeyBlob, _, err := parseSSHKey(key)
	if err != nil {
		return "", err
	}
	if sshKeyType != "ssh-rsa" {
		return "", errors.Errorf("ssh-rsa only supported")
	}
	c.sshKeyBlob, err = base64.StdEncoding.DecodeString(sshKeyBlob)
	if err != nil {
		return "", errors.Wrap(err, "base64 decode failed")
	}

	returnedEmail, _, _, err := c.WhoAmI()
	if err != nil {
		return "", err
	}
	if returnedEmail != email {
		return "", errors.Errorf("login email missmatch")
	}
	return email, nil
}

func (c *client) DisableSSHKeyGuessing() {
	c.disableSSHKeyGuessing = true
}

func (c *client) SetAuthTokenDir(path string) {
	c.authTokenDir = path
}

func (c *client) DeleteCachedCredentials() error {
	c.email = ""
	c.password = ""
	c.authToken = ""
	tokenPath, err := c.getAuthTokenPath(false)
	if err != nil {
		return err
	}
	if !fileutil.FileExists(tokenPath) {
		return nil
	}
	err = os.Remove(tokenPath)
	if err != nil {
		return errors.Wrapf(err, "failed to delete %s", tokenPath)
	}
	return nil
}

func (c *client) SetLoginSSH(email, sshKey string) error {
	sshKeyType, sshKeyBlob, _, err := parseSSHKey(sshKey)
	if err != nil {
		return err
	}

	c.password = ""
	c.authToken = ""
	c.email = email

	c.sshKeyBlob, err = base64.StdEncoding.DecodeString(sshKeyBlob)
	if err != nil {
		return errors.Wrap(err, "base64 decode failed")
	}

	authedEmail, _, _, err := c.WhoAmI()
	if err != nil {
		return err
	}
	if authedEmail != email {
		return errors.Errorf("failed to set correct email") // shouldn't happen
	}
	return c.saveToken(email, sshKeyType, sshKeyBlob)
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

func (c *client) FindSSHAuth() (map[string][]string, error) {
	keys, err := c.GetPublicKeys()
	if err != nil {
		return nil, err
	}

	challenge, err := c.getChallenge()
	if err != nil {
		return nil, err
	}

	foundKeys := map[string][]string{}

	for _, key := range keys {
		email, _, err := c.tryAuth(challenge, key)
		if err == ErrUnauthorized {
			continue // try next key
		} else if err != nil {
			return nil, err
		}
		foundKeys[email] = append(foundKeys[email], key.String())
	}
	return foundKeys, nil
}

func (c *client) RedeemOAuthToken(oauthToken string) (string, error) {
	connectRequest := api.OAuthConnectRequest{
		Token: oauthToken,
	}
	status, body, err := c.doCall("PUT", "/api/v0/oauth/connect", withJSONBody(&connectRequest))
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return "", errors.Errorf("failed to connect OAuth login: %s", msg)
	}

	var connectResponse api.OAuthConnectResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &connectResponse)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal OAuth connect response")
	}

	email, err := c.SetLoginToken(connectResponse.Token)
	if err != nil {
		return "", errors.Wrap(err, "failed to login with OAuth generated token")
	}

	return email, nil
}
