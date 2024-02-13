package cloud

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	logsapi "github.com/earthly/cloud-api/logs"
	"github.com/pkg/errors"
)

func (c *Client) UploadLog(ctx context.Context, pathOnDisk string) (string, error) {
	status, body, err := c.doCall(ctx, http.MethodPost, "/api/v0/logs", withAuth(), withFileBody(pathOnDisk), withHeader("Content-Type", "application/gzip"))
	if err != nil {
		return "", err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return "", errors.Errorf("failed to upload log: %s", msg)
	}

	var uploadBundleResponse logsapi.UploadLogBundleResponse
	err = c.jum.Unmarshal(body, &uploadBundleResponse)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal uploadbundle response")
	}

	return fmt.Sprintf(uploadBundleResponse.ViewURL), nil
}
