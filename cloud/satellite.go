package cloud

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	pb "github.com/earthly/cloud-api/compute"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/earthly/earthly/internal/version"
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

const (
	SatelliteSizeXSmall  = "xsmall"
	SatelliteSizeSmall   = "small"
	SatelliteSizeMedium  = "medium"
	SatelliteSizeLarge   = "large"
	SatelliteSizeXLarge  = "xlarge"
	SatelliteSize2XLarge = "2xlarge"
	SatelliteSize3XLarge = "3xlarge"
	SatelliteSize4XLarge = "4xlarge"
)

const (
	SatellitePlatformAMD64 = "linux/amd64"
	SatellitePlatformARM64 = "linux/arm64"
)

const DefaultSatelliteSize = SatelliteSizeMedium

// SatelliteInstance contains details about a remote Buildkit instance.
type SatelliteInstance struct {
	Name                    string
	Org                     string
	State                   string
	Platform                string
	Size                    string
	Version                 string
	VersionPinned           bool
	FeatureFlags            []string
	MaintenanceWindowStart  string
	MaintenanceWindowEnd    string
	MaintenanceWeekendsOnly bool
	RevisionID              int32
	Hidden                  bool
	LastUsed                time.Time
	CacheRetention          time.Duration
	Address                 string
	IsManaged               bool
	Certificate             *pb.TLSCertificate
	CloudName               string
}

func (c *Client) ListSatellites(ctx context.Context, orgName string, includeHidden bool) ([]SatelliteInstance, error) {
	orgID, err := c.GetOrgID(ctx, orgName)
	if err != nil {
		return nil, errors.Wrap(err, "failed listing satellites")
	}
	resp, err := c.compute.ListSatellites(c.withAuth(ctx), &pb.ListSatellitesRequest{
		OrgId:         orgID,
		IncludeHidden: includeHidden,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed listing satellites")
	}
	instances := make([]SatelliteInstance, len(resp.Instances))
	for i, s := range resp.Instances {
		instances[i] = SatelliteInstance{
			Name:           s.Name,
			Org:            orgID,
			Platform:       s.Platform,
			Size:           s.Size,
			State:          satelliteStatus(s.Status),
			Version:        s.Version,
			Hidden:         s.Hidden,
			LastUsed:       s.LastUsed.AsTime(),
			CacheRetention: s.CacheRetention.AsDuration(),
			CloudName:      s.CloudName,
		}
	}
	return instances, nil
}

func (c *Client) GetSatellite(ctx context.Context, name, orgName string) (*SatelliteInstance, error) {
	orgID, err := c.GetOrgID(ctx, orgName)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting satellites")
	}
	resp, err := c.compute.GetSatellite(c.withAuth(ctx), &pb.GetSatelliteRequest{
		OrgId: orgID,
		Name:  name,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed getting satellite")
	}
	return &SatelliteInstance{
		Name:                    name,
		Org:                     orgID,
		State:                   satelliteStatus(resp.Status),
		Platform:                resp.Platform,
		Size:                    resp.Size,
		Version:                 resp.Version,
		VersionPinned:           resp.VersionPinned,
		FeatureFlags:            resp.FeatureFlags,
		MaintenanceWindowStart:  resp.MaintenanceWindowStart,
		MaintenanceWindowEnd:    resp.MaintenanceWindowEnd,
		MaintenanceWeekendsOnly: resp.MaintenanceWeekendsOnly,
		RevisionID:              resp.RevisionId,
		Hidden:                  resp.Hidden,
		LastUsed:                resp.LastUsed.AsTime(),
		CacheRetention:          resp.CacheRetention.AsDuration(),
		IsManaged:               resp.IsManaged,
		Address:                 resp.SatelliteAddress,
		Certificate:             resp.Certificate,
	}, nil
}

func (c *Client) DeleteSatellite(ctx context.Context, name, orgName string, force bool) error {
	orgID, err := c.GetOrgID(ctx, orgName)
	if err != nil {
		return errors.Wrap(err, "failed deleting satellite")
	}
	_, err = c.compute.DeleteSatellite(c.withAuth(ctx), &pb.DeleteSatelliteRequest{
		OrgId: orgID,
		Name:  name,
		Force: force,
	})
	if err != nil {
		return errors.Wrap(err, "failed deleting satellite")
	}
	return nil
}

type LaunchSatelliteOpt struct {
	Name                    string
	OrgName                 string
	Size                    string
	Platform                string
	PinnedVersion           string
	MaintenanceWindowStart  string
	MaintenanceWeekendsOnly bool
	FeatureFlags            []string
	CloudName               string
}

func (c *Client) LaunchSatellite(ctx context.Context, opt LaunchSatelliteOpt) error {
	orgID, err := c.GetOrgID(ctx, opt.OrgName)
	if err != nil {
		return errors.Wrap(err, "failed launching satellite")
	}
	req := &pb.LaunchSatelliteRequest{
		OrgId:                   orgID,
		Name:                    opt.Name,
		Platform:                opt.Platform,
		Size:                    opt.Size,
		FeatureFlags:            opt.FeatureFlags,
		Version:                 opt.PinnedVersion,
		MaintenanceWindowStart:  opt.MaintenanceWindowStart,
		MaintenanceWeekendsOnly: opt.MaintenanceWeekendsOnly,
		CloudName:               opt.CloudName,
	}
	_, err = c.compute.LaunchSatellite(c.withAuth(ctx), req)
	if err != nil {
		return errors.Wrap(err, "failed launching satellite")
	}
	return nil
}

type SatelliteStatusUpdate struct {
	State string
	Err   error
}

func (c *Client) ReserveSatellite(ctx context.Context, name, orgName, gitAuthor, gitConfigEmail string, isCI bool) (out chan SatelliteStatusUpdate) {
	out = make(chan SatelliteStatusUpdate)
	go func() {
		orgID, err := c.GetOrgID(ctx, orgName)
		if err != nil {
			out <- SatelliteStatusUpdate{Err: errors.Wrap(err, "failed reserving satellite")}
			return
		}
		// Usually satellites reserve in 1-15 seconds, however, in some edge cases it will take longer.
		// It can take a minute if the satellite is actively falling asleep (it needs to finish, then wake back up).
		// In extreme cases, if a satellite update is running, the satellite can take around 6 minutes to finish.
		ctxTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
		defer cancel()
		defer close(out)
		// Note: we were having issues with the stream closing unexpectedly,
		// so we have wrapped this in a retry loop. This may be a temporary solution.
		const numRetries = 5
		var retriedError error
		for i := 1; i <= numRetries; i++ {
			stream, err := c.compute.ReserveSatellite(c.withRetryCount(c.withAuth(ctxTimeout), i), &pb.ReserveSatelliteRequest{
				OrgId:          orgID,
				Name:           name,
				CommitEmail:    gitAuthor,
				GitConfigEmail: gitConfigEmail,
				IsCi:           isCI,
				Metadata: &pb.ReserveSatelliteRequest_Metadata{
					EnvEntryNames: getEnvEntriesNames(),
					CliVersion:    version.Version,
				},
			})
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "retrying connection [attempt %d/%d]\n", i, numRetries)
				time.Sleep(time.Duration(i) * 2 * time.Second)
				continue
			}
			var lastStatus string
			for {
				update, err := stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					if isRetryable(err) {
						retriedError = err
						_, _ = fmt.Fprintf(os.Stderr, "retrying connection [attempt %d/%d]\n", i, numRetries)
						time.Sleep(time.Duration(i) * 2 * time.Second)
						break
					}
					out <- SatelliteStatusUpdate{Err: errors.Wrap(err, "failed receiving satellite reserve update")}
					return
				}
				lastStatus = satelliteStatus(update.Status)
				if lastStatus == SatelliteStatusFailed {
					out <- SatelliteStatusUpdate{Err: errors.New("satellite is in a failed state")}
					return
				}
				out <- SatelliteStatusUpdate{State: satelliteStatus(update.Status)}
			}
		}
		// max retries consumed
		out <- SatelliteStatusUpdate{Err: errors.Wrap(retriedError, "failed to retrieve satellite status")}
	}()
	return out
}

func isRetryable(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded:
		return true
	case codes.Internal:
		return true
	case codes.Unavailable:
		return true
	case codes.ResourceExhausted:
		return true
	case codes.FailedPrecondition:
		return true
	case codes.DataLoss:
		return true
	default:
		return false
	}
}

func (c *Client) WakeSatellite(ctx context.Context, name, orgName string) (out chan SatelliteStatusUpdate) {
	out = make(chan SatelliteStatusUpdate)
	go func() {
		orgID, err := c.GetOrgID(ctx, orgName)
		if err != nil {
			out <- SatelliteStatusUpdate{Err: errors.Wrap(err, "failed waking satellite")}
			return
		}
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
			out <- SatelliteStatusUpdate{State: status}
		}
	}()
	return out
}

func (c *Client) SleepSatellite(ctx context.Context, name, orgName string) (out chan SatelliteStatusUpdate) {
	out = make(chan SatelliteStatusUpdate)
	go func() {
		orgID, err := c.GetOrgID(ctx, orgName)
		if err != nil {
			out <- SatelliteStatusUpdate{Err: errors.Wrap(err, "failed opening satellite sleep stream")}
			return
		}
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
			out <- SatelliteStatusUpdate{State: status}
		}
	}()
	return out
}

type UpdateSatelliteOpt struct {
	Name                    string
	OrgName                 string
	PinnedVersion           string
	Size                    string
	Platform                string
	MaintenanceWindowStart  string
	MaintenanceWeekendsOnly bool
	DropCache               bool
	FeatureFlags            []string
}

func (c *Client) UpdateSatellite(ctx context.Context, opt UpdateSatelliteOpt) error {
	orgID, err := c.GetOrgID(ctx, opt.OrgName)
	if err != nil {
		return errors.Wrap(err, "failed listing satellites")
	}

	req := &pb.UpdateSatelliteRequest{
		OrgId:                   orgID,
		Name:                    opt.Name,
		Version:                 opt.PinnedVersion,
		DropCache:               opt.DropCache,
		FeatureFlags:            opt.FeatureFlags,
		MaintenanceWindowStart:  opt.MaintenanceWindowStart,
		MaintenanceWeekendsOnly: opt.MaintenanceWeekendsOnly,
		Size:                    opt.Size,
		Platform:                opt.Platform,
	}
	_, err = c.compute.UpdateSatellite(c.withAuth(ctx), req)
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

var validSizes = map[string]bool{
	SatelliteSizeXSmall:  true,
	SatelliteSizeSmall:   true,
	SatelliteSizeMedium:  true,
	SatelliteSizeLarge:   true,
	SatelliteSizeXLarge:  true,
	SatelliteSize2XLarge: true,
	SatelliteSize3XLarge: true,
	SatelliteSize4XLarge: true,
}

func ValidSatelliteSize(size string) bool {
	return validSizes[size]
}

var validPlatforms = map[string]bool{
	SatellitePlatformAMD64: true,
	SatellitePlatformARM64: true,
}

func ValidSatellitePlatform(size string) bool {
	return validPlatforms[size]
}

func getEnvEntriesNames() []string {
	environ := os.Environ()
	ret := make([]string, len(environ))
	for i := 0; i < len(environ); i++ {
		ret[i] = strings.SplitN(environ[i], "=", 2)[0]
	}
	sort.Strings(ret)
	return ret
}
