package cliutil

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/earthly/earthly/util/fileutil"
	"github.com/pkg/errors"
)

var earthlyDir string
var earthlyDirErr error
var earthlyDirOnce sync.Once

// GetEarthlyDir returns the .earthly dir. (Usually ~/.earthly).
func GetEarthlyDir() (string, error) {
	earthlyDirOnce.Do(func() {
		earthlyDir, earthlyDirErr = makeEarthlyDir()
	})
	return earthlyDir, earthlyDirErr
}

func makeEarthlyDir() (string, error) {
	homeDir, sudoUser, err := detectHomeDir()
	if err != nil {
		return "", err
	}
	earthlyDir := filepath.Join(homeDir, ".earthly")
	if !fileutil.DirExists(earthlyDir) {
		err := os.MkdirAll(earthlyDir, 0755)
		if err != nil {
			return "", errors.Wrapf(err, "unable to create dir %s", earthlyDir)
		}
		if sudoUser != nil {
			// Attempt to chown the created dir to belong to the sudo user.
			uid, err := strconv.Atoi(sudoUser.Uid)
			if err != nil {
				// Swallow error.
				return earthlyDir, nil
			}
			gid := 0
			if sudoUser.Gid != "" {
				// If cannot convert will use gid 0.
				gid, _ = strconv.Atoi(sudoUser.Gid)
			}
			err = os.Chown(earthlyDir, uid, gid)
			if err != nil {
				// Swallow error.
				return earthlyDir, nil
			}
		}
	}
	return earthlyDir, nil
}

func detectHomeDir() (homeDir string, sudoUser *user.User, err error) {
	if runtime.GOOS == "windows" {
		homeDir, err := os.UserHomeDir()
		return homeDir, nil, err
	}
	// See if SUDO_USER exists. Use that user's home dir.
	sudoUserName, ok := os.LookupEnv("SUDO_USER")
	if ok {
		sudoUser, err := user.Lookup(sudoUserName)
		if err == nil && sudoUser.HomeDir != "" {
			return sudoUser.HomeDir, sudoUser, nil
		}
	}
	// Try to use current user's home dir.
	homeDir, err = os.UserHomeDir()
	if err != nil {
		// Try $HOME.
		homeDir, ok := os.LookupEnv("HOME")
		if ok {
			return homeDir, nil, nil
		}
		// No home dir available - use /etc instead.
		return "/etc", nil, nil
	}
	return homeDir, nil, nil
}
