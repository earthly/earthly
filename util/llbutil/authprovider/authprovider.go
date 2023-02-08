package authprovider

import (
	"context"
	"errors"
	"fmt"

	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrAuthServerNoResponse = fmt.Errorf("AuthServerNoResponse")

func NewAuthProvider(authServers []auth.AuthServer) session.Attachable {
	return &authProvider{
		authServers: authServers,
	}
}

type authProvider struct {
	authServers []auth.AuthServer
}

func (ap *authProvider) Register(server *grpc.Server) {
	auth.RegisterAuthServer(server, ap)
}

func (ap *authProvider) FetchToken(ctx context.Context, req *auth.FetchTokenRequest) (rr *auth.FetchTokenResponse, err error) {
	for _, as := range ap.authServers {
		a, err := as.FetchToken(ctx, req)
		if errors.Is(err, ErrAuthServerNoResponse) {
			continue
		}
		return a, err
	}
	return nil, status.Errorf(codes.Unavailable, "no configured auth servers in the list of client-side configs responded")
}

func (ap *authProvider) Credentials(ctx context.Context, req *auth.CredentialsRequest) (*auth.CredentialsResponse, error) {
	for _, as := range ap.authServers {
		a, err := as.Credentials(ctx, req)
		if errors.Is(err, ErrAuthServerNoResponse) {
			continue
		}
		return a, err
	}
	return nil, status.Errorf(codes.Unavailable, "no configured auth servers in the list of client-side configs responded")
}

func (ap *authProvider) GetTokenAuthority(ctx context.Context, req *auth.GetTokenAuthorityRequest) (*auth.GetTokenAuthorityResponse, error) {
	for _, as := range ap.authServers {
		a, err := as.GetTokenAuthority(ctx, req)
		if errors.Is(err, ErrAuthServerNoResponse) {
			continue
		}
		return a, err
	}
	return nil, status.Errorf(codes.Unavailable, "no configured auth servers in the list of client-side configs responded")
}

func (ap *authProvider) VerifyTokenAuthority(ctx context.Context, req *auth.VerifyTokenAuthorityRequest) (*auth.VerifyTokenAuthorityResponse, error) {
	for _, as := range ap.authServers {
		a, err := as.VerifyTokenAuthority(ctx, req)
		if errors.Is(err, ErrAuthServerNoResponse) {
			continue
		}
		return a, err
	}
	return nil, status.Errorf(codes.Unavailable, "no configured auth servers in the list of client-side configs responded")
}
