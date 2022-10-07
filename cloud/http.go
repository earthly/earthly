package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
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
		marshaler := jsonpb.Marshaler{}
		encodedBody, err := marshaler.MarshalToString(body)
		if err != nil {
			return err
		}

		r.hasBody = true
		r.body = []byte(encodedBody)
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

func withBody(body string) requestOpt {
	return func(r *request) error {
		r.hasBody = true
		r.body = []byte(body)
		return nil
	}
}

func (c *client) doCall(ctx context.Context, method, url string, opts ...requestOpt) (int, string, error) {
	const maxAttempt = 10
	const maxSleepBeforeRetry = time.Second * 3

	var r request
	for _, opt := range opts {
		err := opt(&r)
		if err != nil {
			return 0, "", err
		}
	}

	alreadyReAuthed := false
	if r.hasAuth && time.Now().UTC().After(c.authTokenExpiry) {
		if err := c.Authenticate(ctx); err != nil {
			if errors.Is(err, ErrUnauthorized) {
				return 0, "", ErrUnauthorized
			}
			return 0, "", errors.Wrap(err, "failed refreshing expired auth token")
		}
		alreadyReAuthed = true
	}

	var status int
	var body string
	var callErr error
	duration := time.Millisecond * 100
	for attempt := 0; attempt < maxAttempt; attempt++ {
		status, body, callErr = c.doCallImp(ctx, r, method, url, opts...)
		retry, err := shouldRetry(status, body, callErr, c.warnFunc)
		if err != nil {
			return status, body, err
		}
		if !retry {
			return status, body, nil
		}

		if status == http.StatusUnauthorized {
			if !r.hasAuth || alreadyReAuthed {
				return status, body, ErrUnauthorized
			}
			err := c.Authenticate(ctx)
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

func shouldRetry(status int, body string, callErr error, warnFunc func(string, ...interface{})) (bool, error) {
	if status == http.StatusUnauthorized {
		return true, nil
	}
	if 500 <= status && status <= 599 {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			warnFunc("retrying http request due to unexpected status code %v", status)
		} else {
			warnFunc("retrying http request due to unexpected status code %v: %v", status, msg)
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
		warnFunc("retrying http request due to unexpected error %v", callErr)
		return true, nil
	}
}

func (c *client) doCallImp(ctx context.Context, r request, method, url string, opts ...requestOpt) (int, string, error) {
	var bodyReader io.Reader
	var bodyLen int64
	if r.hasBody {
		bodyReader = bytes.NewReader(r.body)
		bodyLen = int64(len(r.body))
	}

	req, err := http.NewRequestWithContext(ctx, method, c.httpAddr+url, bodyReader)
	if err != nil {
		return 0, "", err
	}
	if bodyReader != nil {
		req.ContentLength = bodyLen
	}
	if r.hasHeaders {
		req.Header = r.headers.Clone()
	}
	if r.hasAuth {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}

	respBody, err := readAllWithContext(ctx, resp.Body)
	if err != nil {
		return 0, "", err
	}
	return resp.StatusCode, string(respBody), nil
}

func readAllWithContext(ctx context.Context, r io.Reader) ([]byte, error) {
	var dt []byte
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
