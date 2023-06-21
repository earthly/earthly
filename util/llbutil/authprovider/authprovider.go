package authprovider

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrAuthProviderNoResponse = fmt.Errorf("AuthServerNoResponse")

func NewAuthProvider(authServers []auth.AuthServer) session.Attachable {
	return &authProvider{
		authServers:     authServers,
		foundAuthServer: map[string]auth.AuthServer{},
		skipAuthServer:  map[string][]auth.AuthServer{},
	}
}

type authProvider struct {
	authServers []auth.AuthServer

	mu sync.Mutex

	// once an authServer has responded succcessfully, only that auth server will be used
	// for all subsequent calls -- this is to prevent accidentally mixing credentials and using them inconsistently
	foundAuthServer map[string]auth.AuthServer

	// if an authServer returns a ErrAuthProviderNoResponse, dont call it again for this host
	skipAuthServer map[string][]auth.AuthServer
}

func (ap *authProvider) Register(server *grpc.Server) {
	auth.RegisterAuthServer(server, ap)
}

func (ap *authProvider) getAuthServers(host string) []auth.AuthServer {
	as, ok := ap.foundAuthServer[host]
	if ok {
		return []auth.AuthServer{as}
	}
	res := []auth.AuthServer{}
	for _, as := range ap.authServers {
		if !ap.shouldSkip(host, as) {
			res = append(res, as)
		}
	}
	return res
}
func (ap *authProvider) setSkipAuthServer(host string, as auth.AuthServer) {
	if ap.shouldSkip(host, as) {
		return // already exists in skip list
	}
	ap.skipAuthServer[host] = append(ap.skipAuthServer[host], as)
}

func (ap *authProvider) shouldSkip(host string, as auth.AuthServer) bool {
	for _, x := range ap.skipAuthServer[host] {
		if x == as {
			return true
		}
	}
	return false
}

func (ap *authProvider) FetchToken(ctx context.Context, req *auth.FetchTokenRequest) (rr *auth.FetchTokenResponse, err error) {
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
		ap.foundAuthServer[req.Host] = as
		return a, nil
	}
	return nil, status.Errorf(codes.Unavailable, "no configured auth servers in the list of client-side configs responded")
}

func (ap *authProvider) Credentials(ctx context.Context, req *auth.CredentialsRequest) (*auth.CredentialsResponse, error) {
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
		ap.foundAuthServer[req.Host] = as
		return a, nil
	}
	return nil, status.Errorf(codes.Unavailable, "no configured auth servers in the list of client-side configs responded")
}

func (ap *authProvider) GetTokenAuthority(ctx context.Context, req *auth.GetTokenAuthorityRequest) (*auth.GetTokenAuthorityResponse, error) {
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
		ap.foundAuthServer[req.Host] = as
		return a, nil
	}
	return nil, status.Errorf(codes.Unavailable, "no configured auth servers in the list of client-side configs responded")
}

func (ap *authProvider) VerifyTokenAuthority(ctx context.Context, req *auth.VerifyTokenAuthorityRequest) (*auth.VerifyTokenAuthorityResponse, error) {
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
		ap.foundAuthServer[req.Host] = as
		return a, nil
	}
	return nil, status.Errorf(codes.Unavailable, "no configured auth servers in the list of client-side configs responded")
}
