package cloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/earthly/cloud-api/secrets"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
)

const (
	testToken = "ey123.abcdefg"
	testEmail = "test@earthly.dev"
	testPass  = "s3cr3t"
)

var testTokenExp = time.Now().Add(24 * time.Hour).UTC()

func TestClient_Authenticate(t *testing.T) {
	srv := mockServer(t)
	cc := &client{
		httpAddr: srv.URL,
		email:    testEmail,
		password: testPass,
		authDir:  "/tmp",
		jm: &jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		},
	}
	ctx := context.Background()

	if err := cc.Authenticate(ctx); err != nil {
		t.Fatalf("unexpected authentication error: %+v", err)
	}

	if cc.authToken != testToken {
		t.Errorf("expected [%s] got [%s]", testToken, cc.authToken)
	}
	if !cc.authTokenExpiry.Equal(testTokenExp) {
		t.Errorf("expected [%s] got [%s]", testTokenExp, cc.authTokenExpiry)
	}
}

func TestClient_loadAuthStorage(t *testing.T) {
	cc := &client{
		authToken:       testToken,
		authTokenExpiry: testTokenExp,
		email:           testEmail,
		password:        testPass,
		authDir:         "/tmp",
	}
	ctx := context.Background()
	cc.saveToken()
	cc.savePasswordCredentials(ctx, cc.email, cc.password)

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
		pbTime, _ := ptypes.TimestampProto(testTokenExp)
		resp := &api.LoginResponse{
			Token:  testToken,
			Expiry: pbTime,
		}
		marshaler := jsonpb.Marshaler{}
		encodedBody, err := marshaler.MarshalToString(resp)
		if err != nil {
			t.Fatal("could not marshal mock response")
		}
		w.Write([]byte(encodedBody))
	}))
}
