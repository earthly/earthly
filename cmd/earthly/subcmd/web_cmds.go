package subcmd

import (
	"fmt"
	"net/url"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/cmd/earthly/helper"

	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type Web struct {
	cli CLI

	loginProvider string
}

func NewWeb(cli CLI) *Web {
	return &Web{
		cli: cli,
	}
}

func (a *Web) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "web",
			Usage:     "*beta* Access the web UI via your default browser and print the url",
			UsageText: "earthly web (--provider=github)",
			Description: `*beta* Prints a url for entering the CI application and attempts to open your default browser with that url.

		If the provider argument is given the CI application will automatically begin an OAuth flow with the given provider.
		If you are logged into the CLI the url will contain a token used to link your OAuth credentials to your Earthly user.`,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "provider",
					EnvVars:     []string{"EARTHLY_LOGIN_PROVIDER"},
					Usage:       "The provider to use when logging into the web ui",
					Required:    false,
					Destination: &a.loginProvider,
				},
			},
			Action: a.action,
		},
	}
}

func (a *Web) action(cliCtx *cli.Context) error {
	urlToOpen, err := url.Parse(fmt.Sprintf("%s/login", a.cli.CIHost()))
	if err != nil {
		return errors.Wrapf(err, "failed to parse url")
	}
	query := urlToOpen.Query()
	if a.loginProvider != "" {
		query.Set("provider", a.loginProvider)
	}

	client, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return errors.Wrap(err, "failed to initialize cloud client")
	}

	token, err := client.GetAuthToken(cliCtx.Context)
	if err != nil {
		if !errors.Is(err, cloud.ErrUnauthorized) {
			return errors.Wrapf(err, "failed to get auth token")
		}
		a.cli.Console().VerbosePrintf("failed to get token %s", err.Error())
	}

	if token != "" {
		query.Set("token", token)
	}

	urlToOpen.RawQuery = query.Encode()
	urlString := urlToOpen.String()

	err = browser.OpenURL(urlString)
	if err != nil {
		err := errors.Wrapf(err, "failed to open UI in browser")
		a.cli.Console().Printf(err.Error())
	}

	a.cli.Console().Printf("Visit UI at %s", urlString)
	return nil
}
