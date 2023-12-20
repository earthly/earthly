package authprovider_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/earthly/earthly/util/llbutil/authprovider"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth"
	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
)

const (
	authFmt = `
{
  "auths": {
    "%s": {
      "auth": "%s"
    }
  }
}
`
)

func TestPodmanProvider(topT *testing.T) {
	type testCtx struct {
		t      *testing.T
		expect expect.Expectation
		os     *mockOS
		stderr *mockWriter
		result chan session.Attachable
	}

	type credentials interface {
		Credentials(ctx context.Context, req *auth.CredentialsRequest) (*auth.CredentialsResponse, error)
	}

	o := onpar.BeforeEach(onpar.New(topT), func(t *testing.T) testCtx {
		tt := testCtx{
			t:      t,
			expect: expect.New(t),
			os:     newMockOS(t, mockTimeout),
			stderr: newMockWriter(t, mockTimeout),
			result: make(chan session.Attachable),
		}
		go func() {
			defer close(tt.result)
			tt.result <- authprovider.NewPodman(tt.stderr, authprovider.WithOS(tt.os))
		}()
		return tt
	})
	defer o.Run()

	o.AfterEach(func(tt testCtx) {
		_, ok := <-tt.result
		tt.expect(ok).To(beFalse()) // Ensure that the channel was closed
	})

	type authFile struct {
		path   any // can be a string or a matcher
		host   string
		user   string
		secret string
	}

	type entry struct {
		envs []string
		auth *authFile
	}

	onpar.TableSpec(o, func(tt testCtx, e entry) {
		for _, env := range e.envs {
			name, val, ok := strings.Cut(env, "=")
			tt.expect(ok).To(beTrue())
			tt.expect(tt.os).To(haveMethodExecuted("Getenv",
				within(timeout),
				withArgs(name),
				returning(val),
			))
		}
		if e.auth == nil {
			// The code should fall back to the default docker auth provider,
			// which we can't really mock out - but we can at least verify that
			// the return value is the correct type.
			tt.expect(tt.os).To(haveMethodExecuted("Open", within(timeout), returning(nil, fs.ErrNotExist)))
			select {
			case res := <-tt.result:
				_, ok := res.(credentials)
				tt.expect(ok).To(beTrue())
			case <-time.After(timeout):
				tt.t.Fatalf("timed out waiting to fall back to the default docker auth provider")
			}
			return
		}
		creds := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", e.auth.user, e.auth.secret)))
		authFile := io.NopCloser(bytes.NewBufferString(fmt.Sprintf(authFmt, e.auth.host, creds)))
		tt.expect(tt.os).To(haveMethodExecuted("Open",
			within(timeout),
			withArgs(e.auth.path),
			returning(authFile, nil)))
		select {
		case res := <-tt.result:
			creds, ok := res.(credentials)
			tt.expect(ok).To(beTrue())
			req := &auth.CredentialsRequest{
				Host: e.auth.host,
			}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			resp, err := creds.Credentials(ctx, req)
			tt.expect(err).To(not(haveOccurred()))
			tt.expect(resp.Username).To(equal(e.auth.user))
			tt.expect(resp.Secret).To(equal(e.auth.secret))
		case <-time.After(timeout):
			tt.t.Fatalf("timed out waiting for a podman auth provider")
		}
	}).
		Entry("it prefers REGISTRY_AUTH_FILE", entry{
			envs: []string{
				"REGISTRY_AUTH_FILE=/path/to/someFile",
			},
			auth: &authFile{
				path:   "/path/to/someFile",
				host:   "foo.bar",
				user:   "foo",
				secret: "bar",
			},
		}).
		Entry("it falls back to XDG_RUNTIME_DIR/containers/auth.json", entry{
			envs: []string{
				"REGISTRY_AUTH_FILE=",
				"XDG_RUNTIME_DIR=/path/to/some/dir",
			},
			auth: &authFile{
				path:   "/path/to/some/dir/containers/auth.json",
				host:   "bacon.eggs",
				user:   "eggs",
				secret: "bacon",
			},
		}).
		Entry("it checks the root runtime dir last", entry{
			envs: []string{
				"REGISTRY_AUTH_FILE=",
				"XDG_RUNTIME_DIR=",
			},
			auth: &authFile{
				path:   matchRegexp("/run/containers/[0-9]*/auth.json"),
				host:   "foo",
				user:   "bar",
				secret: "baz",
			},
		}).
		Entry("it returns a provider even when no podman auth file exists", entry{
			envs: []string{
				"REGISTRY_AUTH_FILE=",
				"XDG_RUNTIME_DIR=",
			},
			auth: nil,
		})
}
