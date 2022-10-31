package delta2cons

import (
	"context"
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/conslogging"
)

type logstreamConsole struct {
	conslogging.ConsoleLogger
	cmd  *command
	lock sync.Mutex
}

type logstreamWriter struct {
	ctx       context.Context
	buildID   string
	commandID string
	targetID  string
	stream    int32
	timestamp time.Time
	client    logstream.LogStreamClient
}

func (lw *logstreamWriter) Write(b []byte) (int, error) {
	stream, err := lw.client.StreamLogs(lw.ctx)
	if err != nil {
		return 0, err
	}

	err = stream.Send(&logstream.StreamLogRequest{
		BuildId: lw.buildID,
		Deltas: []*logstream.Delta{
			{
				DeltaTypeOneof: &logstream.Delta_DeltaLog{DeltaLog: &logstream.DeltaLog{
					TargetId:           lw.targetID,
					CommandId:          lw.commandID,
					Stream:             lw.stream,
					TimestampUnixNanos: uint64(lw.timestamp.UnixNano()),
					Data:               b,
					// Some field to denote this is a 'readable log'?
				}},
			},
		},
	})

	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (l logstreamConsole) WithCommand(cmd *command) Console {
	return logstreamConsole{
		ConsoleLogger: l.ConsoleLogger,
		cmd:           cmd,
	}
}

func (l logstreamConsole) WithPrefixAndSalt(prefix, salt string) Console {
	l2 := l.clone()
	l2.ConsoleLogger = l.ConsoleLogger.WithPrefixAndSalt(prefix, salt)
	return l2
}

func (l logstreamConsole) PrintBytes(b []byte) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.ConsoleLogger = conslogging.NewBufferedLogger()
	//TODO implement me
	panic("implement me")
}

func (l logstreamConsole) WithPrefix(prefix string) Console {
	//TODO implement me
	panic("implement me")
}

func (l logstreamConsole) WithFailed(b bool) Console {
	//TODO implement me
	panic("implement me")
}

func (l logstreamConsole) WithCached(b bool) Console {
	//TODO implement me
	panic("implement me")
}

func (l logstreamConsole) WithMetadataMode(b bool) Console {
	//TODO implement me
	panic("implement me")
}

func (l logstreamConsole) Printf(str string, format ...any) error {
	//TODO implement me
	/**
	if cl.logLevel < Info {
		return
	}
	cl.mu.Lock()
	defer cl.mu.Unlock()
	c := cl.color(noColor)
	if cl.metadataMode {
		c = cl.color(metadataModeColor)
	}
	text := fmt.Sprintf(format, args...)
	text = strings.TrimSuffix(text, "\n")
	for _, line := range strings.Split(text, "\n") {
		cl.printPrefix()
		c.Fprintf(cl.errW, "%s", line)

		// Don't use a background color for \n.
		noColor.Fprintf(cl.errW, "\n")
	}
	*/
	panic("implement me")
}

func (l logstreamConsole) VerbosePrintf(str string, format ...any) {
	//TODO implement me
	panic("implement me")
}

func newLogstreamConsole() Console {
	return &logstreamConsole{}
}

func (l logstreamConsole) clone() logstreamConsole {
	return logstreamConsole{ConsoleLogger: l.ConsoleLogger}
}
