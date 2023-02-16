package cloud

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/earthly/cloud-api/analytics"
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
	ErrNotFound               = errors.Errorf("not found")
	ErrMalformedSecretPath    = errors.Errorf("malformed secret path")
)

const (
	tokenExpiryLayout    = "2006-01-02 15:04:05.999999999 -0700 MST"
	satelliteMgmtTimeout = "5M" // 5 minute timeout when launching or deleting a Satellite
	requestID            = "request-id"
)

type Client struct {
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
	analytics             analytics.AnalyticsClient
	requestID             string
	installationName      string
}

// NewClient provides a new Earthly Cloud client
func NewClient(httpAddr, grpcAddr string, useInsecure bool, agentSockPath, authCredsOverride, authJWTOverride, installationName, requestID string, warnFunc func(string, ...interface{})) (*Client, error) {
	c := &Client{
		httpAddr: httpAddr,
		sshAgent: &lazySSHAgent{
			sockPath: agentSockPath,
		},
		warnFunc:         warnFunc,
		jum:              &protojson.UnmarshalOptions{DiscardUnknown: true},
		installationName: installationName,
		requestID:        requestID,
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
	c.analytics = analytics.NewAnalyticsClient(conn)
	return c, nil
}

func (c *Client) getRequestID() string {
	if c.requestID != "" {
		return c.requestID
	}
	return uuid.NewString()
}
