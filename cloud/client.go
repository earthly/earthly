package cloud

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/earthly/cloud-api/compute"
	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/cloud-api/pipelines"
	"github.com/google/uuid"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/agent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
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
	requestID            = "request-id"
)

// Client contains gRPC and REST endpoints to the Earthly Cloud backend.
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
	AcceptInvite(ctx context.Context, inviteCode string) error
	ListInvites(ctx context.Context, org string) ([]*OrgInvitation, error)
	ListOrgs(ctx context.Context) ([]*OrgDetail, error)
	ListOrgPermissions(ctx context.Context, path string) ([]*OrgPermissions, error)
	ListOrgMembers(ctx context.Context, orgName string) ([]*OrgMember, error)
	UpdateOrgMember(ctx context.Context, orgName, userEmail, permission string) error
	RemoveOrgMember(ctx context.Context, orgName, userEmail string) error
	RevokePermission(ctx context.Context, path, user string) error
	ListPublicKeys(ctx context.Context) ([]string, error)
	AddPublicKey(ctx context.Context, key string) error
	RemovePublicKey(ctx context.Context, key string) error
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
	LaunchSatellite(ctx context.Context, name, org string, features []string) error
	GetOrgID(ctx context.Context, name string) (string, error)
	ListSatellites(ctx context.Context, orgID string) ([]SatelliteInstance, error)
	GetSatellite(ctx context.Context, name, orgID string) (*SatelliteInstance, error)
	DeleteSatellite(ctx context.Context, name, orgID string) error
	ReserveSatellite(ctx context.Context, name, orgID, gitAuthor, gitConfigEmail string, isCI bool) chan SatelliteStatusUpdate
	WakeSatellite(ctx context.Context, name, orgID string) chan SatelliteStatusUpdate
	SleepSatellite(ctx context.Context, name, orgID string) chan SatelliteStatusUpdate
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
	AccountResetRequestToken(ctx context.Context, userEmail string) error
	AccountReset(ctx context.Context, userEmail, token, password string) error
	StreamLogs(ctx context.Context, buildID string, deltas chan []*logstream.Delta) error
}

type client struct {
	httpAddr              string
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
	jum                   *protojson.UnmarshalOptions
	pipelines             pipelines.PipelinesClient
	compute               compute.ComputeClient
	logstream             logstream.LogStreamClient
	requestID             string
	installationName      string
}

var _ Client = &client{}

// NewClient provides a new Earthly Cloud client
func NewClient(httpAddr, grpcAddr string, useInsecure bool, agentSockPath, authCredsOverride, installationName, requestID string, warnFunc func(string, ...interface{})) (Client, error) {
	c := &client{
		httpAddr: httpAddr,
		sshAgent: &lazySSHAgent{
			sockPath: agentSockPath,
		},
		warnFunc:         warnFunc,
		jum:              &protojson.UnmarshalOptions{DiscardUnknown: true},
		installationName: installationName,
		requestID:        requestID,
	}
	if authCredsOverride != "" {
		c.authCredToken = authCredsOverride
	} else {
		if err := c.loadAuthStorage(); err != nil {
			return nil, err
		}
	}
	ctx := context.Background()
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(10),
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithCodes(codes.Internal, codes.Unavailable),
	}
	dialOpts := []grpc.DialOption{
		grpc.WithChainStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...), c.StreamInterceptor()),
		grpc.WithChainUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...), c.UnaryInterceptor()),
	}
	if useInsecure {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsConfig := credentials.NewTLS(&tls.Config{})
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(tlsConfig))
	}
	conn, err := grpc.DialContext(ctx, grpcAddr, dialOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed dialing pipelines grpc")
	}
	c.pipelines = pipelines.NewPipelinesClient(conn)
	c.compute = compute.NewComputeClient(conn)
	c.logstream = logstream.NewLogStreamClient(conn)
	return c, nil
}

func (c *client) getRequestID() string {
	if c.requestID != "" {
		return c.requestID
	}
	return uuid.NewString()
}
