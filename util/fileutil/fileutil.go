package fileutil

import (
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

// FileExists returns true if the file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists returns true if the directory exists
func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func EnsureUserOwned(dir string, owner *user.User) {
	if DirExists(dir) {
		uid, _ := strconv.Atoi(owner.Uid)

		gid := 0
		if owner.Gid != "" {
			// If cannot convert will use gid 0.
			gid, _ = strconv.Atoi(owner.Gid)
		}

		_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			return os.Chown(path, uid, gid)
		})
	}
}
