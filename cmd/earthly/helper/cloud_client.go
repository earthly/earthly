package helper

import (
	"github.com/pkg/errors"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/cmd/earthly/flag"
	"github.com/earthly/earthly/conslogging"
)

type CLI interface {
	Flags() *flag.Global
	Console() conslogging.ConsoleLogger
}

func NewCloudClient(cli CLI, opts ...cloud.ClientOpt) (*cloud.Client, error) {
	cloudClient, err := cloud.NewClient(cli.Flags().CloudHTTPAddr, cli.Flags().CloudGRPCAddr,
		cli.Flags().CloudGRPCInsecure, cli.Flags().SSHAuthSock, "",
		"", cli.Flags().InstallationName, "",
		cli.Console().Warnf, cli.Console().DebugPrintf, cli.Flags().ServerConnTimeout, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloud client")
	}
	return cloudClient, nil
}
