//go:build windows
// +build windows

package fileutil

import (
	"os"
)

// IsDirWritable returns if the path is a directory that the user can write to
func IsDirWritable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	err = nil
	if !info.IsDir() {
		return false, nil
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, nil
	}

	return true, nil
}
