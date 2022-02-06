package producer

import (
	"io"

	"github.com/earthly/earthly/logstream/api"
)

var _ io.Writer = &logWriter{}

type logWriter struct {
	deltaCh       chan api.Delta
	targetID      string
	nextSeekIndex int64
}

func (lw *logWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	lw.deltaCh <- api.Delta{
		Version: api.VersionNumber,
		DeltaLogs: []*api.DeltaLog{
			{
				TargetID:  lw.targetID,
				SeekIndex: lw.nextSeekIndex,
				Data:      p,
			},
		},
	}
	lw.nextSeekIndex += int64(len(p))
	return len(p), nil
}
