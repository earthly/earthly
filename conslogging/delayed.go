package conslogging

import (
	"fmt"
)

// DelayedLogger is a logger that queues up messages until Flush is called.
type DelayedLogger struct {
	cl    *ConsoleLogger
	queue []string
}

// NewDelayedLogger creates a new DelayedLogger.
func NewDelayedLogger(cl *ConsoleLogger) *DelayedLogger {
	return &DelayedLogger{
		cl: cl,
	}
}

// Printf prints a formatted string to the delayed console.
func (dl *DelayedLogger) Printf(format string, v ...interface{}) {
	dl.queue = append(dl.queue, fmt.Sprintf(format, v...))
}

// Flush prints the queued up messages to the underlying console.
func (dl *DelayedLogger) Flush() {
	for _, s := range dl.queue {
		dl.cl.Printf("%s", s)
	}
	dl.queue = nil
}
