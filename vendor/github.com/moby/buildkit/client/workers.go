package client

import (
	"context"
	"time"

	controlapi "github.com/moby/buildkit/api/services/control"
	apitypes "github.com/moby/buildkit/api/types"
	"github.com/moby/buildkit/solver/pb"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// WorkerInfo contains information about a worker
type WorkerInfo struct {
	ID              string              `json:"id"`
	Labels          map[string]string   `json:"labels"`
	Platforms       []ocispecs.Platform `json:"platforms"`
	GCPolicy        []PruneInfo         `json:"gcPolicy"`
	BuildkitVersion BuildkitVersion     `json:"buildkitVersion"`

	// Earthly-specific.
	ParallelismCurrent int `json:"parallelismCurrent"`
	ParallelismMax     int `json:"parallelismMax"`
	ParallelismWaiting int `json:"parallelismWaiting"`
	GCAnalytics        GCAnalytics
}

type GCAnalytics struct {
	// Summary of last numRuns.
	NumRuns           int
	NumFailures       int
	AvgDuration       time.Duration
	AvgRecordsCleared int64
	AvgSizeCleared    int64
	AvgRecordsBefore  int64
	AvgSizeBefore     int64
	// All-time summary.
	AllTimeRuns        int64
	AllTimeMaxDuration time.Duration
	AllTimeDuration    time.Duration
	// Current run (if one is ongoing).
	CurrentStartTime        *time.Time
	CurrentNumRecordsBefore int64
	CurrentSizeBefore       int64
	// Last run.
	LastStartTime         *time.Time
	LastEndTime           *time.Time
	LastNumRecordsBefore  int64
	LastSizeBefore        int64
	LastNumRecordsCleared int64
	LastSizeCleared       int64
	LastSuccess           bool
}

// ListWorkers lists all active workers
func (c *Client) ListWorkers(ctx context.Context, opts ...ListWorkersOption) ([]*WorkerInfo, error) {
	info := &ListWorkersInfo{}
	for _, o := range opts {
		o.SetListWorkersOption(info)
	}

	req := &controlapi.ListWorkersRequest{Filter: info.Filter}
	resp, err := c.ControlClient().ListWorkers(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list workers")
	}

	var wi []*WorkerInfo

	for _, w := range resp.Record {
		var currentStartTime *time.Time
		if w.GetGCAnalytics().GetCurrentStartTimeSecEpoch() != 0 {
			t := time.Unix(w.GetGCAnalytics().GetCurrentStartTimeSecEpoch(), 0)
			currentStartTime = &t
		}
		var lastStartTime *time.Time
		if w.GetGCAnalytics().GetLastStartTimeSecEpoch() != 0 {
			t := time.Unix(w.GetGCAnalytics().GetLastStartTimeSecEpoch(), 0)
			lastStartTime = &t
		}
		var lastEndTime *time.Time
		if w.GetGCAnalytics().GetLastEndTimeSecEpoch() != 0 {
			t := time.Unix(w.GetGCAnalytics().GetLastEndTimeSecEpoch(), 0)
			lastEndTime = &t
		}
		wi = append(wi, &WorkerInfo{
			ID:              w.ID,
			Labels:          w.Labels,
			Platforms:       pb.ToSpecPlatforms(w.Platforms),
			GCPolicy:        fromAPIGCPolicy(w.GCPolicy),
			BuildkitVersion: fromAPIBuildkitVersion(w.BuildkitVersion),

			ParallelismCurrent: int(w.ParallelismCurrent),
			ParallelismMax:     int(w.ParallelismMax),
			ParallelismWaiting: int(w.ParallelismWaiting),

			GCAnalytics: GCAnalytics{
				NumRuns:                 int(w.GetGCAnalytics().GetNumRuns()),
				NumFailures:             int(w.GetGCAnalytics().GetNumFailures()),
				AvgDuration:             time.Duration(w.GetGCAnalytics().GetAvgDurationMs()) * time.Millisecond,
				AvgRecordsCleared:       w.GetGCAnalytics().GetAvgRecordsCleared(),
				AvgSizeCleared:          w.GetGCAnalytics().GetAvgSizeCleared(),
				AvgRecordsBefore:        w.GetGCAnalytics().GetAvgRecordsBefore(),
				AvgSizeBefore:           w.GetGCAnalytics().GetAvgSizeBefore(),
				AllTimeRuns:             w.GetGCAnalytics().GetAllTimeRuns(),
				AllTimeMaxDuration:      time.Duration(w.GetGCAnalytics().GetAllTimeMaxDurationMs()) * time.Millisecond,
				AllTimeDuration:         time.Duration(w.GetGCAnalytics().GetAllTimeDurationMs()) * time.Millisecond,
				CurrentStartTime:        currentStartTime,
				CurrentNumRecordsBefore: w.GetGCAnalytics().GetCurrentNumRecordsBefore(),
				CurrentSizeBefore:       w.GetGCAnalytics().GetCurrentSizeBefore(),
				LastStartTime:           lastStartTime,
				LastEndTime:             lastEndTime,
				LastNumRecordsBefore:    w.GetGCAnalytics().GetLastNumRecordsBefore(),
				LastSizeBefore:          w.GetGCAnalytics().GetLastSizeBefore(),
				LastNumRecordsCleared:   w.GetGCAnalytics().GetLastNumRecordsCleared(),
				LastSizeCleared:         w.GetGCAnalytics().GetLastSizeCleared(),
				LastSuccess:             w.GetGCAnalytics().GetLastSuccess(),
			},
		})
	}

	return wi, nil
}

// ListWorkersOption is an option for a worker list query
type ListWorkersOption interface {
	SetListWorkersOption(*ListWorkersInfo)
}

// ListWorkersInfo is a payload for worker list query
type ListWorkersInfo struct {
	Filter []string
}

func fromAPIGCPolicy(in []*apitypes.GCPolicy) []PruneInfo {
	out := make([]PruneInfo, 0, len(in))
	for _, p := range in {
		out = append(out, PruneInfo{
			All:          p.All,
			Filter:       p.Filters,
			KeepDuration: time.Duration(p.KeepDuration),
			KeepBytes:    p.KeepBytes,
		})
	}
	return out
}
