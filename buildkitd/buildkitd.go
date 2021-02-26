package buildkitd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

var (
	// ErrBuildkitCrashed is an error returned when buildkit has terminated unexpectedly.
	ErrBuildkitCrashed = errors.New("buildkitd crashed")

	// ErrBuildkitStartFailure is an error returned when buildkit has failed to start in time.
	ErrBuildkitStartFailure = errors.New("buildkitd failed to start (in time)")
)

// Address is the address at which the daemon is available.
var Address = fmt.Sprintf("docker-container://%s", ContainerName)

// TODO: Implement all this properly with the docker client.

// NewClient returns a new buildkitd client.
func NewClient(ctx context.Context, console conslogging.ConsoleLogger, image string, settings Settings, opTimeout time.Duration, opts ...client.ClientOpt) (*client.Client, error) {
	if !isDockerAvailable(ctx) {
		console.WithPrefix("buildkitd").Printf("Is docker installed and running? Are you part of the docker group?\n")
		return nil, errors.New("docker not available")
	}
	address, err := MaybeStart(ctx, console, image, settings, opTimeout)
	if err != nil {
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
	err := CheckCompatibility(ctx, settings)
	if len(settings.AdditionalArgs) == 0 && err != nil {
		return errors.Wrap(err, "compatibility")
	}

	settingsHash, err := settings.Hash()
	if err != nil {
		return errors.Wrap(err, "settings hash")
	}
	err = RemoveExited(ctx)
	if err != nil {
		return err
	}
	// Pulling is not strictly needed, but it helps display some progress status to the user in
	// case the image is not available locally.
	err = MaybePull(ctx, image)
	if err != nil {
		fmt.Printf("Error: %s. Attempting to start buildkitd anyway...\n", err.Error())
		// Keep going - it might still work.
	}
	env := os.Environ()
	args := []string{
		"run",
		"-d",
		"-v", fmt.Sprintf("%s:/tmp/earthly:rw", VolumeName),
		"-e", fmt.Sprintf("BUILDKIT_DEBUG=%t", settings.Debug),
		"-e", fmt.Sprintf("EARTHLY_ADDITIONAL_BUILDKIT_CONFIG=%s", settings.AdditionalConfig),
		"--label", fmt.Sprintf("dev.earthly.settingshash=%s", settingsHash),
		"--name", ContainerName,
		"--privileged",
	}
	args = append(args, settings.AdditionalArgs...)
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

	if supportsPlatform(ctx) {
		args = append(args, platformFlag())
	}

	args = append(args,
		"-e", fmt.Sprintf("CACHE_SIZE_MB=%d", settings.CacheSizeMb),
		"-e", fmt.Sprintf("GIT_URL_INSTEAD_OF=%s", settings.GitURLInsteadOf),
	)

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
	// First, wait for the container to be marked as started.
	ctxTimeout1, cancel1 := context.WithTimeout(ctx, opTimeout)
	defer cancel1()
ContainerRunningLoop:
	for {
		select {
		case <-time.After(1 * time.Second):
			isRunning, err := isContainerRunning(ctxTimeout1)
			if err != nil {
				// Has not yet started. Keep waiting.
				continue
			}
			if !isRunning {
				return ErrBuildkitCrashed
			}
			if isRunning {
				break ContainerRunningLoop
			}

		case <-ctxTimeout1.Done():
			return errors.Errorf("timeout %s: buildkitd container did not start", opTimeout)
		}
	}

	// Wait for the connection to be available.
	ctxTimeout2, cancel2 := context.WithTimeout(ctx, opTimeout)
	defer cancel2()
	for {
		select {
		case <-time.After(1 * time.Second):
			// Make sure that it has not crashed on startup.
			isRunning, err := isContainerRunning(ctxTimeout2)
			if err != nil {
				return err
			}
			if !isRunning {
				return ErrBuildkitCrashed
			}
			// Attempt to connect.
			bkClient, err := client.New(ctxTimeout2, address)
			if err != nil {
				// Try again.
				continue
			}
			_, err = bkClient.ListWorkers(ctxTimeout2)
			if err != nil {
				// Try again.
				continue
			}
			err = bkClient.Close()
			if err != nil {
				return errors.Wrap(err, "close buildkit client")
			}
			return nil
		case <-ctxTimeout2.Done():
			return errors.Wrapf(ErrBuildkitStartFailure, "timeout %s: buildkitd did not make connection available after start", opTimeout)
		}
	}
}

// MaybePull checks whether an image is available locally and pulls it if it is not.
func MaybePull(ctx context.Context, image string) error {
	cmd := exec.CommandContext(ctx, "docker", "image", "inspect", image)
	_, err := cmd.CombinedOutput()
	if err == nil {
		// We found the image locally - no need to pull.
		return nil
	}
	args := []string{"pull"}
	if supportsPlatform(ctx) {
		args = append(args, platformFlag())
	}
	args = append(args, image)
	cmd = exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "docker pull %s", image)
	}
	return nil
}

// PrintLogs prints the buildkitd logs to stderr.
func PrintLogs(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "docker", "logs", ContainerName)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "docker logs %s", ContainerName)
	}
	return nil
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
			isRunning, err := isContainerRunning(ctxTimeout)
			if err != nil {
				// The container can no longer be found at all.
				return nil
			}
			if !isRunning {
				return nil
			}
		case <-ctxTimeout.Done():
			return errors.Errorf("timeout %s: buildkitd did not stop", opTimeout)
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

// CheckCompatibility runs all avaliable compatibility checks before starting the buildkitd daemon.
func CheckCompatibility(ctx context.Context, settings Settings) error {
	isNamespaced, err := isNamespacedDocker(ctx)
	if isNamespaced {
		return errors.New(`user namespaces are enabled, set "buildkit_additional_args" in ~/.earthly/config.yml to ["--userns", "host"] to disable`)
	} else if err != nil {
		return errors.Wrap(err, "failed compatibilty check")
	}

	isRootless, err := isRootlessDocker(ctx)
	if isRootless {
		return errors.New(`rootless docker detected. Compatibility is limited. Configure "buildkit_additional_args" in ~/.earthly/config.yml with some additional arguments like ["--log-opt"] to give it a shot`)
	} else if err != nil {
		return errors.Wrap(err, "failed compatibilty check")
	}

	return nil
}

func isNamespacedDocker(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx,
		"docker", "info", "--format={{.SecurityOptions}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.Wrap(err, "get docker security info")
	}

	return strings.Contains(string(output), "name=userns"), nil
}

func isRootlessDocker(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx,
		"docker", "info", "--format={{.SecurityOptions}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.Wrap(err, "get docker security info")
	}

	return strings.Contains(string(output), "rootless"), nil
}

func supportsPlatform(ctx context.Context) bool {
	// We can't run scratch, but the error is different depending on whether
	// --platform is supported or not. This is faster than attempting to run
	// an actual image which may require downloading.
	cmd := exec.CommandContext(ctx,
		"docker", "run", "--rm", platformFlag(), "scratch")
	output, _ := cmd.CombinedOutput()
	return bytes.Contains(output, []byte("Unable to find image"))
}

func platformFlag() string {
	arch := runtime.GOARCH
	if runtime.GOARCH == "arm" {
		arch = "arm/v7"
	}
	return fmt.Sprintf("--platform=linux/%s", arch)
}

func isContainerRunning(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(
		ctx, "docker", "inspect", "--format={{.State.Running}}", ContainerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.Wrapf(err, "docker inspect running")
	}
	isRunning, err := strconv.ParseBool(strings.TrimSpace(string(output)))
	if err != nil {
		return false, errors.Wrapf(err, "cannot interpret output %s", output)
	}
	return isRunning, nil
}

func isDockerAvailable(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "docker", "ps")
	err := cmd.Run()
	return err == nil
}
