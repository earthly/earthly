package fileutil

import (
	"strings"
)

// ExpandPath expands the tilde in the path
func ExpandPath(s string) string {
	if !strings.HasPrefix(s, "~") {
		return s
	}
	homeDir, err := HomeDir()
	if err != nil || homeDir == "" {
		return s // best effort
	}

	if s == "~" {
		return homeDir
	}
	parts := strings.SplitN(s, "/", 2)
	if parts[0] != "~" {
		user := parts[0][1:]
		users, err := GetUserHomeDirs()
		if err != nil {
			return s // best effort
		}
		homeDir = users[user]
		if homeDir == "" {
			return s // best effort
		}
	}
	if len(parts) == 1 {
		return homeDir
	}
	return homeDir + "/" + parts[1]
}
