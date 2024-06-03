package autocomplete

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/util/fileutil"

	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/urfave/cli/v2"
)

var (
	errCompPointOutOfBounds = fmt.Errorf("COMP_POINT out of bounds")
)

func isLocalPath(path string) bool {
	for _, prefix := range []string{".", "..", "/", "~"} {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func isCWDRoot() bool {
	path, err := os.Getwd()
	if err != nil {
		return false
	}
	return path == "/"
}

func containsDirectories(path string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, f := range files {
		if f.IsDir() {
			return true
		}
	}
	return false
}

func getPotentialPaths(ctx context.Context, resolver *buildcontext.Resolver, gwClient gwclient.Client, prefix string) ([]string, error) {
	if prefix == "." {
		potentials := []string{}
		if containsDirectories(".") {
			potentials = append(potentials, "./")
		}
		if !isCWDRoot() {
			potentials = append(potentials, "../")
		}
		return potentials, nil
	}
	if prefix == ".." {
		return []string{"../"}, nil
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

		targetToParse := prefix
		if strings.HasSuffix(targetToParse, "+") {
			targetToParse += "base"
		}
		target, err := domain.ParseTarget(targetToParse)
		if err != nil {
			return nil, err
		}

		targets, err := earthfile2llb.GetTargets(ctx, resolver, gwClient, target)
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

	usePrefixAsDir, err := fileutil.DirExists(prefix)
	if err != nil {
		return nil, err
	}

	if usePrefixAsDir && !strings.HasSuffix(prefix, "/") {
		p, f := path.Split(prefix)
		if hasEarthfile(f) {
			usePrefixAsDir = false
		}
		files, err := os.ReadDir(p)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			fileName := file.Name()
			if strings.HasPrefix(fileName, f) && fileName != f {
				usePrefixAsDir = false
				break
			}
		}
	}

	var f, dir string
	if usePrefixAsDir {
		dir = prefix
	} else {
		dir, f = path.Split(prefix)
	}

	files, err := os.ReadDir(dir)
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

	if prefix != "./" && !strings.HasSuffix(prefix, "/./") { // dont suggest parent directory for "./"
		if abs, _ := filepath.Abs(dir); abs != "/" { // if Abs fails, we will suggest "/../"
			if strings.HasSuffix(prefix, "/.") {
				paths = append(paths, prefix+"./")
			}
			if strings.HasSuffix(prefix, "/") {
				paths = append(paths, prefix+"../")
			}
		}
	}

	return paths, nil
}

func getPotentialTargetBuildArgs(ctx context.Context, resolver *buildcontext.Resolver, gwClient gwclient.Client, targetStr string) ([]string, error) {
	target, err := domain.ParseTarget(targetStr)
	if err != nil {
		return nil, err
	}
	envArgs, err := earthfile2llb.GetTargetArgs(ctx, resolver, gwClient, target)
	if err != nil {
		return nil, err
	}
	return envArgs, nil
}

func getPotentialArtifactBuildArgs(ctx context.Context, resolver *buildcontext.Resolver, gwClient gwclient.Client, artifactStr string) ([]string, error) {
	artifact, err := domain.ParseArtifact(artifactStr)
	if err != nil {
		return nil, err
	}
	return getPotentialTargetBuildArgs(ctx, resolver, gwClient, artifact.Target.String())
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

func getVisibleFlags(flags []cli.Flag) []string {
	visibleFlags := []string{}
	for _, f := range flags {
		if isVisibleFlag(f) {
			for _, n := range f.Names() {
				if len(n) > 1 {
					visibleFlags = append(visibleFlags, n)
				}
			}
		}
	}
	return visibleFlags
}

func isBooleanFlag(flags []cli.Flag, flagName string) (isBool bool, flagFound bool) {
	if flagName == "" {
		return false, false
	}

	isShort := true
	if strings.HasPrefix(flagName, "--") {
		flagName = flagName[2:]
		isShort = false
	} else {
		flagName = flagName[1:]
	}
	_ = isShort // short flags are not suggested; perhaps one day?

	for _, f := range flags {
		for _, n := range f.Names() {
			if n == flagName {
				_, ok := f.(*cli.BoolFlag)
				return ok, true
			}
		}
	}
	return false, false
}

func isFlagValidAndRequiresValue(flags []cli.Flag, flagName string) bool {
	isBool, ok := isBooleanFlag(flags, flagName)
	return ok && !isBool
}

// padStrings takes an array of strings and returns a new array where each
// string element has been padded with a prefix and suffix
func padStrings(flags []string, prefix, suffix string) []string {
	padded := make([]string, len(flags))
	for i, s := range flags {
		padded[i] = prefix + s + suffix
	}
	return padded
}

func getVisibleCommands(commands []*cli.Command) []string {
	visibleCommands := []string{}
	for _, cmd := range commands {
		if !cmd.Hidden {
			visibleCommands = append(visibleCommands, cmd.Name)
		}
	}
	return visibleCommands
}

type completeState int

const (
	unknownState          completeState = iota
	rootState                           // 1
	flagState                           // 2
	flagValueState                      // 3
	commandState                        // 4
	targetState                         // 5
	targetFlagState                     // 6
	artifactFlagState                   // 7
	endOfSuggestionsState               // 8
)

type FlagValuePotentialFn func(ctx context.Context, prefix string) []string

// GetPotentials returns a list of potential arguments for shell auto completion
// NOTE: you can cause earthly to run this command with:
//
//	COMP_LINE="earthly -" COMP_POINT=$(echo -n $COMP_LINE | wc -c) go run cmd/earthly/main.go
func GetPotentials(ctx context.Context, resolver *buildcontext.Resolver, gwClient gwclient.Client, compLine string, compPoint int, app *cli.App, cloudClient cloudListClient) ([]string, error) {
	if compPoint > len(compLine) {
		return nil, errCompPointOutOfBounds
	}
	compLine = compLine[:compPoint]
	subCommands := app.Commands

	// TODO all the urfave/cli commands need to be moved out of the main package
	// so they could be directly referenced rather than storing a list of strings of seen commands
	commandValues := []string{}
	getPrevCommand := func() string {
		n := len(commandValues) - 2
		if n >= 0 {
			return commandValues[n]
		}
		return ""
	}

	flagValues := map[string]string{}
	flagValuePotentialFuncs := map[string]FlagValuePotentialFn{
		"--org": func(ctx context.Context, prefix string) []string {
			if cloudClient == nil {
				return []string{}
			}

			orgs, err := cloudClient.ListOrgs(ctx)
			if err != nil {
				return []string{}
			}
			potentials := []string{}
			for _, org := range orgs {
				potentials = append(potentials, org.Name)
			}
			return potentials
		},
		"--project": func(ctx context.Context, prefix string) []string {
			if cloudClient == nil {
				return []string{}
			}

			org, ok := flagValues["--org"]
			if !ok {
				return []string{}
			}

			projects, err := cloudClient.ListProjects(ctx, org)
			if err != nil {
				return []string{}
			}
			potentials := []string{}
			for _, project := range projects {
				potentials = append(potentials, project.Name)
			}
			return potentials
		},
		"--satellite": func(ctx context.Context, prefix string) []string {
			if cloudClient == nil {
				return []string{}
			}

			org, ok := flagValues["--org"]
			if !ok {
				return []string{}
			}

			satellites, err := cloudClient.ListSatellites(ctx, org)
			if err != nil {
				return []string{}
			}
			potentials := []string{}
			for _, sat := range satellites {
				potentials = append(potentials, sat.Name)
			}
			return potentials
		},
	}

	// getWord returns the next word and a boolean if it is valid
	// TODO this function does not handle escaped space, e.g.
	// earthly --build-arg key="value with space" +mytarget will fail
	hasNextWord := len(compLine) > 0
	getWord := func() (string, bool) {
		if !hasNextWord {
			return "", false
		}
		i := strings.Index(compLine, " ")
		if i < 0 {
			word := compLine
			compLine = ""
			hasNextWord = false
			return word, true
		}
		word := compLine[:i]
		compLine = compLine[(i + 1):]
		return word, true
	}

	// remove first word which is most likely "earthly", or "/some/path/to/earthly", etc.
	prevWord, _ := getWord()

	state := rootState
	var target string
	var artifactMode bool
	var flag string

	var cmd *cli.Command
	getFlags := func() []cli.Flag {
		if cmd != nil {
			return cmd.Flags
		}
		return app.Flags
	}

	for {
		w, ok := getWord()
		if !ok {
			break
		}

		if state == flagValueState {
			flagValues[flag] = prevWord
			prevWord = ""
			state = flagState
		}

		if state == flagState && isFlagValidAndRequiresValue(getFlags(), prevWord) {
			state = flagValueState
			flag = prevWord
		} else if state == rootState || state == commandState || state == flagState {
			if strings.HasPrefix(w, "-") {
				if w == "-a" || w == "--artifact" {
					artifactMode = true
				}
				state = flagState
			} else {
				// targets only work under the root command
				if cmd == nil && (isLocalPath(w) || strings.HasPrefix(w, "+")) { // TODO switch to strings.Contains when remote resolving works
					state = targetState
					target = w
				} else {
					// must be under a command
					foundCmd := getCmd(w, subCommands)
					if foundCmd != nil {
						subCommands = foundCmd.Subcommands
						cmd = foundCmd

						// TODO once urfave/cli commands are moved out of main, this should be removed (and instead the cmd pointer could simply be compared to determine which command we are referencing)
						commandValues = append(commandValues, cmd.Name)
					}
					state = commandState
				}
			}
		} else if state == targetState || state == targetFlagState {
			if !strings.HasPrefix(w, "-") {
				state = endOfSuggestionsState
			} else {
				if strings.HasSuffix(w, "=") {
					state = endOfSuggestionsState
				} else {
					if artifactMode {
						state = artifactFlagState
					} else {
						state = targetFlagState
					}
				}
			}
		}

		prevWord = w
	}

	var potentials []string

	switch state {
	case flagState:
		if cmd != nil {
			potentials = getVisibleFlags(cmd.Flags)
		} else {
			potentials = getVisibleFlags(app.Flags)
			// append flags that urfav/cli automatically include
			potentials = append(potentials, "version", "help")
		}
		potentials = padStrings(potentials, "--", " ")

	case flagValueState:
		if fn, ok := flagValuePotentialFuncs[flag]; ok {
			potentials = append(potentials, fn(ctx, prevWord)...)
		}

	case rootState, commandState:
		if cmd != nil {
			potentials = getVisibleCommands(cmd.Subcommands)
			potentials = padStrings(potentials, "", " ")

			// TODO this should be tied to the instance of the command (and not just command Name value); but that means moving lots out of the main package
			if getPrevCommand() == "satellite" && (cmd.Name == "inspect" || cmd.Name == "rm" || cmd.Name == "select" || cmd.Name == "sleep" || cmd.Name == "update" || cmd.Name == "wake") {
				potentials = append(potentials, flagValuePotentialFuncs["--satellite"](ctx, "")...)
			}
		} else {
			potentials = getVisibleCommands(app.Commands)
			potentials = padStrings(potentials, "", " ")
			if containsDirectories(".") {
				potentials = append(potentials, "./")
			}
			if fileutil.FileExistsBestEffort("Earthfile") {
				potentials = append(potentials, "+")
			}
		}

	case targetState:
		var err error
		potentials, err = getPotentialPaths(ctx, resolver, gwClient, prevWord)
		if err != nil {
			return nil, err
		}

	case targetFlagState:
		var err error
		potentials, err = getPotentialTargetBuildArgs(ctx, resolver, gwClient, target)
		if err != nil {
			return nil, err
		}
		potentials = padStrings(potentials, "--", "=")

	case artifactFlagState:
		var err error
		potentials, err = getPotentialArtifactBuildArgs(ctx, resolver, gwClient, target)
		if err != nil {
			return nil, err
		}
		potentials = padStrings(potentials, "--", "=")
	}

	filteredPotentials := []string{}
	for _, s := range potentials {
		if strings.HasPrefix(s, prevWord) {
			filteredPotentials = append(filteredPotentials, s)
		}
	}

	sort.Strings(filteredPotentials)
	return filteredPotentials, nil
}

func hasEarthfile(dirPath string) bool {
	info, err := os.Stat(path.Join(dirPath, "Earthfile"))
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func hasSubDirs(dirPath string) bool {
	files, err := os.ReadDir(dirPath)
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
