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

	client, err := app.newCloudClient()
	if err != nil {
		return errors.Wrap(err, "failed to initialize cloud client")
	}

	token, err := client.GetAuthToken(cliCtx.Context)
	if err != nil {
		if !errors.Is(err, cloud.ErrUnauthorized) {
			return errors.Wrapf(err, "failed to get auth token")
		}
		app.console.VerbosePrintf("failed to get token %s", err.Error())
	}

	if token != "" {
		query.Set("token", token)
	}

	urlToOpen.RawQuery = query.Encode()
	urlString := urlToOpen.String()

	err = browser.OpenURL(urlString)
	if err != nil {
		err := errors.Wrapf(err, "failed to open UI in browser")
		app.console.Printf(err.Error())
	}

	app.console.Printf("Visit UI at %s", urlString)
	return nil
}
