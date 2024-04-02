package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type request struct {
	hasBody bool
	body    []byte
	headers http.Header

	hasAuth    bool
	hasHeaders bool
}

type requestOpt func(*request) error

func withAuth() requestOpt {
	return func(r *request) error {
		r.hasAuth = true
		return nil
	}
}

func withHeader(key, value string) requestOpt {
	return func(r *request) error {
		r.hasHeaders = true
		r.headers = http.Header{}
		r.headers.Add(key, value)
		return nil
	}
}

func withJSONBody(body proto.Message) requestOpt {
	return func(r *request) error {
		encodedBody, err := protojson.Marshal(body)
		if err != nil {
			return err
		}

		r.hasBody = true
		r.body = encodedBody
		return nil
	}
}

func withFileBody(pathOnDisk string) requestOpt {
	return func(r *request) error {
		_, err := os.Stat(pathOnDisk)
		if err != nil {
			return errors.Wrapf(err, "could not stat file at %s", pathOnDisk)
		}

		contents, err := os.ReadFile(pathOnDisk)
		if err != nil {
			return errors.Wrapf(err, "could not add file %s to request body", pathOnDisk)
		}

		r.hasBody = true
		r.body = contents
		return nil
	}
}

func withBody(body []byte) requestOpt {
	return func(r *request) error {
		r.hasBody = true
		r.body = append([]byte{}, body...)
		return nil
	}
}

func (c *Client) doCall(ctx context.Context, method, url string, opts ...requestOpt) (int, []byte, error) {
	const maxAttempt = 10
	const maxSleepBeforeRetry = time.Second * 3

	var r request
	for _, opt := range opts {
		err := opt(&r)
		if err != nil {
			return 0, nil, err
		}
	}

	alreadyReAuthed := false
	if r.hasAuth && time.Now().UTC().After(c.authTokenExpiry) {
		if _, err := c.Authenticate(ctx); err != nil {
			if errors.Is(err, ErrUnauthorized) {
				return 0, nil, ErrUnauthorized
			}
			if errors.Is(err, ErrAuthTokenExpired) {
				return 0, nil, ErrAuthTokenExpired
			}
			return 0, nil, errors.Wrap(err, "failed refreshing expired auth token")
		}
		alreadyReAuthed = true
	}

	var status int
	body := []byte{}
	var callErr error
	duration := time.Millisecond * 100
	reqID := c.getRequestID()
	for attempt := 0; attempt < maxAttempt; attempt++ {
		status, body, callErr = c.doCallImp(ctx, r, method, url, reqID, opts...)
		retry, err := shouldRetry(status, body, callErr, c.debugFunc, reqID)
		if err != nil {
			return status, body, err
		}
		if !retry {
			return status, body, nil
		}

		if status == http.StatusUnauthorized {
			if !r.hasAuth || alreadyReAuthed {
				msg, err := getMessageFromJSON(bytes.NewReader(body))
				if err != nil || msg != tokenExpiredServerError {
					return status, body, ErrUnauthorized
				}
				return status, body, ErrAuthTokenExpired
			}
			_, err := c.Authenticate(ctx)
			if err != nil {
				return status, body, errors.Wrap(err, "auth credentials not valid")
			}
			alreadyReAuthed = true
		}

		if duration > maxSleepBeforeRetry {
			duration = maxSleepBeforeRetry
		}

		time.Sleep(duration)
		duration *= 2
	}

	return status, body, callErr
}

func shouldRetry(status int, body []byte, callErr error, debugFunc func(string, ...interface{}), reqID string) (bool, error) {
	if status == http.StatusUnauthorized {
		return true, nil
	}
	if 500 <= status && status <= 599 {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			debugFunc("retrying http request due to unexpected status code %v {reqID: %s}", status, reqID)
		} else {
			debugFunc("retrying http request due to unexpected status code %v: %v {reqID: %s}", status, msg, reqID)
		}
		return true, nil
	}
	switch {
	case callErr == nil:
		return false, nil
	case errors.Is(callErr, ErrNoAuthorizedPublicKeys):
		return false, callErr
	case errors.Is(callErr, ErrNoSSHAgent):
		return false, callErr
	case errors.Is(callErr, context.Canceled):
		return false, callErr
	case errors.Is(callErr, context.DeadlineExceeded):
		return false, callErr
	case strings.Contains(callErr.Error(), "failed to connect to ssh-agent"):
		return false, callErr
	default:
		debugFunc("retrying http request due to unexpected error %v {reqID: %s}", callErr, reqID)
		return true, nil
	}
}

func (c *Client) doCallImp(ctx context.Context, r request, method, url, reqID string, opts ...requestOpt) (int, []byte, error) {
	var bodyReader io.Reader
	var bodyLen int64
	if r.hasBody {
		bodyReader = bytes.NewReader(r.body)
		bodyLen = int64(len(r.body))
	}

	req, err := http.NewRequestWithContext(ctx, method, c.httpAddr+url, bodyReader)
	if err != nil {
		return 0, nil, err
	}
	if bodyReader != nil {
		req.ContentLength = bodyLen
	}
	if r.hasHeaders {
		req.Header = r.headers.Clone()
	}
	if r.hasAuth {
		if c.authToken == "" {
			return 0, nil, ErrUnauthorized
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	}
	req.Header.Add(requestID, reqID)

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: c.serverConnTimeout,
			}).DialContext,
			Proxy: http.ProxyFromEnvironment,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	respBody, err := readAllWithContext(ctx, resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, respBody, nil
}

func readAllWithContext(ctx context.Context, r io.Reader) ([]byte, error) {
	dt := []byte{}
	var readErr error
	ch := make(chan struct{})
	go func() {
		dt, readErr = io.ReadAll(r)
		close(ch)
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-ch:
		return dt, readErr
	}
}

func getMessageFromJSON(r io.Reader) (string, error) {
	decoder := json.NewDecoder(r)
	msg := struct {
		Message string `json:"message"`
	}{}
	err := decoder.Decode(&msg)
	if err != nil {
		return "", err
	}
	return msg.Message, nil
}
