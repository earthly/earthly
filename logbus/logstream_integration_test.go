package logbus

import (
	"context"
	"testing"

	"github.com/moby/buildkit/client"
	"github.com/stretchr/testify/require"
)

const (
	buildID = "test-build-id"
	ciHost  = "https://ci-beta.staging.earthly.dev"
)

type nullPrinter struct{}

func (n *nullPrinter) Printf(format string, args ...interface{}) {}
func (n *nullPrinter) Warnf(format string, args ...interface{})  {}

func defaultArgs(t *testing.T) *LogstreamArgs {
	return &LogstreamArgs{
		BuildID:                    buildID,
		CIHost:                     ciHost,
		Debug:                      false,
		Verbose:                    false,
		ForceColor:                 false,
		NoColor:                    false,
		DisableOngoingUpdates:      false,
		UseLogstream:               true,
		UploadLogstream:            true,
		LogstreamDebugFile:         "",
		LogstreamDebugManifestFile: "",
		ConsolePrinter:             &nullPrinter{},
	}
}

func Test_Sanity(t *testing.T) {
	ctx := context.Background()
	l, err := LogstreamFactory(ctx, defaultArgs(t))
	require.NoError(t, err)
	require.NotNil(t, l)

	ch := make(chan *client.SolveStatus)

	go func() {
		err := l.MonitorProgress(ctx, ch)
		require.NoError(t, err)
	}()

	require.NoError(t, err)

	ch <- &client.SolveStatus{
		Vertexes: []*client.Vertex{},
		Statuses: []*client.VertexStatus{},
		Logs:     []*client.VertexLog{},
		Warnings: []*client.VertexWarning{},
	}
}

// TODO: Converter touches targets by setting end
