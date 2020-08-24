package autocomplete

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/earthly/earthly/earthfile2llb"
)

type CompletionType int

const (
	Flag CompletionType = iota
	Command
	Path
	Unknown
)

// ParseLine parses a bash COMP_LINE and COMP_POINT variables into the argument to expand
// e.g. line="earth --argum", cursorLoc=10; this will return "--ar"
func ParseLine(line string, cursorLoc int) string {
	i := cursorLoc
	if i > 0 {
		i--
	}
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

func isFlag(prefix string) bool {
	if len(prefix) == 1 && prefix[0] == '-' {
		return true
	}
	return strings.HasPrefix(prefix, "--")
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
	for _, prefix := range []string{".", "..", "/"} {
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

	targets, err := earthfile2llb.GetTargets(path.Join(dirPath, "Earthfile"))
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

func getPotentialPaths(prefix string) ([]string, error) {
	if prefix == "." {
		return []string{"./", "../"}, nil
	}
	dir, f := path.Split(prefix)

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
					paths = append(paths, s+"+")
				}
				if hasSubDirs(s) {
					paths = append(paths, s+"/")
				}
			}
		}
	}
	return paths, nil
}

// GetPotentials returns a list of potential arguments for shell auto completion
func GetPotentials(prefix string, flags, commands []string) ([]string, error) {
	potentials := []string{}
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
