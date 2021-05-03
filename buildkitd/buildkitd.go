package buildkitd

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/earthly/earthly/conslogging"
	"github.com/fatih/color"
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

// TODO: Implement all this properly with the docker client.

// NewClient returns a new buildkitd client.
func NewClient(ctx context.Context, console conslogging.ConsoleLogger, image string, settings Settings, opts ...client.ClientOpt) (*client.Client, error) {
	if isLocal(settings.BuildkitAddress) {
		err := waitForConnection(ctx, settings.BuildkitAddress, settings.Timeout)
		if err != nil {
			return nil, errors.Wrap(err, "connect provided buildkit")
		}

		bkClient, err := client.New(ctx, settings.BuildkitAddress, opts...)
		if err != nil {
			return nil, errors.Wrap(err, "start provided buildkit")
		}

		return bkClient, nil
	}

	if !isDockerAvailable(ctx) {
		console.WithPrefix("buildkitd").Printf("Is docker installed and running? Are you part of the docker group?\n")
		return nil, errors.New("docker not available")
	}
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

// ResetCache restarts the buildkitd daemon with the reset command.
func ResetCache(ctx context.Context, console conslogging.ConsoleLogger, image string, settings Settings) error {
	// Prune by resetting container.
	if settings.BuildkitAddress != "" {
		return errors.New("cannot reset cache of a provided buildkit-host setting")
	}

	console.
		WithPrefix("buildkitd").
		Printf("Restarting buildkit daemon with reset command...\n")

	// Use twice the restart timeout for reset operations
	// (needs extra time to also remove the files).
	settings.Timeout *= 2

	isStarted, err := IsStarted(ctx)
	if err != nil {
		return errors.Wrap(err, "check is started buildkitd")
	}
	if isStarted {
		err = Stop(ctx)
		if err != nil {
			return err
		}
		err = WaitUntilStopped(ctx, settings.Timeout)
		if err != nil {
			return err
		}
	}
	err = Start(ctx, console, image, settings, true)
	if err != nil {
		return err
	}
	err = WaitUntilStarted(ctx, console, settings.BuildkitAddress, settings.Timeout)
	if err != nil {
		return err
	}
	console.
		WithPrefix("buildkitd").
		Printf("... Done. Future runs will be faster.\n")
	return nil
}

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
		err := Start(ctx, console, image, settings, false)
		if err != nil {
			return "", errors.Wrap(err, "start")
		}
		err = WaitUntilStarted(ctx, console, settings.BuildkitAddress, settings.Timeout)
		if err != nil {
			return "", errors.Wrap(err, "wait until started")
		}
		console.
			WithPrefix("buildkitd").
			Printf("...Done\n")
	}
	return settings.BuildkitAddress, nil
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
			Printf("Updated image available. Restarting buildkit daemon...\n")
	}

	// Replace.
	err = Stop(ctx)
	if err != nil {
		return err
	}
	err = WaitUntilStopped(ctx, settings.Timeout)
	if err != nil {
		return err
	}
	err = Start(ctx, console, image, settings, false)
	if err != nil {
		return err
	}
	err = WaitUntilStarted(ctx, console, settings.BuildkitAddress, settings.Timeout)
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
func Start(ctx context.Context, console conslogging.ConsoleLogger, image string, settings Settings, reset bool) error {
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
	err = MaybePull(ctx, console, image)
	if err != nil {
		console.
			WithPrefix("buildkitd-pull").
			Printf("Error: %s. Attempting to start buildkitd anyway...\n", err.Error())
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
		// Debugger, and buildkit connection only supported in top-most earthly.
		// TODO: Main reason for this is port clash. This could be improved in the future,
		//       if needed.

		// These are controlled by us and should have been validated already
		bkURL, _ := url.Parse(settings.BuildkitAddress)
		dbURL, _ := url.Parse(settings.DebuggerAddress)

		args = append(args,
			"-p", fmt.Sprintf("127.0.0.1:%d:8373", dbURL.Port()),
			"-p", fmt.Sprintf("127.0.0.1:%d:8372", bkURL.Port()))
	}

	if supportsPlatform(ctx) {
		args = append(args, platformFlag())
	}

	if settings.CniMtu > 0 {
		args = append(args, "-e", fmt.Sprintf("CNI_MTU=%v", settings.CniMtu))
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
func WaitUntilStarted(ctx context.Context, console conslogging.ConsoleLogger, address string, opTimeout time.Duration) error {
	// First, wait for the container to be marked as started.
	ctxTimeout, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
ContainerRunningLoop:
	for {
		select {
		case <-time.After(1 * time.Second):
			isRunning, err := isContainerRunning(ctxTimeout)
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

		case <-ctxTimeout.Done():
			return errors.Errorf("timeout %s: buildkitd container did not start", opTimeout)
		}
	}

	// Wait for the connection to be available.
	err := waitForConnection(ctx, address, opTimeout)
	if err != nil {
		if !errors.Is(err, ErrBuildkitStartFailure) {
			return err
		}
		// We timed out. Check if the user has a lot of cache and give buildkit another chance.
		cacheSize, cacheSizeErr := getCacheSize(ctx)
		if cacheSizeErr != nil {
			console.
				WithPrefix("buildkitd").
				Printf("Warning: Could not detect buildkit cache size: %v\n", cacheSizeErr)
			return err
		}
		cacheGigs := cacheSize / 1024 / 1024
		if cacheGigs >= 30 || (cacheGigs >= 10 && runtime.GOOS == "darwin") {
			console.
				WithPrefix("buildkitd").
				Printf("Detected cache size %d GiB. It could take a while for buildkit to start up. Waiting for another %s before giving up...\n", cacheGigs, opTimeout)
			console.
				WithPrefix("buildkitd").
				Printf("To reduce the size of the cache, you can run\n" +
					"\t\tearthly config 'global.cache_size_mb' <new-size>\n" +
					"This sets the BuildKit GC target to a specific value. For more information see " +
					"the Earthly config reference page: https://docs.earthly.dev/configuration/earthly-config\n")
			return waitForConnection(ctx, address, opTimeout)
		}
		return err
	}
	return nil
}

func waitForConnection(ctx context.Context, address string, opTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	for {
		select {
		case <-time.After(1 * time.Second):
			if address == "" {
				// Make sure that our managed buildkit has not crashed on startup.
				isRunning, err := isContainerRunning(ctxTimeout)
				if err != nil {
					return err
				}

				if !isRunning {
					return ErrBuildkitCrashed
				}
			}

			err := checkConnection(ctxTimeout, address)
			if err != nil {
				// Try again.
				continue
			}
			return nil
		case <-ctxTimeout.Done():
			// Try one last time.
			err := checkConnection(ctx, address)
			if err != nil {
				// We give up.
				return errors.Wrapf(ErrBuildkitStartFailure, "timeout %s: buildkitd did not make connection available after start", opTimeout)
			}
			return nil
		}
	}
}

func checkConnection(ctx context.Context, address string) error {
	// Each attempt has limited time to succeed, to prevent hanging for too long
	// here.
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	var connErrMu sync.Mutex
	var connErr error = errors.New("timeout")
	go func() {
		defer cancel()
		bkClient, err := client.New(ctxTimeout, address)
		if err != nil {
			connErrMu.Lock()
			connErr = err
			connErrMu.Unlock()
			return
		}
		defer bkClient.Close()
		_, err = bkClient.ListWorkers(ctxTimeout)
		if err != nil {
			connErrMu.Lock()
			connErr = err
			connErrMu.Unlock()
			return
		}
		// Success.
		connErrMu.Lock()
		connErr = nil
		connErrMu.Unlock()
	}()
	<-ctxTimeout.Done()
	connErrMu.Lock()
	err := connErr
	connErrMu.Unlock()
	return err
}

// MaybePull checks whether an image is available locally and pulls it if it is not.
func MaybePull(ctx context.Context, console conslogging.ConsoleLogger, image string) error {
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
	console.
		WithPrefix("buildkitd-pull").
		Printf("Pulling buildkitd image...\n")
	err = cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "docker pull %s", image)
	}
	console.
		WithPrefix("buildkitd-pull").
		Printf("...Done\n")
	return nil
}

// PrintLogs prints the buildkitd logs to stderr.
func PrintLogs(ctx context.Context, settings Settings, console conslogging.ConsoleLogger) error {
	if !isLocal(settings.BuildkitAddress) {
		return nil
	}

	console.PrintBar(color.New(color.FgHiRed), "Buildkit Logs", "")

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
func GetContainerIP(ctx context.Context, settings Settings) (string, error) {
	if !isLocal(settings.BuildkitAddress) {
		return "", nil // Remote buildkitd is not an error,  but we don't know its IP
	}

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

// getCacheSize returns the size of the earthly cache in KiB.
func getCacheSize(ctx context.Context) (int, error) {
	cmd := exec.CommandContext(
		ctx, "docker", "volume", "inspect", VolumeName, "--format", "{{.Mountpoint}}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, errors.Wrapf(err, "get volume %s mount point", VolumeName)
	}
	mountpoint := string(bytes.TrimSpace(out))

	cmd = exec.CommandContext(
		ctx, "docker", "run", "--privileged", "--pid=host", "--rm", "busybox",
		"nsenter", "-t", "1", "-m", "-u", "-n", "-i", "--",
		"du", "-d", "0", "--", mountpoint)
	out, cmdErr := cmd.Output() // can exit with 1 if there are warnings
	parts := bytes.SplitN(bytes.TrimSpace(out), []byte("\t"), 2)
	size, err := strconv.ParseInt(string(parts[0]), 10, 64)
	if err != nil {
		if cmdErr != nil {
			return 0, errors.Wrapf(cmdErr, "parse cache size \"%s\"", parts[0])
		}
		return 0, errors.Wrapf(err, "parse cache size %s", parts[0])
	}
	return int(size), nil
}

func isLocal(addr string) bool {
	// We consider it local when the address matches one of the ones we allow in our generated GRPC certificates
	parsed, err := url.Parse(addr)
	if err != nil {
		return false
	}

	hostname := parsed.Hostname()
	return hostname == "127.0.0.1" || // The only IP v4 Loopback we honor. Because we need to include it in the TLS certificates.
		hostname == net.IPv6loopback.String() ||
		hostname == "localhost" // Convention. Users hostname omitted; this is only really here for convenience.
}
