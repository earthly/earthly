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
var earthlyDirOnce sync.Once
var earthlyDirSudoUser *user.User

var earthlyDirCreateOnce sync.Once
var earthlyDirCreateErr error

// GetEarthlyDir returns the .earthly dir. (Usually ~/.earthly).
// This function will not attempt to create the directory if missing, for that functionality use to the GetOrCreateEarthlyDir function.
func GetEarthlyDir() string {
	earthlyDirOnce.Do(func() {
		earthlyDir, earthlyDirSudoUser = getEarthlyDirAndUser()
	})
	return earthlyDir
}

func getEarthlyDirAndUser() (string, *user.User) {
	homeDir, u := DetectHomeDir()
	earthlyDir := filepath.Join(homeDir, ".earthly")
	return earthlyDir, u
}

// GetOrCreateEarthlyDir returns the .earthly dir. (Usually ~/.earthly).
// if the directory does not exist, it will attempt to create it.
func GetOrCreateEarthlyDir() (string, error) {
	_ = GetEarthlyDir() // ensure global vars get created so we can reference them below.

	earthlyDirCreateOnce.Do(func() {
		if !fileutil.DirExists(earthlyDir) {
			err := os.MkdirAll(earthlyDir, 0755)
			if err != nil {
				earthlyDirCreateErr = errors.Wrapf(err, "unable to create dir %s", earthlyDir)
				return
			}
			if earthlyDirSudoUser != nil {
				fileutil.EnsureUserOwned(earthlyDir, earthlyDirSudoUser)
			}
		}
	})

	return earthlyDir, earthlyDirCreateErr
}

func getHomeFromSudoUser() (string, *user.User, bool) {
	sudoUserName, ok := os.LookupEnv("SUDO_USER")
	if !ok {
		return "", nil, false
	}
	u, err := user.Lookup(sudoUserName)
	if err != nil {
		return "", nil, false
	}
	if u.HomeDir == "" {
		return "", nil, false
	}
	return u.HomeDir, u, true
}

func getHomeFromHomeEnv() (string, *user.User, bool) {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		return "", nil, false
	}
	return home, nil, true
}

func getHomeFromUserCurrent() (string, *user.User, bool) {
	u, err := user.Current()
	if err != nil {
		return "", nil, false
	}
	if u.HomeDir == "" {
		return "", nil, false
	}
	return u.HomeDir, u, true
}

// DetectHomeDir returns the home directory of the current user, together with
// the user object who owns it. If SUDO_USER is detected, then that user's
// home directory will be used instead.
func DetectHomeDir() (string, *user.User) {
	for _, fn := range []func() (string, *user.User, bool){
		getHomeFromSudoUser,
		getHomeFromHomeEnv,
		getHomeFromUserCurrent,
	} {
		home, u, ok := fn()
		if ok {
			return home, u
		}
	}
	return "/etc", nil
}

// IsBootstrapped provides a tentatively correct guess about the state of our bootstrapping.
func IsBootstrapped() bool {
	return fileutil.DirExists(GetEarthlyDir())
}

// EnsurePermissions changes the permissions of all earthly files to be owned by the user and their group.
func EnsurePermissions() error {
	earthlyDir, sudoUser := getEarthlyDirAndUser()
	if sudoUser != nil {
		fileutil.EnsureUserOwned(earthlyDir, sudoUser)
	}
	return nil
}
