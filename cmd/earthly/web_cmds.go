package main

import (
	"fmt"
	"net/url"

	"github.com/earthly/earthly/cloud"

	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (app *earthlyApp) webUI(cliCtx *cli.Context) error {
	urlToOpen, err := url.Parse(fmt.Sprintf("%s/login", app.getCIHost()))
	if err != nil {
		return errors.Wrapf(err, "failed to parse url")
	}
	query := urlToOpen.Query()
	if app.loginProvider != "" {
		query.Set("provider", app.loginProvider)
	}
	if app.loginFinal != "" {
		query.Set("final", app.loginFinal)
	}

	client, err := app.newCloudClient()
	if err != nil {
		return errors.Wrap(err, "failed to initialize cloud client")
	}

	token, err := client.GetAuthToken(cliCtx.Context)
	if err != nil && !errors.Is(err, cloud.ErrUnauthorized) {
		return errors.Wrapf(err, "failed to get auth token")
	}

	if token != "" {
		query.Set("token", token)
	}

	urlToOpen.RawQuery = query.Encode()
	urlString := urlToOpen.String()

	err = browser.OpenURL(urlString)
	if err != nil {
		app.console.Printf("failed to open UI in browser")
	}

	app.console.Printf("Visit UI at %s", urlString)
	return nil
}
