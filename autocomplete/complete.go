package autocomplete

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"reflect"
	"strings"

	"github.com/earthly/earthly/earthfile2llb"

	"github.com/urfave/cli/v2"
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

func getPotentialPaths(prefix string) ([]string, error) {
	if prefix == "." {
		return []string{"./", "../"}, nil
	}
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	// TODO expand this logic to support other users
	if prefix == "~" {
		prefix = "~/"
	}

	expandedHomeLen := 0
	if strings.HasPrefix(prefix, "~") {
		expandedHomeLen = len(currentUser.HomeDir) + 1
		if len(prefix) > 2 {
			prefix = currentUser.HomeDir + "/" + prefix[2:]
		} else {
			prefix = currentUser.HomeDir + "/"
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

// isHidden returns if a flag is hidden or not
// this code comes from https://github.com/urfave/cli/blob/d648edd48d89ef3a841b1ec75c2ebbd4de5f748f/flag.go#L136
func isVisibleFlag(fl cli.Flag) bool {
	fv := reflect.ValueOf(fl)
	for fv.Kind() == reflect.Ptr {
		fv = reflect.Indirect(fv)
	}
	field := fv.FieldByName("Hidden")
	if !field.IsValid() || !field.Bool() {
		return true
	}
	return false
}

// GetPotentials returns a list of potential arguments for shell auto completion
func GetPotentials(compLine string, compPoint int, app *cli.App) ([]string, error) {
	potentials := []string{}

	flags := []string{}
	for _, f := range app.Flags {
		if isVisibleFlag(f) {
			for _, n := range f.Names() {
				if len(n) > 1 {
					flags = append(flags, n)
				}
			}
		}
	}

	commands := []string{}
	for _, cmd := range app.Commands {
		if !cmd.Hidden {
			commands = append(commands, cmd.Name)
		}
	}

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
