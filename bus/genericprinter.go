package bus

import (
	"fmt"
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GenericPrinter is a generic printer for build output unrelated to a specific target.
type GenericPrinter struct {
	bp       *BuildPrinter
	category string
	mu       sync.Mutex
}

func newGenericPrinter(bp *BuildPrinter, category string) *GenericPrinter {
	return &GenericPrinter{
		bp:       bp,
		category: category,
	}
}

// Write writes the given bytes to the generic printer.
func (gp *GenericPrinter) Write(dt []byte) (int, error) {
	return gp.WriteWithTimestamp(dt, time.Now())
}

// WriteWithTimestamp writes the given bytes to the generic printer.
func (gp *GenericPrinter) WriteWithTimestamp(dt []byte, ts time.Time) (int, error) {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	gp.bp.b.RawDelta(&logstream.Delta{
		DeltaTypeOneof: &logstream.Delta_DeltaLog{
			DeltaLog: &logstream.DeltaLog{
				TargetId:  fmt.Sprintf("_generic:%s", gp.category),
				Timestamp: timestamppb.New(ts),
				Log:       dt,
			},
		},
	})
	return len(dt), nil
}
