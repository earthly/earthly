package cloud

import (
	"context"

	pb "github.com/earthly/cloud-api/compute"
	"github.com/pkg/errors"
)

const (
	CloudStatusConnected = "Connected"
	CloudStatusActive    = "Active"
	CloudStatusProblem   = "Problem"
)

type Installation struct {
	Name          string
	Org           string
	Status        string
	NumSatellites int
	IsDefault     bool
}

func (c *Client) ConfigureCloud(ctx context.Context, orgID, cloudName string, setDefault bool) (*Installation, error) {
	resp, err := c.compute.ConfigureCloud(c.withAuth(ctx), &pb.ConfigureCloudRequest{
		OrgId:      orgID,
		Name:       cloudName,
		SetDefault: setDefault,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error from ConfigureCloud API")
	}
	return &Installation{
		Name:   cloudName,
		Org:    orgID,
		Status: installationStatus(resp.Status),
	}, nil
}

func (c *Client) ListClouds(ctx context.Context, orgID string) ([]Installation, error) {
	resp, err := c.compute.ListCloud(c.withAuth(ctx), &pb.ListCloudRequest{
		OrgId: orgID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error from ListCloud API")
	}
	var installations []Installation
	for _, i := range resp.Clouds {
		installations = append(installations, Installation{
			Name:          i.CloudName,
			Org:           orgID,
			Status:        installationStatus(i.Status),
			NumSatellites: int(i.NumSatellites),
			IsDefault:     i.IsDefault,
		})
	}
	return installations, nil
}

func (c *Client) DeleteCloud(ctx context.Context, orgID, cloudName string) error {
	_, err := c.compute.DeleteCloud(c.withAuth(ctx), &pb.DeleteCloudRequest{
		Name:  cloudName,
		OrgId: orgID,
	})
	if err != nil {
		return errors.Wrap(err, "error from DeleteCloud API")
	}
	return nil
}

func installationStatus(status pb.CloudStatus) string {
	internalStatus := "UNKNOWN"
	switch status {
	case pb.CloudStatus_CLOUD_STATUS_ACCOUNT_ACTIVE:
		internalStatus = CloudStatusActive
	case pb.CloudStatus_CLOUD_STATUS_ACCOUNT_CONNECTED:
		internalStatus = CloudStatusConnected
	case pb.CloudStatus_CLOUD_STATUS_PROBLEM:
		internalStatus = CloudStatusProblem
	}
	return internalStatus
}
