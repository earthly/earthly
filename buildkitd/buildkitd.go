package buildkitd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/moby/buildkit/client"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/pkg/errors"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/fileutil"
)

var (
	// ErrBuildkitCrashed is an error returned when buildkit has terminated unexpectedly.
	ErrBuildkitCrashed = errors.New("buildkitd crashed")

	// ErrBuildkitStartFailure is an error returned when buildkit has failed to start in time.
	ErrBuildkitStartFailure = errors.New("buildkitd failed to start (in time)")
)

// NewClient returns a new buildkitd client, together with a boolean specifying whether the buildkit is local.
func NewClient(ctx context.Context, console conslogging.ConsoleLogger, image, containerName string, fe containerutil.ContainerFrontend, settings Settings, opts ...client.ClientOpt) (*client.Client, error) {
	opts, err := addRequiredOpts(settings, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "add required client opts")
	}

	if !containerutil.IsLocal(settings.BuildkitAddress) {
		if settings.SatelliteName != "" {
			console.WithPrefix("satellite").Printf("Connecting to %s", settings.SatelliteName)
		}

		err := waitForConnection(ctx, containerName, settings.BuildkitAddress, settings.Timeout, fe, opts...)
		if err != nil {
			return nil, errors.Wrap(err, "connect provided buildkit")
		}

		bkClient, err := client.New(ctx, settings.BuildkitAddress, opts...)
		if err != nil {
			return nil, errors.Wrap(err, "start provided buildkit")
		}

		return bkClient, nil
	}

	if !isDockerAvailable(ctx, fe) {
		console.WithPrefix("buildkitd").Printf("Is %[1]s installed and running? Are you part of any needed groups?\n", fe.Config().Binary)
		return nil, fmt.Errorf("%s not available", fe.Config().Binary)
	}
	address, err := MaybeStart(ctx, console, image, containerName, fe, settings, opts...)
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
func ResetCache(ctx context.Context, console conslogging.ConsoleLogger, image, containerName string, fe containerutil.ContainerFrontend, settings Settings, opts ...client.ClientOpt) error {
	// Prune by resetting container.
	if !containerutil.IsLocal(settings.BuildkitAddress) {
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

	isStarted, err := IsStarted(ctx, containerName, fe)
	if err != nil {
		return errors.Wrap(err, "check is started buildkitd")
	}
	if isStarted {
		err = Stop(ctx, containerName, fe)
		if err != nil {
			return err
		}
		err = WaitUntilStopped(ctx, containerName, settings.Timeout, fe)
		if err != nil {
			return err
		}
	}
	err = Start(ctx, console, image, containerName, fe, settings, true)
	if err != nil {
		return err
	}
	err = WaitUntilStarted(ctx, console, containerName, settings.VolumeName, settings.BuildkitAddress, settings.Timeout, fe, opts...)
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
func MaybeStart(ctx context.Context, console conslogging.ConsoleLogger, image, containerName string, fe containerutil.ContainerFrontend, settings Settings, opts ...client.ClientOpt) (string, error) {
	isStarted, err := IsStarted(ctx, containerName, fe)
	if err != nil {
		return "", errors.Wrap(err, "check is started buildkitd")
	}
	if isStarted {
		console.
			WithPrefix("buildkitd").
			Printf("Found buildkit daemon as %s container (%s)\n", fe.Config().Binary, containerName)
		err := MaybeRestart(ctx, console, image, containerName, fe, settings, opts...)
		if err != nil {
			return "", errors.Wrap(err, "maybe restart")
		}
	} else {
		console.
			WithPrefix("buildkitd").
			Printf("Starting buildkit daemon as a %s container (%s)...\n", fe.Config().Binary, containerName)
		err := Start(ctx, console, image, containerName, fe, settings, false)
		if err != nil {
			return "", errors.Wrap(err, "start")
		}
		err = WaitUntilStarted(ctx, console, containerName, settings.VolumeName, settings.BuildkitAddress, settings.Timeout, fe, opts...)
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
func MaybeRestart(ctx context.Context, console conslogging.ConsoleLogger, image, containerName string, fe containerutil.ContainerFrontend, settings Settings, opts ...client.ClientOpt) error {
	containerImageID, err := GetContainerImageID(ctx, containerName, fe)
	if err != nil {
		return err
	}
	availableImageID, err := GetAvailableImageID(ctx, image, fe)
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
		hash, err := GetSettingsHash(ctx, containerName, fe)
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
	err = Stop(ctx, containerName, fe)
	if err != nil {
		return err
	}
	err = WaitUntilStopped(ctx, containerName, settings.Timeout, fe)
	if err != nil {
		return err
	}
	err = Start(ctx, console, image, containerName, fe, settings, false)
	if err != nil {
		return err
	}
	err = WaitUntilStarted(ctx, console, containerName, settings.VolumeName, settings.BuildkitAddress, settings.Timeout, fe, opts...)
	if err != nil {
		return err
	}
	console.
		WithPrefix("buildkitd").
		Printf("...Done\n")
	return nil
}

// RemoveExited removes any stopped or exited buildkitd containers
func RemoveExited(ctx context.Context, fe containerutil.ContainerFrontend, containerName string) error {
	infos, err := fe.ContainerInfo(ctx, containerName)
	if err != nil {
		return errors.Wrapf(err, "get info to remove exited %s", containerName)
	}
	containerInfo, ok := infos[containerName]
	if !ok || containerInfo.Status == containerutil.StatusMissing {
		return nil
	}

	err = fe.ContainerRemove(ctx, false, containerName)
	if err != nil {
		return errors.Wrapf(err, "remove exited %s", containerName)
	}

	return nil
}

// Start starts the buildkitd daemon.
func Start(ctx context.Context, console conslogging.ConsoleLogger, image, containerName string, fe containerutil.ContainerFrontend, settings Settings, reset bool) error {
	settingsHash, err := settings.Hash()
	if err != nil {
		return errors.Wrap(err, "settings hash")
	}
	err = RemoveExited(ctx, fe, containerName)
	if err != nil {
		return err
	}
	// Pulling is not strictly needed, but it helps display some progress status to the user in
	// case the image is not available locally.
	err = MaybePull(ctx, console, image, fe)
	if err != nil {
		console.
			WithPrefix("buildkitd-pull").
			Printf("Error: %s. Attempting to start buildkitd anyway...\n", err.Error())
		// Keep going - it might still work.
	}

	envOpts := map[string]string{
		"BUILDKIT_DEBUG":                 strconv.FormatBool(settings.Debug),
		"BUILDKIT_TCP_TRANSPORT_ENABLED": strconv.FormatBool(settings.UseTCP),
		"BUILDKIT_TLS_ENABLED":           strconv.FormatBool(settings.UseTCP && settings.UseTLS),
		"BUILDKIT_MAX_PARALLELISM":       strconv.Itoa(settings.MaxParallelism),
	}

	labelOpts := map[string]string{
		"dev.earthly.settingshash": settingsHash,
	}

	volumeOpts := containerutil.MountOpt{
		containerutil.Mount{
			Type:     containerutil.MountVolume,
			Source:   settings.VolumeName,
			Dest:     "/tmp/earthly",
			ReadOnly: false,
		},
	}

	portOpts := containerutil.PortOpt{}

	if settings.AdditionalConfig != "" {
		envOpts["EARTHLY_ADDITIONAL_BUILDKIT_CONFIG"] = settings.AdditionalConfig
	}

	if settings.IPTables != "" {
		envOpts["IP_TABLES"] = settings.IPTables
	}

	if os.Getenv("EARTHLY_WITH_DOCKER") == "1" {
		// Add /sys/fs/cgroup if it's earthly-in-earthly.
		volumeOpts = append(volumeOpts, containerutil.Mount{
			Type:   containerutil.MountBind,
			Source: "/sys/fs/cgroup",
			Dest:   "/sys/fs/cgroup",
		})
	} else {
		// TCP ports only supported in top-most earthly.
		// TODO: Main reason for this is port clash. This could be improved in the future,
		//       if needed.
		// These are controlled by us and should have been validated already - hence panics.

		dbURL, err := url.Parse(settings.DebuggerAddress)
		if err != nil {
			panic("Debugger address was not a URL when attempting to start buildkit")
		}

		hostPort, err := strconv.Atoi(dbURL.Port())
		if err != nil {
			panic("Local registry host port was not a number when attempting to start buildkit")
		}

		portOpts = append(portOpts, containerutil.Port{
			IP:            "127.0.0.1",
			HostPort:      hostPort,
			ContainerPort: 8373,
			Protocol:      containerutil.ProtocolTCP,
		})

		if settings.LocalRegistryAddress != "" {
			lrURL, err := url.Parse(settings.LocalRegistryAddress)
			if err != nil {
				panic("Local registry address was not a URL when attempting to start buildkit")
			}
			hostPort, err := strconv.Atoi(lrURL.Port())
			if err != nil {
				panic("Local registry host port was not a number when attempting to start buildkit")
			}
			portOpts = append(portOpts, containerutil.Port{
				IP:            "127.0.0.1",
				HostPort:      hostPort,
				ContainerPort: 8371,
				Protocol:      containerutil.ProtocolTCP,
			})
		}

		bkURL, err := url.Parse(settings.BuildkitAddress)
		if err != nil {
			panic("Buildkit address was not a URL when attempting to start buildkit")
		}
		if settings.UseTCP {
			hostPort, err := strconv.Atoi(bkURL.Port())
			if err != nil {
				panic("Local registry host port was not a number when attempting to start buildkit")
			}
			portOpts = append(portOpts, containerutil.Port{
				IP:            "127.0.0.1",
				HostPort:      hostPort,
				ContainerPort: 8372,
				Protocol:      containerutil.ProtocolTCP,
			})
			if settings.UseTLS {
				if settings.TLSCA != "" {
					caPath, err := makeTLSPath(settings.TLSCA)
					if err != nil {
						return errors.Wrap(err, "start buildkitd")
					}
					volumeOpts = append(volumeOpts, containerutil.Mount{
						Type:     containerutil.MountBind,
						Source:   caPath,
						Dest:     "/etc/ca.pem",
						ReadOnly: true,
					})
				}

				if settings.ServerTLSCert != "" {
					certPath, err := makeTLSPath(settings.ServerTLSCert)
					if err != nil {
						return errors.Wrap(err, "start buildkitd")
					}
					volumeOpts = append(volumeOpts, containerutil.Mount{
						Type:     containerutil.MountBind,
						Source:   certPath,
						Dest:     "/etc/cert.pem",
						ReadOnly: true,
					})
				}

				if settings.ServerTLSKey != "" {
					keyPath, err := makeTLSPath(settings.ServerTLSKey)
					if err != nil {
						return errors.Wrap(err, "start buildkitd")
					}
					volumeOpts = append(volumeOpts, containerutil.Mount{
						Type:     containerutil.MountBind,
						Source:   keyPath,
						Dest:     "/etc/key.pem",
						ReadOnly: true,
					})
				}
			}
		}
	}

	if settings.CniMtu > 0 {
		envOpts["CNI_MTU"] = strconv.FormatUint(uint64(settings.CniMtu), 10)
	}

	if settings.CacheSizeMb > 0 {
		envOpts["CACHE_SIZE_MB"] = strconv.FormatInt(int64(settings.CacheSizeMb), 10)
	}

	if settings.CacheSizePct > 0 {
		envOpts["CACHE_SIZE_PCT"] = strconv.FormatInt(int64(settings.CacheSizePct), 10)
	}

	// Apply reset.
	if reset {
		envOpts["EARTHLY_RESET_TMP_DIR"] = "true"
	}

	// Execute.
	err = fe.ContainerRun(ctx, containerutil.ContainerRun{
		NameOrID:       containerName,
		ImageRef:       image,
		Privileged:     true,
		Envs:           envOpts,
		Labels:         labelOpts,
		Mounts:         volumeOpts,
		Ports:          portOpts,
		AdditionalArgs: settings.AdditionalArgs,
	})
	if err != nil {
		return errors.Wrap(err, "could not start buildkit")
	}

	return nil
}

// Stop stops the buildkitd container.
func Stop(ctx context.Context, containerName string, fe containerutil.ContainerFrontend) error {
	return fe.ContainerStop(ctx, 10, containerName)
}

// IsStarted checks if the buildkitd container has been started.
func IsStarted(ctx context.Context, containerName string, fe containerutil.ContainerFrontend) (bool, error) {
	infos, err := fe.ContainerInfo(ctx, containerName)
	if err != nil {
		return false, err
	}
	containerInfo, ok := infos[containerName]
	if !ok {
		return false, err
	}
	return containerInfo.Status == containerutil.StatusRunning, nil
}

// WaitUntilStarted waits until the buildkitd daemon has started and is healthy.
func WaitUntilStarted(ctx context.Context, console conslogging.ConsoleLogger, containerName, volumeName, address string, opTimeout time.Duration, fe containerutil.ContainerFrontend, opts ...client.ClientOpt) error {
	// First, wait for the container to be marked as started.
	ctxTimeout, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
ContainerRunningLoop:
	for {
		select {
		case <-time.After(1 * time.Second):
			isRunning, err := isContainerRunning(ctxTimeout, containerName, fe)
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
	err := waitForConnection(ctx, containerName, address, opTimeout, fe, opts...)
	if err != nil {
		if !errors.Is(err, ErrBuildkitStartFailure) {
			return err
		}
		// We timed out. Check if the user has a lot of cache and give buildkit another chance.
		cacheSize, cacheSizeErr := getCacheSize(ctx, volumeName, fe)
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
				Printf("To reduce the size of the cache, you can run one of\n" +
					"\t\tearthly config 'global.cache_size_mb' <new-size>\n" +
					"\t\tearthly config 'global.cache_size_pct' <new-percent>\n" +
					"These set the BuildKit GC target to a specific value. For more information see " +
					"the Earthly config reference page: https://docs.earthly.dev/docs/earthly-config\n")
			return waitForConnection(ctx, containerName, address, opTimeout, fe)
		}
		return err
	}
	return nil
}

func waitForConnection(ctx context.Context, containerName, address string, opTimeout time.Duration, fe containerutil.ContainerFrontend, opts ...client.ClientOpt) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	for {
		select {
		case <-time.After(1 * time.Second):
			if address == "" {
				// Make sure that our managed buildkit has not crashed on startup.
				isRunning, err := isContainerRunning(ctxTimeout, containerName, fe)
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
				return errors.Wrapf(ErrBuildkitStartFailure, "timeout %s: buildkitd did not make connection available after start with error: %s", opTimeout, err.Error())
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
func MaybePull(ctx context.Context, console conslogging.ConsoleLogger, image string, fe containerutil.ContainerFrontend) error {
	infos, err := fe.ImageInfo(ctx, image)
	if err != nil {
		return errors.Wrap(err, "could not get container info")
	}
	if len(infos) > 0 { // the presence of an item implies its local
		return nil
	}

	console.
		WithPrefix("buildkitd-pull").
		Printf("Pulling buildkitd image...\n")
	err = fe.ImagePull(ctx, image)
	if err != nil {
		return errors.Wrapf(err, "could not pull %s", image)
	}
	console.
		WithPrefix("buildkitd-pull").
		Printf("...Done\n")
	return nil
}

// GetDockerVersion returns the docker version command output
func GetDockerVersion(ctx context.Context, fe containerutil.ContainerFrontend) (string, error) {
	info, err := fe.Information(ctx)
	if err != nil {
		return "", errors.Wrap(err, "get info from frontend")
	}

	return fmt.Sprintf("%#v", info), nil
}

// GetLogs returns earthly-buildkitd logs
func GetLogs(ctx context.Context, containerName string, fe containerutil.ContainerFrontend, settings Settings) (string, error) {
	if !containerutil.IsLocal(settings.BuildkitAddress) {
		return "", nil
	}

	logs, err := fe.ContainerLogs(ctx, containerName)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	if containerLogs, ok := logs[containerName]; ok {
		return containerLogs.Stdout, nil
	}

	return "", fmt.Errorf("logs for container %s were not found", containerName)
}

// GetContainerIP returns the IP of the buildkit container.
func GetContainerIP(ctx context.Context, containerName string, fe containerutil.ContainerFrontend, settings Settings) (string, error) {
	if !containerutil.IsLocal(settings.BuildkitAddress) {
		return "", nil // Remote buildkitd is not an error,  but we don't know its IP
	}

	infos, err := fe.ContainerInfo(ctx, containerName)
	if err != nil {
		return "", errors.Wrap(err, "could not get container info to determine ip")
	}

	if containerInfo, ok := infos[containerName]; ok {
		// default is bridge. If someone has a weirdo setup this should be able to handle it with some config option.
		return containerInfo.IPs["bridge"], nil
	}

	return "", fmt.Errorf("ip for container %s was not found", containerName)
}

// WaitUntilStopped waits until the buildkitd daemon has stopped.
func WaitUntilStopped(ctx context.Context, containerName string, opTimeout time.Duration, fe containerutil.ContainerFrontend) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	for {
		select {
		case <-time.After(1 * time.Second):
			isRunning, err := isContainerRunning(ctxTimeout, containerName, fe)
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
func GetSettingsHash(ctx context.Context, containerName string, fe containerutil.ContainerFrontend) (string, error) {
	infos, err := fe.ContainerInfo(ctx, containerName)
	if err != nil {
		return "", errors.Wrap(err, "get container info for settings")
	}

	if containerInfo, ok := infos[containerName]; ok {
		return containerInfo.Labels["dev.earthly.settingshash"], nil
	}

	return "", fmt.Errorf("settings hash for container %s was not found", containerName)
}

// GetContainerImageID fetches the ID of the image used for the running buildkitd container.
func GetContainerImageID(ctx context.Context, containerName string, fe containerutil.ContainerFrontend) (string, error) {
	infos, err := fe.ContainerInfo(ctx, containerName)
	if err != nil {
		return "", errors.Wrap(err, "get container info for current container image ID")
	}

	if containerInfo, ok := infos[containerName]; ok {
		return containerInfo.ImageID, nil
	}

	return "", fmt.Errorf("image id for container %s was not found", containerName)

}

// GetAvailableImageID fetches the ID of the image buildkitd image available.
func GetAvailableImageID(ctx context.Context, image string, fe containerutil.ContainerFrontend) (string, error) {
	infos, err := fe.ImageInfo(ctx, image)
	if err != nil {
		return "", errors.Wrap(err, "get output for available image ID")
	}
	return infos[image].ID, nil
}

func isContainerRunning(ctx context.Context, containerName string, fe containerutil.ContainerFrontend) (bool, error) {
	infos, err := fe.ContainerInfo(ctx, containerName)
	if err != nil {
		return false, errors.Wrap(err, "failed to get container info while checking if running")
	}

	if containerInfo, ok := infos[containerName]; ok {
		return containerInfo.Status == containerutil.StatusRunning, nil
	}

	return false, fmt.Errorf("status for container %s was not found", containerName)
}

func isDockerAvailable(ctx context.Context, fe containerutil.ContainerFrontend) bool {
	return fe.IsAvaliable(ctx)
}

// getCacheSize returns the size of the earthly cache in KiB.
func getCacheSize(ctx context.Context, volumeName string, fe containerutil.ContainerFrontend) (int, error) {
	infos, err := fe.VolumeInfo(ctx, volumeName)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get volume info for cache size %s", volumeName)
	}

	return int(infos[volumeName].Size), nil
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

	exists, err := fileutil.FileExists(fullPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to check if %s exists", fullPath)
	}
	if !exists {
		return "", fmt.Errorf("path '%s' does not exist", path)
	}

	return fullPath, nil
}

func addRequiredOpts(settings Settings, opts ...client.ClientOpt) ([]client.ClientOpt, error) {
	if settings.SatelliteName != "" {
		return append(opts, client.WithAdditionalMetadataContext(
			"satellite_name", settings.SatelliteName,
			"satellite_org", settings.SatelliteOrgID,
			"satellite_token", settings.SatelliteToken),
		), nil
	}

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
