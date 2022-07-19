package secretprovider

import (
	"context"
	"errors"
	"net/url"
	"os"
	"os/exec"
	"syscall"

	"github.com/earthly/earthly/debugger/common"

	"github.com/alessio/shellescape"
	"github.com/moby/buildkit/session/secrets"
)

type cmdStore struct {
	cmd string
}

// NewSecretProviderCmd returns a SecretStore that shells out to a user-supplied command
func NewSecretProviderCmd(cmd string) (secrets.SecretStore, error) {
	return &cmdStore{
		cmd: cmd,
	}, nil
}

// GetSecret gets a secret from the map store
func (c *cmdStore) GetSecret(ctx context.Context, id string) ([]byte, error) {
	if len(c.cmd) == 0 {
		return nil, secrets.ErrNotFound
	}
	q, err := url.ParseQuery(id)
	if err != nil {
		return nil, errors.New("invalid secret ID format")
	}
	name := q.Get("name")
	if name == common.DebuggerSettingsSecretsKey {
		// the interactive debugger passes config values by abusing secrets,
		// we must not call the user's secret provider in this case.
		return nil, secrets.ErrNotFound
	}
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", c.cmd+" "+shellescape.Quote(name))
	envs := os.Environ()
	if q.Get("v") == "1" {
		envs = append(envs, "EARTHLY_ORG="+shellescape.Quote(q.Get("org")))
		envs = append(envs, "EARTHLY_PROJECT="+shellescape.Quote(q.Get("project")))
	}
	cmd.Env = envs
	cmd.Stderr = os.Stderr
	dt, err := cmd.Output()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitStatus := status.ExitStatus()
				if exitStatus == 2 {
					// exit code of 2 indicates secret not found (and earthly should continue looking in other stores)
					return nil, secrets.ErrNotFound
				}
			}
		}
		return nil, err
	}
	return dt, nil
}
