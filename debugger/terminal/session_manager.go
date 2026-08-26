package terminal

import (
	"context"
	"io"
	"sync"

	"github.com/earthly/earthly/conslogging"
)

type connectTermFunc func(context.Context, io.ReadWriteCloser, conslogging.ConsoleLogger) error

// SessionManager ensures that only one of the concurrent debugger connections
// controls the terminal during a build. Its zero value is ready to use.
type SessionManager struct {
	selectSession sync.Once
	sessionDone   chan struct{}
	connect       connectTermFunc
}

// ConnectTerm connects the first debugger session and keeps any concurrent
// sessions from competing for the shared terminal.
func (m *SessionManager) ConnectTerm(ctx context.Context, conn io.ReadWriteCloser, console conslogging.ConsoleLogger) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	isPrimary := false
	m.selectSession.Do(func() {
		isPrimary = true
		m.sessionDone = make(chan struct{})
		if m.connect == nil {
			m.connect = ConnectTerm
		}
	})

	if isPrimary {
		defer close(m.sessionDone)
		return m.connect(ctx, conn, console)
	}

	// Keep other failing commands blocked until the interactive session ends.
	// Returning earlier would let them fail the build and cancel that session.
	select {
	case <-m.sessionDone:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
