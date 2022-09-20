package cloud

import (
	"context"
	"io"
	"time"

	pb "github.com/earthly/cloud-api/pipelines"
	"github.com/pkg/errors"
)

const (
	SatelliteStatusOperational = "Operational"
	SatelliteStatusSleep       = "Sleeping"
	SatelliteStatusStarting    = "Starting"
	SatelliteStatusCreating    = "Creating"
	SatelliteStatusFailed      = "Failed"
	SatelliteStatusDestroying  = "Destroying"
	SatelliteStatusOffline     = "Offline"
	SatelliteStatusUnknown     = "Unknown"
)

// SatelliteInstance contains details about a remote Buildkit instance.
type SatelliteInstance struct {
	Name     string
	Org      string
	Status   string
	Platform string
}

func (c *client) ListSatellites(ctx context.Context, orgID string) ([]SatelliteInstance, error) {
	resp, err := c.pipelines.ListSatellites(c.withAuth(ctx), &pb.ListSatellitesRequest{
		OrgId: orgID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed listing satellites")
	}
	instances := make([]SatelliteInstance, len(resp.Instances))
	for i, s := range resp.Instances {
		instances[i] = SatelliteInstance{
			Name:     s.Name,
			Org:      orgID,
			Platform: s.Platform,
			Status:   satelliteStatus(s.Status),
		}
	}
	return instances, nil
}

func (c *client) GetSatellite(ctx context.Context, name, orgID string) (*SatelliteInstance, error) {
	resp, err := c.pipelines.GetSatellite(c.withAuth(ctx), &pb.GetSatelliteRequest{
		OrgId: orgID,
		Name:  name,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed getting satellite")
	}
	return &SatelliteInstance{
		Name:     name,
		Org:      orgID,
		Status:   satelliteStatus(resp.Status),
		Platform: resp.Platform,
	}, nil
}

func (c *client) DeleteSatellite(ctx context.Context, name, orgID string) error {
	_, err := c.pipelines.DeleteSatellite(c.withAuth(ctx), &pb.DeleteSatelliteRequest{
		OrgId: orgID,
		Name:  name,
	})
	if err != nil {
		return errors.Wrap(err, "failed deleting satellite")
	}
	return nil
}

func (c *client) LaunchSatellite(ctx context.Context, name, orgID string, features []string) error {
	_, err := c.pipelines.LaunchSatellite(c.withAuth(ctx), &pb.LaunchSatelliteRequest{
		OrgId:        orgID,
		Name:         name,
		Platform:     "linux/amd64", // TODO support arm64 as well
		FeatureFlags: features,
	})
	if err != nil {
		return errors.Wrap(err, "failed launching satellite")
	}
	return nil
}

func (c *client) ReserveSatellite(ctx context.Context, name, orgID string, out chan<- string) error {
	defer close(out)
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	stream, err := c.pipelines.ReserveSatellite(c.withAuth(ctxTimeout), &pb.ReserveSatelliteRequest{
		OrgId: orgID,
		Name:  name,
	})
	if err != nil {
		return errors.Wrap(err, "failed opening satellite reserve stream")
	}
	var update *pb.ReserveSatelliteResponse
	for {
		update, err = stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "failed receiving satellite reserve update")
		}
		out <- satelliteStatus(update.Status)
	}
}

func satelliteStatus(status pb.SatelliteStatus) string {
	switch status {
	case pb.SatelliteStatus_SATELLITE_STATUS_OPERATIONAL:
		return SatelliteStatusOperational
	case pb.SatelliteStatus_SATELLITE_STATUS_SLEEP:
		return SatelliteStatusSleep
	case pb.SatelliteStatus_SATELLITE_STATUS_STARTING:
		return SatelliteStatusStarting
	case pb.SatelliteStatus_SATELLITE_STATUS_CREATING:
		return SatelliteStatusCreating
	case pb.SatelliteStatus_SATELLITE_STATUS_FAILED:
		return SatelliteStatusFailed
	case pb.SatelliteStatus_SATELLITE_STATUS_DESTROYING:
		return SatelliteStatusDestroying
	case pb.SatelliteStatus_SATELLITE_STATUS_OFFLINE:
		return SatelliteStatusOffline
	default:
		return SatelliteStatusUnknown
	}
}
