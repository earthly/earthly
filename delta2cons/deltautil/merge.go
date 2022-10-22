package deltautil

import (
	"fmt"

	pb "github.com/earthly/cloud-api/logstream"
)

const Version = 2

// SimplifyDeltas takes a delta and flattens it.
func SimplifyDeltas(delta *pb.Delta, nextManifestOrderID int64) (*pb.Delta, int64) {
	if delta.GetVersion() != Version {
		panic(fmt.Errorf("unsupported delta version %d", delta.GetVersion()))
	}
	ret := &pb.Delta{
		Version: Version,
	}
	ret.DeltaManifests, nextManifestOrderID = mergeDeltaManifests(
		delta.GetDeltaManifests(), nextManifestOrderID)
	ret.DeltaLogs = mergeDeltaLogs(delta.GetDeltaLogs())
	return ret, nextManifestOrderID
}

func mergeDeltaManifests(dms []*pb.DeltaManifest, nextManifestOrderID int64) ([]*pb.DeltaManifest, int64) {
	if len(dms) == 0 {
		// Optimization.
		return nil, nextManifestOrderID
	}
	var ret []*pb.DeltaManifest
	var reset *pb.BuildManifest
	postResetIndex := 0
	for index, dm := range dms {
		switch dm.GetDeltaManifestOneof().(type) {
		case *pb.DeltaManifest_ResetAll:
			reset = dm.GetResetAll()
			postResetIndex = index + 1
		default:
		}
	}
	if reset != nil {
		ret = append(ret, &pb.DeltaManifest{
			OrderId:            nextManifestOrderID,
			DeltaManifestOneof: &pb.DeltaManifest_ResetAll{ResetAll: reset},
		})
		nextManifestOrderID++
	}
	if postResetIndex == len(dms) {
		return ret, nextManifestOrderID
	}
	nextIsEof := dms[postResetIndex].GetEof()
	if nextIsEof {
		ret = append(ret, &pb.DeltaManifest{
			OrderId:            nextManifestOrderID,
			DeltaManifestOneof: &pb.DeltaManifest_Eof{Eof: true},
		})
		nextManifestOrderID++
		return ret, nextManifestOrderID
	}
	mergedF := &pb.DeltaManifest_FieldsDelta{}
	merged := &pb.DeltaManifest{
		OrderId: nextManifestOrderID,
		DeltaManifestOneof: &pb.DeltaManifest_Fields{
			Fields: mergedF,
		},
	}
	nextManifestOrderID++
	ret = append(ret, merged)
	for index := postResetIndex; index < len(dms); index++ {
		dm := dms[index]
		switch dm.GetDeltaManifestOneof().(type) {
		case *pb.DeltaManifest_Eof:
			eofDelta := &pb.DeltaManifest{
				OrderId:            nextManifestOrderID,
				DeltaManifestOneof: &pb.DeltaManifest_Eof{Eof: true},
			}
			nextManifestOrderID++
			ret = append(ret, eofDelta)
		case *pb.DeltaManifest_Fields:
			f := dm.GetFields()
			if f.GetStartedAt() != 0 {
				mergedF.StartedAt = f.GetStartedAt()
			}
			if f.GetFinishedAt() != 0 {
				mergedF.FinishedAt = f.GetFinishedAt()
			}
			if f.GetStatus() != pb.BuildStatus_BUILD_STATUS_UNKNOWN {
				mergedF.Status = f.GetStatus()
			}
			if f.GetFailedTarget() != "" {
				mergedF.FailedTarget = f.GetFailedTarget()
			}
			if f.GetFailedSummary() != "" {
				mergedF.FailedSummary = f.GetFailedSummary()
			}
			if f.GetTargets() != nil {
				if mergedF.Targets == nil {
					mergedF.Targets = make(map[string]*pb.DeltaTargetManifest)
				}
				for targetID, t := range f.GetTargets() {
					mt, found := mergedF.Targets[targetID]
					if !found {
						mt = new(pb.DeltaTargetManifest)
						mergedF.Targets[targetID] = mt
					}
					if t.GetName() != "" {
						mt.Name = t.GetName()
					}
					if len(t.GetOverrideArgs()) != 0 {
						mt.OverrideArgs = append([]string{}, t.GetOverrideArgs()...)
					}
					if t.GetPlatform() != "" {
						mt.Platform = t.GetPlatform()
					}
					if t.GetStatus() != pb.BuildStatus_BUILD_STATUS_UNKNOWN {
						mt.Status = t.GetStatus()
					}
					if t.GetStartedAt() != 0 {
						mt.StartedAt = t.GetStartedAt()
					}
					if t.GetFinishedAt() != 0 {
						mt.FinishedAt = t.GetFinishedAt()
					}
					if t.GetCommands() != nil {
						if mt.Commands == nil {
							mt.Commands = make(map[int32]*pb.DeltaCommandManifest)
						}
						for commandID, c := range t.GetCommands() {
							mc, found := mt.Commands[commandID]
							if !found {
								mc = new(pb.DeltaCommandManifest)
								mt.Commands[commandID] = mc
							}
							if c.GetName() != "" {
								mc.Name = c.GetName()
							}
							if c.GetStatus() != pb.BuildStatus_BUILD_STATUS_UNKNOWN {
								mc.Status = c.GetStatus()
							}
							if c.GetHasCached() {
								mc.HasCached = true
								mc.IsCached = c.GetIsCached()
							}
							if c.GetHasPush() {
								mc.HasPush = true
								mc.IsPush = c.GetIsPush()
							}
							if c.GetHasLocal() {
								mc.HasLocal = true
								mc.IsLocal = c.GetIsLocal()
							}
							if c.GetStartedAt() != 0 {
								mc.StartedAt = c.GetStartedAt()
							}
							if c.GetFinishedAt() != 0 {
								mc.FinishedAt = c.GetFinishedAt()
							}
							if c.GetHasHasProgress() {
								mc.HasHasProgress = true
								mc.HasProgress = c.GetHasProgress()
							}
							if c.GetProgress() != 0 {
								mc.Progress = c.GetProgress()
							}
						}
					}
				}
			}
		}
	}
	return ret, nextManifestOrderID
}

func mergeDeltaLogs(dls []*pb.DeltaLog) []*pb.DeltaLog {
	targets := make(map[string]*pb.DeltaLog)
	eof := make(map[string]bool)
	for _, dl := range dls {
		tdl, found := targets[dl.GetTargetId()]
		if !found {
			tdl = &pb.DeltaLog{
				TargetId:      dl.GetTargetId(),
				SeekIndex:     dl.GetSeekIndex(),
				DeltaLogOneof: &pb.DeltaLog_Data{},
			}
			targets[dl.GetTargetId()] = tdl
		}
		dataOneof := tdl.DeltaLogOneof.(*pb.DeltaLog_Data)
		switch dl.GetDeltaLogOneof().(type) {
		case *pb.DeltaLog_Eof:
			eof[dl.GetTargetId()] = true
		case *pb.DeltaLog_Data:
			dataOneof.Data = append(dataOneof.Data, dl.GetData()...)
		}
	}
	ret := make([]*pb.DeltaLog, 0, len(targets))
	for targetID, tdl := range targets {
		if len(tdl.GetData()) > 0 {
			ret = append(ret, tdl)
		}
		if eof[targetID] {
			ret = append(ret, &pb.DeltaLog{
				TargetId:      targetID,
				SeekIndex:     tdl.GetSeekIndex() + int64(len(tdl.GetData())),
				DeltaLogOneof: &pb.DeltaLog_Eof{Eof: true},
			})
		}
	}
	return ret
}
