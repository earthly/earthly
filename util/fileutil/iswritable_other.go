//go:build !windows
// +build !windows

package fileutil

import (
	"os"
	"syscall"
)

// IsDirWritable returns if the path is a directory that the user can write to
func IsDirWritable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if !info.IsDir() {
		return false, nil
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, nil
	}

	var stat syscall.Stat_t
	if err = syscall.Stat(path, &stat); err != nil {
		return false, err
	}

	return true, nil
}
