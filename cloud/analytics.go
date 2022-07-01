package cloud

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// EarthlyAnalytics is the payload used in SendAnalytics.
// It contains information about the command that was run,
// the environment it was run in, and the result of the command.
type EarthlyAnalytics struct {
	Key              string                    `json:"key"`
	InstallID        string                    `json:"install_id"`
	Version          string                    `json:"version"`
	Platform         string                    `json:"platform"`
	GitSHA           string                    `json:"git_sha"`
	ExitCode         int                       `json:"exit_code"`
	CI               string                    `json:"ci_name"`
	RepoHash         string                    `json:"repo_hash"`
	ExecutionSeconds float64                   `json:"execution_seconds"`
	Terminal         bool                      `json:"terminal"`
	Counts           map[string]map[string]int `json:"counts"`
}

// SendAnalytics send an analytics event to the Cloud server.
func (c *client) SendAnalytics(ctx context.Context, data *EarthlyAnalytics) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}
	opts := []requestOpt{
		withBody(string(payload)),
		withHeader("Content-Type", "application/json; charset=utf-8"),
	}
	if c.IsLoggedIn(ctx) {
		opts = append(opts, withAuth())
	}
	status, _, err := c.doCall(ctx, "PUT", "/analytics", opts...)
	if err != nil {
		return errors.Wrap(err, "failed sending analytics")
	}
	if status != http.StatusCreated {
		return errors.Errorf("unexpected response from analytics server: %d", status)
	}
	return nil
}
