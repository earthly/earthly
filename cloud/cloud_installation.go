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
	StatusMessage string
	NumSatellites int
	IsDefault     bool
}

type CloudConfigurationOpt struct {
	Name               string
	SetDefault         bool
	SshKeyName         string
	ComputeRoleArn     string
	AccountId          string
	AllowedSubnetIds   []string
	SecurityGroupId    string
	Region             string
	InstanceProfileArn string
}

func (c *Client) ConfigureCloud(ctx context.Context, orgID string, configuration *CloudConfigurationOpt) (*Installation, error) {
	resp, err := c.compute.ConfigureCloud(c.withAuth(ctx), &pb.ConfigureCloudRequest{
		OrgId:              orgID,
		Name:               configuration.Name,
		SetDefault:         configuration.SetDefault,
		SshKeyName:         configuration.SshKeyName,
		ComputeRoleArn:     configuration.ComputeRoleArn,
		AccountId:          configuration.AccountId,
		AllowedSubnetIds:   configuration.AllowedSubnetIds,
		SecurityGroupId:    configuration.SecurityGroupId,
		Region:             configuration.Region,
		InstanceProfileArn: configuration.InstanceProfileArn,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error from ConfigureCloud API")
	}
	return &Installation{
		Name:          configuration.Name,
		Org:           orgID,
		Status:        installationStatus(resp.Status),
		StatusMessage: resp.Message,
	}, nil
}

func (c *Client) ListClouds(ctx context.Context, orgID string) ([]Installation, error) {
	resp, err := c.compute.ListClouds(c.withAuth(ctx), &pb.ListCloudsRequest{
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
