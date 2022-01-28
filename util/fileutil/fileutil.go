package fileutil

import (
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

// FileExists returns true if the file exists
func FileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, "unable to stat %s", filename)
	}
	return !info.IsDir(), nil
}

// FileExistsBestEffort returns true if the directory exists and ignores errors
func FileExistsBestEffort(filename string) bool {
	ok, _ := FileExists(filename)
	return ok
}

// DirExists returns true if the directory exists
func DirExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, "unable to stat %s", filename)
	}
	return info.IsDir(), nil
}

// DirExistsBestEffort returns true if the directory exists and ignores errors
func DirExistsBestEffort(filename string) bool {
	ok, _ := DirExists(filename)
	return ok
}

// EnsureUserOwned changes the files in the directory to be owned by the use and their group, as specified by the provided user.
func EnsureUserOwned(dir string, owner *user.User) error {
	exists, err := DirExists(dir)
	if err != nil || !exists {
		return err
	}

	uid, _ := strconv.Atoi(owner.Uid)

	gid := 0
	if owner.Gid != "" {
		// If cannot convert will use gid 0.
		gid, _ = strconv.Atoi(owner.Gid)
	}

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		return os.Chown(path, uid, gid)
	})
}
