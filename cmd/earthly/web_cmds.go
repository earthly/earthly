package main

import (
	"fmt"
	"net/url"

	"github.com/earthly/earthly/cloud"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (app *earthlyApp) webUI(cliCtx *cli.Context) error {
	urlToOpen, err := url.Parse(fmt.Sprintf("%s/login", app.cloudHTTPAddr))
	if err != nil {
		return errors.Wrapf(err, "failed to parse url")
	}
	if app.loginProvider != "" {
		urlToOpen.Query().Set("provider", app.loginProvider)
	}
	if app.loginFinal != "" {
		urlToOpen.Query().Set("final", app.loginFinal)
	}

	client, err := app.newCloudClient()
	if err != nil {
		return errors.Wrap(err, "failed to initialize cloud client")
	}

	token, err := client.GetAuthToken(cliCtx.Context)
	if err != nil && err != cloud.ErrUnauthorized {
		return errors.Wrapf(err, "failed to get auth token")
	}

	if token != "" {
		urlToOpen.Query().Set("token", token)
	}

	app.console.Printf("Visit UI at %s", urlToOpen.String())
	return nil
}
