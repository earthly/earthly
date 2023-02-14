//go:generate hel --output helheim_mocks_test.go

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
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/agent"
	"google.golang.org/grpc"
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

// GRPCConn represents a grpc connection. It's effectively a copy of
// grpc.ClientConnInterface, copied here to make it easier to see what is
// required by this package without having to dig through the grpc docs.
type GRPCConn interface {
	Invoke(ctx context.Context, method string, args, reply any, opts ...grpc.CallOption) error
	NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error)
}

// Client is a client for communicating with Earthly Cloud.
type Client struct {
	httpAddr              string
	sshKeyBlob            []byte // sshKey to use
	forceSSHKey           bool   // if true only use the above ssh key, don't attempt to guess others
	sshAgent              agent.ExtendedAgent
	warnFunc              func(string, ...any)
	email                 string
	password              string
	authToken             string
	authTokenExpiry       time.Time
	authCredToken         string
	authDir               string
	disableSSHKeyGuessing bool
	jum                   *protojson.UnmarshalOptions
	grpcConn              GRPCConn
	pipelines             pipelines.PipelinesClient
	compute               compute.ComputeClient
	logstream             logstream.LogStreamClient
	analytics             analytics.AnalyticsClient
	requestID             string
	installationName      string

	// postAllocationOpts are options that can't be applied until the Client has
	// been fully allocated and all Opts have been applied.
	//
	// ... Which ... is a problem, because if a postAllocationOpt depends on
	// another postAllocationOpt...
	//
	// It's not a big deal for now, but we probably should break the cycle when
	// we get a chance to descend into that rabbit hole. Maybe `Client` is doing
	// too much and should be broken up into a few types (e.g. Authenticator,
	// Interceptor, Client).
	postAllocationOpts []Opt
}

// NewClient provides a new Earthly Cloud client
func NewClient(httpAddr, grpcAddr string, useInsecure bool, agentSockPath, authCredsOverride, authJWTOverride, installationName, requestID string, warnFunc func(string, ...interface{})) (*Client, error) {
	opts := []Opt{
		WithHTTPAddr(httpAddr),
		WithAgentSockPath(agentSockPath),
		WithInstallation(installationName),
		WithRequestID(requestID),
		WithWarnCallback(warnFunc),
	}

	// NOTE: retry and interceptor options are dealt with in the WithGRPCDial
	// option.
	var dialOpts []grpc.DialOption
	if useInsecure {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsConfig := credentials.NewTLS(&tls.Config{})
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(tlsConfig))
	}
	opts = append(opts, WithGRPCDial(grpcAddr, dialOpts...))

	if authJWTOverride != "" {
		opts = append(opts, WithAuthOverride(JWT(authJWTOverride)))
	} else if authCredsOverride != "" {
		opts = append(opts, WithAuthOverride(Creds(authCredsOverride)))
	}

	return NewClientOpts(context.Background(), opts...)
}

// NewClientOpts provides a new Earthly Cloud client using the provided opts to
// control how it is configured.
func NewClientOpts(ctx context.Context, opts ...Opt) (*Client, error) {
	c := &Client{
		jum: &protojson.UnmarshalOptions{DiscardUnknown: true},
	}
	for _, opt := range opts {
		newC, err := opt(ctx, c)
		if err != nil {
			return nil, err
		}
		c = newC
	}
	if c.authToken == "" && c.authCredToken == "" {
		if err := c.loadAuthStorage(); err != nil {
			return nil, err
		}
	}

	for _, opt := range c.postAllocationOpts {
		newC, err := opt(ctx, c)
		if err != nil {
			return nil, err
		}
		c = newC
	}
	c.pipelines = pipelines.NewPipelinesClient(c.grpcConn)
	c.compute = compute.NewComputeClient(c.grpcConn)
	c.logstream = logstream.NewLogStreamClient(c.grpcConn)
	c.analytics = analytics.NewAnalyticsClient(conn)

	return c, nil
}

func (c *Client) getRequestID() string {
	if c.requestID != "" {
		return c.requestID
	}
	return uuid.NewString()
}
