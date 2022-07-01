package cloud

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func (c *client) Remove(ctx context.Context, path string) error {
	if path == "" || path[0] != '/' {
		return errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "DELETE", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to remove secret: %s", msg)
	}
	return nil
}

func (c *client) List(ctx context.Context, path string) ([]string, error) {
	if path != "" && !strings.HasSuffix(path, "/") {
		return nil, errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "GET", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to list secrets: %s", msg)
	}
	if body == "" {
		return []string{}, nil
	}
	return strings.Split(body, "\n"), nil
}

func (c *client) Get(ctx context.Context, path string) ([]byte, error) {
	if path == "" || path[0] != '/' || strings.HasSuffix(path, "/") {
		return nil, errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "GET", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to get secret: %s", msg)
	}
	return []byte(body), nil
}

func (c *client) Set(ctx context.Context, path string, data []byte) error {
	if path == "" || path[0] != '/' {
		return errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "PUT", fmt.Sprintf("/api/v0/secrets%s", path), withAuth(), withBody(string(data)))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to set secret: %s", msg)
	}
	return nil
}
