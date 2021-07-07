package buildkitd

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/fatih/color"
	"github.com/moby/buildkit/client"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/pkg/errors"
)

var (
	// ErrBuildkitCrashed is an error returned when buildkit has terminated unexpectedly.
	ErrBuildkitCrashed = errors.New("buildkitd crashed")

	// ErrBuildkitStartFailure is an error returned when buildkit has failed to start in time.
	ErrBuildkitStartFailure = errors.New("buildkitd failed to start (in time)")
)

// TCPAddress is the address at which the daemon s available when using TCP.
var TCPAddress = "tcp://127.0.0.1:8372"

// TODO: Implement all this properly with the docker client.

// NewClient returns a new buildkitd client, together with a boolean specifying whether the buildkit is local.
func NewClient(ctx context.Context, console conslogging.ConsoleLogger, image, containerName string, settings Settings, opts ...client.ClientOpt) (*client.Client, error) {
	opts, err := addRequiredOpts(settings, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "add required client opts")
	}

	if !IsLocal(settings.BuildkitAddress) {
		err := waitForConnection(ctx, containerName, settings.BuildkitAddress, settings.Timeout, opts...)
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
	address, err := MaybeStart(ctx, console, image, containerName, settings, opts...)
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
func ResetCache(ctx context.Context, console conslogging.ConsoleLogger, image, containerName string, settings Settings, opts ...client.ClientOpt) error {
	// Prune by resetting container.
	if !IsLocal(settings.BuildkitAddress) {
		return errors.New("cannot reset cache of a provided buildkit-host setting")
	}

	opts, err := addRequiredOpts(settings, opts...)
	if err != nil {
		return errors.Wrap(err, "add required client opts")
	}

	console.
		WithPrefix("buildkitd").
		Printf("Restarting buildkit daemon with reset command...\n")

	// Use twice the restart timeout for reset operations
	// (needs extra time to also remove the files).
	settings.Timeout *= 2

	isStarted, err := IsStarted(ctx, containerName)
	if err != nil {
		return errors.Wrap(err, "check is started buildkitd")
	}
	if isStarted {
		err = Stop(ctx, containerName)
		if err != nil {
			return err
		}
		err = WaitUntilStopped(ctx, containerName, settings.Timeout)
		if err != nil {
			return err
		}
	}
	err = Start(ctx, console, image, containerName, settings, true)
	if err != nil {
		return err
	}
	err = WaitUntilStarted(ctx, console, containerName, settings.VolumeName, settings.BuildkitAddress, settings.Timeout, opts...)
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
func MaybeStart(ctx context.Context, console conslogging.ConsoleLogger, image, containerName string, settings Settings, opts ...client.ClientOpt) (string, error) {
	isStarted, err := IsStarted(ctx, containerName)
	if err != nil {
		return "", errors.Wrap(err, "check is started buildkitd")
	}
	if isStarted {
		console.
			WithPrefix("buildkitd").
			Printf("Found buildkit daemon as docker container (%s)\n", containerName)
		err := MaybeRestart(ctx, console, image, containerName, settings, opts...)
		if err != nil {
			return "", errors.Wrap(err, "maybe restart")
		}
	} else {
		console.
			WithPrefix("buildkitd").
			Printf("Starting buildkit daemon as a docker container (%s)...\n", containerName)
		err := Start(ctx, console, image, containerName, settings, false)
		if err != nil {
			return "", errors.Wrap(err, "start")
		}
		err = WaitUntilStarted(ctx, console, containerName, settings.VolumeName, settings.BuildkitAddress, settings.Timeout, opts...)
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
func MaybeRestart(ctx context.Context, console conslogging.ConsoleLogger, image, containerName string, settings Settings, opts ...client.ClientOpt) error {
	containerImageID, err := GetContainerImageID(ctx, containerName)
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
	console.
		WithPrefix("buildkitd").
		VerbosePrintf("Comparing running container image (%q) with available image (%q)\n", containerImageID, availableImageID)
	if containerImageID == availableImageID {
		// Images are the same. Check settings hash.
		hash, err := GetSettingsHash(ctx, containerName)
		if err != nil {
			return err
		}
		ok, err := settings.VerifyHash(hash)
		if err != nil {
			return errors.Wrap(err, "verify hash")
		}
		if ok {
			// No need to replace: images are the same and settings are the same.
			console.
				WithPrefix("buildkitd").
				VerbosePrintf("Settings hashes match (%q), no restart required\n", hash)
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
	err = Stop(ctx, containerName)
	if err != nil {
		return err
	}
	err = WaitUntilStopped(ctx, containerName, settings.Timeout)
	if err != nil {
		return err
	}
	err = Start(ctx, console, image, containerName, settings, false)
	if err != nil {
		return err
	}
	err = WaitUntilStarted(ctx, console, containerName, settings.VolumeName, settings.BuildkitAddress, settings.Timeout, opts...)
	if err != nil {
		return err
	}
	console.
		WithPrefix("buildkitd").
		Printf("...Done\n")
	return nil
}

// RemoveExited removes any stopped or exited buildkitd containers
func RemoveExited(ctx context.Context, containerName string) error {
	cmd := exec.CommandContext(ctx, "docker", "ps", "-a", "-q", "-f", fmt.Sprintf("name=%s", containerName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "get combined output")
	}
	if len(output) == 0 {
		return nil
	}
	return exec.CommandContext(ctx, "docker", "rm", containerName).Run()
}

// Start starts the buildkitd daemon.
func Start(ctx context.Context, console conslogging.ConsoleLogger, image, containerName string, settings Settings, reset bool) error {
	err := CheckCompatibility(ctx, settings)
	if len(settings.AdditionalArgs) == 0 && err != nil {
		return errors.Wrap(err, "compatibility")
	}

	settingsHash, err := settings.Hash()
	if err != nil {
		return errors.Wrap(err, "settings hash")
	}
	err = RemoveExited(ctx, containerName)
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
		"-v", fmt.Sprintf("%s:/tmp/earthly:rw", settings.VolumeName),
		"-e", fmt.Sprintf("BUILDKIT_DEBUG=%t", settings.Debug),
		"-e", fmt.Sprintf("EARTHLY_ADDITIONAL_BUILDKIT_CONFIG=%s", settings.AdditionalConfig),
		"-e", fmt.Sprintf("BUILDKIT_TCP_TRANSPORT_ENABLED=%t", settings.UseTCP),
		"-e", fmt.Sprintf("BUILDKIT_TLS_ENABLED=%t", settings.UseTCP && settings.UseTLS),
		"--label", fmt.Sprintf("dev.earthly.settingshash=%s", settingsHash),
		"--name", containerName,
		"--privileged",
	}
	args = append(args, settings.AdditionalArgs...)
	if os.Getenv("EARTHLY_WITH_DOCKER") == "1" {
		// Add /sys/fs/cgroup if it's earthly-in-earthly.
		args = append(args, "-v", "/sys/fs/cgroup:/sys/fs/cgroup")
	} else {
		// TCP ports only supported in top-most earthly.
		// TODO: Main reason for this is port clash. This could be improved in the future,
		//       if needed.
		// These are controlled by us and should have been validated already - hence panics.

		dbURL, err := url.Parse(settings.DebuggerAddress)
		if err != nil {
			panic("Debugger address was not a URL when attempting to start buildkit")
		}
		args = append(args, "-p", fmt.Sprintf("127.0.0.1:%s:8373", dbURL.Port()))

		if settings.LocalRegistryAddress != "" {
			lrURL, err := url.Parse(settings.LocalRegistryAddress)
			if err != nil {
				panic("Local registry address was not a URL when attempting to start buildkit")
			}
			args = append(args, "-p", fmt.Sprintf("127.0.0.1:%s:8371", lrURL.Port()))
			args = append(args, "-e", "BUILDKIT_LOCAL_REGISTRY_LISTEN_PORT=8371")
		}

		bkURL, err := url.Parse(settings.BuildkitAddress)
		if err != nil {
			panic("Buildkit address was not a URL when attempting to start buildkit")
		}
		if settings.UseTCP {
			args = append(args, "-p", fmt.Sprintf("127.0.0.1:%s:8372", bkURL.Port()))

			if settings.UseTLS {
				if settings.TLSCA != "" {
					caPath, err := makeTLSPath(settings.TLSCA)
					if err != nil {
						return errors.Wrap(err, "start buildkitd")
					}
					args = append(args, "-v", fmt.Sprintf("%s:/etc/ca.pem", caPath))
				}

				if settings.ServerTLSCert != "" {
					certPath, err := makeTLSPath(settings.ServerTLSCert)
					if err != nil {
						return errors.Wrap(err, "start buildkitd")
					}
					args = append(args, "-v", fmt.Sprintf("%s:/etc/cert.pem", certPath))
				}

				if settings.ServerTLSKey != "" {
					keyPath, err := makeTLSPath(settings.ServerTLSKey)
					if err != nil {
						return errors.Wrap(err, "start buildkitd")
					}
					args = append(args, "-v", fmt.Sprintf("%s:/etc/key.pem", keyPath))
				}
			}
		}
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
func Stop(ctx context.Context, containerName string) error {
	cmd := exec.CommandContext(ctx, "docker", "stop", containerName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "get combined output")
	}
	return nil
}

// IsStarted checks if the buildkitd container has been started.
func IsStarted(ctx context.Context, containerName string) (bool, error) {
	cmd := exec.CommandContext(ctx, "docker", "ps", "-q", "-f", fmt.Sprintf("name=%s", containerName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.Wrap(err, "get combined output")
	}
	return (len(output) != 0), nil
}

// WaitUntilStarted waits until the buildkitd daemon has started and is healthy.
func WaitUntilStarted(ctx context.Context, console conslogging.ConsoleLogger, containerName, volumeName, address string, opTimeout time.Duration, opts ...client.ClientOpt) error {
	// First, wait for the container to be marked as started.
	ctxTimeout, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
ContainerRunningLoop:
	for {
		select {
		case <-time.After(1 * time.Second):
			isRunning, err := isContainerRunning(ctxTimeout, containerName)
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
	err := waitForConnection(ctx, containerName, address, opTimeout, opts...)
	if err != nil {
		if !errors.Is(err, ErrBuildkitStartFailure) {
			return err
		}
		// We timed out. Check if the user has a lot of cache and give buildkit another chance.
		cacheSize, cacheSizeErr := getCacheSize(ctx, volumeName)
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
			return waitForConnection(ctx, containerName, address, opTimeout)
		}
		return err
	}
	return nil
}

func waitForConnection(ctx context.Context, containerName, address string, opTimeout time.Duration, opts ...client.ClientOpt) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	for {
		select {
		case <-time.After(1 * time.Second):
			if address == "" {
				// Make sure that our managed buildkit has not crashed on startup.
				isRunning, err := isContainerRunning(ctxTimeout, containerName)
				if err != nil {
					return err
				}

				if !isRunning {
					return ErrBuildkitCrashed
				}
			}

			err := checkConnection(ctxTimeout, address, opts...)
			if err != nil {
				// Try again.
				continue
			}
			return nil
		case <-ctxTimeout.Done():
			// Try one last time.
			err := checkConnection(ctx, address, opts...)
			if err != nil {
				// We give up.
				return errors.Wrapf(ErrBuildkitStartFailure, "timeout %s: buildkitd did not make connection available after start", opTimeout)
			}
			return nil
		}
	}
}

func checkConnection(ctx context.Context, address string, opts ...client.ClientOpt) error {
	// Each attempt has limited time to succeed, to prevent hanging for too long
	// here.
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	var connErrMu sync.Mutex
	var connErr error = errors.New("timeout")
	go func() {
		defer cancel()
		bkClient, err := client.New(ctxTimeout, address, opts...)
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
func PrintLogs(ctx context.Context, containerName string, settings Settings, console conslogging.ConsoleLogger) error {
	if !IsLocal(settings.BuildkitAddress) {
		return nil
	}

	console.PrintBar(color.New(color.FgHiRed), "Buildkit Logs", "")

	cmd := exec.CommandContext(ctx, "docker", "logs", containerName)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "docker logs %s", containerName)
	}
	return nil
}

// GetContainerIP returns the IP of the buildkit container.
func GetContainerIP(ctx context.Context, containerName string, settings Settings) (string, error) {
	if !IsLocal(settings.BuildkitAddress) {
		return "", nil // Remote buildkitd is not an error,  but we don't know its IP
	}

	cmd := exec.CommandContext(ctx, "docker", "inspect", "-f", "{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "get combined output ip")
	}
	return string(bytes.TrimSpace(output)), nil
}

// WaitUntilStopped waits until the buildkitd daemon has stopped.
func WaitUntilStopped(ctx context.Context, containerName string, opTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	for {
		select {
		case <-time.After(1 * time.Second):
			isRunning, err := isContainerRunning(ctxTimeout, containerName)
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
func GetSettingsHash(ctx context.Context, containerName string) (string, error) {
	cmd := exec.CommandContext(ctx,
		"docker", "inspect",
		"--format={{index .Config.Labels \"dev.earthly.settingshash\"}}",
		containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "get output for settings hash")
	}
	return string(output), nil
}

// GetContainerImageID fetches the ID of the image used for the running buildkitd container.
func GetContainerImageID(ctx context.Context, containerName string) (string, error) {
	cmd := exec.CommandContext(ctx,
		"docker", "inspect", "--format={{index .Image}}", containerName)
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

func isContainerRunning(ctx context.Context, containerName string) (bool, error) {
	cmd := exec.CommandContext(
		ctx, "docker", "inspect", "--format={{.State.Running}}", containerName)
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
func getCacheSize(ctx context.Context, volumeName string) (int, error) {
	cmd := exec.CommandContext(
		ctx, "docker", "volume", "inspect", volumeName, "--format", "{{.Mountpoint}}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, errors.Wrapf(err, "get volume %s mount point", volumeName)
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

// IsLocal parses a URL and returns whether it is considered a local buildkit host + port that we
// need to manage ourselves.
func IsLocal(addr string) bool {
	parsed, err := url.Parse(addr)
	if err != nil {
		return false
	}

	hostname := parsed.Hostname()
	// These need to match what we put in our certificates.
	return hostname == "127.0.0.1" || // The only IP v4 Loopback we honor. Because we need to include it in the TLS certificates.
		hostname == net.IPv6loopback.String() ||
		hostname == "localhost" || // Convention. Users hostname omitted; this is only really here for convenience.
		parsed.Scheme == "docker-container" // Accomodate feature flagging during transition. This will have omitted TLS?
}

func makeTLSPath(path string) (string, error) {
	fullPath := path

	if !filepath.IsAbs(path) {
		earthlyDir, err := cliutil.GetOrCreateEarthlyDir()
		if err != nil {
			return "", err
		}

		fullPath = filepath.Join(earthlyDir, path)
	}

	if !fileutil.FileExists(fullPath) {
		return "", fmt.Errorf("path '%s' does not exist", path)
	}

	return fullPath, nil
}

func addRequiredOpts(settings Settings, opts ...client.ClientOpt) ([]client.ClientOpt, error) {
	if !settings.UseTCP || !settings.UseTLS {
		return opts, nil
	}

	server, err := url.Parse(settings.BuildkitAddress)
	if err != nil {
		return []client.ClientOpt{}, errors.Wrap(err, "invalid buildkit url")
	}

	caPath, err := makeTLSPath(settings.TLSCA)
	if err != nil {
		return []client.ClientOpt{}, errors.Wrap(err, "caPath")
	}

	certPath, err := makeTLSPath(settings.ClientTLSCert)
	if err != nil {
		return []client.ClientOpt{}, errors.Wrap(err, "certPath")
	}

	keyPath, err := makeTLSPath(settings.ClientTLSKey)
	if err != nil {
		return []client.ClientOpt{}, errors.Wrap(err, "keyPath")
	}

	return append(opts, client.WithCredentials(server.Hostname(), caPath, certPath, keyPath)), nil
}
