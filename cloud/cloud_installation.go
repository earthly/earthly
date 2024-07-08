package cloud

import (
	"context"

	pb "github.com/earthly/cloud-api/compute"
	"github.com/pkg/errors"
)

const (
	CloudStatusGreen   = "Green"
	CloudStatusYellow  = "Yellow"
	CloudStatusRed     = "Red"
	CloudStatusUnknown = "Unknown"
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
	AddressResolution  string
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
		//AddressResolution:  configuration.AddressResolution,
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

func (c *Client) UseCloud(ctx context.Context, orgID string, configuration *CloudConfigurationOpt) (*Installation, error) {
	resp, err := c.compute.UseCloud(c.withAuth(ctx), &pb.UseCloudRequest{
		OrgId: orgID,
		Name:  configuration.Name,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error from UseCloud API")
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
			StatusMessage: i.StatusContext,
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
	switch status {
	case pb.CloudStatus_CLOUD_STATUS_GREEN:
		return CloudStatusGreen
	case pb.CloudStatus_CLOUD_STATUS_YELLOW:
		return CloudStatusYellow
	case pb.CloudStatus_CLOUD_STATUS_RED:
		return CloudStatusRed
	default:
		return CloudStatusUnknown
	}
}
