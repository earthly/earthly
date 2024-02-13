package cloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/earthly/cloud-api/secrets"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	testToken = "ey123.abcdefg"
	testEmail = "test@earthly.dev"
	testPass  = "s3cr3t"
)

var testTokenExp = time.Now().Add(24 * time.Hour).UTC()

func TestClient_Authenticate(t *testing.T) {
	srv := mockServer(t)
	cc := &Client{
		httpAddr: srv.URL,
		email:    testEmail,
		password: testPass,
		authDir:  "/tmp",
		jum:      &protojson.UnmarshalOptions{DiscardUnknown: true},
	}
	ctx := context.Background()

	authMethod, err := cc.Authenticate(ctx)
	if err != nil {
		t.Fatalf("unexpected authentication error: %+v", err)
	}
	if authMethod != AuthMethodPassword {
		t.Errorf("expected [%s] got [%s]", AuthMethodPassword, authMethod)
	}

	if cc.authToken != testToken {
		t.Errorf("expected [%s] got [%s]", testToken, cc.authToken)
	}
	if !cc.authTokenExpiry.Equal(testTokenExp) {
		t.Errorf("expected [%s] got [%s]", testTokenExp, cc.authTokenExpiry)
	}
}

func TestClient_loadAuthStorage(t *testing.T) {
	cc := &Client{
		authToken:       testToken,
		authTokenExpiry: testTokenExp,
		email:           testEmail,
		password:        testPass,
		authDir:         "/tmp",
		jum:             &protojson.UnmarshalOptions{DiscardUnknown: true},
	}
	ctx := context.Background()
	err := cc.saveToken()
	assert.NoError(t, err)

	err = cc.savePasswordCredentials(ctx, cc.email, cc.password)
	assert.NoError(t, err)

	cc.email = ""
	cc.password = ""
	cc.authToken = ""
	cc.authTokenExpiry = time.Now()

	if err := cc.loadAuthStorage(); err != nil {
		t.Errorf("could not reload auth storage: %+v", err)
	}
	if cc.authToken != testToken {
		t.Errorf("expected [%s] got [%s]", testToken, cc.authToken)
	}
	if cc.authTokenExpiry != testTokenExp {
		t.Errorf("expected [%s] got [%s]", testTokenExp, cc.authTokenExpiry)
	}
	if cc.email != testEmail {
		t.Errorf("expected [%s] got [%s]", testEmail, cc.email)
	}
	if cc.password != testPass {
		t.Errorf("expected [%s] got [%s]", testPass, cc.password)
	}
}

func mockServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &api.LoginResponse{
			Token:  testToken,
			Expiry: timestamppb.New(testTokenExp),
		}
		encodedBody, err := protojson.Marshal(resp)
		if err != nil {
			t.Fatal("could not marshal mock response")
		}
		_, err = w.Write(encodedBody)
		assert.NoError(t, err)
	}))
}
