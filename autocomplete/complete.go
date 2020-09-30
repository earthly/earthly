package autocomplete

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/earthly/earthly/earthfile2llb"
)

func hasTargetOrCommand(line string) bool {
	splits := strings.Split(line, " ")
	for i, s := range splits {
		if i == 0 {
			continue // skip earth command
		}
		if len(s) == 0 {
			continue // skip empty commands
		}
		if s[0] == '-' {
			continue // skip flags
		}
		return true // found a command or target
	}
	return false
}

// parseLine parses a bash COMP_LINE and COMP_POINT variables into the argument to expand
// e.g. line="earth --argum", cursorLoc=10; this will return "--ar"
func parseLine(line string, cursorLoc int) string {
	var i int
	for i = cursorLoc; i > 0; i-- {
		if line[i-1] == ' ' {
			break
		}
	}
	if i >= cursorLoc {
		return ""
	}
	return line[i:cursorLoc]
}

func trimFlag(prefix string) (string, bool) {
	if len(prefix) == 1 && prefix[0] == '-' {
		return "", true
	}
	if strings.HasPrefix(prefix, "--") {
		return prefix[2:], true
	}
	return "", false
}

func isLocalPath(path string) bool {
	for _, prefix := range []string{".", "..", "/", "~"} {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func getPotentialTarget(prefix string) ([]string, error) {
	splits := strings.SplitN(prefix, "+", 2)
	if len(splits) < 2 {
		return []string{}, nil
	}
	dirPath := splits[0]

	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	realDirPath := dirPath
	if strings.HasPrefix(prefix, "~/") {
		realDirPath = currentUser.HomeDir + "/" + dirPath[2:]
	}

	targets, err := earthfile2llb.GetTargets(path.Join(realDirPath, "Earthfile"))
	if err != nil {
		return nil, err
	}

	potentials := []string{}
	for _, target := range targets {
		s := dirPath + "+" + target + " "
		if strings.HasPrefix(s, prefix) {
			potentials = append(potentials, s)
		}
	}

	return potentials, nil
}

func lookupUsers(prefix string) ([]string, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	lPath := strings.Split(u.HomeDir, "/")
	if len(lPath) < 2 {
		return nil, errors.New("users not found")
	}

	directoryList, err := ioutil.ReadDir(fmt.Sprintf("/%s", lPath[1]))
	if err != nil {
		return nil, err
	}
	var res []string
	for _, s := range directoryList {
		if prefix != "" && !strings.HasPrefix(s.Name(), prefix) {
			continue
		}

		if !s.IsDir() {
			continue
		}

		u, err := user.Lookup(s.Name())
		if err != nil {
			continue
		}
		res = append(res, fmt.Sprintf("~%s", u.Username))
	}

	if prefix == "~" {
		return res, nil
	}

	return res, nil
}

func getUserFromPrefix(prefix string) (*user.User, error) {
	username := strings.Replace(prefix, "~", "", 1)
	spl := strings.Split(username, "/")
	if len(spl) != 1 {
		username = spl[0]
	}
	return user.Lookup(username)
}

func filterUsers(prefix string, lUsers []string) []string {
	var filter []string
	for _, l := range lUsers {
		if strings.HasPrefix(l, prefix) {
			filter = append(filter, l)
		}
	}
	return filter
}

func getPotentialPaths(prefix string) ([]string, error) {
	switch {
	case prefix == ".":
		return []string{"./", "../"}, nil
	case prefix == "~":
		return lookupUsers(strings.Replace(prefix, "~", "", 1))
	case strings.HasPrefix(prefix, "~") && !strings.Contains(prefix, "/"):
		lUser, err := lookupUsers("")
		if err != nil {
			return nil, err
		}
		return filterUsers(prefix, lUser), nil
	}

	expandedHomeLen := 0
	if strings.HasPrefix(prefix, "~") {
		// get user from prefix
		currentUser, err := getUserFromPrefix(prefix)
		if err != nil {
			return nil, err
		}

		expandedHomeLen = len(currentUser.HomeDir) + 1
		prefix = strings.Replace(prefix, "~", strings.Replace(currentUser.HomeDir, currentUser.Username, "", 1), 1)
		if strings.HasSuffix(prefix, currentUser.Username) {
			prefix = fmt.Sprintf("%s/", prefix)
		}
	}

	dir, f := path.Split(prefix)

	replaceHomePrefix := func(s string) string {
		if expandedHomeLen == 0 {
			return s
		}
		return "~/" + s[expandedHomeLen:]
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	paths := []string{}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), f) {
			if file.IsDir() {
				s := path.Join(dir, file.Name())
				if strings.HasPrefix(prefix, "./") {
					s = "./" + s
				}
				if hasEarthfile(s) {
					paths = append(paths, replaceHomePrefix(s)+"+")
				}
				if hasSubDirs(s) {
					paths = append(paths, replaceHomePrefix(s)+"/")
				}
			}
		}
	}
	return paths, nil
}

// GetPotentials returns a list of potential arguments for shell auto completion
func GetPotentials(compLine string, compPoint int, flags, commands []string) ([]string, error) {
	potentials := []string{}

	prefix := parseLine(compLine, compPoint)

	// already has a full command or target (we're done now)
	if hasTargetOrCommand(compLine) && prefix == "" {
		return potentials, nil
	}

	if flagPrefix, ok := trimFlag(prefix); ok {
		for _, s := range flags {
			if strings.HasPrefix(s, flagPrefix) {
				potentials = append(potentials, "--"+s+" ")
			}
		}
		return potentials, nil
	}

	if isLocalPath(prefix) || strings.HasPrefix(prefix, "+") {
		if strings.Contains(prefix, "+") {
			return getPotentialTarget(prefix)
		}
		return getPotentialPaths(prefix)
	}

	if prefix == "" {
		if hasEarthfile(".") {
			potentials = append(potentials, "+")
		}
		potentials = append(potentials, "./")
	}

	for _, cmd := range commands {
		if strings.HasPrefix(cmd, prefix) {
			potentials = append(potentials, cmd+" ")
		}
	}
	return potentials, nil
}

func hasEarthfile(dirPath string) bool {
	info, err := os.Stat(path.Join(dirPath, "Earthfile"))
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func hasSubDirs(dirPath string) bool {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return false
	}

	for _, file := range files {
		if file.IsDir() {
			return true
		}
	}
	return false
}
