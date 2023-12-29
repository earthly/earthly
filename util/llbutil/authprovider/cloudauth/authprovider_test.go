package cloudauth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	dockerauth "github.com/containerd/containerd/remotes/docker/auth"
	"github.com/docker/cli/cli/config/types"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/moby/buildkit/session/auth"
	"github.com/stretchr/testify/require"
)

func newRetryClient() *retryablehttp.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.CheckRetry = checkRetryFunc

	// Use backoff of 0 to limit impact to test performance.
	retryClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		return 0
	}

	return retryClient
}

func newAuthProvider(host string, httpClient *http.Client) *authProvider {
	return &authProvider{
		httpClient: httpClient,
		authConfigCache: map[string]*authConfig{
			host: {
				loc: "user",
				ac: &types.AuthConfig{
					Username: "user",
					Password: "pass",
				},
			},
		},
	}
}

func TestFetchTokenRetry(t *testing.T) {

	mux := http.NewServeMux()

	count := 0

	mux.Handle("/v2/token", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		if count <= 3 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		resp := &dockerauth.OAuthTokenResponse{
			AccessToken:  "token",
			RefreshToken: "refresh",
			ExpiresIn:    0,
			IssuedAt:     time.Now(),
		}

		_ = json.NewEncoder(w).Encode(resp)
	}))

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	retryClient := newRetryClient()

	attempts := 0
	retryClient.RequestLogHook = func(logger retryablehttp.Logger, r *http.Request, i int) {
		attempts = i + 1
	}

	u, _ := url.Parse(testServer.URL)

	p := newAuthProvider(u.Host, retryClient.StandardClient())

	ctx := context.Background()

	_, err := p.FetchToken(ctx, &auth.FetchTokenRequest{
		ClientID: "client-id",
		Host:     u.Host,
		Realm:    testServer.URL + "/v2/token",
		Service:  "service",
		Scopes:   []string{"foo", "bar"},
	})

	r := require.New(t)
	r.NoError(err)
	r.Equal(4, attempts, "4 total retries (3 failures & 1 success)")
}

type eofTransport struct {
	count    int
	times    int
	fallback http.RoundTripper
}

func (t *eofTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.count++
	if t.count < t.times {
		return nil, io.EOF
	}
	return t.fallback.RoundTrip(req)
}

func TestFetchTokenRetryEOF(t *testing.T) {

	retryClient := newRetryClient()
	retryClient.HTTPClient = &http.Client{
		Transport: &eofTransport{fallback: http.DefaultTransport, times: 2},
	}

	mux := http.NewServeMux()

	mux.Handle("/v2/token", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &dockerauth.OAuthTokenResponse{
			AccessToken:  "token",
			RefreshToken: "refresh",
			ExpiresIn:    0,
			IssuedAt:     time.Now(),
		}

		_ = json.NewEncoder(w).Encode(resp)
	}))

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	attempts := 0
	retryClient.RequestLogHook = func(logger retryablehttp.Logger, r *http.Request, i int) {
		attempts++
	}

	u, _ := url.Parse(testServer.URL)

	p := newAuthProvider(u.Host, retryClient.StandardClient())

	ctx := context.Background()

	_, err := p.FetchToken(ctx, &auth.FetchTokenRequest{
		ClientID: "client-id",
		Host:     u.Host,
		Realm:    testServer.URL + "/v2/token",
		Service:  "service",
		Scopes:   []string{"foo", "bar"},
	})

	r := require.New(t)
	r.NoError(err)
	r.Equal(2, attempts, "2 total retries (1 EOF & 1 success)")
}

func TestFetchTokenRetry404(t *testing.T) {

	mux := http.NewServeMux()

	count := 0

	mux.Handle("/v2/token", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost: // Initial Post request to /v2/token
			w.WriteHeader(http.StatusNotFound)
			return
		case http.MethodGet: // Follow-on Get requests to /v2/token (with failures)
			count++
			if count <= 3 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			resp := &dockerauth.OAuthTokenResponse{
				AccessToken:  "token",
				RefreshToken: "refresh",
				ExpiresIn:    0,
				IssuedAt:     time.Now(),
			}

			_ = json.NewEncoder(w).Encode(resp)
		}
	}))

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	retryClient := newRetryClient()

	attempts := 0
	retryClient.RequestLogHook = func(logger retryablehttp.Logger, r *http.Request, i int) {
		attempts++
	}

	u, _ := url.Parse(testServer.URL)

	p := newAuthProvider(u.Host, retryClient.StandardClient())

	ctx := context.Background()

	_, err := p.FetchToken(ctx, &auth.FetchTokenRequest{
		ClientID: "client-id",
		Host:     u.Host,
		Realm:    testServer.URL + "/v2/token",
		Service:  "service",
		Scopes:   []string{"foo", "bar"},
	})

	r := require.New(t)
	r.NoError(err)
	r.Equal(5, attempts, "5 total attempts (1 Post 404, 3 Get 500s, & 1 Get success)")
}
