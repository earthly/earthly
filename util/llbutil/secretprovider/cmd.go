package secretprovider

import (
	"context"
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
	if id == common.DebuggerSettingsSecretsKey {
		// the interactive debugger passes config values by abusing secrets,
		// we must not call the user's secret provider in this case.
		return nil, secrets.ErrNotFound
	}
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", c.cmd+" "+shellescape.Quote(id))
	cmd.Env = os.Environ()
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
