// Package cloudstoredauthprovider was forked from buildkit/session/auth/authprovider in order to allow using
// registry credentials that are stored on a server rather than local filesystem.
// This package is distributed under the original file's license, The Apache License, which is defined under
// https://github.com/moby/buildkit/blob/7c3e9fdd48c867f48a07a80cde64cc2d578cb332/LICENSE

package cloudauth

import (
	"context"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	authutil "github.com/containerd/containerd/remotes/docker/auth"
	remoteserrors "github.com/containerd/containerd/remotes/errors"
	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/llbutil/authprovider"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth"
	"github.com/moby/buildkit/util/progress/progresswriter"
	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/sign"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultExpiration = 60
	maxAuthRetries    = 5
)

var ErrNoCredentialsFound = fmt.Errorf("no credentials found")

type ProjectBasedAuthProvider interface {
	AddProject(org, project string)
}

func NewProvider(cfg *configfile.ConfigFile, cloudClient *cloud.Client, console conslogging.ConsoleLogger) session.Attachable {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = maxAuthRetries
	retryClient.CheckRetry = checkRetryFunc
	retryClient.Logger = nil
	return &authProvider{
		authConfigCache: map[string]*authConfig{},
		config:          cfg,
		seeds:           &tokenSeeds{dir: config.Dir()},
		loggerCache:     map[string]struct{}{},
		cloudClient:     cloudClient,
		seenOrgProjects: map[string]struct{}{},
		console:         console.WithPrefix("registry auth"),
		httpClient:      retryClient.StandardClient(),
	}
}

func checkRetryFunc(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if errors.Is(err, io.EOF) {
		return true, nil
	}
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}

type orgProject struct {
	org     string
	project string
}

type authConfig struct {
	loc string // where the credentials came from (e.g. user, or org/project)
	ac  *types.AuthConfig
}

type authProvider struct {
	authConfigCache map[string]*authConfig
	config          *configfile.ConfigFile
	seeds           *tokenSeeds
	logger          progresswriter.Logger
	loggerCache     map[string]struct{}
	httpClient      *http.Client

	// earthly-add on
	orgProjects     []orgProject
	seenOrgProjects map[string]struct{}
	cloudClient     *cloud.Client
	console         conslogging.ConsoleLogger

	// The need for this mutex is not well understood.
	// Without it, the docker cli on OS X hangs when
	// reading credentials from docker-credential-osxkeychain.
	// See issue https://github.com/docker/cli/issues/1862
	mu sync.Mutex
}

func (ap *authProvider) SetLogger(l progresswriter.Logger) {
	ap.mu.Lock()
	ap.logger = l
	ap.mu.Unlock()
}

func (ap *authProvider) Register(server *grpc.Server) {
	auth.RegisterAuthServer(server, ap)
}

func secondsToDuration(s int) time.Duration {
	return time.Duration(s) * time.Second
}

func (ap *authProvider) FetchToken(ctx context.Context, req *auth.FetchTokenRequest) (rr *auth.FetchTokenResponse, err error) {
	ac, err := ap.getAuthConfig(ctx, req.Host)
	if err != nil {
		return nil, err
	}

	// check for statically configured bearer token
	if ac.ac.RegistryToken != "" {
		return toTokenResponse(ac.ac.RegistryToken, time.Time{}, 0), nil
	}

	creds, err := ap.credentials(ctx, req.Host)
	if err != nil {
		return nil, err
	}

	if creds.Username != "" {
		ap.console.VerbosePrintf("attempting to login to %s with username %s (using credentials from %s)", req.Host, creds.Username, ac.loc)
	}

	to := authutil.TokenOptions{
		Realm:    req.Realm,
		Service:  req.Service,
		Scopes:   req.Scopes,
		Username: creds.Username,
		Secret:   creds.Secret,
	}

	if creds.Secret != "" {
		done := func(progresswriter.SubLogger) error {
			return err
		}
		defer func() {
			err = errors.Wrap(err, "failed to fetch oauth token")
		}()
		ap.mu.Lock()
		name := fmt.Sprintf("[auth] %v token for %s", strings.Join(trimScopePrefix(req.Scopes), " "), req.Host)
		if _, ok := ap.loggerCache[name]; !ok {
			_ = progresswriter.Wrap(name, ap.logger, done)
		}
		ap.mu.Unlock()
		// credential information is provided, use oauth POST endpoint
		resp, err := authutil.FetchTokenWithOAuth(ctx, ap.httpClient, nil, "buildkit-client", to)
		if err != nil {
			var errStatus remoteserrors.ErrUnexpectedStatus
			if errors.As(err, &errStatus) {
				// Registries without support for POST may return 404 for POST /v2/token.
				// As of September 2017, GCR is known to return 404.
				// As of February 2018, JFrog Artifactory is known to return 401.
				if (errStatus.StatusCode == 405 && to.Username != "") || errStatus.StatusCode == 404 || errStatus.StatusCode == 401 {
					resp, err := authutil.FetchToken(ctx, ap.httpClient, nil, to)
					if err != nil {
						ap.console.Warnf("failed to login to %s with username %s (using credentials from %s): %s", req.Host, creds.Username, ac.loc, err)
						return nil, err
					}
					ap.console.Printf("logged into %s with username %s (using credentials from %s)", req.Host, creds.Username, ac.loc)
					ap.console.VerbosePrintf("%s token was issued at %s, and expires in %s seconds", req.Host, resp.IssuedAt, secondsToDuration(resp.ExpiresIn))
					return toTokenResponse(resp.Token, resp.IssuedAt, resp.ExpiresIn), nil
				}
			}
			ap.console.Warnf("failed to login to %s with username %s (using credentials from %s): %s", req.Host, creds.Username, ac.loc, err)
			return nil, err
		}
		ap.console.Printf("logged into %s with username %s (using credentials from %s)", req.Host, creds.Username, ac.loc)
		ap.console.VerbosePrintf("%s token was issued at %s, and expires in %s seconds", req.Host, resp.IssuedAt, secondsToDuration(resp.ExpiresIn))
		return toTokenResponse(resp.AccessToken, resp.IssuedAt, resp.ExpiresIn), nil
	}
	// do request anonymously
	resp, err := authutil.FetchToken(ctx, ap.httpClient, nil, to)
	if err != nil {
		ap.console.Warnf("failed to login to %s anonymously: %s", req.Host, err)
		return nil, errors.Wrap(err, "failed to fetch anonymous token")
	}
	ap.console.Printf("logged into %s anonymously", req.Host)
	ap.console.VerbosePrintf("%s token was issued at %s, and expires in %s seconds", req.Host, resp.IssuedAt, secondsToDuration(resp.ExpiresIn))
	return toTokenResponse(resp.Token, resp.IssuedAt, resp.ExpiresIn), nil
}

func (ap *authProvider) credentials(ctx context.Context, host string) (*auth.CredentialsResponse, error) {
	ac, err := ap.getAuthConfig(ctx, host)
	if err != nil {
		return nil, err
	}
	res := &auth.CredentialsResponse{}
	if ac.ac.IdentityToken != "" {
		res.Secret = ac.ac.IdentityToken
	} else {
		res.Username = ac.ac.Username
		res.Secret = ac.ac.Password
	}
	return res, nil
}

func (ap *authProvider) Credentials(ctx context.Context, req *auth.CredentialsRequest) (*auth.CredentialsResponse, error) {
	return ap.credentials(ctx, req.Host)
}

func (ap *authProvider) GetTokenAuthority(ctx context.Context, req *auth.GetTokenAuthorityRequest) (*auth.GetTokenAuthorityResponse, error) {
	key, err := ap.getAuthorityKey(ctx, req.Host, req.Salt)
	if err != nil {
		return nil, err
	}

	return &auth.GetTokenAuthorityResponse{PublicKey: key[32:]}, nil
}

func (ap *authProvider) VerifyTokenAuthority(ctx context.Context, req *auth.VerifyTokenAuthorityRequest) (*auth.VerifyTokenAuthorityResponse, error) {
	key, err := ap.getAuthorityKey(ctx, req.Host, req.Salt)
	if err != nil {
		return nil, err
	}

	priv := new([64]byte)
	copy((*priv)[:], key)

	return &auth.VerifyTokenAuthorityResponse{Signed: sign.Sign(nil, req.Payload, priv)}, nil
}

func (ap *authProvider) projectExists(ctx context.Context, org, project string) (bool, error) {
	_, err := ap.cloudClient.ListSecrets(ctx, fmt.Sprintf("/%s/%s/std/registry", org, project))
	if err != nil {
		if strings.Contains(err.Error(), "resource not found") { // TODO better support for this error is needed by the client/server
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ap *authProvider) getAuthConfigUsernamePassword(ctx context.Context, host, org, project string) (*authConfig, error) {
	usernamePath := getRegistrySecret(host, org, project, "username")
	passwordPath := getRegistrySecret(host, org, project, "password")

	ap.console.VerbosePrintf("looking up %s", usernamePath)
	var username string
	usernameSecret, err := ap.cloudClient.GetUserOrProjectSecret(ctx, usernamePath)
	if err == nil {
		username = usernameSecret.Value
	} else if !errors.Is(err, cloud.ErrNotFound) {
		return nil, err
	}

	ap.console.VerbosePrintf("looking up %s", passwordPath)
	var password string
	passwordSecret, err := ap.cloudClient.GetUserOrProjectSecret(ctx, passwordPath)
	if err == nil {
		password = passwordSecret.Value
	} else if !errors.Is(err, cloud.ErrNotFound) {
		return nil, err
	}
	// TODO look for insecure and http config options (not sure how these options will be propagated to the rest of buildkit)

	if username == "" && password != "" {
		return nil, fmt.Errorf("found %s, but no %s", usernamePath, passwordPath)
	}
	if username != "" && password == "" {
		return nil, fmt.Errorf("found %s, but no %s", passwordPath, usernamePath)
	}

	if username == "" && password == "" {
		return nil, ErrNoCredentialsFound
	}

	return &authConfig{
		ac: &types.AuthConfig{
			ServerAddress: host,
			Username:      username,
			Password:      password,
		},
		loc: getRegistrySecretPrefix(host, org, project),
	}, nil
}

func getRegistrySecretPrefix(host, org, project string) string {
	pathPrefix := fmt.Sprintf("std/registry/%s/", host)
	if org == "" && project == "" {
		return fmt.Sprintf("/user/%s", pathPrefix)
	}
	return fmt.Sprintf("/%s/%s/%s", org, project, pathPrefix)
}

func getRegistrySecret(host, org, project, filename string) string {
	return getRegistrySecretPrefix(host, org, project) + filename
}

func (ap *authProvider) getAuthConfigForProject(ctx context.Context, org, project, host string) (*authConfig, error) {
	credHelperPath := getRegistrySecret(host, org, project, "cred_helper")
	ap.console.VerbosePrintf("looking up %s", credHelperPath)
	var credHelper string
	credHelperSecret, err := ap.cloudClient.GetUserOrProjectSecret(ctx, credHelperPath)
	if err == nil {
		credHelper = strings.TrimSpace(credHelperSecret.Value)
		ap.console.VerbosePrintf("configured cred_helper for %s is %q", host, credHelper)
	} else if errors.Is(err, cloud.ErrNotFound) {
		ap.console.VerbosePrintf("no cred_helper for %q exists, checking for username/password based config", credHelperPath)
	} else {
		return nil, err
	}

	switch credHelper {
	case "", "none": // either empty or ErrNotFound
		return ap.getAuthConfigUsernamePassword(ctx, host, org, project)
	case ECRCredHelper:
		return ap.getAuthConfigECR(ctx, host, org, project)
	case GCPCredHelper:
		return ap.getAuthConfigGCP(ctx, host, org, project)
	default:
		return nil, fmt.Errorf("unsupported cred_helper %s found at %s", credHelper, credHelperPath)
	}
}

// getAuthConfig was re-written to make use of earthly-cloud based credentials
func (ap *authProvider) getAuthConfig(ctx context.Context, host string) (*authConfig, error) {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	// check the cache
	ac, ok := ap.authConfigCache[host]
	if ok {
		return ac, nil
	}

	// check user's secrets
	ac, err := ap.getAuthConfigForProject(ctx, "", "", host)
	if err == nil {
		ap.authConfigCache[host] = ac
		return ac, nil
	}
	if !errors.Is(err, ErrNoCredentialsFound) {
		return nil, fmt.Errorf("failed to lookup credentials for %s: %w", host, err)
	}

	// fall back to project's secrets (starting with the root-level Earthfile's org/project)
	for _, op := range ap.orgProjects {
		exists, err := ap.projectExists(ctx, op.org, op.project)
		if err != nil {
			return nil, err
		}
		if !exists {
			ap.console.Warnf("Warning: PROJECT %s/%s was defined; however it doesn't exist (or you don't have access)", op.org, op.project)
			continue
		}

		ac, err := ap.getAuthConfigForProject(ctx, op.org, op.project, host)
		if err == nil {
			ap.authConfigCache[host] = ac
			return ac, nil
		}
		if !errors.Is(err, ErrNoCredentialsFound) {
			return nil, fmt.Errorf("failed to lookup credentials for %s: %w", host, err)
		}
	}

	return nil, authprovider.ErrAuthProviderNoResponse
}

func (ap *authProvider) getAuthorityKey(ctx context.Context, host string, salt []byte) (ed25519.PrivateKey, error) {
	if v, err := strconv.ParseBool(os.Getenv("BUILDKIT_NO_CLIENT_TOKEN")); err == nil && v {
		return nil, status.Errorf(codes.Unavailable, "client side tokens disabled")
	}

	creds, err := ap.credentials(ctx, host)
	if err != nil {
		return nil, err
	}
	seed, err := ap.seeds.getSeed(host)
	if err != nil {
		return nil, err
	}

	mac := hmac.New(sha256.New, salt)
	if creds.Secret != "" {
		mac.Write(seed)
	}

	sum := mac.Sum(nil)

	return ed25519.NewKeyFromSeed(sum[:ed25519.SeedSize]), nil
}

func (ap *authProvider) AddProject(org, project string) {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	orgProjectKey := fmt.Sprintf("%s/%s", org, project)

	if _, exists := ap.seenOrgProjects[orgProjectKey]; exists {
		return
	}
	ap.seenOrgProjects[orgProjectKey] = struct{}{}
	ap.orgProjects = append(ap.orgProjects, orgProject{
		org:     org,
		project: project,
	})
}

func toTokenResponse(token string, issuedAt time.Time, expires int) *auth.FetchTokenResponse {
	if expires == 0 {
		expires = defaultExpiration
	}
	resp := &auth.FetchTokenResponse{
		Token:     token,
		ExpiresIn: int64(expires),
	}
	if !issuedAt.IsZero() {
		resp.IssuedAt = issuedAt.Unix()
	}
	return resp
}

func trimScopePrefix(scopes []string) []string {
	out := make([]string, len(scopes))
	for i, s := range scopes {
		out[i] = strings.TrimPrefix(s, "repository:")
	}
	return out
}
