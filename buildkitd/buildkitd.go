package buildkitd

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/containerd/containerd/platforms"
	"github.com/docker/go-units"
	"github.com/dustin/go-humanize"
	"github.com/gofrs/flock"
	"github.com/moby/buildkit/client"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/earthly/earthly/ast/hint"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/buildkitutil"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/earthly/earthly/util/semverutil"
)

const minRecommendedCacheSize = 10 << 30 // 10 GiB

var (
	// ErrBuildkitCrashed is an error returned when buildkit has terminated unexpectedly.
	ErrBuildkitCrashed = errors.New("buildkitd crashed")

	// ErrBuildkitConnectionFailure is an error returned when buildkit has failed to respond.
	ErrBuildkitConnectionFailure = errors.New("buildkitd did not respond (in time)")
)

// NewClient returns a new buildkitd client. If the buildkitd daemon is local, this function
// might start one up, if not already started.
func NewClient(ctx context.Context, console conslogging.ConsoleLogger, image, containerName, installationName string, fe containerutil.ContainerFrontend, earthlyVersion string, settings Settings, opts ...client.ClientOpt) (_ *client.Client, retErr error) {
	defer func() {
		if retErr == nil {
			return
		}
		if errors.Is(retErr, os.ErrNotExist) {
			switch fe.Config().Setting {
			case containerutil.FrontendPodman, containerutil.FrontendPodmanShell:
				tlsPaths := []string{
					settings.TLSCA,
					settings.ServerTLSKey,
					settings.ServerTLSCert,
					settings.ClientTLSKey,
					settings.ClientTLSCert,
				}
				if containsAny(retErr.Error(), tlsPaths...) {
					retErr = hint.Wrap(retErr,
						"podman now requires TLS certs by default - try stopping the earthly-buildkitd container and re-running 'earthly bootstrap'",
						"alternatively, run 'earthly config global.tls_enabled false' to disable TLS",
					)
				}
			default:
			}
			return
		}
		if strings.Contains(retErr.Error(), rsa.ErrVerification.Error()) {
			// verification errors can happen server-side, which means
			// errors.Is() won't work. We use strings.Contains instead to handle
			// that case.
			retErr = hint.Wrap(retErr, "did earthly's certificates get regenerated? you may need to manually stop the earthly-buildkitd container.")
			return
		}
	}()

	opts, err := addRequiredOpts(settings, installationName, fe.Config().Setting == containerutil.FrontendPodmanShell, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "add required client opts")
	}

	isLocal := containerutil.IsLocal(settings.BuildkitAddress)
	if !isLocal {
		remoteConsole := console
		if settings.SatelliteName != "" {
			remoteConsole = console.WithPrefix("satellite")
			if settings.SatelliteIsManaged {
				remoteConsole.Printf("Connecting to %s...", settings.SatelliteDisplayName)
			} else {
				remoteConsole.Printf("Connecting to %s (hosted at: %s)...", settings.SatelliteDisplayName, settings.BuildkitAddress)
			}
		} else {
			remoteConsole = console.WithPrefix("buildkitd")
			remoteConsole.Printf("Connecting to %s...", settings.BuildkitAddress)
		}

		info, workerInfo, err := waitForConnection(ctx, containerName, settings.BuildkitAddress, settings.Timeout, fe, opts...)
		if err != nil {
			return nil, errors.Wrap(err, "connect provided buildkit")
		}
		remoteConsole.Printf("...Done")
		printBuildkitInfo(remoteConsole, info, workerInfo, earthlyVersion, isLocal, settings.HasConfiguredCacheSize())

		bkClient, err := client.New(ctx, settings.BuildkitAddress, opts...)
		if err != nil {
			return nil, errors.Wrap(err, "start provided buildkit")
		}
		return bkClient, nil
	}

	bkCons := console.WithPrefix("buildkitd")
	if !isDockerAvailable(ctx, fe) {
		bkCons.Printf("Is %[1]s installed and running? Are you part of any needed groups?\n", fe.Config().Binary)
		return nil, fmt.Errorf("%s not available", fe.Config().Binary)
	}
	info, workerInfo, err := maybeStart(ctx, console, image, containerName, installationName, fe, settings, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "maybe start buildkitd")
	}
	printBuildkitInfo(bkCons, info, workerInfo, earthlyVersion, isLocal, settings.HasConfiguredCacheSize())
	bkClient, err := client.New(ctx, settings.BuildkitAddress, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "new buildkit client")
	}
	return bkClient, nil
}

// ResetCache restarts the buildkitd daemon with the reset command.
func ResetCache(ctx context.Context, console conslogging.ConsoleLogger, image, containerName, installationName string, fe containerutil.ContainerFrontend, settings Settings, opts ...client.ClientOpt) error {
	// Prune by resetting container.
	if !containerutil.IsLocal(settings.BuildkitAddress) {
		return errors.New("cannot reset cache of a provided buildkit-host setting")
	}

	opts, err := addRequiredOpts(settings, installationName, fe.Config().Setting == containerutil.FrontendPodmanShell, opts...)
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
	err = Start(ctx, console, image, containerName, installationName, fe, settings, true)
	if err != nil {
		return err
	}
	_, _, err = WaitUntilStarted(ctx, console, containerName, settings.VolumeName, settings.BuildkitAddress, settings.Timeout, fe, opts...)
	if err != nil {
		return err
	}
	console.
		WithPrefix("buildkitd").
		Printf("... Done")
	return nil
}

// maybeStart ensures that the buildkitd daemon is started. It returns the URL
// that can be used to connect to it.
func maybeStart(ctx context.Context, console conslogging.ConsoleLogger, image, containerName, installationName string, fe containerutil.ContainerFrontend, settings Settings, opts ...client.ClientOpt) (cinfo *client.Info, winfo *client.WorkerInfo, finalErr error) {
	if settings.StartUpLockPath != "" {
		var tryLockDone atomic.Bool
		go func() {
			time.Sleep(3 * time.Second)
			if !tryLockDone.Load() {
				console.Warnf("waiting on other instance of earthly to start buildkitd (as indicated by %q existing)", settings.StartUpLockPath)
			}
		}()
		startLock := flock.New(settings.StartUpLockPath)
		timeoutCtx, cancel := context.WithTimeout(ctx, opTimeout)
		defer cancel()
		_, err := startLock.TryLockContext(timeoutCtx, 200*time.Millisecond)
		tryLockDone.Store(true)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return nil, nil, errors.Errorf("timeout waiting for other instance of earthly to start buildkitd")
			}
			return nil, nil, errors.Wrapf(err, "try flock context %s", settings.StartUpLockPath)
		}
		defer func() {
			err := startLock.Unlock()
			if err != nil {
				console.Warnf("Failed to unlock %s: %v", settings.StartUpLockPath, err)
				if finalErr == nil {
					finalErr = err
				}
				return
			}
		}()
	}
	isStarted, err := IsStarted(ctx, containerName, fe)
	if err != nil {
		return nil, nil, errors.Wrap(err, "check is started buildkitd")
	}
	if isStarted {
		console.
			WithPrefix("buildkitd").
			Printf("Found buildkit daemon as %s container (%s)\n", fe.Config().Binary, containerName)
		info, workerInfo, err := maybeRestart(ctx, console, image, containerName, installationName, fe, settings, opts...)
		if err != nil {
			return nil, nil, errors.Wrap(err, "maybe restart")
		}
		return info, workerInfo, nil
	}
	console.
		WithPrefix("buildkitd").
		Printf("Starting buildkit daemon as a %s container (%s)...\n", fe.Config().Binary, containerName)
	err = Start(ctx, console, image, containerName, installationName, fe, settings, false)
	if err != nil {
		return nil, nil, errors.Wrap(err, "start")
	}
	info, workerInfo, err := WaitUntilStarted(ctx, console, containerName, settings.VolumeName, settings.BuildkitAddress, settings.Timeout, fe, opts...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "wait until started")
	}

	// check arch is correct
	runningContainerInfo, err := GetContainerInfo(ctx, containerName, fe)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetContainerInfo %s", containerName)
	}
	currentImageInfo, err := GetImageInfo(ctx, runningContainerInfo.Image, fe)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetImageInfo %s", runningContainerInfo.Image)
	}
	if currentImageInfo.Architecture != runtime.GOARCH {
		console.
			WithPrefix("buildkitd").
			Warnf("Warning: %s was started using architecture %s, but host architecture is %s; is DOCKER_DEFAULT_PLATFORM accidentally set?\n", containerName, currentImageInfo.Architecture, runtime.GOARCH)
	}

	console.
		WithPrefix("buildkitd").
		Printf("...Done\n")
	return info, workerInfo, nil
}

// maybeRestart checks whether the there is a different buildkitd image available locally or if
// settings of the current container are different from the provided settings. In either case,
// the container is restarted.
func maybeRestart(ctx context.Context, console conslogging.ConsoleLogger, image, containerName, installationName string, fe containerutil.ContainerFrontend, settings Settings, opts ...client.ClientOpt) (*client.Info, *client.WorkerInfo, error) {
	bkCons := console.WithPrefix("buildkitd")
	runningContainerInfo, err := GetContainerInfo(ctx, containerName, fe)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get container info")
	}
	currentImageInfo, err := GetImageInfo(ctx, runningContainerInfo.Image, fe)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get image info")
	}
	if currentImageInfo.Architecture != runtime.GOARCH {
		console.
			WithPrefix("buildkitd").
			Warnf("Warning: currently running %s under architecture %s, but host architecture is %s; is DOCKER_DEFAULT_PLATFORM accidentally set?\n", containerName, currentImageInfo.Architecture, runtime.GOARCH)
	}
	containerImageID := runningContainerInfo.ImageID
	availableImageID, err := GetAvailableImageID(ctx, image, fe)
	if err != nil {
		// Could not get available image ID. This happens when a new image tag is given and that
		// tag has not yet been pulled locally. Restarting will cause that tag to be pulled.
		availableImageID = "" // Will cause equality to fail and force a restart.
		// Keep going anyway.
	}
	bkCons.VerbosePrintf("Comparing running container %q image (%q) with available image %q (%q)\n", containerName, containerImageID, image, availableImageID)
	if containerImageID == availableImageID {
		// Images are the same. Check settings hash.
		hash, err := GetSettingsHash(ctx, containerName, fe)
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not get settings hash")
		}
		ok, err := settings.VerifyHash(hash)
		if err != nil {
			return nil, nil, errors.Wrap(err, "verify hash")
		}
		if ok {
			// No need to replace: images are the same and settings are the same.
			bkCons.VerbosePrintf("Settings hashes match (%q), no restart required\n", hash)
			info, workerInfo, err := checkConnection(ctx, settings.BuildkitAddress, 5*time.Second, opts...)
			if err != nil {
				return nil, nil, errors.Wrap(err, "could not connect to buildkitd to shut down container")
			}
			return info, workerInfo, nil
		}
		bkCons.Printf("Settings do not match. Restarting buildkit daemon with updated settings...\n")
	} else {
		if settings.NoUpdate {
			bkCons.Printf("Updated image available. But update was inhibited.\n")
			info, workerInfo, err := checkConnection(ctx, settings.BuildkitAddress, 5*time.Second, opts...)
			if err != nil {
				return nil, nil, errors.Wrap(err, "could not verify connection to buildkitd container")
			}
			return info, workerInfo, nil
		}
		bkCons.Printf("Updated image available. Restarting buildkit daemon...\n")
	}

	// Replace.
	err = Stop(ctx, containerName, fe)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not shut down container %q", containerName)
	}
	err = WaitUntilStopped(ctx, containerName, settings.Timeout, fe)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not wait for container %q to stop", containerName)
	}
	err = Start(ctx, console, image, containerName, installationName, fe, settings, false)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not start container %q", containerName)
	}
	info, workerInfo, err := WaitUntilStarted(ctx, console, containerName, settings.VolumeName, settings.BuildkitAddress, settings.Timeout, fe, opts...)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not wait for container %q to start", containerName)
	}
	bkCons.Printf("...Done\n")
	return info, workerInfo, nil
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
func Start(ctx context.Context, console conslogging.ConsoleLogger, image, containerName, installationName string, fe containerutil.ContainerFrontend, settings Settings, reset bool) error {
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
			return errors.Wrap(err, "error parsing buildkit address url")
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
			if settings.EnableProfiler {
				portOpts = append(portOpts, containerutil.Port{
					IP:            "127.0.0.1",
					HostPort:      6061, // 6060 is reserved for earthly client
					ContainerPort: 6060,
					Protocol:      containerutil.ProtocolTCP,
				})
			}
			if settings.UseTLS {
				if settings.TLSCA != "" {
					if exists, _ := fileutil.FileExists(settings.TLSCA); !exists {
						return errors.Wrapf(os.ErrNotExist, "TLS CA file %q is missing", settings.TLSCA)
					}
					volumeOpts = append(volumeOpts, containerutil.Mount{
						Type:     containerutil.MountBind,
						Source:   settings.TLSCA,
						Dest:     "/etc/ca.pem",
						ReadOnly: true,
					})
				}

				if settings.ServerTLSCert != "" {
					if exists, _ := fileutil.FileExists(settings.ServerTLSCert); !exists {
						return errors.Wrapf(os.ErrNotExist, "TLS certificate %q is missing", settings.ServerTLSCert)
					}
					volumeOpts = append(volumeOpts, containerutil.Mount{
						Type:     containerutil.MountBind,
						Source:   settings.ServerTLSCert,
						Dest:     "/etc/cert.pem",
						ReadOnly: true,
					})
				}

				if settings.ServerTLSKey != "" {
					if exists, _ := fileutil.FileExists(settings.ServerTLSKey); !exists {
						return errors.Wrapf(os.ErrNotExist, "TLS private key %q is missing", settings.ServerTLSKey)
					}
					volumeOpts = append(volumeOpts, containerutil.Mount{
						Type:     containerutil.MountBind,
						Source:   settings.ServerTLSKey,
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

	if settings.CacheKeepDuration > 0 {
		envOpts["CACHE_KEEP_DURATION"] = strconv.FormatInt(int64(settings.CacheKeepDuration), 10)
	}

	if settings.EnableProfiler {
		envOpts["BUILDKIT_PPROF_ENABLED"] = strconv.FormatBool(true)
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
func WaitUntilStarted(ctx context.Context, console conslogging.ConsoleLogger, containerName, volumeName, address string, opTimeout time.Duration, fe containerutil.ContainerFrontend, opts ...client.ClientOpt) (*client.Info, *client.WorkerInfo, error) {
	// First, wait for the container to be marked as started.
	ctxTimeout, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
ContainerRunningLoop:
	for {
		select {
		case <-time.After(200 * time.Millisecond):
			isRunning, err := isContainerRunning(ctxTimeout, containerName, fe)
			if err != nil {
				// Has not yet started. Keep waiting.
				continue
			}
			if !isRunning {
				return nil, nil, ErrBuildkitCrashed
			}
			if isRunning {
				break ContainerRunningLoop
			}

		case <-ctxTimeout.Done():
			return nil, nil, errors.Errorf("timeout %s: buildkitd container did not start", opTimeout)
		}
	}

	// Wait for the connection to be available.
	info, workerInfo, err := waitForConnection(ctx, containerName, address, opTimeout, fe, opts...)
	if err != nil {
		if !errors.Is(err, ErrBuildkitConnectionFailure) {
			return nil, nil, err
		}
		// We timed out. Check if the user has a lot of cache and give buildkit another chance.
		cacheSizeBytes, cacheSizeErr := getCacheSize(ctx, volumeName, fe)
		if cacheSizeErr != nil {
			console.
				WithPrefix("buildkitd").
				Printf("Warning: Could not detect buildkit cache size: %v\n", cacheSizeErr)
			return nil, nil, err
		}
		cacheGigs := cacheSizeBytes / 1024 / 1024 / 1024
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
			info, workerInfo, err := waitForConnection(ctx, containerName, address, opTimeout, fe, opts...)
			if err != nil {
				return nil, nil, err
			}
			return info, workerInfo, nil
		}
		return nil, nil, err
	}
	return info, workerInfo, nil
}

func waitForConnection(ctx context.Context, containerName, address string, opTimeout time.Duration, fe containerutil.ContainerFrontend, opts ...client.ClientOpt) (*client.Info, *client.WorkerInfo, error) {
	retryInterval := 200 * time.Millisecond
	if !containerutil.IsLocal(address) {
		retryInterval = 1 * time.Second
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	attemptTimeout := 500 * time.Millisecond
	if !containerutil.IsLocal(address) {
		attemptTimeout = 1 * time.Second
	}
	for {
		select {
		case <-time.After(retryInterval):
			if containerutil.IsLocal(address) {
				// Make sure that our managed buildkit has not crashed on startup.
				isRunning, err := isContainerRunning(ctxTimeout, containerName, fe)
				if err != nil {
					return nil, nil, err
				}
				if !isRunning {
					return nil, nil, ErrBuildkitCrashed
				}
			}

			info, workerInfo, err := checkConnection(ctxTimeout, address, attemptTimeout, opts...)
			if err != nil {
				// Try again.
				attemptTimeout *= 2
				// keep timeout reasonable
				if attemptTimeout > opTimeout {
					attemptTimeout = opTimeout
				}
				continue
			}
			return info, workerInfo, nil
		case <-ctxTimeout.Done():
			// Try one last time.
			info, workerInfo, err := checkConnection(ctx, address, attemptTimeout, opts...)
			if err != nil {
				// We give up.
				return nil, nil, errors.Wrapf(ErrBuildkitConnectionFailure, "timeout %s: could not connect to buildkit: %s", opTimeout, err.Error())
			}
			return info, workerInfo, nil
		}
	}
}

func checkConnection(ctx context.Context, address string, timeout time.Duration, opts ...client.ClientOpt) (*client.Info, *client.WorkerInfo, error) {
	// Each attempt has limited time to succeed, to prevent hanging for too long
	// here.
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	var (
		mu         sync.Mutex // protects the vars below
		connErr    error      = errors.New("timeout")
		info       *client.Info
		workerInfo *client.WorkerInfo
	)
	go func() {
		defer cancel()
		bkClient, err := client.New(ctxTimeout, address, opts...)
		if err != nil {
			mu.Lock()
			connErr = err
			mu.Unlock()
			return
		}
		defer bkClient.Close()
		// Use ListWorkers for backwards compatibility. (Info is relatively new)
		ws, err := bkClient.ListWorkers(ctxTimeout)
		if err != nil {
			mu.Lock()
			connErr = err
			mu.Unlock()
			return
		}
		if len(ws) == 0 {
			mu.Lock()
			connErr = errors.New("no workers")
			mu.Unlock()
			return
		}

		// Success.
		mu.Lock()
		defer mu.Unlock()
		connErr = nil
		workerInfo = ws[0]
		info, err = bkClient.Info(ctxTimeout)
		if err != nil {
			s, ok := status.FromError(errors.Cause(err))
			if ok && s.Code() == codes.Unimplemented {
				// Degrade gracefully.
				info = &client.Info{
					BuildkitVersion: client.BuildkitVersion{
						Version:  "unknown",
						Package:  "unknown",
						Revision: "unknown",
					},
				}
			} else {
				connErr = err
				return
			}
		}
	}()
	<-ctxTimeout.Done() // timeout or goroutine finished
	mu.Lock()
	defer mu.Unlock()
	if connErr != nil {
		return nil, nil, connErr
	}
	return info, workerInfo, nil
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
		case <-time.After(200 * time.Millisecond):
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

// GetContainerInfo inspects the running container (running under containerName)
func GetContainerInfo(ctx context.Context, containerName string, fe containerutil.ContainerFrontend) (*containerutil.ContainerInfo, error) {
	infos, err := fe.ContainerInfo(ctx, containerName)
	if err != nil {
		return nil, errors.Wrap(err, "get container info for current container image ID")
	}

	if containerInfo, ok := infos[containerName]; ok {
		return containerInfo, nil
	}

	return nil, fmt.Errorf("info for container %s was not found", containerName)
}

// GetImageInfo inspects an image
func GetImageInfo(ctx context.Context, image string, fe containerutil.ContainerFrontend) (*containerutil.ImageInfo, error) {
	infos, err := fe.ImageInfo(ctx, image)
	if err != nil {
		return nil, errors.Wrapf(err, "get image info %s", image)
	}

	if info, ok := infos[image]; ok {
		return info, nil
	}

	return nil, fmt.Errorf("info for image %s was not found", image)
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
	return fe.IsAvailable(ctx)
}

func printBuildkitInfo(bkCons conslogging.ConsoleLogger, info *client.Info, workerInfo *client.WorkerInfo, earthlyVersion string, isLocal, hasConfiguredCacheSize bool) {
	// Print most of this stuff only for remote buildkits / satellites.
	printFun := bkCons.Printf
	if isLocal {
		printFun = bkCons.VerbosePrintf
	}
	if info.BuildkitVersion.Version != "unknown" {
		printFun(
			"Version %s %s %s",
			info.BuildkitVersion.Package, info.BuildkitVersion.Version, info.BuildkitVersion.Revision)
		if info.BuildkitVersion.Package != "github.com/earthly/buildkit" {
			bkCons.Warnf("Using a non-Earthly version of Buildkit. This is not supported.")
		} else {
			if info.BuildkitVersion.Version != earthlyVersion {
				if isLocal {
					// For local buildkits we expect perfect version match.
					bkCons.Warnf(
						"Warning: Buildkit version (%s) is different from Earthly version (%s)",
						info.BuildkitVersion.Version, earthlyVersion)
				} else {
					compatible := true
					bkVersion, err := semverutil.Parse(info.BuildkitVersion.Version)
					if err != nil {
						bkCons.VerbosePrintf("Warning: could not parse buildkit version: %v", err)
						compatible = false
					}
					earthlyVersion, err := semverutil.Parse(earthlyVersion)
					if err != nil {
						bkCons.VerbosePrintf("Warning: could not parse earthly version: %v", err)
						compatible = false
					}
					compatible = compatible && semverutil.IsCompatible(bkVersion, earthlyVersion)
					if !compatible {
						bkCons.Warnf("Warning: Buildkit version (%s) is not compatible with Earthly version (%s)",
							info.BuildkitVersion.Version, earthlyVersion)
					} else {
						bkCons.VerbosePrintf("Buildkit version (%s) is compatible with Earthly version (%s)",
							info.BuildkitVersion.Version, earthlyVersion)
					}
				}
			}
		}
	} else {
		bkCons.Warnf(
			"Warning: Buildkit version is unknown. This usually means that " +
				"it's from a version lower than Earthly Buildkit v0.6.20")
	}
	ps := make([]string, len(workerInfo.Platforms))
	for i, p := range workerInfo.Platforms {
		ps[i] = platforms.Format(p)
	}
	if len(ps) > 0 {
		printFun("Platforms: %s (native) %s", ps[0], strings.Join(ps[1:], " "))
	}
	load := workerInfo.ParallelismCurrent + workerInfo.ParallelismWaiting
	printFun(buildkitutil.FormatUtilization(info.NumSessions, load, workerInfo.ParallelismMax))
	switch {
	case workerInfo.ParallelismWaiting > 5:
		bkCons.Warnf("Warning: Currently under heavy load. Performance will be affected")
	case workerInfo.ParallelismWaiting > 0:
		bkCons.Printf("Note: Currently under significant load. Performance will be affected")
	default:
	}
	ld := time.Duration(0)
	if workerInfo.GCAnalytics.LastEndTime != nil &&
		workerInfo.GCAnalytics.LastStartTime != nil {
		ld = workerInfo.GCAnalytics.LastEndTime.Sub(*workerInfo.GCAnalytics.LastStartTime)
	}
	printFun(
		"GC stats: %s cache, avg GC duration %v, all-time GC duration %v, last GC duration %v, last cleared %v",
		humanize.Bytes(uint64(workerInfo.GCAnalytics.LastSizeBefore)),
		workerInfo.GCAnalytics.AvgDuration,
		workerInfo.GCAnalytics.AllTimeDuration,
		ld,
		humanize.Bytes(uint64(workerInfo.GCAnalytics.LastSizeCleared)))
	if workerInfo.GCAnalytics.CurrentStartTime != nil {
		d := time.Since(*workerInfo.GCAnalytics.CurrentStartTime).Round(time.Second)
		switch {
		case d > 5*time.Minute:
			bkCons.Warnf("Warning: GC has been running for a long time, started %v ago", d)
		case d > 1*time.Minute:
			bkCons.Printf("GC currently ongoing, started %v ago", d)
		default:
		}
	}

	if isLocal && !hasConfiguredCacheSize {
		if size, ok := getGCPolicySize(workerInfo); ok && size < minRecommendedCacheSize {
			bkCons.Warnf("Configured cache size of %s is smaller than the minimum recommended size of %s",
				units.HumanSize(float64(size)), units.HumanSize(minRecommendedCacheSize))
			bkCons.Warnf("Please consider increasing the cache size: https://docs.earthly.dev/docs/caching/managing-cache")
		}
	}
}

func getGCPolicySize(workerInfo *client.WorkerInfo) (int64, bool) {
	for _, p := range workerInfo.GCPolicy {
		if p.All {
			return p.KeepBytes, true
		}
	}
	return 0, false
}

// getCacheSize returns the size of the earthly cache in bytes.
func getCacheSize(ctx context.Context, volumeName string, fe containerutil.ContainerFrontend) (int, error) {
	infos, err := fe.VolumeInfo(ctx, volumeName)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get volume info for cache size %s", volumeName)
	}

	return int(infos[volumeName].SizeBytes), nil
}

func addRequiredOpts(settings Settings, installationName string, isUsingPodman bool, opts ...client.ClientOpt) ([]client.ClientOpt, error) {
	server, err := url.Parse(settings.BuildkitAddress)
	if err != nil {
		return []client.ClientOpt{}, errors.Wrapf(err, "failed to parse buildkit url %s", settings.BuildkitAddress)
	}

	if settings.SatelliteName != "" {
		opts = append(opts,
			client.WithDefaultGRPCDialer(),
			client.WithAdditionalMetadataContext(
				"satellite_name", settings.SatelliteName,
				"satellite_org", settings.SatelliteOrgID,
				"satellite_token", settings.SatelliteToken),
		)
	}

	if !settings.UseTCP || !settings.UseTLS {
		return opts, nil
	}

	if settings.TLSCA == "" && settings.ClientTLSCert == "" && settings.ClientTLSKey == "" {
		return append(opts, client.WithServerConfigSystem("")), nil
	}

	opts = append(opts,
		client.WithCredentials(settings.ClientTLSCert, settings.ClientTLSKey),
		client.WithServerConfig(server.Hostname(), settings.TLSCA),
	)

	return opts, nil
}

// PrintSatelliteInfo prints the instance's details,
// including its Buildkit version, current workload, and garbage collection.
func PrintSatelliteInfo(ctx context.Context, console conslogging.ConsoleLogger, earthlyVersion string, settings Settings, installationName string) error {
	console.Printf("Connecting to %s...", settings.SatelliteDisplayName)
	opts, err := addRequiredOpts(settings, installationName, false)
	if err != nil {
		return errors.Wrap(err, "add required client opts")
	}
	info, workerInfo, err := waitForConnection(ctx, "", settings.BuildkitAddress, settings.Timeout, nil, opts...)
	if err != nil {
		return errors.Wrap(err, "connect provided buildkit")
	}
	printBuildkitInfo(console, info, workerInfo, earthlyVersion, false, false)
	return nil
}

func containsAny(hs string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(hs, n) {
			return true
		}
	}
	return false
}
