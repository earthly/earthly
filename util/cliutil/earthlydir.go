package cliutil

import (
	"os"
	"os/user"
	"path/filepath"
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
		fileutil.EnsureUserOwned(earthlyDir, sudoUser)
	}
	return earthlyDir, nil
}

// DetectHomeDir returns the home directory of the current user, together with
// the user object who owns it.
func DetectHomeDir() (string, *user.User, error) {
	u, isSudo, err := currentNonSudoUser()
	if err != nil {
		return "", nil, errors.Wrap(err, "lookup user for homedir")
	}
	if !isSudo {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			return homeDir, u, nil
		}
	}

	if u.HomeDir == "" {
		return "/etc", u, nil
	}
	return u.HomeDir, u, nil
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

	return true
}

// EnsurePermissions changes the permissions of all earthly files to be owned by the user and their group.
func EnsurePermissions() error {
	_, err := GetEarthlyDir()
	if err != nil {
		return err
	}

	u, _, err := currentNonSudoUser()
	if err != nil {
		return errors.Wrap(err, "get non-sudo user")
	}

	fileutil.EnsureUserOwned(earthlyDir, u)
	return nil
}

func currentNonSudoUser() (*user.User, bool, error) {
	if sudoUserName, ok := os.LookupEnv("SUDO_USER"); ok {
		sudoUser, err := user.Lookup(sudoUserName)
		if err == nil {
			return sudoUser, true, nil
		}
	}

	u, err := user.Current()
	if err != nil {
		return nil, false, err
	}
	return u, false, nil
}
