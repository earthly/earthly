package fileutil

import (
	"os"
	"os/user"
	"runtime"
)

// HomeDir returns the home directory of the current user, an additional sudoUser
// is returned if the user is currently running as root
func HomeDir() (homeDir string, sudoUser *user.User, err error) {
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
		return "", nil, nil
	}
	return homeDir, nil, nil
}
