package cloud

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/earthly/cloud-api/secrets"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/durationpb"
)

func (c *Client) Remove(ctx context.Context, path string) error {
	if path == "" || path[0] != '/' {
		return errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "DELETE", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to remove secret: %s", msg)
	}
	return nil
}

func (c *Client) List(ctx context.Context, path string) ([]string, error) {
	if path != "" && !strings.HasSuffix(path, "/") {
		return nil, errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "GET", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to list secrets: %s", msg)
	}
	if len(body) == 0 {
		return []string{}, nil
	}
	return strings.Split(string(body), "\n"), nil
}

func (c *Client) Get(ctx context.Context, path string) ([]byte, error) {
	if path == "" || path[0] != '/' || strings.HasSuffix(path, "/") {
		return nil, errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "GET", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to get secret: %s", msg)
	}
	return body, nil
}

func (c *Client) Set(ctx context.Context, path string, data []byte) error {
	if path == "" || path[0] != '/' {
		return errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "PUT", fmt.Sprintf("/api/v0/secrets%s", path), withAuth(), withBody(data))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to set secret: %s", msg)
	}
	return nil
}

func (c *Client) GetAWSCredentials(ctx context.Context, sessionName string, roleARN string, orgName string, projectName string, region string, sessionDuration *time.Duration) (*secrets.GetAWSCredentialsResponse, error) {
	if orgName == "" {
		return nil, errors.New("org must be set in order to use AWS OIDC")
	}
	if projectName == "" {
		return nil, errors.New("project must be set in order to use AWS OIDC")
	}
	var duration *durationpb.Duration
	if sessionDuration != nil {
		duration = durationpb.New(*sessionDuration)
	}
	response, err := c.secrets.GetAWSCredentials(c.withAuth(ctx), &secrets.GetAWSCredentialsRequest{
		RoleArn:         roleARN,
		SessionName:     sessionName,
		SessionDuration: duration,
		Region:          region,
		OrgName:         orgName,
		ProjectName:     projectName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get aws credentials via oidc provider")
	}
	return response, nil
}
