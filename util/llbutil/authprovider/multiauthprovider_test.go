package authprovider_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/pkg/pers"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/llbutil/authprovider"
	"github.com/moby/buildkit/session/auth"
	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newConsLogger() conslogging.ConsoleLogger {
	return conslogging.New(os.Stderr, &sync.Mutex{}, conslogging.NoColor, 0, conslogging.Info, false)
}

func TestMultiAuth(t *testing.T) {
	type testCtx struct {
		*testing.T
		expect   expect.Expectation
		children []*mockChild
		multi    *authprovider.MultiAuthProvider
	}

	o := onpar.BeforeEach(onpar.New(t), func(t *testing.T) testCtx {
		children := []*mockChild{
			newMockChild(t, mockTimeout),
			newMockChild(t, mockTimeout),
		}
		var srv []authprovider.Child
		for _, c := range children {
			srv = append(srv, c)
		}
		return testCtx{
			T:        t,
			expect:   expect.New(t),
			children: children,
			multi:    authprovider.New(newConsLogger(), srv),
		}
	})
	defer o.Run()

	type fetchResult struct {
		resp *auth.FetchTokenResponse
		err  error
	}

	o.Spec("it calls child ProjectAdders", func(t testCtx) {
		type projectProvider struct {
			*mockChild
			*mockProjectAdder
		}
		p := projectProvider{
			mockChild:        newMockChild(t, mockTimeout),
			mockProjectAdder: newMockProjectAdder(t, mockTimeout),
		}
		t.multi = authprovider.New(newConsLogger(), []authprovider.Child{p})
		t.multi.AddProject("foo", "bar")
		t.expect(p.mockProjectAdder).To(haveMethodExecuted("AddProject", withArgs("foo", "bar")))
	})

	o.Spec("it does not continue to contact servers with no credentials for a given host", func(t testCtx) {
		const host = "foo.bar"
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		res := make(chan fetchResult)
		go func() {
			resp, err := t.multi.FetchToken(ctx, &auth.FetchTokenRequest{Host: host})
			res <- fetchResult{resp, err}
		}()

		for _, c := range t.children {
			t.expect(c).To(haveMethodExecuted(
				"FetchToken",
				within(timeout),
				withArgs(pers.Any, equal(&auth.FetchTokenRequest{Host: host})),
				returning(nil, authprovider.ErrAuthProviderNoResponse),
			))
		}

		select {
		case result := <-res:
			t.expect(result.resp).To(beNil())
			t.expect(status.Code(result.err)).To(equal(codes.Unavailable))
		case <-time.After(timeout):
			t.Fatal("timed out waiting for FetchToken to return")
		}

		go func() {
			resp, err := t.multi.FetchToken(ctx, &auth.FetchTokenRequest{Host: host})
			res <- fetchResult{resp, err}
		}()

		for _, c := range t.children {
			t.expect(c).To(not(haveMethodExecuted(
				"FetchToken",
				within(10*time.Millisecond),
			)))
		}

		select {
		case result := <-res:
			t.expect(result.resp).To(beNil())
			t.expect(status.Code(result.err)).To(equal(codes.Unavailable))
		case <-time.After(timeout):
			t.Fatal("timed out waiting for FetchToken to return")
		}
	})

	o.Spec("it resets its knowledge of which servers it should contact after a project is added", func(t testCtx) {
		const host = "foo.bar"
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		res := make(chan fetchResult)
		go func() {
			resp, err := t.multi.FetchToken(ctx, &auth.FetchTokenRequest{Host: host})
			res <- fetchResult{resp, err}
		}()

		for i, c := range t.children {
			ret := []any{
				nil,
				authprovider.ErrAuthProviderNoResponse,
			}
			if i == len(t.children)-1 {
				// ensure one child responds so that we can prove that
				// successful results are also cache-busted.
				ret = []any{
					&auth.FetchTokenResponse{},
					nil,
				}
			}
			t.expect(c).To(haveMethodExecuted(
				"FetchToken",
				within(timeout),
				withArgs(pers.Any, equal(&auth.FetchTokenRequest{Host: host})),
				returning(ret...),
			))
		}

		select {
		case result := <-res:
			t.expect(status.Code(result.err)).To(not(haveOccurred()))
		case <-time.After(timeout):
			t.Fatal("timed out waiting for FetchToken to return")
		}

		t.multi.AddProject("foo", "bar")

		go func() {
			resp, err := t.multi.FetchToken(ctx, &auth.FetchTokenRequest{Host: host})
			res <- fetchResult{resp, err}
		}()

		for _, c := range t.children {
			t.expect(c).To(haveMethodExecuted(
				"FetchToken",
				within(timeout),
				withArgs(pers.Any, equal(&auth.FetchTokenRequest{Host: host})),
				returning(nil, authprovider.ErrAuthProviderNoResponse),
			))
		}

		select {
		case result := <-res:
			t.expect(result.resp).To(beNil())
			t.expect(status.Code(result.err)).To(equal(codes.Unavailable))
		case <-time.After(timeout):
			t.Fatal("timed out waiting for FetchToken to return")
		}
	})
}
