package fileutil

import (
	"bufio"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// GetUserHomeDirs returns a map of all users and their homedirs
func GetUserHomeDirs() (map[string]string, error) {
	users := map[string]string{}
	if runtime.GOOS == "darwin" {
		currentUser, err := user.Current()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get current user")
		}
		home := filepath.Dir(currentUser.HomeDir)
		directoryList, err := os.ReadDir(home)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read dir")
		}
		for _, s := range directoryList {
			if !s.IsDir() {
				continue
			}
			u, err := user.Lookup(s.Name())
			if err != nil {
				continue
			}
			users[u.Username] = path.Join(home, u.Username)
		}
	} else {
		fp, err := os.Open("/etc/passwd")
		if err != nil {
			return nil, errors.Wrap(err, "failed to open /etc/passwd")
		}
		reader := bufio.NewReader(fp)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, errors.Wrap(err, "failed to read line")
			}
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.Split(line, ":")
			if len(parts) >= 6 {
				user := parts[0]
				home := parts[5]
				users[user] = home
			}
		}
	}
	return users, nil
}
