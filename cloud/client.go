package cloud

import (
	"bytes"
	"context"
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

	logsapi "github.com/earthly/cloud-api/logs"
	pipelinesapi "github.com/earthly/cloud-api/pipelines"
	secretsapi "github.com/earthly/cloud-api/secrets"

	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/fileutil"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/agent"
)

// ErrUnauthorized occurs when a user is unauthorized to access a resource
var ErrUnauthorized = errors.New("unauthorized")

// ErrNoAuthorizedPublicKeys occurs when no authorized public keys are found
var ErrNoAuthorizedPublicKeys = errors.New("no authorized public keys found")

const (
	tokenExpiryLayout    = "2006-01-02 15:04:05.999999999 -0700 MST"
	satelliteMgmtTimeout = "5M" // 5 minute timeout when launching or deleting a Satellite
)

// OrgDetail contains an organization and details
type OrgDetail struct {
	ID    string
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

// SatelliteInstance contains details about a remote Buildkit instance.
type SatelliteInstance struct {
	Name     string
	Org      string
	Status   string
	Version  string
	Platform string
}

// Client provides a client to the shared secrets service
type Client interface {
	RegisterEmail(ctx context.Context, email string) error
	CreateAccount(ctx context.Context, email, verificationToken, password, publicKey string, termsConditionsPrivacy bool) error
	Authenticate(ctx context.Context) error
	Get(ctx context.Context, path string) ([]byte, error)
	Remove(ctx context.Context, path string) error
	Set(ctx context.Context, path string, data []byte) error
	List(ctx context.Context, path string) ([]string, error)
	GetPublicKeys(ctx context.Context) ([]*agent.Key, error)
	CreateOrg(ctx context.Context, org string) error
	Invite(ctx context.Context, org, user string, write bool) error
	ListOrgs(ctx context.Context) ([]*OrgDetail, error)
	ListOrgPermissions(ctx context.Context, path string) ([]*OrgPermissions, error)
	RevokePermission(ctx context.Context, path, user string) error
	ListPublicKeys(ctx context.Context) ([]string, error)
	AddPublickKey(ctx context.Context, key string) error
	RemovePublickKey(ctx context.Context, key string) error
	CreateToken(context.Context, string, bool, *time.Time) (string, error)
	ListTokens(ctx context.Context) ([]*TokenDetail, error)
	RemoveToken(ctx context.Context, token string) error
	WhoAmI(ctx context.Context) (string, string, bool, error)
	UploadLog(ctx context.Context, pathOnDisk string) (string, error)
	SetPasswordCredentials(context.Context, string, string) error
	SetTokenCredentials(ctx context.Context, token string) (string, error)
	SetSSHCredentials(ctx context.Context, email, sshKey string) error
	FindSSHCredentials(ctx context.Context, emailToFind string) error
	DeleteAuthCache(ctx context.Context) error
	DeleteCachedToken(ctx context.Context) error
	DisableSSHKeyGuessing(ctx context.Context)
	SetAuthTokenDir(ctx context.Context, path string)
	SendAnalytics(ctx context.Context, data *EarthlyAnalytics) error
	IsLoggedIn(ctx context.Context) bool
	GetAuthToken(ctx context.Context) (string, error)
	LaunchSatellite(ctx context.Context, name, org string) (*SatelliteInstance, error)
	GetOrgID(ctx context.Context, name string) (string, error)
	ListSatellites(ctx context.Context, orgID string) ([]SatelliteInstance, error)
	GetSatellite(ctx context.Context, name, orgID string) (*SatelliteInstance, error)
	DeleteSatellite(ctx context.Context, name, orgID string) error
}

type request struct {
	hasBody bool
	body    []byte
	headers http.Header

	hasAuth    bool
	hasHeaders bool
}
type requestOpt func(*request) error

func withAuth() requestOpt {
	return func(r *request) error {
		r.hasAuth = true
		return nil
	}
}

func withHeader(key, value string) requestOpt {
	return func(r *request) error {
		r.hasHeaders = true
		r.headers = http.Header{}
		r.headers.Add(key, value)
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

func withFileBody(pathOnDisk string) requestOpt {
	return func(r *request) error {
		_, err := os.Stat(pathOnDisk)
		if err != nil {
			return errors.Wrapf(err, "could not stat file at %s", pathOnDisk)
		}

		contents, err := os.ReadFile(pathOnDisk)
		if err != nil {
			return errors.Wrapf(err, "could not add file %s to request body", pathOnDisk)
		}

		r.hasBody = true
		r.body = contents
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

func (c *client) doCall(ctx context.Context, method, url string, opts ...requestOpt) (int, string, error) {
	const maxAttempt = 10
	const maxSleepBeforeRetry = time.Second * 3

	var r request
	for _, opt := range opts {
		err := opt(&r)
		if err != nil {
			return 0, "", err
		}
	}

	alreadyReAuthed := false
	if r.hasAuth && time.Now().UTC().After(c.authTokenExpiry) {
		if err := c.Authenticate(ctx); err != nil {
			if errors.Is(err, ErrUnauthorized) {
				return 0, "", ErrUnauthorized
			}
			return 0, "", errors.Wrap(err, "failed refreshing expired auth token")
		}
		alreadyReAuthed = true
	}

	var status int
	var body string
	var err error
	duration := time.Millisecond * 100
	for attempt := 0; attempt < maxAttempt; attempt++ {
		status, body, err = c.doCallImp(ctx, r, method, url, opts...)

		if !shouldRetry(status, body, err, c.warnFunc) {
			return status, body, err
		}

		if status == http.StatusUnauthorized {
			if !r.hasAuth || alreadyReAuthed {
				return status, body, ErrUnauthorized
			}
			if err = c.Authenticate(ctx); err != nil {
				return status, body, errors.Wrap(err, "auth credentials not valid")
			}
			alreadyReAuthed = true
		}

		if duration > maxSleepBeforeRetry {
			duration = maxSleepBeforeRetry
		}

		time.Sleep(duration)
		duration *= 2
	}

	return status, body, err
}

func shouldRetry(status int, body string, err error, warnFunc func(string, ...interface{})) bool {
	if status == http.StatusUnauthorized {
		return true
	}
	if 500 <= status && status <= 599 {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			warnFunc("retrying http request due to unexpected status code %v", status)
		} else {
			warnFunc("retrying http request due to unexpected status code %v: %v", status, msg)
		}
		return true
	}
	if err != nil {
		if errors.Cause(err) == ErrNoAuthorizedPublicKeys ||
			errors.Cause(err) == ErrNoSSHAgent ||
			strings.Contains(err.Error(), "failed to connect to ssh-agent") {
			return false
		}
		warnFunc("retrying http request due to unexpected error %v", err)
		return true
	}
	return false
}

func (c *client) doCallImp(ctx context.Context, r request, method, url string, opts ...requestOpt) (int, string, error) {
	var bodyReader io.Reader
	var bodyLen int64
	if r.hasBody {
		bodyReader = bytes.NewReader(r.body)
		bodyLen = int64(len(r.body))
	}

	req, err := http.NewRequestWithContext(ctx, method, c.host+url, bodyReader)
	if err != nil {
		return 0, "", err
	}
	if bodyReader != nil {
		req.ContentLength = bodyLen
	}
	if r.hasHeaders {
		req.Header = r.headers
	}
	if r.hasAuth {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}

	respBody, err := readAllWithContext(ctx, resp.Body)
	if err != nil {
		return 0, "", err
	}
	return resp.StatusCode, string(respBody), nil
}

func readAllWithContext(ctx context.Context, r io.Reader) ([]byte, error) {
	var dt []byte
	var readErr error
	ch := make(chan struct{})
	go func() {
		dt, readErr = io.ReadAll(r)
		close(ch)
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-ch:
		return dt, readErr
	}
}

type client struct {
	host                  string
	sshKeyBlob            []byte // sshKey to use
	forceSSHKey           bool   // if true only use the above ssh key, don't attempt to guess others
	sshAgent              agent.ExtendedAgent
	warnFunc              func(string, ...interface{})
	email                 string
	password              string
	authToken             string
	authTokenExpiry       time.Time
	authCredToken         string
	authDir               string
	disableSSHKeyGuessing bool
	jm                    *jsonpb.Unmarshaler
}

// NewClient provides a new Earthly Cloud client
func NewClient(host, agentSockPath, authCredsOverride string, warnFunc func(string, ...interface{})) (Client, error) {
	c := &client{
		host: host,
		sshAgent: &lazySSHAgent{
			sockPath: agentSockPath,
		},
		warnFunc: warnFunc,
		jm: &jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		},
	}
	if authCredsOverride != "" {
		c.authCredToken = authCredsOverride
	} else {
		if err := c.loadAuthStorage(); err != nil {
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

func (c *client) GetPublicKeys(ctx context.Context) ([]*agent.Key, error) {
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

func (c *client) RegisterEmail(ctx context.Context, email string) error {
	status, body, err := c.doCall(ctx, "PUT", fmt.Sprintf("/api/v0/account/create/%s", url.QueryEscape(email)))
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
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &resp)
	if err != nil {
		return "", zero, errors.Wrap(err, "failed to unmarshal login response")
	}
	return resp.Token, resp.Expiry.AsTime().UTC(), nil
}

// ping calls the ping endpoint on the server,
// which is used to both test an auth token and retrieve the associated email address.
func (c *client) ping(ctx context.Context) (email string, writeAccess bool, err error) {
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
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &resp)
	if err != nil {
		return "", false, errors.Wrap(err, "failed to unmarshal challenge response")
	}
	return resp.Email, resp.WriteAccess, nil
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

func (c *client) CreateAccount(ctx context.Context, email, verificationToken, password, publicKey string, termsConditionsPrivacy bool) error {
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
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
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

func (c *client) getChallenge(ctx context.Context) (string, error) {
	status, body, err := c.doCall(ctx, "GET", "/api/v0/account/auth-challenge")
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

	var challengeResponse secretsapi.AuthChallengeResponse
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

func (c *client) CreateOrg(ctx context.Context, org string) error {
	status, body, err := c.doCall(ctx, "PUT", fmt.Sprintf("/api/v0/admin/organizations/%s", url.QueryEscape(org)), withAuth())
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

func (c *client) Remove(ctx context.Context, path string) error {
	if path == "" || path[0] != '/' {
		return errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "DELETE", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
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

func (c *client) List(ctx context.Context, path string) ([]string, error) {
	if path != "" && !strings.HasSuffix(path, "/") {
		return nil, errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "GET", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
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

func (c *client) Get(ctx context.Context, path string) ([]byte, error) {
	if path == "" || path[0] != '/' || strings.HasSuffix(path, "/") {
		return nil, errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "GET", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
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

func (c *client) Set(ctx context.Context, path string, data []byte) error {
	if path == "" || path[0] != '/' {
		return errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "PUT", fmt.Sprintf("/api/v0/secrets%s", path), withAuth(), withBody(string(data)))
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

func (c *client) Invite(ctx context.Context, path, user string, write bool) error {
	orgName, ok := getOrgFromPath(path)
	if !ok {
		return errors.Errorf("invalid path")
	}

	permission := secretsapi.OrgPermissions{
		Path:  path,
		Email: user,
		Write: write,
	}

	status, body, err := c.doCall(ctx, "PUT", fmt.Sprintf("/api/v0/admin/organizations/%s/permissions", url.QueryEscape(orgName)), withAuth(), withJSONBody(&permission))
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

func (c *client) RevokePermission(ctx context.Context, path, user string) error {
	orgName, ok := getOrgFromPath(path)
	if !ok {
		return errors.Errorf("invalid path")
	}

	permission := secretsapi.OrgPermissions{
		Path:  path,
		Email: user,
	}

	status, body, err := c.doCall(ctx, "DELETE", fmt.Sprintf("/api/v0/admin/organizations/%s/permissions", url.QueryEscape(orgName)), withAuth(), withJSONBody(&permission))
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

func (c *client) ListOrgPermissions(ctx context.Context, path string) ([]*OrgPermissions, error) {
	orgName, ok := getOrgFromPath(path)
	if !ok {
		return nil, errors.Errorf("invalid path")
	}

	status, body, err := c.doCall(ctx, "GET", fmt.Sprintf("/api/v0/admin/organizations/%s/permissions", url.QueryEscape(orgName)), withAuth())
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

	var listOrgPermissionsResponse secretsapi.ListOrgPermissionsResponse
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

func (c *client) ListOrgs(ctx context.Context) ([]*OrgDetail, error) {
	status, body, err := c.doCall(ctx, "GET", "/api/v0/admin/organizations", withAuth())
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

	var listOrgsResponse secretsapi.ListOrgsResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &listOrgsResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal list orgs response")
	}

	res := []*OrgDetail{}
	for _, org := range listOrgsResponse.Details {
		res = append(res, &OrgDetail{
			ID:    org.Id,
			Name:  org.Name,
			Admin: org.Admin,
		})
	}

	return res, nil
}

func (c *client) ListPublicKeys(ctx context.Context) ([]string, error) {
	status, body, err := c.doCall(ctx, "GET", "/api/v0/account/keys", withAuth())
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

func (c *client) AddPublickKey(ctx context.Context, key string) error {
	key = strings.TrimSpace(key) + "\n"
	status, body, err := c.doCall(ctx, "PUT", "/api/v0/account/keys", withAuth(), withBody(key))
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

func (c *client) RemovePublickKey(ctx context.Context, key string) error {
	key = strings.TrimSpace(key) + "\n"
	status, body, err := c.doCall(ctx, "DELETE", "/api/v0/account/keys", withAuth(), withBody(key))
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

func (c *client) CreateToken(ctx context.Context, name string, write bool, expiry *time.Time) (string, error) {
	name = url.QueryEscape(name)

	expiryPB, err := ptypes.TimestampProto(expiry.UTC())
	if err != nil {
		return "", errors.Wrap(err, "TimestampProto failed")
	}

	authToken := secretsapi.AuthToken{
		Write:  write,
		Expiry: expiryPB,
	}
	status, body, err := c.doCall(ctx, "PUT", "/api/v0/account/token/"+name, withAuth(), withJSONBody(&authToken))
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

func (c *client) ListTokens(ctx context.Context) ([]*TokenDetail, error) {
	status, body, err := c.doCall(ctx, "GET", "/api/v0/account/tokens", withAuth())
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

	var listTokensResponse secretsapi.ListAuthTokensResponse
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

func (c *client) RemoveToken(ctx context.Context, name string) error {
	name = url.QueryEscape(name)
	status, body, err := c.doCall(ctx, "DELETE", "/api/v0/account/token/"+name, withAuth())
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

func (c *client) WhoAmI(ctx context.Context) (string, string, bool, error) {
	email, writeAccess, err := c.ping(ctx)
	if err != nil {
		return "", "", false, err
	}

	authType := "ssh"
	if c.password != "" {
		authType = "password"
	} else if c.authCredToken != "" {
		authType = "token"
	}

	return email, authType, writeAccess, nil
}

func (c *client) LaunchSatellite(ctx context.Context, name, orgID string) (*SatelliteInstance, error) {
	req := pipelinesapi.LaunchSatelliteRequest{
		OrgId:    orgID,
		Name:     name,
		Platform: "linux/amd64", // TODO support arm64 as well
	}
	status, body, err := c.doCall(ctx, "POST", "/api/v0/satellites",
		withAuth(), withHeader("Grpc-Timeout", satelliteMgmtTimeout), withJSONBody(&req))
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.Errorf("failed launching satellite: %s", body)
	}
	var resp pipelinesapi.LaunchSatelliteResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal LaunchSatellite response")
	}
	return &SatelliteInstance{
		Name:     name,
		Org:      orgID,
		Status:   resp.Status.String(),
		Version:  resp.Version,
		Platform: "linux/amd64",
	}, nil
}

func (c *client) GetOrgID(ctx context.Context, orgName string) (string, error) {
	orgs, err := c.ListOrgs(ctx)
	if err != nil {
		return "", err
	}
	for _, o := range orgs {
		if o.Name == orgName {
			return o.ID, nil
		}
	}
	return "", errors.Errorf("org not found: %s", orgName)
}

func (c *client) ListSatellites(ctx context.Context, orgID string) ([]SatelliteInstance, error) {
	url := fmt.Sprintf("/api/v0/satellites?orgId=%s", url.QueryEscape(orgID))
	status, body, err := c.doCall(ctx, "GET", url, withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.Errorf("failed listing satellites: %s", body)
	}
	var resp pipelinesapi.ListSatellitesResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal listTokens response")
	}
	instances := make([]SatelliteInstance, len(resp.Instances))
	for i, s := range resp.Instances {
		instances[i] = SatelliteInstance{
			Name:     s.Name,
			Org:      orgID,
			Version:  s.Version,
			Platform: s.Platform,
			Status:   s.Status.String(),
		}
	}
	return instances, nil
}

func (c *client) GetSatellite(ctx context.Context, name, orgID string) (*SatelliteInstance, error) {
	url := fmt.Sprintf("/api/v0/satellites/%s?orgId=%s", name, url.QueryEscape(orgID))
	status, body, err := c.doCall(ctx, "GET", url, withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.Errorf("failed listing satellites: %s", body)
	}
	var resp pipelinesapi.GetSatelliteResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal listTokens response")
	}
	return &SatelliteInstance{
		Name:     name,
		Org:      orgID,
		Status:   resp.Status.String(),
		Version:  resp.Version,
		Platform: resp.Platform,
	}, nil
}

func (c *client) DeleteSatellite(ctx context.Context, name, orgID string) error {
	url := fmt.Sprintf("/api/v0/satellites/%s?orgId=%s", name, url.QueryEscape(orgID))
	status, body, err := c.doCall(ctx, "DELETE", url,
		withAuth(), withHeader("Grpc-Timeout", satelliteMgmtTimeout))
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return errors.Errorf("failed listing satellites: %s", body)
	}
	return nil
}

// EarthlyAnalytics is the payload used in SendAnalytics.
// It contains information about the command that was run,
// the environment it was run in, and the result of the command.
type EarthlyAnalytics struct {
	Key              string                    `json:"key"`
	InstallID        string                    `json:"install_id"`
	Version          string                    `json:"version"`
	Platform         string                    `json:"platform"`
	GitSHA           string                    `json:"git_sha"`
	ExitCode         int                       `json:"exit_code"`
	CI               string                    `json:"ci_name"`
	RepoHash         string                    `json:"repo_hash"`
	ExecutionSeconds float64                   `json:"execution_seconds"`
	Terminal         bool                      `json:"terminal"`
	Counts           map[string]map[string]int `json:"counts"`
}

// SendAnalytics send an analytics event to the Cloud server.
func (c *client) SendAnalytics(ctx context.Context, data *EarthlyAnalytics) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}
	opts := []requestOpt{
		withBody(string(payload)),
		withHeader("Content-Type", "application/json; charset=utf-8"),
	}
	if c.IsLoggedIn(ctx) {
		opts = append(opts, withAuth())
	}
	status, _, err := c.doCall(ctx, "PUT", "/analytics", opts...)
	if err != nil {
		return errors.Wrap(err, "failed sending analytics")
	}
	if status != http.StatusCreated {
		return errors.Errorf("unexpected response from analytics server: %d", status)
	}
	return nil
}

func (c *client) IsLoggedIn(ctx context.Context) bool {
	return c.authToken != "" || c.authCredToken != ""
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

// loads the following files:
//  * ~/.earthly/auth.credentials
//  * ~/.earthly/auth.jwt
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

// IsValidEmail returns true if email is valid
func IsValidEmail(ctx context.Context, email string) bool {
	if strings.Contains(email, " ") {
		return false
	}
	parts := strings.Split(email, "@")
	return len(parts) == 2
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

func (c *client) UploadLog(ctx context.Context, pathOnDisk string) (string, error) {
	status, body, err := c.doCall(ctx, http.MethodPost, "/api/v0/logs", withAuth(), withFileBody(pathOnDisk), withHeader("Content-Type", "application/gzip"))
	if err != nil {
		return "", err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return "", errors.Errorf("failed to upload log: %s", msg)
	}

	var uploadBundleResponse logsapi.UploadLogBundleResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &uploadBundleResponse)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal uploadbundle response")
	}

	return fmt.Sprintf(uploadBundleResponse.ViewURL), nil
}

func (c *client) GetAuthToken(ctx context.Context) (string, error) {
	err := c.Authenticate(ctx) // Ensure the current token is valid
	if err != nil {
		return "", errors.Wrap(err, "could not authenticate")
	}
	return c.authToken, nil
}
