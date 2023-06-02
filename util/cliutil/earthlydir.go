package cliutil

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/earthly/earthly/util/fileutil"
	"github.com/pkg/errors"
)

var earthlyDir string
var earthlyDirOnce sync.Once
var earthlyDirSudoUser *user.User

var earthlyDirCreateOnce sync.Once
var earthlyDirCreateErr error

// GetEarthlyDir returns the .earthly dir. (Usually ~/.earthly).
// This function will not attempt to create the directory if missing, for that functionality use to the GetOrCreateEarthlyDir function.
func GetEarthlyDir(installName string) string {
	if installName == "" {
		// if GetEarthlyDir is called by the autocomplete code, this may not be set
		installName = "earthly"
	}
	earthlyDirOnce.Do(func() {
		earthlyDir, earthlyDirSudoUser = getEarthlyDirAndUser(installName)
	})
	return earthlyDir
}

func getEarthlyDirAndUser(installName string) (string, *user.User) {
	homeDir, u := fileutil.HomeDir()
	earthlyDir := filepath.Join(homeDir, fmt.Sprintf(".%s", installName))
	return earthlyDir, u
}

// GetOrCreateEarthlyDir returns the .earthly dir. (Usually ~/.earthly).
// if the directory does not exist, it will attempt to create it.
func GetOrCreateEarthlyDir(installName string) (string, error) {
	_ = GetEarthlyDir(installName) // ensure global vars get created so we can reference them below.

	earthlyDirCreateOnce.Do(func() {
		earthlyDirExists, err := fileutil.DirExists(earthlyDir)
		if err != nil {
			earthlyDirCreateErr = errors.Wrapf(err, "unable to create dir %s", earthlyDir)
			return
		}
		if !earthlyDirExists {
			err := os.MkdirAll(earthlyDir, 0755)
			if err != nil {
				earthlyDirCreateErr = errors.Wrapf(err, "unable to create dir %s", earthlyDir)
				return
			}
			if earthlyDirSudoUser != nil {
				err := fileutil.EnsureUserOwned(earthlyDir, earthlyDirSudoUser)
				if err != nil {
					earthlyDirCreateErr = errors.Wrapf(err, "failed to ensure %s is owned by %s", earthlyDir, earthlyDirSudoUser)
				}
			}
		}
	})

	return earthlyDir, earthlyDirCreateErr
}

// IsBootstrapped provides a tentatively correct guess about the state of our bootstrapping.
func IsBootstrapped(installName string) bool {
	exists, _ := fileutil.DirExists(GetEarthlyDir(installName))
	return exists
}

// EnsurePermissions changes the permissions of all earthly files to be owned by the user and their group.
func EnsurePermissions(installName string) error {
	earthlyDir, sudoUser := getEarthlyDirAndUser(installName)
	if sudoUser != nil {
		err := fileutil.EnsureUserOwned(earthlyDir, sudoUser)
		if err != nil {
			return errors.Wrapf(err, "failed to ensure %s is owned by %s", earthlyDir, sudoUser)
		}
	}
	return nil
}
