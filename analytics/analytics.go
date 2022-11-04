package analytics

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/syncutil"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// DetectCI determines if Earthly is being run from a CI environment. It returns
// the name of the CI tool and true if we detect one.
func DetectCI() (string, bool) {
	for k, v := range map[string]string{
		"GITHUB_WORKFLOW": "github-actions",
		"CIRCLECI":        "circle-ci",
		"JENKINS_HOME":    "jenkins",
		"BUILDKITE":       "buildkite",
		"DRONE_BRANCH":    "drone",
		"TRAVIS":          "travis",
		"GITLAB_CI":       "gitlab",
		"EARTHLY_IMAGE":   "earthly-image",
		"AGENT_WORKDIR":   "jenkins", // https://github.com/jenkinsci/docker-agent/blob/master/11/alpine/Dockerfile#L35
	} {
		if _, ok := os.LookupEnv(k); ok {
			return v, true
		}
	}

	// another way to tell if it's Jenkins
	_, err := os.Stat("/home/jenkins")
	if err == nil {
		// /home/jenkins exists.
		return "jenkins", true
	}

	// default catch-all
	if v, ok := os.LookupEnv("CI"); ok {
		isCI, err := strconv.ParseBool(v)
		if err == nil && isCI {
			return "ci-env-var-set", true
		}
	}

	return "false", false
}

func getRepo() string {
	if isGitInstalled() {
		if !isGitDir() {
			return ""
		}

		cmd := exec.Command("git", "config", "--get", "remote.origin.url")
		out, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
	}

	for _, k := range []string{
		"GITHUB_REPOSITORY",
		"CIRCLE_PROJECT_REPONAME",
		"GIT_URL",
		"BUILDKITE_REPO",
		"DRONE_REPO",
		"TRAVIS_REPO_SLUG",
		"EARTHLY_GIT_ORIGIN_URL",
		"CI_REPOSITORY_URL",
	} {
		if v, ok := os.LookupEnv(k); ok {
			return strings.TrimSpace(v)
		}
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			if strings.Contains(pair[1], "git") {
				return strings.TrimSpace(pair[1])
			}
		}
	}

	return "unknown"
}

func isGitInstalled() bool {
	cmd := exec.Command("git", "--version")
	err := cmd.Run()
	return (err == nil)
}

func isGitDir() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return (err == nil)
}

func getRepoHash() string {
	return RepoHashFromCloneURL(getRepo())
}

// RepoHashFromCloneURL returns the repoHash of a ref
func RepoHashFromCloneURL(repo string) string {
	if repo == "unknown" || repo == "" {
		return repo
	}
	consistentRepo, err := gitutil.ParseGitRemoteURL(repo)
	if err == nil {
		repo = consistentRepo
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(repo)))
}

func getInstallID() (string, error) {
	earthlyDir, err := cliutil.GetOrCreateEarthlyDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(earthlyDir, "install_id")
	pathExists, err := fileutil.FileExists(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to check if %s exists", path)
	}
	if !pathExists {
		u, err := uuid.NewUUID()
		if err != nil {
			u, err = uuid.NewRandom()
			if err != nil {
				return "", errors.Wrap(err, "failed to generate uuid")
			}
		}
		ID := u.String()
		err = os.WriteFile(path, []byte(ID), 0644)
		if err != nil {
			return "", errors.Wrapf(err, "failed to write %q", path)
		}
		return ID, nil
	}

	s, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read %q", path)
	}
	return string(s), nil
}

func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// Count increases the global count of (subsystem, key) which then gets reported when CollectAnalytics is called.
func Count(subsystem, key string) {
	counts.Count(subsystem, key)
}

// Meta holds metadata about the current run of the program.
type Meta struct {
	Version          string
	Platform         string
	BuildkitPlatform string
	UserPlatform     string
	GitSHA           string
	CommandName      string
	ExitCode         int
	IsSatellite      bool
	SatelliteVersion string
	IsRemoteBuildkit bool
	Realtime         time.Duration
}

// CollectAnalytics sends analytics to api.earthly.dev
func CollectAnalytics(ctx context.Context, cloudClient cloud.Client, displayErrors bool, meta Meta) {
	var err error
	ciName, ci := DetectCI()
	repoHash := getRepoHash()
	installID, overrideInstallID := os.LookupEnv("EARTHLY_INSTALL_ID")
	if !overrideInstallID {
		if ci {
			if repoHash == "unknown" {
				installID = "unknown"
			} else {
				installID = fmt.Sprintf("%x", sha256.Sum256([]byte(ciName+repoHash)))
			}
		} else {
			installID, err = getInstallID()
			if err != nil {
				if displayErrors {
					fmt.Fprintf(os.Stderr, "Failed to get install ID: %s\n", err.Error())
				}
				installID = "unknown"
			}
		}
	}

	key := "cli-" + meta.CommandName

	var wg sync.WaitGroup

	// send data to api.earthly.dev
	wg.Add(1)
	go func() {
		defer wg.Done()

		countsMap, countsMapUnlock := counts.getMap()
		defer countsMapUnlock()

		err := cloudClient.SendAnalytics(ctx, &cloud.EarthlyAnalytics{
			Key:              key,
			InstallID:        installID,
			Version:          meta.Version,
			Platform:         meta.Platform,
			BuildkitPlatform: meta.BuildkitPlatform,
			UserPlatform:     meta.UserPlatform,
			GitSHA:           meta.GitSHA,
			ExitCode:         meta.ExitCode,
			CI:               ciName,
			IsSatellite:      meta.IsSatellite,
			SatelliteVersion: meta.SatelliteVersion,
			IsRemoteBuildkit: meta.IsRemoteBuildkit,
			RepoHash:         repoHash,
			ExecutionSeconds: meta.Realtime.Seconds(),
			Terminal:         isTerminal(),
			Counts:           countsMap,
		})
		if err != nil && displayErrors {
			fmt.Fprintf(os.Stderr, "error while sending analytics to earthly: %s\n", err.Error())
		}
	}()

	ok := syncutil.WaitContext(ctx, &wg)
	if !ok && displayErrors {
		fmt.Fprintf(os.Stderr, "Warning: timed out while sending analytics\n")
	}
}
