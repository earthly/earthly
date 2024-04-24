package authprovider

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/earthly/earthly/conslogging"
	"github.com/moby/buildkit/session/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrAuthProviderNoResponse = fmt.Errorf("AuthServerNoResponse")

// ProjectAdder is an optional interface that auth servers may implement. If
// they do, the MultiAuthProvider will call their AddProject method when its
// AddProject method is called.
type ProjectAdder interface {
	AddProject(org, project string)
}

// Child is the interface that child auth providers need to implement for
// MultiAuthProvider.
type Child interface {
	Credentials(context.Context, *auth.CredentialsRequest) (*auth.CredentialsResponse, error)
	FetchToken(context.Context, *auth.FetchTokenRequest) (*auth.FetchTokenResponse, error)
	GetTokenAuthority(context.Context, *auth.GetTokenAuthorityRequest) (*auth.GetTokenAuthorityResponse, error)
	VerifyTokenAuthority(context.Context, *auth.VerifyTokenAuthorityRequest) (*auth.VerifyTokenAuthorityResponse, error)
}

// New returns a new MultiAuthProvider, wrapping up multiple child auth providers.
func New(console conslogging.ConsoleLogger, authServers []Child) *MultiAuthProvider {
	return &MultiAuthProvider{
		console:         console,
		authServers:     authServers,
		foundAuthServer: map[string]Child{},
		skipAuthServer:  map[string][]Child{},
	}
}

// MultiAuthProvider is an auth provider that delegates authentication to
// multiple child auth providers.
type MultiAuthProvider struct {
	console     conslogging.ConsoleLogger
	authServers []Child
	mu          sync.Mutex

	// once an authServer has responded succcessfully, only that auth server
	// will be used for all subsequent calls -- this is to prevent accidentally
	// mixing credentials and using them inconsistently
	foundAuthServer map[string]Child

	// if an authServer returns an ErrAuthProviderNoResponse, dont call it again
	// for this host unless AddProject is called.
	skipAuthServer map[string][]Child
}

// Register registers ap against server.
func (ap *MultiAuthProvider) Register(server *grpc.Server) {
	auth.RegisterAuthServer(server, ap)
}

func (ap *MultiAuthProvider) getAuthServers(host string) []Child {
	as, ok := ap.foundAuthServer[host]
	if ok {
		return []Child{as}
	}
	res := []Child{}
	for _, as := range ap.authServers {
		if !ap.shouldSkip(host, as) {
			res = append(res, as)
		}
	}
	return res
}

func (ap *MultiAuthProvider) setAuthServer(host string, as Child) {
	ap.console.VerbosePrintf("using %T for %s", as, host)
	ap.foundAuthServer[host] = as
}

func (ap *MultiAuthProvider) setSkipAuthServer(host string, as Child) {
	if ap.shouldSkip(host, as) {
		return // already exists in skip list
	}
	ap.skipAuthServer[host] = append(ap.skipAuthServer[host], as)
}

func (ap *MultiAuthProvider) shouldSkip(host string, as Child) bool {
	for _, x := range ap.skipAuthServer[host] {
		if x == as {
			return true
		}
	}
	return false
}

// AddProject searches for any children implementing ProjectAdder and calls
// them, then invalidates its cached auth server responses.
func (ap *MultiAuthProvider) AddProject(org, proj string) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	ap.foundAuthServer = make(map[string]Child)
	ap.skipAuthServer = make(map[string][]Child)
	for _, s := range ap.authServers {
		adder, ok := s.(ProjectAdder)
		if !ok {
			continue
		}
		adder.AddProject(org, proj)
	}
}

// FetchToken calls child FetchToken methods until one of ap's children
// succeeds.
func (ap *MultiAuthProvider) FetchToken(ctx context.Context, req *auth.FetchTokenRequest) (rr *auth.FetchTokenResponse, err error) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	for _, as := range ap.getAuthServers(req.Host) {
		a, err := as.FetchToken(ctx, req)
		if err != nil {
			if errors.Is(err, ErrAuthProviderNoResponse) {
				ap.setSkipAuthServer(req.Host, as)
				continue
			}
			return nil, err
		}
		if a.Anonymous {
			ap.console.Warnf("Warning: you are not logged into %s, you may experience rate-limitting when pulling images\n", req.Host)
		}
		ap.setAuthServer(req.Host, as)
		return a, nil
	}
	return nil, status.Errorf(codes.Unavailable, "no configured auth servers in the list of client-side configs responded")
}

// Credentials calls child Credentials methods until one of ap's children
// succeeds.
func (ap *MultiAuthProvider) Credentials(ctx context.Context, req *auth.CredentialsRequest) (*auth.CredentialsResponse, error) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	for _, as := range ap.getAuthServers(req.Host) {
		a, err := as.Credentials(ctx, req)
		if err != nil {
			if errors.Is(err, ErrAuthProviderNoResponse) {
				ap.setSkipAuthServer(req.Host, as)
				continue
			}
			return nil, err
		}
		ap.setAuthServer(req.Host, as)
		return a, nil
	}
	return nil, status.Errorf(codes.Unavailable, "no configured auth servers in the list of client-side configs responded")
}

// GetTokenAuthority calls child GetTokenAuthority methods until one of ap's
// children succeeds.
func (ap *MultiAuthProvider) GetTokenAuthority(ctx context.Context, req *auth.GetTokenAuthorityRequest) (*auth.GetTokenAuthorityResponse, error) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	for _, as := range ap.getAuthServers(req.Host) {
		a, err := as.GetTokenAuthority(ctx, req)
		if err != nil {
			if errors.Is(err, ErrAuthProviderNoResponse) {
				ap.setSkipAuthServer(req.Host, as)
				continue
			}
			return nil, err
		}
		ap.setAuthServer(req.Host, as)
		return a, nil
	}
	return nil, status.Errorf(codes.Unavailable, "no configured auth servers in the list of client-side configs responded")
}

// VerifyTokenAuthority calls child VerifyTokenAuthority methods until one of
// ap's children succeeds.
func (ap *MultiAuthProvider) VerifyTokenAuthority(ctx context.Context, req *auth.VerifyTokenAuthorityRequest) (*auth.VerifyTokenAuthorityResponse, error) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	for _, as := range ap.getAuthServers(req.Host) {
		a, err := as.VerifyTokenAuthority(ctx, req)
		if err != nil {
			if errors.Is(err, ErrAuthProviderNoResponse) {
				ap.setSkipAuthServer(req.Host, as)
				continue
			}
			return nil, err
		}
		ap.setAuthServer(req.Host, as)
		return a, nil
	}
	return nil, status.Errorf(codes.Unavailable, "no configured auth servers in the list of client-side configs responded")
}
