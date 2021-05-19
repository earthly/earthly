package cliutil

import (
	"os"
	"os/user"
	"path/filepath"
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
	homeDir, sudoUser, err := DetectHomeDir()
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

// DetectHomeDir returns the home directory of the current user, an additional sudoUser
// is returned if the user is currently running as root
func DetectHomeDir() (homeDir string, sudoUser *user.User, err error) {
	homeDir, sudoUser, err = fileutil.HomeDir()
	if err != nil {
		return
	}
	if homeDir == "" {
		homeDir = "/etc" // No home dir available - use /etc instead.
	}
	return
}

// IsBootstrapped provides a tentatively correct guess about the state of our bootstrapping.
func IsBootstrapped() bool {
	homeDir, _, err := DetectHomeDir()
	if err != nil {
		return false
	}

	earthlyDir := filepath.Join(homeDir, ".earthly")
	if !fileutil.DirExists(earthlyDir) {
		return false
	}

	installID := filepath.Join(homeDir, ".earthly", "install_id")
	if !fileutil.FileExists(installID) {
		return false
	}

	return true
}
