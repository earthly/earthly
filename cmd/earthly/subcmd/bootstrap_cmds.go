package subcmd

import (
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/adrg/xdg"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cmd/earthly/common"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/earthly/earthly/util/termutil"
)

type BootstrapInterface interface {
	NewBootstrap(CLI) *Bootstrap
}

type Bootstrap struct {
	cli CLI

	homebrewSource   string
	noBuildkit       bool
	genCerts         bool
	withAutocomplete bool
	certsHostName    string
}

func NewBootstrap(cli CLI) *Bootstrap {
	return &Bootstrap{
		cli: cli,
	}
}

func (b *Bootstrap) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "bootstrap",
			Usage:       "Bootstraps earthly installation including buildkit image download and optionally shell autocompletion",
			UsageText:   "earthly [options] bootstrap [--no-buildkit, --with-autocomplete, --certs-hostname]",
			Description: "Bootstraps earthly installation including buildkit image download and optionally shell autocompletion.",
			Action:      b.Action,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "source",
					Usage:       "Output source file (for use in homebrew install)",
					Hidden:      true, // only meant for use with homebrew formula
					Destination: &b.homebrewSource,
				},
				&cli.BoolFlag{
					Name:        "no-buildkit",
					Usage:       "Skips setting up the BuildKit container",
					Destination: &b.noBuildkit,
				},
				&cli.BoolFlag{
					Name:        "with-autocomplete",
					Usage:       "Install shell autocompletions during bootstrap",
					Destination: &b.withAutocomplete,
				},
				&cli.BoolFlag{
					Name:        "force-certificate-generation",
					Usage:       "Force the generation of self-signed TLS certificates, even when no BuildKit container is started",
					Destination: &b.genCerts,
				},
				&cli.StringFlag{
					Name:        "certs-hostname",
					Usage:       "Hostname to generate certificates for",
					EnvVars:     []string{"EARTHLY_CERTS_HOSTNAME"},
					Value:       "localhost",
					Destination: &b.certsHostName,
				},
			},
		},
	}
}

func (a *Bootstrap) Action(cliCtx *cli.Context) error {
	a.cli.SetCommandName("actionbootstrap")

	switch a.homebrewSource {
	case "bash":
		compEntry, err := bashCompleteEntry()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to enable bash-completion: %s\n", err)
			return nil // zsh-completion isn't available, silently fail.
		}
		fmt.Print(compEntry)
		return nil
	case "zsh":
		compEntry, err := zshCompleteEntry()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to bootstrap zsh-completion: %s\n", err)
			return nil // zsh-completion isn't available, silently fail.
		}
		fmt.Print(compEntry)
		return nil
	case "":
		break
	default:
		return errors.Errorf("unhandled source %q", a.homebrewSource)
	}

	return a.bootstrap(cliCtx)
}

func (a *Bootstrap) bootstrap(cliCtx *cli.Context) error {
	var err error
	console := a.cli.Console().WithPrefix("bootstrap")
	defer func() {
		// cliutil.IsBootstrapped() determines if bootstrapping was done based
		// on the existence of ~/.earthly; therefore we must ensure it's created.
		_, err := cliutil.GetOrCreateEarthlyDir(a.cli.Flags().InstallationName)
		if err != nil {
			console.Warnf("Warning: Failed to create Earthly Dir: %v", err)
			// Keep going.
		}
		err = cliutil.EnsurePermissions(a.cli.Flags().InstallationName)
		if err != nil {
			console.Warnf("Warning: Failed to ensure permissions: %v", err)
			// Keep going.
		}
	}()

	if a.withAutocomplete {
		// Because this requires sudo, it should warn and not fail the rest of it.
		err = a.insertBashCompleteEntry()
		if err != nil {
			console.Warnf("Warning: %s\n", err.Error())
			// Keep going.
		}
		err = a.insertZSHCompleteEntry()
		if err != nil {
			console.Warnf("Warning: %s\n", err.Error())
			// Keep going.
		}

		console.Printf("You may have to restart your shell for autocomplete to get initialized (e.g. run \"exec $SHELL\")\n")
	}
	err = symlinkEarthlyToEarth()
	if err != nil {
		console.Warnf("Warning: %s\n", err.Error())
		err = nil
	}

	if !a.noBuildkit || a.genCerts {
		bkURL, err := url.Parse(a.cli.Flags().BuildkitHost)
		if err != nil {
			return errors.Wrapf(err, "invalid buildkit_host: %s", a.cli.Flags().BuildkitHost)
		}
		if bkURL.Scheme == "tcp" && a.cli.Cfg().Global.TLSEnabled {
			err := buildkitd.GenCerts(*a.cli.Cfg(), a.certsHostName)
			if err != nil {
				return errors.Wrap(err, "failed to generate TLS certs")
			}
		}
	}

	if !a.noBuildkit {
		// connect to local buildkit instance (to trigger pulling and running the earthly/buildkitd image)
		bkClient, err := a.cli.GetBuildkitClient(cliCtx)
		if err != nil {
			console.Warnf("Warning: Bootstrapping buildkit failed: %v", err)
			// Keep going.
		} else {
			defer bkClient.Close()
		}
	}

	console.Printf("Bootstrapping successful.\n")
	return nil
}

func (a *Bootstrap) insertBashCompleteEntry() error {
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
	ok, err := a.insertBashCompleteEntryAt(path)
	if err != nil {
		return err
	}
	if ok {
		a.cli.Console().VerbosePrintf("Successfully enabled bash-completion at %s\n", path)
	} else {
		a.cli.Console().VerbosePrintf("Bash-completion already present at %s\n", path)
	}
	return nil
}

func (a *Bootstrap) insertBashCompleteEntryAt(path string) (bool, error) {
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

// If debugging this, it might be required to run `rm ~/.zcompdump*` to remove the cache
func (a *Bootstrap) insertZSHCompleteEntry() error {
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
			return a.insertZSHCompleteEntryUnderPath(dirPath)
		}
	}

	fmt.Fprintf(os.Stderr, "Warning: unable to enable zsh-completion: none of %s does not exist\n", strings.Join(potentialPaths, ", "))
	return nil // zsh-completion isn't available, silently fail.
}

func (a *Bootstrap) insertZSHCompleteEntryUnderPath(dirPath string) error {
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

	return a.deleteZcompdump()
}

func (a *Bootstrap) deleteZcompdump() error {
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

func symlinkEarthlyToEarth() error {
	binPath, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "failed to get current executable path")
	}

	baseName := path.Base(binPath)
	if baseName != "earthly" {
		return nil
	}

	earthPath := path.Join(path.Dir(binPath), "earth")

	earthPathExists, err := fileutil.FileExists(earthPath)
	if err != nil {
		return errors.Wrapf(err, "failed to check if %q exists", earthPath)
	}
	if !earthPathExists && termutil.IsTTY() {
		return nil // legacy earth binary doesn't exist, don't create it (unless we're under a non-tty system e.g. CI)
	}

	if !common.IsEarthlyBinary(earthPath) {
		return nil // file exists but is not an earthly binary, leave it alone.
	}

	// otherwise legacy earth command has been detected, remove it and symlink
	// to the new earthly command.
	err = os.Remove(earthPath)
	if err != nil {
		return errors.Wrapf(err, "failed to remove old install at %s", earthPath)
	}
	err = os.Symlink(binPath, earthPath)
	if err != nil {
		return errors.Wrapf(err, "failed to symlink %s to %s", binPath, earthPath)
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
