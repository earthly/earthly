package cloud

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"time"

	pb "github.com/earthly/cloud-api/compute"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	// SatelliteStatusOperational indicates an on satellite that is ready to accept connections.
	SatelliteStatusOperational = "Operational"
	// SatelliteStatusSleep indicates a satellite that is in a sleep state.
	SatelliteStatusSleep = "Sleeping"
	// SatelliteStatusStarting indicates a satellite that is waking from a sleep state.
	SatelliteStatusStarting = "Starting"
	// SatelliteStatusStopping indicates a new satellite that is currently going to sleep.
	SatelliteStatusStopping = "Stopping"
	// SatelliteStatusCreating indicates a new satellite that is currently being launched.
	SatelliteStatusCreating = "Creating"
	// SatelliteStatusUpdating indicates a satellite that is upgrading to a new version, either manually or via maintenance window.
	SatelliteStatusUpdating = "Updating"
	// SatelliteStatusFailed indicates a satellite that has crashed and cannot be used.
	SatelliteStatusFailed = "Failed"
	// SatelliteStatusDestroying indicates a satellite that is actively being deleted.
	SatelliteStatusDestroying = "Destroying"
	// SatelliteStatusOffline indicates a satellite that has been stopped and will not be woken up normally via build.
	SatelliteStatusOffline = "Offline"
	// SatelliteStatusUnknown is used when an unexpected satellite status is returned by the server.
	SatelliteStatusUnknown = "Unknown"
)

// SatelliteInstance contains details about a remote Buildkit instance.
type SatelliteInstance struct {
	Name              string
	Org               string
	State             string
	Platform          string
	Size              string
	Version           string
	FeatureFlags      []string
	MaintenanceWindow string
}

func (c *client) ListSatellites(ctx context.Context, orgID string) ([]SatelliteInstance, error) {
	resp, err := c.compute.ListSatellites(c.withAuth(ctx), &pb.ListSatellitesRequest{
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
			Size:     s.Size,
			State:    satelliteStatus(s.Status),
			Version:  s.Version,
		}
	}
	return instances, nil
}

func (c *client) GetSatellite(ctx context.Context, name, orgID string) (*SatelliteInstance, error) {
	resp, err := c.compute.GetSatellite(c.withAuth(ctx), &pb.GetSatelliteRequest{
		OrgId: orgID,
		Name:  name,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed getting satellite")
	}
	return &SatelliteInstance{
		Name:              name,
		Org:               orgID,
		State:             satelliteStatus(resp.Status),
		Platform:          resp.Platform,
		Size:              resp.Size,
		Version:           resp.Version,
		FeatureFlags:      resp.FeatureFlags,
		MaintenanceWindow: resp.MaintenanceWindowStart,
	}, nil
}

func (c *client) DeleteSatellite(ctx context.Context, name, orgID string) error {
	_, err := c.compute.DeleteSatellite(c.withAuth(ctx), &pb.DeleteSatelliteRequest{
		OrgId: orgID,
		Name:  name,
	})
	if err != nil {
		return errors.Wrap(err, "failed deleting satellite")
	}
	return nil
}

func (c *client) LaunchSatellite(ctx context.Context, name, orgID, platform, size, version, maintenanceWindow string, features []string) error {
	req := &pb.LaunchSatelliteRequest{
		OrgId:                  orgID,
		Name:                   name,
		Platform:               platform,
		Size:                   size,
		FeatureFlags:           features,
		Version:                version,
		MaintenanceWindowStart: maintenanceWindow,
	}
	_, err := c.compute.LaunchSatellite(c.withAuth(ctx), req)
	if err != nil {
		return errors.Wrap(err, "failed launching satellite")
	}
	return nil
}

type SatelliteStatusUpdate struct {
	State string
	Err   error
}

func (c *client) ReserveSatellite(ctx context.Context, name, orgID, gitAuthor, gitConfigEmail string, isCI bool) (out chan SatelliteStatusUpdate) {
	out = make(chan SatelliteStatusUpdate)
	go func() {
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		defer close(out)
		stream, err := c.compute.ReserveSatellite(c.withAuth(ctxTimeout), &pb.ReserveSatelliteRequest{
			OrgId:          orgID,
			Name:           name,
			CommitEmail:    gitAuthor,
			GitConfigEmail: gitConfigEmail,
			IsCi:           isCI,
		})
		if err != nil {
			out <- SatelliteStatusUpdate{Err: errors.Wrap(err, "failed opening satellite reserve stream")}
			return
		}
		for {
			update, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				out <- SatelliteStatusUpdate{Err: errors.Wrap(err, "failed receiving satellite reserve update")}
				return
			}
			status := satelliteStatus(update.Status)
			if status == SatelliteStatusFailed {
				out <- SatelliteStatusUpdate{Err: errors.New("satellite is in a failed state")}
				return
			}
			out <- SatelliteStatusUpdate{State: satelliteStatus(update.Status)}
		}
	}()
	return out
}

func (c *client) WakeSatellite(ctx context.Context, name, orgID string) (out chan SatelliteStatusUpdate) {
	out = make(chan SatelliteStatusUpdate)
	go func() {
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		defer close(out)
		stream, err := c.compute.WakeSatellite(c.withAuth(ctxTimeout), &pb.WakeSatelliteRequest{
			OrgId: orgID,
			Name:  name,
		})
		if err != nil {
			out <- SatelliteStatusUpdate{Err: errors.Wrap(err, "failed opening satellite wake stream")}
			return
		}
		for {
			update, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				out <- SatelliteStatusUpdate{Err: errors.Wrap(err, "failed receiving satellite wake update")}
				return
			}
			status := satelliteStatus(update.Status)
			if status == SatelliteStatusFailed {
				out <- SatelliteStatusUpdate{Err: errors.New("satellite is in a failed state")}
				return
			}
			out <- SatelliteStatusUpdate{State: satelliteStatus(update.Status)}
		}
	}()
	return out
}

func (c *client) SleepSatellite(ctx context.Context, name, orgID string) (out chan SatelliteStatusUpdate) {
	out = make(chan SatelliteStatusUpdate)
	go func() {
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		defer close(out)
		stream, err := c.compute.SleepSatellite(c.withAuth(ctxTimeout), &pb.SleepSatelliteRequest{
			OrgId:          orgID,
			Name:           name,
			UpdateInterval: durationpb.New(10 * time.Second),
		})
		if err != nil {
			out <- SatelliteStatusUpdate{Err: errors.Wrap(err, "failed opening satellite sleep stream")}
			return
		}
		for {
			update, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				out <- SatelliteStatusUpdate{Err: errors.Wrap(err, "failed receiving satellite sleep update")}
				return
			}
			status := satelliteStatus(update.Status)
			if status == SatelliteStatusFailed {
				out <- SatelliteStatusUpdate{Err: errors.New("satellite is in a failed state")}
				return
			}
			out <- SatelliteStatusUpdate{State: satelliteStatus(update.Status)}
		}
	}()
	return out
}

func (c *client) UpdateSatellite(ctx context.Context, name, orgID, version, maintenanceWindow string, dropCache bool, featureFlags []string) error {
	req := &pb.UpdateSatelliteRequest{
		OrgId:                  orgID,
		Name:                   name,
		Version:                version,
		DropCache:              dropCache,
		FeatureFlags:           featureFlags,
		MaintenanceWindowStart: maintenanceWindow,
	}
	_, err := c.compute.UpdateSatellite(c.withAuth(ctx), req)
	if err != nil {
		return errors.Wrap(err, "failed getting satellite")
	}
	return nil
}

var maintenanceWindowRx = regexp.MustCompile(`[0-9]{2}:[0-9]{2}\z`)

// LocalMaintenanceWindowToUTC checks if the provided maintenance window is valid
// and returns a new maintenance window converted from local time to UTC format.
func LocalMaintenanceWindowToUTC(window string, loc *time.Location) (string, error) {
	if !maintenanceWindowRx.MatchString(window) {
		return "", errors.New("maintenance window must be in the format HH:MM (24hr)")
	}
	t, err := time.ParseInLocation("15:04:05", fmt.Sprintf("%s:00", window), loc)
	if err != nil {
		return "", errors.Wrap(err, "failed parsing maintenance window")
	}
	return t.UTC().Format("15:04"), nil
}

// UTCMaintenanceWindowToLocal checks if the provided maintenance window is valid
// and returns a new maintenance window converted from local time to UTC format.
func UTCMaintenanceWindowToLocal(window string, loc *time.Location) (string, error) {
	t, err := time.ParseInLocation("15:04:05", fmt.Sprintf("%s:00", window), time.UTC)
	if err != nil {
		return "", errors.Wrap(err, "failed parsing maintenance window")
	}
	return t.In(loc).Format("15:04"), nil
}

func satelliteStatus(status pb.SatelliteStatus) string {
	switch status {
	case pb.SatelliteStatus_SATELLITE_STATUS_OPERATIONAL:
		return SatelliteStatusOperational
	case pb.SatelliteStatus_SATELLITE_STATUS_SLEEP:
		return SatelliteStatusSleep
	case pb.SatelliteStatus_SATELLITE_STATUS_STARTING:
		return SatelliteStatusStarting
	case pb.SatelliteStatus_SATELLITE_STATUS_STOPPING:
		return SatelliteStatusStopping
	case pb.SatelliteStatus_SATELLITE_STATUS_CREATING:
		return SatelliteStatusCreating
	case pb.SatelliteStatus_SATELLITE_STATUS_UPDATING:
		return SatelliteStatusUpdating
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
