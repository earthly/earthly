package analytics

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/earthly/earthly/fileutils"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/pkg/errors"
	"gopkg.in/segmentio/analytics-go.v3"
)

func detectCI() (string, bool) {
	for k, v := range map[string]string{
		"GITHUB_WORKFLOW": "github-actions",
		"CIRCLECI":        "circle-ci",
		"EARTHLY":         "earthly",
		"JENKINS_HOME":    "jenkins",
		"BUILDKITE":       "buildkite",
		"DRONE_BRANCH":    "drone",
		"TRAVIS":          "travis",
	} {
		if _, ok := os.LookupEnv(k); ok {
			return v, true
		}
	}

	// default catch-all
	if v, ok := os.LookupEnv("CI"); ok {
		if strings.ToLower(v) == "true" {
			return "unknown", true
		}
		return v, true
	}

	return "false", false
}

func getRepo() string {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err == nil {
		return string(out)
	}

	for _, k := range []string{
		"GITHUB_REPOSITORY",
		"CIRCLE_PROJECT_REPONAME",
		"GIT_URL",
		"BUILDKITE_REPO",
		"DRONE_REPO",
		"TRAVIS_REPO_SLUG",
		"EARTHLY_GIT_ORIGIN_URL",
	} {
		if v, ok := os.LookupEnv(k); ok {
			return v
		}
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			if strings.Contains(pair[1], "git") {
				return pair[1]
			}
		}
	}

	return "unknown"
}

func getRepoHash() string {
	repo := getRepo()
	if repo == "unknown" || repo == "" {
		return repo
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(repo)))
}

func getInstallID() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user home dir")
	}

	path := filepath.Join(homeDir, ".earthly", "install_id")
	if !fileutils.FileExists(path) {
		u, err := uuid.NewV4()
		if err != nil {
			return "", errors.Wrap(err, "failed to generate uuid")
		}
		ID := u.String()
		err = ioutil.WriteFile(path, []byte(ID), 0644)
		if err != nil {
			return "", errors.Wrapf(err, "failed to write %q", path)
		}
		return ID, nil
	}

	s, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read %q", path)
	}
	return string(s), nil
}

// CollectAnalytics sends analytics to segment.io
func CollectAnalytics(version, gitSha, commandName string, exitCode int, realtime time.Duration) {
	var err error
	ciName, ci := detectCI()
	installID, overrideInstallID := os.LookupEnv("EARTHLY_INSTALL_ID")
	repoHash := getRepoHash()
	if !overrideInstallID {
		if repoHash == "unknown" {
			installID = "unknown"
		} else {
			if ci {
				installID = fmt.Sprintf("%x", sha256.Sum256([]byte(ciName+repoHash)))
			} else {
				installID, err = getInstallID()
				if err != nil {
					installID = "unknown"
				}
			}
		}
	}
	segmentClient := analytics.New("RtwJaMBswcW3CNMZ7Ops79dV6lEZqsXf")
	segmentClient.Enqueue(analytics.Track{
		Event:  "cli-" + commandName,
		UserId: installID,
		Properties: analytics.NewProperties().
			Set("version", version).
			Set("gitsha", gitSha).
			Set("exitcode", exitCode).
			Set("ci", ciName).
			Set("repohash", repoHash).
			Set("realtime", realtime.Seconds()),
	})
	done := make(chan bool, 1)
	go func() {
		segmentClient.Close()
		done <- true
	}()
	select {
	case <-time.After(time.Millisecond * 500):
	case <-done:
	}
}
