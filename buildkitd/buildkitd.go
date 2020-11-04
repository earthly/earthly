package buildkitd

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/earthly/earthly/conslogging"
	"github.com/moby/buildkit/client"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/pkg/errors"
)

const (
	// ContainerName is the name of the buildkitd container.
	ContainerName = "earthly-buildkitd"
	// VolumeName is the name of the docker volume used for storing the cache.
	VolumeName = "earthly-cache"
)

// Address is the address at which the daemon is available.
var Address = fmt.Sprintf("docker-container://%s", ContainerName)

// TODO: Implement all this properly with the docker client.

// NewClient returns a new buildkitd client.
func NewClient(ctx context.Context, console conslogging.ConsoleLogger, image string, settings Settings, opTimeout time.Duration, opts ...client.ClientOpt) (*client.Client, error) {
	address, err := MaybeStart(ctx, console, image, settings, opTimeout)
	if err != nil {
		console.WithPrefix("buildkitd").Printf("Is docker installed and running? Are you part of the docker group?\n")
		return nil, errors.Wrap(err, "maybe start buildkitd")
	}
	bkClient, err := client.New(ctx, address, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "new buildkit client")
	}
	return bkClient, nil
}

// ResetCache restarts the buildkitd daemon with the reset command.
func ResetCache(ctx context.Context, console conslogging.ConsoleLogger, image string, settings Settings, opTimeout time.Duration) error {
	console.
		WithPrefix("buildkitd").
		Printf("Restarting buildkit daemon with reset command...\n")
	isStarted, err := IsStarted(ctx)
	if err != nil {
		return errors.Wrap(err, "check is started buildkitd")
	}
	if isStarted {
		err = Stop(ctx)
		if err != nil {
			return err
		}
		err = WaitUntilStopped(ctx, opTimeout)
		if err != nil {
			return err
		}
	}
	err = Start(ctx, image, settings, true)
	if err != nil {
		return err
	}
	err = WaitUntilStarted(ctx, Address, opTimeout)
	if err != nil {
		return err
	}
	console.
		WithPrefix("buildkitd").
		Printf("...Done\n")
	return nil
}

// MaybeStart ensures that the buildkitd daemon is started. It returns the URL
// that can be used to connect to it.
func MaybeStart(ctx context.Context, console conslogging.ConsoleLogger, image string, settings Settings, opTimeout time.Duration) (string, error) {
	isStarted, err := IsStarted(ctx)
	if err != nil {
		return "", errors.Wrap(err, "check is started buildkitd")
	}
	if isStarted {
		console.
			WithPrefix("buildkitd").
			Printf("Found buildkit daemon as docker container (%s)\n", ContainerName)
		err := MaybeRestart(ctx, console, image, settings, opTimeout)
		if err != nil {
			return "", errors.Wrap(err, "maybe restart")
		}
	} else {
		console.
			WithPrefix("buildkitd").
			Printf("Starting buildkit daemon as a docker container (%s)...\n", ContainerName)
		err := Start(ctx, image, settings, false)
		if err != nil {
			return "", errors.Wrap(err, "start")
		}
		err = WaitUntilStarted(ctx, Address, opTimeout)
		if err != nil {
			return "", errors.Wrap(err, "wait until started")
		}
		console.
			WithPrefix("buildkitd").
			Printf("...Done\n")
	}
	return Address, nil
}

// MaybeRestart checks whether the there is a different buildkitd image available locally or if
// settings of the current container are different from the provided settings. In either case,
// the container is restarted.
func MaybeRestart(ctx context.Context, console conslogging.ConsoleLogger, image string, settings Settings, opTimeout time.Duration) error {
	containerImageID, err := GetContainerImageID(ctx)
	if err != nil {
		return err
	}
	availableImageID, err := GetAvailableImageID(ctx, image)
	if err != nil {
		// Could not get available image ID. This happens when a new image tag is given and that
		// tag has not yet been pulled locally. Restarting will cause that tag to be pulled.
		availableImageID = "" // Will cause equality to fail and force a restart.
		// Keep going anyway.
	}
	if containerImageID == availableImageID {
		// Images are the same. Check settings hash.
		hash, err := GetSettingsHash(ctx)
		if err != nil {
			return err
		}
		ok, err := settings.VerifyHash(hash)
		if err != nil {
			return errors.Wrap(err, "verify hash")
		}
		if ok {
			// No need to replace: images are the same and settings are the same.
			return nil
		}

		console.
			WithPrefix("buildkitd").
			Printf("Settings do not match. Restarting buildkit daemon with updated settings...\n")
	} else {
		console.
			WithPrefix("buildkitd").
			Printf("Newer image available. Restarting buildkit daemon...\n")
	}

	// Replace.
	err = Stop(ctx)
	if err != nil {
		return err
	}
	err = WaitUntilStopped(ctx, opTimeout)
	if err != nil {
		return err
	}
	err = Start(ctx, image, settings, false)
	if err != nil {
		return err
	}
	err = WaitUntilStarted(ctx, Address, opTimeout)
	if err != nil {
		return err
	}
	console.
		WithPrefix("buildkitd").
		Printf("...Done\n")
	return nil
}

// RemoveExited removes any stopped or exited buildkitd containers
func RemoveExited(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "docker", "ps", "-a", "-q", "-f", fmt.Sprintf("name=%s", ContainerName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "get combined output")
	}
	if len(output) == 0 {
		return nil
	}
	return exec.CommandContext(ctx, "docker", "rm", ContainerName).Run()
}

// Start starts the buildkitd daemon.
func Start(ctx context.Context, image string, settings Settings, reset bool) error {
	settingsHash, err := settings.Hash()
	if err != nil {
		return errors.Wrap(err, "settings hash")
	}
	err = RemoveExited(ctx)
	if err != nil {
		return err
	}
	env := os.Environ()
	runMount := fmt.Sprintf("%s:/run/earthly:consistent", settings.RunDir)
	args := []string{
		"run",
		"-d",
		"-v", fmt.Sprintf("%s:/tmp/earthly:rw", VolumeName),
		"-v", runMount,
		"-e", fmt.Sprintf("BUILDKIT_DEBUG=%t", settings.Debug),
		"--label", fmt.Sprintf("dev.earthly.settingshash=%s", settingsHash),
		"--name", ContainerName,
		"--privileged",
	}
	if os.Getenv("EARTHLY_WITH_DOCKER") == "1" {
		// Add /sys/fs/cgroup if it's earthly-in-earthly.
		args = append(args, "-v", "/sys/fs/cgroup:/sys/fs/cgroup")
	} else {
		// Debugger only supported in top-most earthly.
		// TODO: Main reason for this is port clash. This could be improved in the future,
		//       if needed.
		args = append(args,
			"-p", fmt.Sprintf("127.0.0.1:%d:8373", settings.DebuggerPort))
	}
	// Apply some git-related settings.
	if settings.SSHAuthSock != "" {
		args = append(args,
			"-v", fmt.Sprintf("%s:/ssh-agent.sock", settings.SSHAuthSock),
			"-e", "SSH_AUTH_SOCK=/ssh-agent.sock",
		)
	}

	args = append(args,
		"-e", fmt.Sprintf("CACHE_SIZE_MB=%d", settings.CacheSizeMb),
		"-e", "EARTHLY_GIT_CONFIG",
		"-e", fmt.Sprintf("GIT_URL_INSTEAD_OF=%s", settings.GitURLInsteadOf),
	)
	env = append(env,
		fmt.Sprintf("EARTHLY_GIT_CONFIG=%s", base64.StdEncoding.EncodeToString([]byte(settings.GitConfig))),
	)

	for i, data := range settings.GitCredentials {
		args = append(args,
			"-e", fmt.Sprintf("GIT_CREDENTIALS_%d", i),
		)
		env = append(env,
			fmt.Sprintf("GIT_CREDENTIALS_%d=%s", i, base64.StdEncoding.EncodeToString([]byte("password="+data))),
		)
	}

	// Apply reset.
	if reset {
		args = append(args, "-e", "EARTHLY_RESET_TMP_DIR=true")
	}
	// Execute.
	args = append(args, image)
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "docker run %s: %s", image, string(output))
	}
	return nil
}

// Stop stops the buildkitd container.
func Stop(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "docker", "stop", ContainerName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "get combined output")
	}
	return nil
}

// IsStarted checks if the buildkitd container has been started.
func IsStarted(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx, "docker", "ps", "-q", "-f", fmt.Sprintf("name=%s", ContainerName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.Wrap(err, "get combined output")
	}
	return (len(output) != 0), nil
}

// WaitUntilStarted waits until the buildkitd daemon has started and is healthy.
func WaitUntilStarted(ctx context.Context, address string, opTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	for {
		select {
		case <-time.After(1 * time.Second):
			bkClient, err := client.New(ctxTimeout, address)
			if err != nil {
				// Try again.
				continue
			}
			_, err = bkClient.ListWorkers(ctxTimeout)
			if err != nil {
				// Try again.
				continue
			}
			err = bkClient.Close()
			if err != nil {
				return errors.Wrap(err, "close buildkit client")
			}
			return nil
		case <-ctxTimeout.Done():
			return errors.New("Timeout: Buildkitd did not start")
		}
	}
}

// GetContainerIP returns the IP of the buildkit container.
func GetContainerIP(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "inspect", "-f", "{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}", ContainerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "get combined output ip")
	}
	return string(bytes.TrimSpace(output)), nil
}

// WaitUntilStopped waits until the buildkitd daemon has stopped.
func WaitUntilStopped(ctx context.Context, opTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	for {
		select {
		case <-time.After(1 * time.Second):
			cmd := exec.CommandContext(
				ctx, "docker", "inspect", "--format={{.State.Running}}", ContainerName)
			output, err := cmd.CombinedOutput()
			if err != nil {
				// The container can no longer be found at all.
				return nil
			}
			isRunning, err := strconv.ParseBool(strings.TrimSpace(string(output)))
			if err != nil {
				return errors.Wrapf(err, "cannot interpret output %s", output)
			}
			if !isRunning {
				return nil
			}
		case <-ctxTimeout.Done():
			return errors.New("Timeout: Buildkitd did not start")
		}
	}
}

// GetSettingsHash fetches the hash of the currently running buildkitd container.
func GetSettingsHash(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx,
		"docker", "inspect",
		"--format={{index .Config.Labels \"dev.earthly.settingshash\"}}",
		ContainerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "get output for settings hash")
	}
	return string(output), nil
}

// GetContainerImageID fetches the ID of the image used for the running buildkitd container.
func GetContainerImageID(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx,
		"docker", "inspect", "--format={{index .Image}}", ContainerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "get output for container image ID")
	}
	return string(output), nil
}

// GetAvailableImageID fetches the ID of the image buildkitd image available.
func GetAvailableImageID(ctx context.Context, image string) (string, error) {
	cmd := exec.CommandContext(ctx,
		"docker", "inspect", "--format={{index .Id}}", image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "get output for available image ID")
	}
	return string(output), nil
}
