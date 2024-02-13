//go:build windows
// +build windows

package inodeutil

import (
	"os"
	"syscall"
)

// GetInodeBestEffort returns an inode if available, or 0 on failure
func GetInodeBestEffort(path string) uint64 {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()

	var info syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(syscall.Handle(f.Fd()), &info); err != nil {
		return 0
	}
	inode := uint64(info.FileIndexHigh)
	inode <<= 32
	inode += uint64(info.FileIndexLow)
	return inode
}
