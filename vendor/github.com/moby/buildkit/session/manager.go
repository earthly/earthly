package session

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Caller can invoke requests on the session
type Caller interface {
	Context() context.Context
	Supports(method string) bool
	Conn() *grpc.ClientConn
	Name() string
	SharedKey() string
}

type client struct {
	Session
	cc        *grpc.ClientConn
	supported map[string]struct{}
}

type History struct {
	Start time.Time
	End   time.Time
}

// Manager is a controller for accessing currently active sessions
type Manager struct {
	sessions        map[string]*client
	mu              sync.Mutex
	updateCondition *sync.Cond
	healthCfg       ManagerHealthCfg

	stop            bool // Earthly-specific.
	shutdownCh      chan struct{}
	idleAt          time.Time           // Earthly-specific
	history         map[string]*History // Earthly-specific
	historyDuration time.Duration       // Earthly-specific
}

// ManagerHealthCfg is the healthcheck configuration for gRPC healthchecks
type ManagerHealthCfg struct {
	frequency       time.Duration
	timeout         time.Duration
	allowedFailures int
}

// ManagerOpt is earthly-specific, and required for custom health-check overrides
type ManagerOpt struct {
	HealthFrequency       time.Duration
	HealthTimeout         time.Duration
	HealthAllowedFailures int
	ShutdownCh            chan struct{}
}

// NewManager returns a new Manager
// earthly-specific: opt param is required for our custom health config
func NewManager(opt *ManagerOpt) (*Manager, error) {
	var historyDuration time.Duration
	if dur, ok := os.LookupEnv("BUILDKIT_SESSION_HISTORY_DURATION"); ok {
		var err error
		historyDuration, err = time.ParseDuration(dur)
		if err != nil {
			return nil, errors.Wrapf(err, "could not parse session history duration value of '%s'", dur)
		}
	}
	sm := &Manager{
		sessions: make(map[string]*client),
		healthCfg: ManagerHealthCfg{
			frequency:       opt.HealthFrequency,
			timeout:         opt.HealthTimeout,
			allowedFailures: opt.HealthAllowedFailures,
		},
		shutdownCh:      opt.ShutdownCh,
		idleAt:          time.Now(),
		history:         make(map[string]*History),
		historyDuration: historyDuration,
	}
	sm.updateCondition = sync.NewCond(&sm.mu)
	return sm, nil
}

// NumSessions returns the number of active sessions.
// earthly-specific
func (sm *Manager) NumSessions() (sessions int, durationIdle time.Duration) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sessions = len(sm.sessions)
	if sessions == 0 {
		durationIdle = time.Now().Sub(sm.idleAt)
	}
	return sessions, durationIdle
}

// StopIfIdle stops the manager if there are no active sessions.
// earthly-specific
func (sm *Manager) StopIfIdle() (bool, int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	isIdle := sm.idleAt.Add(time.Minute).Before(time.Now())
	if isIdle && len(sm.sessions) == 0 {
		close(sm.shutdownCh)
		sm.stop = true
		return true, 0
	}
	return false, len(sm.sessions)
}

// Reserve signals intent to start a build.
// It resets the idleAt counter so that the buildkit will not shutdown.
// earthly-specific
func (sm *Manager) Reserve() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.stop {
		return errors.New("already shutting down")
	}
	sm.idleAt = time.Now()
	return nil
}

// earthly-specific
func (sm *Manager) recordSessionStart(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.history[sessionID] = &History{Start: time.Now()}
}

// earthly-specific
func (sm *Manager) recordSessionEnd(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if h := sm.history[sessionID]; h != nil {
		h.End = time.Now()
	}
	for id, history := range sm.history {
		if time.Since(history.Start) > sm.historyDuration {
			delete(sm.history, id)
		}
	}
}

// GetSessionHistory returns a map of session ID to History entries
// earthly-specific
func (sm *Manager) GetSessionHistory() map[string]*History {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	historyCopy := make(map[string]*History)
	for id, h := range sm.history {
		history := *h
		historyCopy[id] = &history
	}
	return historyCopy
}

// HandleHTTPRequest handles an incoming HTTP request
func (sm *Manager) HandleHTTPRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return errors.New("handler does not support hijack")
	}

	id := r.Header.Get(headerSessionID)

	proto := r.Header.Get("Upgrade")

	sm.mu.Lock()
	if _, ok := sm.sessions[id]; ok {
		sm.mu.Unlock()
		return errors.Errorf("session %s already exists", id)
	}

	if proto == "" {
		sm.mu.Unlock()
		return errors.New("no upgrade proto in request")
	}

	if proto != "h2c" {
		sm.mu.Unlock()
		return errors.Errorf("protocol %s not supported", proto)
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		sm.mu.Unlock()
		return errors.Wrap(err, "failed to hijack connection")
	}

	resp := &http.Response{
		StatusCode: http.StatusSwitchingProtocols,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
	}
	resp.Header.Set("Connection", "Upgrade")
	resp.Header.Set("Upgrade", proto)

	// set raw mode
	conn.Write([]byte{})
	resp.Write(conn)

	return sm.handleConn(ctx, conn, r.Header)
}

// HandleConn handles an incoming raw connection
func (sm *Manager) HandleConn(ctx context.Context, conn net.Conn, opts map[string][]string) error {
	sm.mu.Lock()
	return sm.handleConn(ctx, conn, opts)
}

// caller needs to take lock, this function will release it
func (sm *Manager) handleConn(ctx context.Context, conn net.Conn, opts map[string][]string) error {
	if sm.stop {
		sm.mu.Unlock()
		return errors.New("shutting down")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts = canonicalHeaders(opts)

	h := http.Header(opts)
	id := h.Get(headerSessionID)
	name := h.Get(headerSessionName)
	sharedKey := h.Get(headerSessionSharedKey)

	ctx, cc, err := grpcClientConn(ctx, conn, sm.healthCfg)
	if err != nil {
		sm.mu.Unlock()
		return err
	}

	c := &client{
		Session: Session{
			id:        id,
			name:      name,
			sharedKey: sharedKey,
			ctx:       ctx,
			cancelCtx: cancel,
			done:      make(chan struct{}),
		},
		cc:        cc,
		supported: make(map[string]struct{}),
	}

	for _, m := range opts[headerSessionMethod] {
		c.supported[strings.ToLower(m)] = struct{}{}
	}
	sm.sessions[id] = c
	sm.updateCondition.Broadcast()
	sm.mu.Unlock()
	sm.recordSessionStart(id) // earthly-specific

	defer func() {
		sm.mu.Lock()
		delete(sm.sessions, id)
		if len(sm.sessions) == 0 {
			sm.idleAt = time.Now() // earthly-specific
		}
		sm.mu.Unlock()
		sm.recordSessionEnd(id) // earthly-specific
	}()

	<-c.ctx.Done()
	conn.Close()
	close(c.done)

	return nil
}

// Get returns a session by ID
func (sm *Manager) Get(ctx context.Context, id string, noWait bool) (Caller, error) {
	// session prefix is used to identify vertexes with different contexts so
	// they would not collide, but for lookup we don't need the prefix
	if p := strings.SplitN(id, ":", 2); len(p) == 2 && len(p[1]) > 0 {
		id = p[1]
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		sm.mu.Lock()
		sm.updateCondition.Broadcast()
		sm.mu.Unlock()
	}()

	var c *client

	sm.mu.Lock()
	for {
		select {
		case <-ctx.Done():
			sm.mu.Unlock()
			return nil, errors.Wrapf(ctx.Err(), "no active session for %s", id)
		default:
		}
		var ok bool
		c, ok = sm.sessions[id]
		if (!ok || c.closed()) && !noWait {
			sm.updateCondition.Wait()
			continue
		}
		sm.mu.Unlock()
		break
	}

	if c == nil {
		return nil, nil
	}

	return c, nil
}

func (c *client) Context() context.Context {
	return c.context()
}

func (c *client) Name() string {
	return c.name
}

func (c *client) SharedKey() string {
	return c.sharedKey
}

func (c *client) Supports(url string) bool {
	_, ok := c.supported[strings.ToLower(url)]
	return ok
}
func (c *client) Conn() *grpc.ClientConn {
	return c.cc
}

func canonicalHeaders(in map[string][]string) map[string][]string {
	out := map[string][]string{}
	for k := range in {
		out[http.CanonicalHeaderKey(k)] = in[k]
	}
	return out
}
