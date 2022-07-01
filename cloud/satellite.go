package cloud

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	pipelinesapi "github.com/earthly/cloud-api/pipelines"
	"github.com/pkg/errors"
)

// SatelliteInstance contains details about a remote Buildkit instance.
type SatelliteInstance struct {
	Name     string
	Org      string
	Status   string
	Version  string
	Platform string
}

func (c *client) ListSatellites(ctx context.Context, orgID string) ([]SatelliteInstance, error) {
	url := fmt.Sprintf("/api/v0/satellites?orgId=%s", url.QueryEscape(orgID))
	status, body, err := c.doCall(ctx, "GET", url, withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.Errorf("failed listing satellites: %s", body)
	}
	var resp pipelinesapi.ListSatellitesResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal listTokens response")
	}
	instances := make([]SatelliteInstance, len(resp.Instances))
	for i, s := range resp.Instances {
		instances[i] = SatelliteInstance{
			Name:     s.Name,
			Org:      orgID,
			Version:  s.Version,
			Platform: s.Platform,
			Status:   s.Status.String(),
		}
	}
	return instances, nil
}

func (c *client) GetSatellite(ctx context.Context, name, orgID string) (*SatelliteInstance, error) {
	url := fmt.Sprintf("/api/v0/satellites/%s?orgId=%s", name, url.QueryEscape(orgID))
	status, body, err := c.doCall(ctx, "GET", url, withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.Errorf("failed listing satellites: %s", body)
	}
	var resp pipelinesapi.GetSatelliteResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal listTokens response")
	}
	return &SatelliteInstance{
		Name:     name,
		Org:      orgID,
		Status:   resp.Status.String(),
		Version:  resp.Version,
		Platform: resp.Platform,
	}, nil
}

func (c *client) DeleteSatellite(ctx context.Context, name, orgID string) error {
	url := fmt.Sprintf("/api/v0/satellites/%s?orgId=%s", name, url.QueryEscape(orgID))
	status, body, err := c.doCall(ctx, "DELETE", url,
		withAuth(), withHeader("Grpc-Timeout", satelliteMgmtTimeout))
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return errors.Errorf("failed listing satellites: %s", body)
	}
	return nil
}

func (c *client) LaunchSatellite(ctx context.Context, name, orgID string) (*SatelliteInstance, error) {
	req := pipelinesapi.LaunchSatelliteRequest{
		OrgId:    orgID,
		Name:     name,
		Platform: "linux/amd64", // TODO support arm64 as well
	}
	status, body, err := c.doCall(ctx, "POST", "/api/v0/satellites",
		withAuth(), withHeader("Grpc-Timeout", satelliteMgmtTimeout), withJSONBody(&req))
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.Errorf("failed launching satellite: %s", body)
	}
	var resp pipelinesapi.LaunchSatelliteResponse
	err = c.jm.Unmarshal(bytes.NewReader([]byte(body)), &resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal LaunchSatellite response")
	}
	return &SatelliteInstance{
		Name:     name,
		Org:      orgID,
		Status:   resp.Status.String(),
		Version:  resp.Version,
		Platform: "linux/amd64",
	}, nil
}
