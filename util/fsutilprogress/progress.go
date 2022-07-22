package fsutilprogress

import (
	"path"
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/earthly/earthly/conslogging"
	"github.com/tonistiigi/fsutil"
)

// ProgressCallback exposes two different levels of callbacks for displaying status on files being sent or received
type ProgressCallback interface {
	Info(numBytes int, last bool)
	Verbose(relPath string, status fsutil.VerboseProgressStatus, numBytes int)
}

type progressCallback struct {
	console     conslogging.ConsoleLogger
	mutex       sync.Mutex
	pathPrefix  string
	numStats    int
	numSent     int
	numReceived int
	filesize    map[string]int
}

// New returns a new verbose progress callback for use with fsutil
func New(pathPrefix string, console conslogging.ConsoleLogger) ProgressCallback {
	return &progressCallback{
		console:    console,
		pathPrefix: pathPrefix,
		filesize:   map[string]int{},
	}
}

func (s *progressCallback) Info(numBytes int, last bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if last {
		s.console.Printf("transferred %d file(s) for context %s (%s, %d file/dir stats)", s.numSent, s.pathPrefix, humanize.Bytes(uint64(numBytes)), s.numStats)
	}
}

func (s *progressCallback) Verbose(relPath string, status fsutil.VerboseProgressStatus, numBytes int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	fullPath := path.Join(s.pathPrefix, relPath)
	switch status {
	case fsutil.StatusStat:
		s.numStats++
		//s.console.VerbosePrintf("sent file stat for %s\n", fullPath) ignored as it is too verbose. TODO add different verbose levels to support ExtraVerbosePrintf
	case fsutil.StatusSent:
		s.console.VerbosePrintf("sent data for %s (%s)\n", fullPath, humanize.Bytes(uint64(numBytes)))
		s.numSent++
	case fsutil.StatusReceiving:
		s.filesize[fullPath] += numBytes
		//ignore
	case fsutil.StatusReceived:
		if numBytes == 0 {
			numBytes = s.filesize[fullPath]
		}
		s.console.VerbosePrintf("received data for %s (%s)\n", fullPath, humanize.Bytes(uint64(numBytes)))
		s.numReceived++
	case fsutil.StatusFailed:
		s.console.VerbosePrintf("sent data for %s failed\n", fullPath)
	case fsutil.StatusSkipped:
		s.console.VerbosePrintf("ignoring %s\n", fullPath)
	default:
		s.console.Warnf("unhandled progress status %v (path=%s, numBytes=%d)\n", status, fullPath, numBytes)
	}
}
