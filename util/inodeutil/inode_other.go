//go:build !windows
// +build !windows

package inodeutil

import (
	"syscall"
)

// GetInodeBestEffort returns an inode if available, or 0 on failure
func GetInodeBestEffort(path string) uint64 {
	var stat syscall.Stat_t
	inode := uint64(0)
	if err := syscall.Stat(path, &stat); err == nil {
		inode = uint64(stat.Ino)
	}
	return inode
}
