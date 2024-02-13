package helper

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"

	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"

	"github.com/earthly/earthly/cmd/earthly/base"

	"github.com/earthly/earthly/autocomplete"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/cliutil"
)

// to enable autocomplete, enter
// complete -o nospace -C "/path/to/earthly" earthly
func AutoComplete(ctx context.Context, cli *base.CLI) {
	_, found := os.LookupEnv("COMP_LINE")
	if !found {
		return
	}

	cli.SetConsole(cli.Console().WithLogLevel(conslogging.Silent))

	err := autoCompleteImp(ctx, cli)
	if err != nil {
		errToLog := err
		logDir, err := cliutil.GetOrCreateEarthlyDir(cli.Flags().InstallationName)
		if err != nil {
			os.Exit(1)
		}
		logFile := filepath.Join(logDir, "autocomplete.log")
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			os.Exit(1)
		}
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			os.Exit(1)
		}
		fmt.Fprintf(f, "error during autocomplete: %s\n", errToLog)
		os.Exit(1)
	}
	os.Exit(0)
}

func autoCompleteImp(ctx context.Context, cli *base.CLI) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("recovered panic in autocomplete %s: %s", r, debug.Stack())
		}
	}()

	compLine := os.Getenv("COMP_LINE")   // full command line
	compPoint := os.Getenv("COMP_POINT") // where the cursor is

	compPointInt, err := strconv.ParseUint(compPoint, 10, 64)
	if err != nil {
		return err
	}
	if !(compPointInt > 0 && compPointInt < math.MaxInt) {
		err = errors.Errorf("compPointInt is out of bounds.")
		return err
	}

	gitLookup := buildcontext.NewGitLookup(cli.Console(), cli.Flags().SSHAuthSock)
	resolver := buildcontext.NewResolver(nil, gitLookup, cli.Console(), "", cli.Flags().GitBranchOverride, "", 0, "")
	var gwClient gwclient.Client // TODO this is a nil pointer which causes a panic if we try to expand a remotely referenced earthfile
	// it's expensive to create this gwclient, so we need to implement a lazy eval which returns it when required.

	potentials, err := autocomplete.GetPotentials(ctx, resolver, gwClient, compLine, int(compPointInt), cli.App(), autocomplete.NewCachedCloudClient(cli.Flags().InstallationName, getCloudClientForAutoCompleter(ctx, cli)))
	if err != nil {
		return err
	}
	for _, p := range potentials {
		fmt.Printf("%s\n", p)
	}

	return nil
}

func getCloudClientForAutoCompleter(ctx context.Context, cli *base.CLI) *cloud.Client {
	// TODO these need to be set outside of urfave/cli; since the auto-competer happens without envoking it
	// maybe we can half-parse them? or just set the defaults?
	cli.Flags().CloudHTTPAddr = "https://api.earthly.dev"
	cli.Flags().CloudGRPCAddr = "ci.earthly.dev:443"
	cli.Flags().SSHAuthSock = os.Getenv("SSH_AUTH_SOCK")

	cloudClient, _ := NewCloudClient(cli) // best effort
	return cloudClient
}
