//go:build !windows
// +build !windows

package gitutil

import (
	"context"
	"os/exec"
)

func detectGitBinary(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", "which git")
	_, err := cmd.Output()
	if err != nil {
		_, isExitError := err.(*exec.ExitError)
		if isExitError {
			return ErrNoGitBinary
		}
		return err
	}
	return nil
}
