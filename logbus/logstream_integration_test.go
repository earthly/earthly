package logbus

import (
	"context"
	"testing"

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
	l, err := LogstreamFactory(context.Background(), defaultArgs(t))
	require.NoError(t, err)
	require.NotNil(t, l)
}
