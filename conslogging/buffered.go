package conslogging

import (
	"fmt"
)

// BufferedLogger is a logger that queues up messages until Flush is called.
type BufferedLogger struct {
	cl    *ConsoleLogger
	queue []string
}

// NewBufferedLogger creates a new BufferedLogger.
func NewBufferedLogger(cl *ConsoleLogger) *BufferedLogger {
	return &BufferedLogger{
		cl: cl,
	}
}

// Printf prints a formatted string to the delayed console.
func (bl *BufferedLogger) Printf(format string, v ...interface{}) {
	bl.queue = append(bl.queue, fmt.Sprintf(format, v...))
}

// Flush prints the queued up messages to the underlying console.
func (bl *BufferedLogger) Flush() {
	for _, s := range bl.queue {
		bl.cl.Printf("%s", s)
	}
	bl.queue = nil
}
