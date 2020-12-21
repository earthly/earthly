package analytics

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/fileutils"
	"github.com/earthly/earthly/syncutil"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/pkg/errors"
	"gopkg.in/segmentio/analytics-go.v3"
)

func detectCI() (string, bool) {
	for k, v := range map[string]string{
		"GITHUB_WORKFLOW": "github-actions",
		"CIRCLECI":        "circle-ci",
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
	consistentRepo, err := buildcontext.ParseGitRemoteURL(repo)
	if err == nil {
		repo = consistentRepo
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

// EarthlyAnalytics contains analytical data which is sent to api.earthly.dev
type EarthlyAnalytics struct {
	Key              string  `json:"key"`
	InstallID        string  `json:"install_id"`
	Version          string  `json:"version"`
	GitSHA           string  `json:"git_sha"`
	ExitCode         int     `json:"exit_code"`
	CI               string  `json:"ci_name"`
	RepoHash         string  `json:"repo_hash"`
	ExecutionSeconds float64 `json:"execution_seconds"`
}

func saveData(server string, data *EarthlyAnalytics) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}

	client := &http.Client{
		Timeout: time.Millisecond * 500,
	}

	// set the HTTP method, url, and request body
	url := server + "/analytics"
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrap(err, "failed to create request for sending analytics")
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send analytics")
	}

	return nil
}

// CollectAnalytics sends analytics to segment.io
func CollectAnalytics(ctx context.Context, earthlyServer string, displayErrors bool, version, gitSha, commandName string, exitCode int, realtime time.Duration) {
	var err error
	ciName, ci := detectCI()
	installID, overrideInstallID := os.LookupEnv("EARTHLY_INSTALL_ID")
	repoHash := getRepoHash()
	if !overrideInstallID {
		if ciName != "false" {
			installID = "ci"
		} else if repoHash == "unknown" {
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

	key := "cli-" + commandName

	var wg sync.WaitGroup

	// send data to api.earthly.dev
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := saveData(earthlyServer, &EarthlyAnalytics{
			Key:              key,
			InstallID:        installID,
			Version:          version,
			GitSHA:           gitSha,
			ExitCode:         exitCode,
			CI:               ciName,
			RepoHash:         repoHash,
			ExecutionSeconds: realtime.Seconds(),
		})
		if err != nil && displayErrors {
			fmt.Fprintf(os.Stderr, "error while sending analytics to earthly: %s\n", err.Error())
		}
	}()

	// send data to segment
	wg.Add(1)
	go func() {
		defer wg.Done()
		segmentClient := analytics.New("RtwJaMBswcW3CNMZ7Ops79dV6lEZqsXf")
		segmentClient.Enqueue(analytics.Track{
			Event:  key,
			UserId: installID,
			Properties: analytics.NewProperties().
				Set("version", version).
				Set("gitsha", gitSha).
				Set("exitcode", exitCode).
				Set("ci", ciName).
				Set("repohash", repoHash).
				Set("realtime", realtime.Seconds()),
		})
		err := segmentClient.Close()
		if err != nil && displayErrors {
			fmt.Fprintf(os.Stderr, "error while sending analytics to segment: %s\n", err.Error())
		}
	}()

	ok := syncutil.WaitContext(ctx, &wg)
	if !ok && displayErrors {
		fmt.Fprintf(os.Stderr, "Warning: timedout while sending analytics\n")
	}
}
