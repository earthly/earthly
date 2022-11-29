package main

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/earthly/earthly/autocomplete"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/fileutil"

	"github.com/adrg/xdg"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

// to enable autocomplete, enter
// complete -o nospace -C "/path/to/earthly" earthly
func (app *earthlyApp) autoComplete(ctx context.Context) {
	_, found := os.LookupEnv("COMP_LINE")
	if !found {
		return
	}

	app.console = app.console.WithLogLevel(conslogging.Silent)

	err := app.autoCompleteImp(ctx)
	if err != nil {
		errToLog := err
		logDir, err := cliutil.GetOrCreateEarthlyDir(app.installationName)
		if err != nil {
			os.Exit(1)
		}
		logFile := filepath.Join(logDir, "autocomplete.log")
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			os.Exit(1)
		}
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			os.Exit(1)
		}
		fmt.Fprintf(f, "error during autocomplete: %s\n", errToLog)
		os.Exit(1)
	}
	os.Exit(0)
}

func (app *earthlyApp) autoCompleteImp(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("recovered panic in autocomplete %s: %s", r, debug.Stack())
		}
	}()

	compLine := os.Getenv("COMP_LINE")   // full command line
	compPoint := os.Getenv("COMP_POINT") // where the cursor is

	compPointInt, err := strconv.ParseUint(compPoint, 10, 64)
	if err != nil {
		return err
	}

	gitLookup := buildcontext.NewGitLookup(app.console, app.sshAuthSock)
	resolver := buildcontext.NewResolver(nil, gitLookup, app.console, "")
	var gwClient gwclient.Client // TODO this is a nil pointer which causes a panic if we try to expand a remotely referenced earthfile
	// it's expensive to create this gwclient, so we need to implement a lazy eval which returns it when required.

	potentials, err := autocomplete.GetPotentials(ctx, resolver, gwClient, compLine, int(compPointInt), app.cliApp)
	if err != nil {
		return err
	}
	for _, p := range potentials {
		fmt.Printf("%s\n", p)
	}

	return err
}

func (app *earthlyApp) insertBashCompleteEntry() error {
	u, err := user.Current()
	if err != nil {
		return errors.Wrapf(err, "could not get current user")
	}
	isRootUser := u.Uid == "0"
	var path string
	// Assume that non-root can't write to the system and that installation
	// to root's home isn't desirable.  One possible exception might be if
	// those directories are on an R/O filesystem, but user can install these
	// manually in that case.
	if isRootUser {
		if runtime.GOOS == "darwin" {
			path = "/usr/local/etc/bash_completion.d/earthly"
		} else {
			path = "/usr/share/bash-completion/completions/earthly"
		}
	} else {
		// https://github.com/scop/bash-completion/blob/master/README.md#faq
		userPath, ok := os.LookupEnv("BASH_COMPLETION_USER_DIR")
		if !ok {
			// This will give a standardized fallback even if XDG isn't active
			userPath = xdg.DataHome
		}
		path = filepath.Join(userPath, "bash-completion/completions/earthly")
	}
	ok, err := app.insertBashCompleteEntryAt(path)
	if err != nil {
		return err
	}
	if ok {
		app.console.VerbosePrintf("Successfully enabled bash-completion at %s\n", path)
	} else {
		app.console.VerbosePrintf("Bash-completion already present at %s\n", path)
	}
	return nil
}

func (app *earthlyApp) insertBashCompleteEntryAt(path string) (bool, error) {
	dirPath := filepath.Dir(path)

	dirPathExists, err := fileutil.DirExists(dirPath)
	if err != nil {
		return false, errors.Wrapf(err, "failed checking if %s exists", dirPath)
	}
	if !dirPathExists {
		return false, errors.New(fmt.Sprintf("%s does not exist", dirPath))
	}

	pathExists, err := fileutil.FileExists(path)
	if err != nil {
		return false, errors.Wrapf(err, "failed checking if %s exists", path)
	}
	if pathExists {
		return false, nil // file already exists, don't update it.
	}

	// create the completion file
	f, err := os.Create(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	bashEntry, err := bashCompleteEntry()
	if err != nil {
		return false, errors.Wrapf(err, "failed to add entry")
	}

	_, err = f.Write([]byte(bashEntry))
	if err != nil {
		return false, errors.Wrapf(err, "failed writing to %s", path)
	}
	return true, nil
}

func (app *earthlyApp) deleteZcompdump() error {
	var homeDir string
	sudoUser, found := os.LookupEnv("SUDO_USER")
	if !found {
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return errors.Wrapf(err, "failed to lookup current user home dir")
		}
	} else {
		currentUser, err := user.Lookup(sudoUser)
		if err != nil {
			return errors.Wrapf(err, "failed to lookup user %s", sudoUser)
		}
		homeDir = currentUser.HomeDir
	}
	files, err := os.ReadDir(homeDir)
	if err != nil {
		return errors.Wrapf(err, "failed to read dir %s", homeDir)
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".zcompdump") {
			path := filepath.Join(homeDir, f.Name())
			err := os.Remove(path)
			if err != nil {
				return errors.Wrapf(err, "failed to remove %s", path)
			}
		}
	}
	return nil
}

func bashCompleteEntry() (string, error) {
	template := "complete -o nospace -C '__earthly__' earthly\n"
	return renderEntryTemplate(template)
}

func zshCompleteEntry() (string, error) {
	template := `#compdef _earthly earthly

function _earthly {
    autoload -Uz bashcompinit
    bashcompinit
    complete -o nospace -C '__earthly__' earthly
}
`
	return renderEntryTemplate(template)
}

func renderEntryTemplate(template string) (string, error) {
	earthlyPath, err := os.Executable()
	if err != nil {
		return "", errors.Wrapf(err, "failed to determine earthly path: %s", err)
	}
	return strings.ReplaceAll(template, "__earthly__", earthlyPath), nil
}

// If debugging this, it might be required to run `rm ~/.zcompdump*` to remove the cache
func (app *earthlyApp) insertZSHCompleteEntry() error {
	potentialPaths := []string{
		"/usr/local/share/zsh/site-functions",
		"/usr/share/zsh/site-functions",
	}
	for _, dirPath := range potentialPaths {
		dirPathExists, err := fileutil.DirExists(dirPath)
		if err != nil {
			return errors.Wrapf(err, "failed to check if %s exists", dirPath)
		}
		if dirPathExists {
			return app.insertZSHCompleteEntryUnderPath(dirPath)
		}
	}

	fmt.Fprintf(os.Stderr, "Warning: unable to enable zsh-completion: none of %s does not exist\n", strings.Join(potentialPaths, ", "))
	return nil // zsh-completion isn't available, silently fail.
}

func (app *earthlyApp) insertZSHCompleteEntryUnderPath(dirPath string) error {
	path := filepath.Join(dirPath, "_earthly")

	pathExists, err := fileutil.FileExists(path)
	if err != nil {
		return errors.Wrapf(err, "failed to check if %s exists", path)
	}
	if pathExists {
		return nil // file already exists, don't update it.
	}

	// create the completion file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	compEntry, err := zshCompleteEntry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: unable to enable zsh-completion: %s\n", err)
		return nil // zsh-completion isn't available, silently fail.
	}

	_, err = f.Write([]byte(compEntry))
	if err != nil {
		return errors.Wrapf(err, "failed writing to %s", path)
	}

	return app.deleteZcompdump()
}
