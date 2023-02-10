// This file was generated by git.sr.ht/~nelsam/hel/v4.  Do not
// edit this code by hand unless you *really* know what you're
// doing.  Expect any changes made manually to be overwritten
// the next time hel regenerates this file.

package logstreamer_test

import (
	"context"
	"time"

	"git.sr.ht/~nelsam/hel/v4/vegr"
	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/cloud"
)

type mockCloudClient struct {
	t                vegr.T
	timeout          time.Duration
	StreamLogsCalled chan bool
	StreamLogsInput  struct {
		Ctx     chan context.Context
		BuildID chan string
		Deltas  chan cloud.Deltas
	}
	StreamLogsOutput struct {
		Ret0 chan error
	}
}

func newMockCloudClient(t vegr.T, timeout time.Duration) *mockCloudClient {
	m := &mockCloudClient{t: t, timeout: timeout}
	m.StreamLogsCalled = make(chan bool, 100)
	m.StreamLogsInput.Ctx = make(chan context.Context, 100)
	m.StreamLogsInput.BuildID = make(chan string, 100)
	m.StreamLogsInput.Deltas = make(chan cloud.Deltas, 100)
	m.StreamLogsOutput.Ret0 = make(chan error, 100)
	return m
}
func (m *mockCloudClient) StreamLogs(ctx context.Context, buildID string, deltas cloud.Deltas) (ret0 error) {
	m.t.Helper()
	m.StreamLogsCalled <- true
	m.StreamLogsInput.Ctx <- ctx
	m.StreamLogsInput.BuildID <- buildID
	m.StreamLogsInput.Deltas <- deltas
	vegr.PopulateReturns(m.t, "StreamLogs", m.timeout, m.StreamLogsOutput, &ret0)
	return ret0
}

type mockContext struct {
	t              vegr.T
	timeout        time.Duration
	DeadlineCalled chan bool
	DeadlineOutput struct {
		Deadline chan time.Time
		Ok       chan bool
	}
	DoneCalled chan bool
	DoneOutput struct {
		Ret0 chan (<-chan struct{})
	}
	ErrCalled chan bool
	ErrOutput struct {
		Ret0 chan error
	}
	ValueCalled chan bool
	ValueInput  struct {
		Key chan any
	}
	ValueOutput struct {
		Ret0 chan any
	}
}

func newMockContext(t vegr.T, timeout time.Duration) *mockContext {
	m := &mockContext{t: t, timeout: timeout}
	m.DeadlineCalled = make(chan bool, 100)
	m.DeadlineOutput.Deadline = make(chan time.Time, 100)
	m.DeadlineOutput.Ok = make(chan bool, 100)
	m.DoneCalled = make(chan bool, 100)
	m.DoneOutput.Ret0 = make(chan (<-chan struct{}), 100)
	m.ErrCalled = make(chan bool, 100)
	m.ErrOutput.Ret0 = make(chan error, 100)
	m.ValueCalled = make(chan bool, 100)
	m.ValueInput.Key = make(chan any, 100)
	m.ValueOutput.Ret0 = make(chan any, 100)
	return m
}
func (m *mockContext) Deadline() (deadline time.Time, ok bool) {
	m.t.Helper()
	m.DeadlineCalled <- true
	vegr.PopulateReturns(m.t, "Deadline", m.timeout, m.DeadlineOutput, &deadline, &ok)
	return deadline, ok
}
func (m *mockContext) Done() (ret0 <-chan struct{}) {
	m.t.Helper()
	m.DoneCalled <- true
	vegr.PopulateReturns(m.t, "Done", m.timeout, m.DoneOutput, &ret0)
	return ret0
}
func (m *mockContext) Err() (ret0 error) {
	m.t.Helper()
	m.ErrCalled <- true
	vegr.PopulateReturns(m.t, "Err", m.timeout, m.ErrOutput, &ret0)
	return ret0
}
func (m *mockContext) Value(key any) (ret0 any) {
	m.t.Helper()
	m.ValueCalled <- true
	m.ValueInput.Key <- key
	vegr.PopulateReturns(m.t, "Value", m.timeout, m.ValueOutput, &ret0)
	return ret0
}

type mockDeltas struct {
	t          vegr.T
	timeout    time.Duration
	NextCalled chan bool
	NextInput  struct {
		Ctx chan context.Context
	}
	NextOutput struct {
		Ret0 chan []*logstream.Delta
		Ret1 chan error
	}
}

func newMockDeltas(t vegr.T, timeout time.Duration) *mockDeltas {
	m := &mockDeltas{t: t, timeout: timeout}
	m.NextCalled = make(chan bool, 100)
	m.NextInput.Ctx = make(chan context.Context, 100)
	m.NextOutput.Ret0 = make(chan []*logstream.Delta, 100)
	m.NextOutput.Ret1 = make(chan error, 100)
	return m
}
func (m *mockDeltas) Next(ctx context.Context) (ret0 []*logstream.Delta, ret1 error) {
	m.t.Helper()
	m.NextCalled <- true
	m.NextInput.Ctx <- ctx
	vegr.PopulateReturns(m.t, "Next", m.timeout, m.NextOutput, &ret0, &ret1)
	return ret0, ret1
}
