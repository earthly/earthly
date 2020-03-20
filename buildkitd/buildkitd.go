package buildkitd

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/moby/buildkit/client"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/pkg/errors"
	"github.com/vladaionescu/earthly/buildkitd/credpass"
	"github.com/vladaionescu/earthly/conslogging"
	"github.com/vladaionescu/earthly/logging"
)

const (
	// ContainerName is the name of the buildkitd container.
	ContainerName = "earthly-buildkitd"
	// TempDir is the directory used for buildkitd cache.
	TempDir = "/tmp/earthly"
)

// Address is the address at which the daemon is available.
var Address = fmt.Sprintf("docker-container://%s", ContainerName)

// TODO: Implement all this properly with the docker client.

// MaybeStart ensures that the buildkitd daemon is started. It returns the URL
// that can be used to connect to it.
func MaybeStart(ctx context.Context, console conslogging.ConsoleLogger, image string, settings Settings) (string, error) {
	isStarted, err := IsStarted(ctx)
	if err != nil {
		return "", errors.Wrap(err, "check is started buildkitd")
	}
	if isStarted {
		console.
			WithPrefix("buildkitd").
			Printf("Found buildkit daemon as docker container (%s)\n", ContainerName)
		err := MaybeRestart(ctx, console, image, settings)
		if err != nil {
			return "", errors.Wrap(err, "maybe restart")
		}
	} else {
		console.
			WithPrefix("buildkitd").
			Printf("Starting buildkit daemon as a docker container (%s)...\n", ContainerName)
		err := Start(ctx, image, settings)
		if err != nil {
			return "", errors.Wrap(err, "start")
		}
		err = WaitUntilStarted(ctx, Address)
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
func MaybeRestart(ctx context.Context, console conslogging.ConsoleLogger, image string, settings Settings) error {
	containerImageID, err := GetContainerImageID(ctx)
	if err != nil {
		return err
	}
	availableImageID, err := GetAvailableImageID(ctx, image)
	if err != nil {
		return err
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
	err = WaitUntilStopped(ctx)
	if err != nil {
		return err
	}
	err = Start(ctx, image, settings)
	if err != nil {
		return err
	}
	err = WaitUntilStarted(ctx, Address)
	if err != nil {
		return err
	}
	console.
		WithPrefix("buildkitd").
		Printf("...Done\n")
	return nil
}

// Start starts the buildkitd daemon.
func Start(ctx context.Context, image string, settings Settings) error {
	settingsHash, err := settings.Hash()
	if err != nil {
		return errors.Wrap(err, "settings hash")
	}
	cacheMount := fmt.Sprintf("%s:%s:delegated", TempDir, TempDir)
	args := []string{
		"run",
		"-d", "--rm",
		"-v", cacheMount,
		"--label", fmt.Sprintf("dev.earthly.settingshash=%s", settingsHash),
		"--name", ContainerName,
		"--privileged",
	}
	// Apply some buildkitd-related settings.
	if settings.CacheSizeMb > 0 {
		args = append(args,
			"-e", fmt.Sprintf("CACHE_SIZE_MB=%d", settings.CacheSizeMb),
		)
	}
	// Apply some git-related settings.
	if settings.SSHAuthSock != "" {
		args = append(args,
			"-v", fmt.Sprintf("%s:/ssh-agent", settings.SSHAuthSock),
			"-e", "SSH_AUTH_SOCK=/ssh-agent",
		)
	}
	if len(settings.GitSettings) > 0 {
		// TODO: Only the first GitSettings entry is used, and it is not bound
		//       to the specified domain.
		args = append(args,
			"-e", fmt.Sprintf("GIT_USERNAME=%s", settings.GitSettings[0].Username),
			"-e", fmt.Sprintf("GIT_PASSWORD=%s", settings.GitSettings[0].Password),
		)
	}
	if settings.GitURLInsteadOf != "" {
		args = append(args,
			"-e", fmt.Sprintf("GIT_URL_INSTEAD_OF=%s", settings.GitURLInsteadOf),
		)
	}
	// Pass docker credentials.
	dockerConfig, dockerConfigDir, err := credpass.Read(ctx)
	if err != nil {
		warningText := fmt.Sprintf("Unable to pass docker credentials to buildkitd: %s", err.Error())
		logging.GetLogger(ctx).Warning(warningText)
		// TODO: Print console warning properly.
		fmt.Printf("Warning: %s\nYour builds will not be able to pull from private logged-in registries!\n", warningText)
		// Keep going anyway.
	} else {
		if dockerConfig != "" {
			// Pass docker config as env var.
			args = append(args,
				"-e", fmt.Sprintf("DOCKER_CONFIG_JSON=%s", dockerConfig),
			)
		} else if dockerConfigDir != "" {
			// Pass docker config as mount.
			args = append(args,
				"-v", fmt.Sprintf("%s:/root/.docker:ro", dockerConfigDir),
			)
		}
	}
	// Execute.
	args = append(args, image)
	cmd := exec.CommandContext(ctx, "docker", args...)
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
func WaitUntilStarted(ctx context.Context, address string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
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

// WaitUntilStopped waits until the buildkitd daemon has stopped.
func WaitUntilStopped(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	for {
		select {
		case <-time.After(1 * time.Second):
			cmd := exec.CommandContext(
				ctx, "docker", "inspect", "--format={{.State.Running}}", ContainerName)
			_, err := cmd.CombinedOutput()
			if err != nil {
				// The container has stopped successfully when this command returns an error
				// (container can no longer be found).
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

// NewClient returns a new buildkitd client.
func NewClient(ctx context.Context, console conslogging.ConsoleLogger, image string, settings Settings, opts ...client.ClientOpt) (*client.Client, error) {
	address, err := MaybeStart(ctx, console, image, settings)
	if err != nil {
		return nil, errors.Wrap(err, "maybe start buildkitd")
	}
	bkClient, err := client.New(ctx, address, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "new buildkit client")
	}
	return bkClient, nil
}
