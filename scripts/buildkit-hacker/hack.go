package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

type opts struct {
	Verbose       bool `long:"verbose" short:"v" description:"Enable verbose logging"`
	Reset         bool `long:"reset" description:"Reset repos"`
	SkipVendor    bool `long:"skip-vendor" description:"skip vendoring (and any fsutils work)"`
	SkipGenerated bool `long:"skip-generated" description:"skip auto-generating files"`
}

func mustGetGitSha(gitPath string) string {
	sha := mustRunOutput(`set -e
cd $gitPath
git rev-parse HEAD`, "gitPath="+gitPath)
	return strings.TrimSpace(sha)
}
func mustGetGitCommitDate(gitPath string) time.Time {
	commitdate := mustRunOutput(`set -e
cd $gitPath
git show -s --format=%ci`, "gitPath="+gitPath)
	commitdate = strings.TrimSpace(commitdate)

	d, err := time.Parse("2006-01-02 15:04:05 -0700", commitdate)
	if err != nil {
		panic(err)
	}
	return d
}

func mustRun(bourneScript string, envs ...string) {
	cmd := exec.Command("sh", "-c", bourneScript)
	cmd.Env = append(os.Environ(), envs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func mustRunOutput(bourneScript string, envs ...string) string {
	cmd := exec.Command("sh", "-c", bourneScript)
	cmd.Env = append(os.Environ(), envs...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return string(stdout.Bytes())
}

func resetRepos(githubuser, hackbranch string) {
	mustRun(`set -e
if [ -d "$HOME/buildkit-hacker" ]; then
  echo "delete $HOME/buildkit-hacker before continuing"
  exit 1
fi
`)

	mustRun(`set -e
mkdir -p $HOME/buildkit-hacker/fsutil
cd $HOME/buildkit-hacker/fsutil

git init .
git config commit.gpgsign false
git remote add earthly git@github.com:earthly/fsutil.git
git remote add tonistiigi git@github.com:tonistiigi/fsutil.git
git remote add "$githubuser" git@github.com:$githubuser/fsutil.git
git fetch --all
git branch --track tonistiigi-master tonistiigi/master
git branch --track earthly-main earthly/earthly-main

git push --force "$githubuser" "earthly-main:$hackbranch"

git branch --track hack $githubuser/$hackbranch
git checkout hack
`, "githubuser="+githubuser, "hackbranch="+hackbranch)

	mustRun(`set -e
mkdir -p $HOME/buildkit-hacker/buildkit
cd $HOME/buildkit-hacker/buildkit

set -e
git init .
git config commit.gpgsign false
git remote add earthly git@github.com:earthly/buildkit.git
git remote add moby git@github.com:moby/buildkit.git
git remote add "$githubuser" git@github.com:$githubuser/buildkit.git
git fetch --all
git branch --track moby-master moby/master
git branch --track earthly-main earthly/earthly-main

git push --force "$githubuser" "earthly-main:$hackbranch"

git branch --track hack $githubuser/$hackbranch
git checkout hack
`, "githubuser="+githubuser, "hackbranch="+hackbranch)

}

func getGoModReplaceStr(githubuser, repo string) string {
	sha1 := mustGetGitSha(os.Getenv("HOME") + "/buildkit-hacker/" + repo)
	commitDate := mustGetGitCommitDate(os.Getenv("HOME") + "/buildkit-hacker/" + repo)
	timestamp := commitDate.UTC().Format("20060102150405")

	repoPath := "github.com/" + githubuser + "/" + repo
	version := "v0.0.0"
	s := repoPath + "@" + version + "-" + timestamp + "-" + sha1[:12]
	return s
}

func build(githubuser, hackbranch string, skipVendoring, skipGeneratedFiles bool) {
	if hackbranch == "main" || hackbranch == "master" || hackbranch == "earthly-main" {
		panic("hackbranch cant be main")
	}

	mustRun(`set -e
cd "$HOME/buildkit-hacker/buildkit"
test "$(git rev-parse --abbrev-ref HEAD)" = "$hackbranch" || (echo expected $hackbranch && exit 1)

cd "$HOME/buildkit-hacker/fsutil"
test "$(git rev-parse --abbrev-ref HEAD)" = "$hackbranch" || (echo expected $hackbranch && exit 1)
`, "hackbranch="+hackbranch)

	var fsutilreplace string
	if !skipVendoring {
		mustRun(`set -e
cd "$HOME/buildkit-hacker/fsutil"
git commit -am wip || true
git push
`)

		fsutilreplace = getGoModReplaceStr(githubuser, "fsutil") // NOTE this must be done after the git commit
		fmt.Printf("fsutil is %s\n", fsutilreplace)

		mustRun(`set -e
cd "$HOME/buildkit-hacker/buildkit"
go mod edit -replace "github.com/tonistiigi/fsutil=$fsutilreplace"
make vendor
`, "fsutilreplace="+fsutilreplace)
	} else {
		fsutilreplace = getGoModReplaceStr(githubuser, "fsutil")
		fmt.Printf("fsutil is %s\n", fsutilreplace)
	}

	if !skipGeneratedFiles {
		mustRun(`set -e
cd "$HOME/buildkit-hacker/buildkit"
make generated-files
`)
	}
	mustRun(`set -e
cd "$HOME/buildkit-hacker/buildkit"
make binaries
git commit -am wip || true
git push
`)

	buildkitreplace := getGoModReplaceStr(githubuser, "buildkit")
	fmt.Printf("buildkit is %s\n", buildkitreplace)

	buildkitgitSha := mustGetGitSha(os.Getenv("HOME") + "/buildkit-hacker/buildkit")
	mustUpdateBuildkitEarthfile(githubuser, buildkitgitSha)

	mustRun(`set -e
go mod edit -replace "github.com/moby/buildkit=$buildkitreplace"
go mod edit -replace "github.com/tonistiigi/fsutil=$fsutilreplace"
go mod tidy
echo building earthly
./earthly --no-sat -VD +for-linux
`, "buildkitreplace="+buildkitreplace, "fsutilreplace="+fsutilreplace)
}

func mustUpdateBuildkitEarthfile(githubuser, buildkitsha string) {
	data, err := os.ReadFile("buildkitd/Earthfile")
	if err != nil {
		panic(err)
	}
	newLines := []string{}
	for _, l := range strings.Split(string(data), "\n") {
		s := strings.TrimSpace(l)
		magic := "ARG BUILDKIT_BASE_IMAGE=github.com"
		if strings.HasPrefix(s, magic) {
			splits := strings.Split(l, "=")
			if len(splits) != 2 {
				panic("I hurt myself doing the splits")
			}
			l = splits[0] + "=github.com/" + githubuser + "/buildkit:" + buildkitsha + "+build"
		}
		newLines = append(newLines, l)
	}
	newEarthfile := strings.TrimSpace(strings.Join(newLines, "\n")) + "\n"
	err = os.WriteFile("buildkitd/Earthfile", []byte(newEarthfile), 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	programName := "buildkit-hacker"
	if len(os.Args) > 0 {
		programName = path.Base(os.Args[0])
	}

	// TODO move to config
	githubuser := "alexcb"
	hackbranch := "hack"

	progOpts := opts{}
	p := flags.NewNamedParser("", flags.PrintErrors|flags.PassDoubleDash|flags.PassAfterNonOption|flags.HelpFlag)
	_, err := p.AddGroup(fmt.Sprintf("%s [options] args", programName), "", &progOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}
	args, err := p.ParseArgs(os.Args[1:])
	if err != nil {
		p.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	if len(args) != 0 {
		p.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	if progOpts.Reset {
		resetRepos(githubuser, hackbranch)
		os.Exit(0)
	}
	build(githubuser, hackbranch, progOpts.SkipVendor, progOpts.SkipGenerated)
}
