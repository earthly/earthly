package main

import (
	"fmt"
	"net/url"

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

	// TODO: Get token
	var token string
	if token != "" {
		urlToOpen.Query().Set("token", token)
	}

	return errors.New("unimplemented")
}
