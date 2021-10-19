package autocomplete

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"reflect"
	"strings"

	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/util/fileutil"

	"github.com/urfave/cli/v2"
)

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

func getPotentialPaths(prefix string) ([]string, error) {
	if prefix == "." {
		return []string{"./", "../"}, nil
	}
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	var user string
	expandedHomeLen := 0
	if strings.HasPrefix(prefix, "~") {
		users, err := fileutil.GetUserHomeDirs()
		if err != nil {
			return nil, err
		}

		// handle username completion
		if !strings.Contains(prefix, "/") {
			potentials := []string{}
			for user := range users {
				if strings.HasPrefix(user, prefix[1:]) {
					potentials = append(potentials, "~"+user+"/")
				}
			}
			return potentials, nil
		}

		// otherwise expand ~ into complete path
		parts := strings.SplitN(prefix[1:], "/", 2)
		user = parts[0]
		rest := parts[1]

		var homeDir string
		if len(user) == 0 {
			homeDir = currentUser.HomeDir
		} else {
			homeDir = users[user]
		}

		expandedHomeLen = len(currentUser.HomeDir) + 1
		prefix = homeDir + "/" + rest
	}

	replaceHomePrefix := func(s string) string {
		if expandedHomeLen == 0 {
			return s
		}
		return "~" + user + "/" + s[expandedHomeLen:]
	}

	// handle targets
	if strings.Contains(prefix, "+") {
		splits := strings.SplitN(prefix, "+", 2)
		if len(splits) < 2 {
			return []string{}, nil
		}
		dirPath := splits[0]

		targets, err := earthfile2llb.GetTargets(path.Join(dirPath, "Earthfile"))
		if err != nil {
			return nil, err
		}
		if len(targets) == 0 {
			// only suggest when Earthfile has no other targets
			targets = append(targets, "base")
		}

		potentials := []string{}
		for _, target := range targets {
			s := dirPath + "+" + target + " "
			if strings.HasPrefix(s, prefix) {
				potentials = append(potentials, replaceHomePrefix(s))
			}
		}

		return potentials, nil
	}

	// handle paths
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

// isVisibleFlag returns if a flag is hidden or not
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

func getCmd(name string, cmds []*cli.Command) *cli.Command {
	for _, c := range cmds {
		if name == c.Name {
			return c
		}
	}
	return nil
}

func getVisibleFlags(flags []cli.Flag, showHidden bool) []string {
	visibleFlags := []string{}
	for _, f := range flags {
		if isVisibleFlag(f) || showHidden {
			for _, n := range f.Names() {
				if len(n) > 1 {
					visibleFlags = append(visibleFlags, n)
				}
			}
		}
	}
	return visibleFlags
}

func getVisibleCommands(commands []*cli.Command, showHidden bool) []string {
	visibleCommands := []string{}
	for _, cmd := range commands {
		if !cmd.Hidden || showHidden {
			visibleCommands = append(visibleCommands, cmd.Name)
		}
	}
	return visibleCommands
}

// GetPotentials returns a list of potential arguments for shell auto completion
func GetPotentials(compLine string, compPoint int, app *cli.App, showHidden bool) ([]string, error) {
	potentials := []string{}

	compLine = compLine[:compPoint]
	subCommands := app.Commands

	// determine sub command
	parts := strings.Split(compLine, " ")
	var cmd *cli.Command
	for _, word := range parts[1 : len(parts)-1] {
		foundCmd := getCmd(word, subCommands)
		if foundCmd != nil {
			subCommands = foundCmd.Subcommands
			cmd = foundCmd
		}
	}
	lastWord := parts[len(parts)-1]

	var flags []string
	if cmd != nil {
		flags = getVisibleFlags(cmd.Flags, showHidden)
	} else {
		flags = getVisibleFlags(app.Flags, showHidden)
		// append flags that urfav/cli automatically include
		flags = append(flags, "version", "help")
	}

	var commands []string
	if cmd != nil {
		commands = getVisibleCommands(cmd.Subcommands, showHidden)
	} else {
		commands = getVisibleCommands(app.Commands, showHidden)
	}

	if flagPrefix, ok := trimFlag(lastWord); ok {
		for _, s := range flags {
			if strings.HasPrefix(s, flagPrefix) {
				potentials = append(potentials, "--"+s+" ")
			}
		}
		return potentials, nil
	}

	if isLocalPath(lastWord) || strings.HasPrefix(lastWord, "+") {
		return getPotentialPaths(lastWord)
	}

	if lastWord == "" && cmd == nil {
		if hasEarthfile(".") {
			potentials = append(potentials, "+")
		}
		if hasSubDirs(".") {
			potentials = append(potentials, "./")
		}
	}

	for _, cmd := range commands {
		if strings.HasPrefix(cmd, lastWord) {
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
