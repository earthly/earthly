package cloud

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/earthly/cloud-api/analytics"
	"github.com/earthly/cloud-api/askv"
	"github.com/earthly/cloud-api/billing"
	"github.com/earthly/cloud-api/compute"
	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/cloud-api/pipelines"
	"github.com/earthly/cloud-api/secrets"

	"github.com/google/uuid"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/agent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	// ErrUnauthorized occurs when a user is unauthorized to access a resource
	ErrUnauthorized     = errors.New("unauthorized")
	ErrAuthTokenExpired = errors.New("auth token expired")
	// ErrNoAuthorizedPublicKeys occurs when no authorized public keys are found
	ErrNoAuthorizedPublicKeys = errors.New("no authorized public keys found")
	ErrNotFound               = errors.Errorf("not found")
	ErrMalformedSecretPath    = errors.Errorf("malformed secret path")
)

const (
	tokenExpiryLayout       = "2006-01-02 15:04:05.999999999 -0700 MST"
	satelliteMgmtTimeout    = "5M" // 5 minute timeout when launching or deleting a Satellite
	requestID               = "request-id"
	retryCount              = "retry-count"
	tokenExpiredServerError = "token expired"
)

type logstreamClient interface {
	StreamLogs(ctx context.Context, opts ...grpc.CallOption) (logstream.LogStream_StreamLogsClient, error)
}

type Client struct {
	httpAddr                 string
	sshKeyBlob               []byte // sshKey to use
	forceSSHKey              bool   // if true only use the above ssh key, don't attempt to guess others
	sshAgent                 agent.ExtendedAgent
	warnFunc                 func(string, ...interface{})
	debugFunc                func(string, ...interface{})
	email                    string
	password                 string
	authToken                string
	authTokenExpiry          time.Time
	authCredToken            string
	authDir                  string
	disableSSHKeyGuessing    bool
	jum                      *protojson.UnmarshalOptions
	pipelines                pipelines.PipelinesClient
	compute                  compute.ComputeClient
	logstream                logstreamClient
	logstreamBackoff         time.Duration
	analytics                analytics.AnalyticsClient
	askv                     askv.AskvClient
	billing                  billing.BillingClient
	secrets                  secrets.SecretsClient
	requestID                string
	installationName         string
	logstreamAddressOverride string
	serverConnTimeout        time.Duration
	orgIDCache               sync.Map // orgName -> orgID
	lastAuthMethod           AuthMethod
	lastAuthMethodExpiry     time.Time
}

// ClientOpt is used to customize the Cloud client.
type ClientOpt func(*Client)

// WithLogstreamGRPCAddressOverride can be used to override the Logstream gRPC address.
func WithLogstreamGRPCAddressOverride(address string) ClientOpt {
	return func(client *Client) {
		client.logstreamAddressOverride = address
	}
}

// WithAuthToken can be used to ignore other authentication mechanisms
// besides the given token. Any previously set JWT tokens are cleared.
func WithAuthToken(token string) ClientOpt {
	return func(client *Client) {
		client.authCredToken = token
		client.authToken = ""
		client.authTokenExpiry = time.Time{}
	}
}

// NewClient provides a new Earthly Cloud client
func NewClient(httpAddr, grpcAddr string, useInsecure bool, agentSockPath, authCredsOverride,
	authJWTOverride, installationName, requestID string, warnFunc func(string, ...interface{}),
	debugFunc func(string, ...interface{}), serverConnTimeout time.Duration, opts ...ClientOpt) (*Client, error) {
	c := &Client{
		httpAddr: httpAddr,
		sshAgent: &lazySSHAgent{
			sockPath: agentSockPath,
		},
		warnFunc:          warnFunc,
		debugFunc:         debugFunc,
		jum:               &protojson.UnmarshalOptions{DiscardUnknown: true},
		installationName:  installationName,
		requestID:         requestID,
		serverConnTimeout: serverConnTimeout,
	}

	if authJWTOverride != "" {
		c.authToken = authJWTOverride
		c.authTokenExpiry = time.Now().Add(24 * 365 * time.Hour) // Never expire when using JWT.
	} else if authCredsOverride != "" {
		c.authCredToken = authCredsOverride
	} else {
		if err := c.loadAuthStorage(); err != nil {
			return nil, err
		}
	}

	for _, opt := range opts {
		opt(c)
	}

	ctx := context.Background()

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(10),
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithCodes(codes.Internal, codes.Unavailable),
	}

	dialOpts := []grpc.DialOption{
		grpc.WithChainStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...), c.StreamInterceptor()),
		grpc.WithChainUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...), c.UnaryInterceptor(WithSkipAuth("/api.public.analytics.Analytics/SendAnalytics"))),
	}

	var transportCreds credentials.TransportCredentials
	if useInsecure {
		transportCreds = insecure.NewCredentials()
	} else {
		transportCreds = credentials.NewTLS(&tls.Config{})
	}

	dialOpts = append(dialOpts, grpc.WithTransportCredentials(transportCreds))
	conn, err := grpc.DialContext(ctx, grpcAddr, dialOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed dialing pipelines grpc")
	}

	c.pipelines = pipelines.NewPipelinesClient(conn)
	c.compute = compute.NewComputeClient(conn)
	c.analytics = analytics.NewAnalyticsClient(conn)
	c.askv = askv.NewAskvClient(conn)
	c.billing = billing.NewBillingClient(conn)
	c.secrets = secrets.NewSecretsClient(conn)

	logstreamAddr := grpcAddr
	if c.logstreamAddressOverride != "" {
		logstreamAddr = c.logstreamAddressOverride
	}

	c.logstreamBackoff = 250 * time.Millisecond
	c.logstream, err = newLogstreamClient(ctx, logstreamAddr, transportCreds)
	if err != nil {
		return nil, errors.Wrap(err, "cloud: could not create logstream client")
	}

	return c, nil
}

func (c *Client) getRequestID() string {
	if c.requestID != "" {
		return c.requestID
	}
	return uuid.NewString()
}

var serviceConfig = `{
	"methodConfig": [{
		"name": [{"service": "` + logstream.LogStream_ServiceDesc.ServiceName + `"}],
		"waitForReady": true,
		"retryPolicy": {
			"MaxAttempts": 10,
			"InitialBackoff": ".5s",
			"MaxBackoff": "10s",
			"BackoffMultiplier": 1.5,
			"RetryableStatusCodes": [ "UNAVAILABLE", "UNKNOWN" ]
		}
	}]
}`

func newLogstreamClient(ctx context.Context, addr string, transportCreds credentials.TransportCredentials) (logstream.LogStreamClient, error) {

	// Use custom dial options for log streaming as it uses long-lived,
	// sometimes idle, connections.
	dialOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(serviceConfig),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: 10 * time.Second}),
		grpc.WithTransportCredentials(transportCreds),
	}

	// Client context cancellation is managed by another process.
	conn, err := grpc.DialContext(context.WithoutCancel(ctx), addr, dialOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "cloud: failed dialing logstream grpc")
	}

	return logstream.NewLogStreamClient(conn), nil
}
