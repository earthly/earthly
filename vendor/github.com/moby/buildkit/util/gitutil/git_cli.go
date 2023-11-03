package gitutil

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/moby/buildkit/util/bklog"
	"github.com/moby/buildkit/util/urlutil"
	"github.com/pkg/errors"
)

// arthlyCtxDebugLevelKey is earthly-specific and is used to pass along the debug level
const EarthlyCtxDebugLevelKey = "EARTHLY_DEBUG_LEVEL" // earthly-specific

// GitLogLevel is earthly-specific
type GitLogLevel int

const (
	GitLogLevelDefault GitLogLevel = iota
	GitLogLevelDebug
	GitLogLevelTrace
)

// GitCLI carries config to pass to the git cli to make running multiple
// commands less repetitive.
type GitCLI struct {
	git  string
	exec func(context.Context, *exec.Cmd) error

	args    []string
	dir     string
	streams StreamFunc

	workTree string
	gitDir   string

	sshAuthSock   string
	sshKnownHosts string
}

// Option provides a variadic option for configuring the git client.
type Option func(b *GitCLI)

// WithGitBinary sets the git binary path.
func WithGitBinary(path string) Option {
	return func(b *GitCLI) {
		b.git = path
	}
}

// WithExec sets the command exec function.
func WithExec(exec func(context.Context, *exec.Cmd) error) Option {
	return func(b *GitCLI) {
		b.exec = exec
	}
}

// WithArgs sets extra args.
func WithArgs(args ...string) Option {
	return func(b *GitCLI) {
		b.args = append(b.args, args...)
	}
}

// WithDir sets working directory.
//
// This should be a path to any directory within a standard git repository.
func WithDir(dir string) Option {
	return func(b *GitCLI) {
		b.dir = dir
	}
}

// WithWorkTree sets the --work-tree arg.
//
// This should be the path to the top-level directory of the checkout. When
// setting this, you also likely need to set WithGitDir.
func WithWorkTree(workTree string) Option {
	return func(b *GitCLI) {
		b.workTree = workTree
	}
}

// WithGitDir sets the --git-dir arg.
//
// This should be the path to the .git directory. When setting this, you may
// also need to set WithWorkTree, unless you are working with a bare
// repository.
func WithGitDir(gitDir string) Option {
	return func(b *GitCLI) {
		b.gitDir = gitDir
	}
}

// WithSSHAuthSock sets the ssh auth sock.
func WithSSHAuthSock(sshAuthSock string) Option {
	return func(b *GitCLI) {
		b.sshAuthSock = sshAuthSock
	}
}

// WithSSHKnownHosts sets the known hosts file.
func WithSSHKnownHosts(sshKnownHosts string) Option {
	return func(b *GitCLI) {
		b.sshKnownHosts = sshKnownHosts
	}
}

type StreamFunc func(context.Context) (io.WriteCloser, io.WriteCloser, func())

// WithStreams configures a callback for getting the streams for a command. The
// stream callback will be called once for each command, and both writers will
// be closed after the command has finished.
func WithStreams(streams StreamFunc) Option {
	return func(b *GitCLI) {
		b.streams = streams
	}
}

// New initializes a new git client
func NewGitCLI(opts ...Option) *GitCLI {
	c := &GitCLI{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// New returns a new git client with the same config as the current one, but
// with the given options applied on top.
func (cli *GitCLI) New(opts ...Option) *GitCLI {
	c := &GitCLI{
		git:           cli.git,
		dir:           cli.dir,
		workTree:      cli.workTree,
		gitDir:        cli.gitDir,
		args:          append([]string{}, cli.args...),
		streams:       cli.streams,
		sshAuthSock:   cli.sshAuthSock,
		sshKnownHosts: cli.sshKnownHosts,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func gitDebug() bool {
	return os.Getenv("BUILDKIT_DEBUG_GIT") == "1"
}

// Run executes a git command with the given args.
func (cli *GitCLI) Run(ctx context.Context, args ...string) (_ []byte, err error) {
	gitBinary := "git"
	if cli.git != "" {
		gitBinary = cli.git
	}

	// earthly-specific: smuggle in the logLevel
	logLevel, ok := ctx.Value(EarthlyCtxDebugLevelKey).(GitLogLevel)
	if !ok {
		bklog.G(ctx).Warnf("failed to extract %s", EarthlyCtxDebugLevelKey)
	}
	if gitDebug() || logLevel >= GitLogLevelDebug { // earthly-specific
		bklog.G(ctx).Infof("GitCLI.Run called with %v", args)
	}

	for {
		var cmd *exec.Cmd
		if cli.exec == nil {
			cmd = exec.CommandContext(ctx, gitBinary)
		} else {
			cmd = exec.Command(gitBinary)
		}

		cmd.Dir = cli.dir
		if cmd.Dir == "" {
			cmd.Dir = cli.workTree
		}

		// Block sneaky repositories from using repos from the filesystem as submodules.
		cmd.Args = append(cmd.Args, "-c", "protocol.file.allow=user")
		if cli.workTree != "" {
			cmd.Args = append(cmd.Args, "--work-tree", cli.workTree)
		}
		if cli.gitDir != "" {
			cmd.Args = append(cmd.Args, "--git-dir", cli.gitDir)
		}
		cmd.Args = append(cmd.Args, cli.args...)
		cmd.Args = append(cmd.Args, args...)

		buf := bytes.NewBuffer(nil)
		errbuf := bytes.NewBuffer(nil)
		cmd.Stdin = nil
		cmd.Stdout = buf
		cmd.Stderr = errbuf
		if cli.streams != nil {
			stdout, stderr, flush := cli.streams(ctx)
			if stdout != nil {
				cmd.Stdout = io.MultiWriter(stdout, cmd.Stdout)
			}
			if stderr != nil {
				cmd.Stderr = io.MultiWriter(stderr, cmd.Stderr)
			}
			defer stdout.Close()
			defer stderr.Close()
			defer func() {
				if err != nil {
					flush()
				}
			}()
		}

		cmd.Env = []string{
			"PATH=" + os.Getenv("PATH"),

			// earthly-specific settings
			"HOME=" + os.Getenv("HOME"), // earthly-specific: we need this for git to read /root/.gitconfig
			"GIT_LFS_SKIP_SMUDGE=1",     // earthly-specific: dont automatically pull large files

			"GIT_TERMINAL_PROMPT=0",
			"GIT_SSH_COMMAND=" + getGitSSHCommand(cli.sshKnownHosts, logLevel),
			//	"GIT_TRACE=1",
			// earthly-specific: Commented out. We do not want to disable reading from gitconfig.
			//"GIT_CONFIG_NOSYSTEM=1", // Disable reading from system gitconfig.
			//"HOME=/dev/null",        // Disable reading from user gitconfig.

			"LC_ALL=C", // Ensure consistent output.
		}

		// earthly-specific
		if logLevel >= GitLogLevelTrace {
			cmd.Env = append(cmd.Env, "GIT_TRACE=1")
		}

		if cli.sshAuthSock != "" {
			cmd.Env = append(cmd.Env, "SSH_AUTH_SOCK="+cli.sshAuthSock)
		}

		if cli.exec != nil {
			// remote git commands spawn helper processes that inherit FDs and don't
			// handle parent death signal so exec.CommandContext can't be used
			err = cli.exec(ctx, cmd)
		} else {
			err = cmd.Run()
		}

		if err != nil {
			if strings.Contains(errbuf.String(), "--depth") || strings.Contains(errbuf.String(), "shallow") {
				if newArgs := argsNoDepth(args); len(args) > len(newArgs) {
					args = newArgs
					continue
				}
			}

			// Earthly-TODO: moby/master didn't use to include the git stderr; however, change it at some point AFTER we added in this EARTHLY_GIT_STDERR hack; however we will keep ours for the time being (since we look specifically for EARTHLY_GIT_STDERR in earthly)
			//return buf.Bytes(), errors.Errorf("git error: %s\nstderr:\n%s", err, errbuf.String())

			// earthly-specific
			if gitDebug() {
				bklog.G(ctx).Infof("knownHosts: %s", cli.sshKnownHosts)
				bklog.G(ctx).Infof("git stdout: %s", buf.String())
				bklog.G(ctx).Infof("git stderr: %s", errbuf.String())
			}
			err = errors.Wrapf(err, "EARTHLY_GIT_STDERR: %s", base64.StdEncoding.EncodeToString([]byte(urlutil.RedactAllCredentials(fmt.Sprintf("git %s\n%s", strings.Join(args, " "), errbuf.String()))))) // earthly-specific
			return buf.Bytes(), err
		}

		return buf.Bytes(), nil
	}
}

func getGitSSHCommand(knownHosts string, logLevel GitLogLevel) string {
	gitSSHCommand := "ssh -F /dev/null"
	if knownHosts != "" {
		gitSSHCommand += " -o UserKnownHostsFile=" + knownHosts
	} else {
		gitSSHCommand += " -o StrictHostKeyChecking=no"
	}
	if gitDebug() || logLevel >= GitLogLevelTrace {
		gitSSHCommand += " -vvvv"
	} else if logLevel >= GitLogLevelDebug {
		gitSSHCommand += " -v"
	}
	return gitSSHCommand
}

func argsNoDepth(args []string) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if a != "--depth=1" {
			out = append(out, a)
		}
	}
	return out
}
