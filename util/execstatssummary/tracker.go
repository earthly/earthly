package execstatssummary

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
)

// Tracker is used for tracking exec stats summary for each RUN command
type Tracker struct {
	mu    sync.Mutex
	stats map[string]*stats
	path  string
}

// NewTracker creates a new exec stats summary tracker
func NewTracker(path string) *Tracker {
	return &Tracker{
		stats: map[string]*stats{},
		path:  path,
	}
}

// Observe records an exec stats payload for a particular (target, command) pair
func (t *Tracker) Observe(target, command string, memory uint64, cpu time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	k := target + "|" + command
	stat, ok := t.stats[k]
	if !ok {
		stat = &stats{
			target:  target,
			command: command,
		}
		t.stats[k] = stat
	}
	if memory > stat.memory {
		stat.memory = memory
	}
	if cpu > stat.cpu {
		stat.cpu = cpu
	}
}

// String returns a summarized table
func (t *Tracker) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	keys := make([]string, 0, len(t.stats))
	for k := range t.stats {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return t.stats[keys[i]].memory < t.stats[keys[j]].memory
	})

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "target\tcommand\tmemory\tcpu\n")
	for _, k := range keys {
		v := t.stats[k]
		fmt.Fprintf(w, "%s\t%s\t%v\t%v\n", v.target, v.command, humanize.Bytes(v.memory), v.cpu)
	}
	w.Flush()
	return buf.String()
}

// Close closes the tracker, and writes the summary to disk (or stdout)
func (t *Tracker) Close(ctx context.Context) error {
	summary := t.String()
	if t.path == "-" {
		fmt.Print(summary)
		return nil
	}
	return os.WriteFile(t.path, []byte(summary), 0644)
}

type stats struct {
	memory  uint64
	cpu     time.Duration
	target  string
	command string
}
