package fileutil

import (
	"os"
	"os/user"
)

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
	writable, _ := IsDirWritable(u.HomeDir)
	if !writable {
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

	// do NOT return the user here, because the user is only
	// required during the SUDO_USER case; where as this case
	// the permissions will belong to the current user and won't need changing.
	return u.HomeDir, nil, true
}

// HomeDir returns the home directory of the current user, together with
// the user object who owns it. If SUDO_USER is detected, then that user's
// home directory will be used instead.
func HomeDir() (string, *user.User) {
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
