package secretsclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/earthly/earthly/secretsclient/api"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/agent"
)

// ErrAccountExists occurs account creation when an account already exists
var ErrAccountExists = fmt.Errorf("account already exists")

// ErrUnauthorized occurs when a user is unauthorized to access a resource
var ErrUnauthorized = fmt.Errorf("unauthorized")

// ErrNoAuthorizedPublicKeys occurs when no authorized public keys are found
var ErrNoAuthorizedPublicKeys = fmt.Errorf("no authorized public keys found")

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

// Client provides a client to the shared secrets service
type Client interface {
	RegisterEmail(email string) error
	CreateAccount(email, verificationToken, password, publicKey string) error
	Get(path string) ([]byte, error)
	Remove(path string) error
	Set(path string, data []byte) error
	List(path string) ([]string, error)
	GetPublicKeys() ([]*agent.Key, error)
	CreateOrg(org string) error
	Invite(org, user string, write bool) error
	ListOrgs() ([]*OrgDetail, error)
	ListOrgPermissions(path string) ([]*OrgPermissions, error)
}

type request struct {
	hasBody bool
	body    []byte

	hasAuth bool
}
type requestOpt func(*request) error

func withPublicKeyAuth() requestOpt {
	return func(r *request) error {
		r.hasAuth = true
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

func (c *client) doCall(method, url string, opts ...requestOpt) (int, string, error) {
	var r request
	for _, opt := range opts {
		err := opt(&r)
		if err != nil {
			return 0, "", err
		}
	}
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
		_, authToken, err := c.getAuthToken()
		if err != nil {
			return 0, "", err
		}
		req.Header.Add("Authorization", authToken)
	}

	client := &http.Client{}

	resp, err := client.Do(req) // TODO add in auto-retry logic for any 500 errors
	if err != nil {
		return 0, "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}
	return resp.StatusCode, string(respBody), nil
}

type client struct {
	secretServer          string
	lastUsedPublicKeyPath string
	sshAgent              agent.ExtendedAgent
}

// NewClient provides a new client
func NewClient(secretServer, agentSockPath string) Client {
	return &client{
		secretServer: secretServer,
		sshAgent: &lazySSHAgent{
			sockPath: agentSockPath,
		},
	}
}

func (c *client) GetPublicKeys() ([]*agent.Key, error) {
	keys, err := c.sshAgent.List()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list ssh keys")
	}

	key, err := c.getLastUsedPublicKey()
	if err != nil {
		// ignore error and return existing list from sshAgent
		return keys, nil
	}
	// otherwise, move last used key to the front of the list
	// to make it more likely to get a valid auth token
	keys2 := []*agent.Key{}
	for _, k := range keys {
		if k.String() == key {
			keys2 = append(keys2, k)
		}
	}
	for _, k := range keys {
		if k.String() != key {
			keys2 = append(keys2, k)
		}
	}
	return keys2, nil
}

func (c *client) RegisterEmail(email string) error {
	status, body, err := c.doCall("PUT", fmt.Sprintf("/api/v0/account/create/%s", email))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return fmt.Errorf("failed to create account registration request: %s", msg)
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
	data, err := ioutil.ReadFile("/tmp/last-used-public-key")
	if err != nil {
		return "", errors.Wrap(err, "failed to read file")
	}
	return string(data), nil
}

func (c *client) savePublicKey(publicKey string) error {
	f, err := os.Create("/tmp/last-used-public-key")
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

func (c *client) CreateAccount(email, verificationToken, password, publicKey string) error {
	createAccountRequest := api.CreateAccountRequest{
		Email:             email,
		VerificationToken: verificationToken,
		PublicKey:         publicKey,
		Password:          password,
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
		return fmt.Errorf("failed to create account: %s", msg)
	}
	c.savePublicKey(publicKey)
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
		return "", fmt.Errorf("failed to get auth challenge: %s", msg)
	}

	var challengeResponse api.AuthChallengeResponse
	err = jsonpb.Unmarshal(bytes.NewReader([]byte(body)), &challengeResponse)
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
	err = jsonpb.Unmarshal(resp.Body, &pingResponse)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to unmarshal challenge response")
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return "", "", ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to get secret: %s", pingResponse.Message)
	}

	return pingResponse.Email, authToken, nil
}

func (c *client) getAuthToken() (*agent.Key, string, error) {
	challenge, err := c.getChallenge()
	if err != nil {
		return nil, "", err
	}

	keys, err := c.GetPublicKeys()
	if err != nil {
		return nil, "", err
	}
	for _, key := range keys {
		_, authToken, err := c.tryAuth(challenge, key)
		if err == ErrUnauthorized {
			continue // try next key
		} else if err != nil {
			return nil, "", err
		}
		c.savePublicKey(key.String())
		return key, authToken, nil
	}
	return nil, "", ErrNoAuthorizedPublicKeys
}

func (c *client) CreateOrg(org string) error {
	status, body, err := c.doCall("PUT", fmt.Sprintf("/api/v0/admin/organizations/%s", org), withPublicKeyAuth())
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return fmt.Errorf("failed to create org: %s", msg)
	}
	return nil
}

func (c *client) Remove(path string) error {
	if path == "" || path[0] != '/' {
		return fmt.Errorf("invalid path")
	}
	status, body, err := c.doCall("DELETE", fmt.Sprintf("/api/v0/secrets%s", path), withPublicKeyAuth())
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return fmt.Errorf("failed to remove secret: %s", msg)
	}
	return nil
}

func (c *client) List(path string) ([]string, error) {
	if path != "" && !strings.HasSuffix(path, "/") {
		return nil, fmt.Errorf("invalid path")
	}
	status, body, err := c.doCall("GET", fmt.Sprintf("/api/v0/secrets%s", path), withPublicKeyAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, fmt.Errorf("failed to list secrets: %s", msg)
	}
	if body == "" {
		return []string{}, nil
	}
	return strings.Split(body, "\n"), nil
}

func (c *client) Get(path string) ([]byte, error) {
	if path == "" || path[0] != '/' || strings.HasSuffix(path, "/") {
		return nil, fmt.Errorf("invalid path")
	}
	status, body, err := c.doCall("GET", fmt.Sprintf("/api/v0/secrets%s", path), withPublicKeyAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, fmt.Errorf("failed to get secret: %s", msg)
	}
	return []byte(body), nil
}

func (c *client) Set(path string, data []byte) error {
	if path == "" || path[0] != '/' {
		return fmt.Errorf("invalid path")
	}
	status, body, err := c.doCall("PUT", fmt.Sprintf("/api/v0/secrets%s", path), withPublicKeyAuth(), withBody(string(data)))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return fmt.Errorf("failed to set secret: %s", msg)
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
		return fmt.Errorf("invalid path")
	}

	permission := api.OrgPermissions{
		Path:  path,
		Email: user,
		Write: write,
	}

	status, body, err := c.doCall("PUT", fmt.Sprintf("/api/v0/admin/organizations/%s/permissions", orgName), withPublicKeyAuth(), withJSONBody(&permission))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return fmt.Errorf("failed to invite user into org: %s", msg)
	}
	return nil
}

func (c *client) ListOrgPermissions(path string) ([]*OrgPermissions, error) {
	orgName, ok := getOrgFromPath(path)
	if !ok {
		return nil, fmt.Errorf("invalid path")
	}

	status, body, err := c.doCall("GET", fmt.Sprintf("/api/v0/admin/organizations/%s/permissions", orgName), withPublicKeyAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, fmt.Errorf("failed to invite user into org: %s", msg)
	}

	var listOrgPermissionsResponse api.ListOrgPermissionsResponse
	err = jsonpb.Unmarshal(bytes.NewReader([]byte(body)), &listOrgPermissionsResponse)
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
	status, body, err := c.doCall("GET", "/api/v0/admin/organizations", withPublicKeyAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, fmt.Errorf("failed to invite user into org: %s", msg)
	}

	var listOrgsResponse api.ListOrgsResponse
	err = jsonpb.Unmarshal(bytes.NewReader([]byte(body)), &listOrgsResponse)
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
