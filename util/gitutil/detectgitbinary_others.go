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
		return ErrNoGitBinary
	}
	return nil
}
