package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/agent"
)

var (
	// ErrUnauthorized occurs when a user is unauthorized to access a resource
	ErrUnauthorized = errors.New("unauthorized")
	// ErrNoAuthorizedPublicKeys occurs when no authorized public keys are found
	ErrNoAuthorizedPublicKeys = errors.New("no authorized public keys found")
)

const (
	tokenExpiryLayout    = "2006-01-02 15:04:05.999999999 -0700 MST"
	satelliteMgmtTimeout = "5M" // 5 minute timeout when launching or deleting a Satellite
)

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
	InviteToOrg(ctx context.Context, invite *OrgInvitation) (string, error)
	ListOrgs(ctx context.Context) ([]*OrgDetail, error)
	ListOrgPermissions(ctx context.Context, path string) ([]*OrgPermissions, error)
	ListOrgMembers(ctx context.Context, orgName string) ([]*OrgMember, error)
	UpdateOrgMember(ctx context.Context, orgName, userEmail, permission string) error
	RemoveOrgMember(ctx context.Context, orgName, userEmail string) error
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
	LaunchSatellite(ctx context.Context, name, org string) error
	GetOrgID(ctx context.Context, name string) (string, error)
	ListSatellites(ctx context.Context, orgID string) ([]SatelliteInstance, error)
	GetSatellite(ctx context.Context, name, orgID string) (*SatelliteInstance, error)
	DeleteSatellite(ctx context.Context, name, orgID string) error
	CreateProject(ctx context.Context, name, orgName string) (*Project, error)
	ListProjects(ctx context.Context, orgName string) ([]*Project, error)
	GetProject(ctx context.Context, orgName, name string) (*Project, error)
	DeleteProject(ctx context.Context, orgName, name string) error
	AddProjectMember(ctx context.Context, orgName, name, userEmail, permission string) error
	UpdateProjectMember(ctx context.Context, orgName, name, userEmail, permission string) error
	ListProjectMembers(ctx context.Context, orgName, name string) ([]*ProjectMember, error)
	RemoveProjectMember(ctx context.Context, orgName, name, userEmail string) error
	ListSecrets(ctx context.Context, path string) ([]*Secret, error)
	SetSecret(ctx context.Context, path string, secret []byte) error
	RemoveSecret(ctx context.Context, path string) error
	ListSecretPermissions(ctx context.Context, path string) ([]*SecretPermission, error)
	SetSecretPermission(ctx context.Context, path, userEmail, permission string) error
	RemoveSecretPermission(ctx context.Context, path, userEmail string) error
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

		if err != nil && strings.Contains(err.Error(), "context canceled") {
			// Some operations can be canceled gracefully, so we can signal this to higher-level functions.
			// Note the actual error caught here is not an instance of context.Canceled,
			// but is an error message that contains the HTTP call plus the string "context canceled".
			return status, body, context.Canceled
		}

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
		req.Header = r.headers.Clone()
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

var _ Client = &client{}

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
