package containerutil_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/containerutil"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

func TestFrontendNew(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
	}{
		{"docker", containerutil.NewDockerShellFrontend},
		{"podman", containerutil.NewPodmanShellFrontend},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)
			NotNil(t, fe)
		})
	}
}

func TestFrontendScheme(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
		scheme  string
	}{
		{"docker", containerutil.NewDockerShellFrontend, "docker-container"},
		{"podman", containerutil.NewPodmanShellFrontend, "podman-container"},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			scheme := fe.Scheme()
			Equal(t, tC.scheme, scheme)
		})
	}
}

func TestFrontendIsAvailable(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
	}{
		{"docker", containerutil.NewDockerShellFrontend},
		{"podman", containerutil.NewPodmanShellFrontend},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			available := fe.IsAvailable(ctx)
			True(t, available)
		})
	}
}

func TestFrontendInformation(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
	}{
		{"docker", containerutil.NewDockerShellFrontend},
		{"podman", containerutil.NewPodmanShellFrontend},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			info, err := fe.Information(ctx)
			NoError(t, err)
			NotNil(t, info)
		})
	}
}

func TestFrontendContainerInfo(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
	}{
		{"docker", containerutil.NewDockerShellFrontend},
		{"podman", containerutil.NewPodmanShellFrontend},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			testContainers := []string{"test-1", "test-2"}
			cleanup, err := spawnTestContainers(ctx, tC.binary, testContainers...)
			t.Cleanup(cleanup)
			NoError(t, err)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			getInfos := append(testContainers, "missing")
			info, err := fe.ContainerInfo(ctx, getInfos...)
			NoError(t, err)
			NotNil(t, info)

			Len(t, info, 3)

			Equal(t, getInfos[0], info[getInfos[0]].Name)
			Equal(t, "docker.io/library/nginx:1.21", info[getInfos[0]].Image)

			Equal(t, getInfos[1], info[getInfos[1]].Name)
			Equal(t, "docker.io/library/nginx:1.21", info[getInfos[1]].Image)

			Equal(t, getInfos[2], info[getInfos[2]].Name)
			Equal(t, containerutil.StatusMissing, info[getInfos[2]].Status)
		})
	}
}

func TestFrontendContainerRemove(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
	}{
		{"docker", containerutil.NewDockerShellFrontend},
		{"podman", containerutil.NewPodmanShellFrontend},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			testContainers := []string{"remove-1", "remove-2"}
			cleanup, err := spawnTestContainers(ctx, tC.binary, testContainers...)
			t.Cleanup(cleanup)
			NoError(t, err)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			info, err := fe.ContainerInfo(ctx, testContainers...)
			NoError(t, err)
			Len(t, info, 2)

			err = fe.ContainerRemove(ctx, true, testContainers...)
			NoError(t, err)

			info, err = fe.ContainerInfo(ctx, testContainers...)
			NoError(t, err)
			Equal(t, containerutil.StatusMissing, info[testContainers[0]].Status)
			Equal(t, containerutil.StatusMissing, info[testContainers[1]].Status)
		})
	}
}

func TestFrontendContainerStop(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
	}{
		{"docker", containerutil.NewDockerShellFrontend},
		{"podman", containerutil.NewPodmanShellFrontend},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			testContainers := []string{"stop-1", "stop-2"}
			cleanup, err := spawnTestContainers(ctx, tC.binary, testContainers...)
			t.Cleanup(cleanup)
			NoError(t, err)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			info, err := fe.ContainerInfo(ctx, testContainers...)
			NoError(t, err)
			Len(t, info, 2)

			err = fe.ContainerStop(ctx, 0, testContainers...)
			NoError(t, err)

			_, err = fe.ContainerInfo(ctx, testContainers...)
			NoError(t, err)
			Len(t, info, 2)
		})
	}
}

func TestFrontendContainerLogs(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
	}{
		{"docker", containerutil.NewDockerShellFrontend},
		{"podman", containerutil.NewPodmanShellFrontend},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			testContainers := []string{"logs-1", "logs-2"}
			cleanup, err := spawnTestContainers(ctx, tC.binary, testContainers...)
			t.Cleanup(cleanup)
			NoError(t, err)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			logs, err := fe.ContainerLogs(ctx, testContainers...)
			NoError(t, err)
			Len(t, logs, 2)

			Equal(t, "output stream\n", logs[testContainers[0]].Stdout)
			Equal(t, "error stream\n", logs[testContainers[0]].Stderr)

			Equal(t, "output stream\n", logs[testContainers[1]].Stdout)
			Equal(t, "error stream\n", logs[testContainers[1]].Stderr)
		})
	}
}

func TestFrontendContainerRun(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
	}{
		{"docker", containerutil.NewDockerShellFrontend},
		{"podman", containerutil.NewPodmanShellFrontend},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			testContainers := []string{"create-1", "create-2"}
			runs := []containerutil.ContainerRun{}
			for _, name := range testContainers {
				runs = append(runs, containerutil.ContainerRun{
					NameOrID:       name,
					ImageRef:       "docker.io/nginx:1.21",
					Privileged:     false,
					Envs:           containerutil.EnvMap{"test": name},
					Labels:         containerutil.LabelMap{"test": name},
					ContainerArgs:  []string{"nginx-debug", "-g", "daemon off;"},
					AdditionalArgs: []string{"--rm"},
					Mounts: containerutil.MountOpt{
						containerutil.Mount{
							Type:     containerutil.MountVolume,
							Source:   fmt.Sprintf("vol-%s", name),
							Dest:     "/test",
							ReadOnly: true,
						},
					},
					Ports: containerutil.PortOpt{
						containerutil.Port{
							IP:            "127.0.0.1",
							HostPort:      0,
							ContainerPort: 5678,
							Protocol:      containerutil.ProtocolTCP,
						},
					},
				})
			}

			defer func() {
				for _, name := range testContainers {
					// Roll our own cleanup since we can't use the spawn test containers helper... since
					// the whole point of this test is to create them with a frontend. Also theres a volume
					cmd := exec.CommandContext(ctx, tC.binary, "rm", "-f", name)
					_ = cmd.Run() // Just best effort

					cmd = exec.CommandContext(ctx, tC.binary, "volume", "rm", "-f", fmt.Sprintf("vol-%s", name))
					_ = cmd.Run()
				}
			}()

			info, err := fe.ContainerInfo(ctx, testContainers...)
			NoError(t, err)
			Equal(t, containerutil.StatusMissing, info[testContainers[0]].Status)
			Equal(t, containerutil.StatusMissing, info[testContainers[1]].Status)

			err = fe.ContainerRun(ctx, runs...)
			NoError(t, err)

			info, err = fe.ContainerInfo(ctx, testContainers...)
			NoError(t, err)
			Equal(t, containerutil.StatusRunning, info[testContainers[0]].Status)
			Equal(t, containerutil.StatusRunning, info[testContainers[1]].Status)
		})
	}
}

func TestFrontendImagePull(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
		refList []string
	}{
		{"docker", containerutil.NewDockerShellFrontend, []string{"nginx:1.21", "alpine:3.18"}},
		{"podman", containerutil.NewPodmanShellFrontend, []string{"docker.io/nginx:1.21", "docker.io/alpine:3.18"}}, // Podman prefers... and exports fully-qualified image tags
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{
				LocalRegistryHostFileValue: "tcp://some-host:5309", // podman pull needs some potentially valid address to check against, otherwise panic
				Console:                    testLogger(),
			})
			NoError(t, err)

			err = fe.ImagePull(ctx, tC.refList...)
			NoError(t, err)

			defer func() {
				for _, ref := range tC.refList {
					cmd := exec.CommandContext(ctx, "docker", "image", "rm", "-f", ref)
					_ = cmd.Run()
				}
			}()
		})
	}
}

func TestFrontendImageInfo(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
		refList []string
	}{
		{"docker", containerutil.NewDockerShellFrontend, []string{"info:1", "info:2"}},
		{"podman", containerutil.NewPodmanShellFrontend, []string{"localhost/info:1", "localhost/info:2"}},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			cleanup, err := spawnTestImages(ctx, tC.binary, tC.refList...)
			NoError(t, err)
			t.Cleanup(cleanup)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			info, err := fe.ImageInfo(ctx, tC.refList...)
			NoError(t, err)

			Len(t, info, 2)

			Contains(t, info[tC.refList[0]].Tags, tC.refList[0])
			Contains(t, info[tC.refList[1]].Tags, tC.refList[1])
		})
	}
}

func TestFrontendImageRemove(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
	}{
		{"docker", containerutil.NewDockerShellFrontend},
		{"podman", containerutil.NewPodmanShellFrontend},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			refList := []string{"remove:1", "remove:2"}
			cleanup, err := spawnTestImages(ctx, tC.binary, refList...)
			NoError(t, err)
			t.Cleanup(cleanup)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			info, err := fe.ImageInfo(ctx, refList...)
			NoError(t, err)
			Len(t, info, 2)

			err = fe.ImageRemove(ctx, true, refList...)
			NoError(t, err)

			info, err = fe.ImageInfo(ctx, refList...)
			NoError(t, err)
			Empty(t, info[refList[0]].ID)
			Empty(t, info[refList[1]].ID)
		})
	}
}

func TestFrontendImageTag(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
		tagList []string
	}{
		{"docker", containerutil.NewDockerShellFrontend, []string{"tag:1", "tag:2"}},
		{"podman", containerutil.NewPodmanShellFrontend, []string{"localhost/tag:1", "localhost/tag:2"}},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			ref := "tag:me"
			cleanup, err := spawnTestImages(ctx, tC.binary, ref)
			NoError(t, err)
			t.Cleanup(cleanup)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			info, err := fe.ImageInfo(ctx, ref)
			NoError(t, err)

			imageID := info[ref].ID
			tags := []containerutil.ImageTag{}
			for _, tagName := range tC.tagList {
				tags = append(tags, containerutil.ImageTag{
					SourceRef: imageID,
					TargetRef: tagName,
				})
			}

			err = fe.ImageTag(ctx, tags...)
			NoError(t, err)

			info, err = fe.ImageInfo(ctx, tC.tagList...)
			NoError(t, err)

			Contains(t, info[tC.tagList[0]].Tags, tC.tagList[0])
			Contains(t, info[tC.tagList[1]].Tags, tC.tagList[1])
		})
	}
}

func TestFrontendImageLoad(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
		ref     string
	}{
		{"docker", containerutil.NewDockerShellFrontend, "load:me"},
		{"podman", containerutil.NewPodmanShellFrontend, "localhost/load:me"},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			cleanup, err := spawnTestImages(ctx, tC.binary, tC.ref)
			NoError(t, err)

			imgBuffer := &bytes.Buffer{}
			cmd := exec.CommandContext(ctx, tC.binary, "image", "save", tC.ref)
			cmd.Stdout = bufio.NewWriter(imgBuffer)
			err = cmd.Run()
			NoError(t, err)

			cleanup()

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			err = fe.ImageLoad(ctx, bufio.NewReader(imgBuffer))
			NoError(t, err)

			defer func() {
				cmd := exec.CommandContext(ctx, tC.binary, "image", "rm", "-f", tC.ref)
				_ = cmd.Run()
			}()

			info, err := fe.ImageInfo(ctx, tC.ref)
			NoError(t, err)
			Contains(t, info[tC.ref].Tags, tC.ref)
		})
	}
}

func TestFrontendImageLoadHybrid(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
		ref     string
	}{
		{"docker", containerutil.NewDockerShellFrontend, "hybrid:test"},
		{"podman", containerutil.NewPodmanShellFrontend, "localhost/hybrid:test"},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			data, err := os.ReadFile("./testdata/hybrid.tar")
			NoError(t, err)
			reader := bytes.NewReader(data)

			err = fe.ImageLoad(ctx, reader)
			NoError(t, err)

			defer func() {
				cmd := exec.CommandContext(ctx, tC.binary, "image", "rm", "-f", tC.ref)
				_ = cmd.Run()
			}()

			info, err := fe.ImageInfo(ctx, tC.ref)
			NoError(t, err)
			Contains(t, info[tC.ref].Tags, tC.ref)
		})
	}
}

func TestFrontendVolumeInfo(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
	}{
		{"docker", containerutil.NewDockerShellFrontend},
		{"podman", containerutil.NewPodmanShellFrontend},
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			volList := []string{"test1", "test2"}
			cleanup, err := spawnTestVolumes(ctx, tC.binary, volList...)
			NoError(t, err)
			t.Cleanup(cleanup)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			NoError(t, err)

			info, err := fe.VolumeInfo(ctx, volList...)
			NoError(t, err)
			Len(t, info, 2)
		})
	}
}

func onlyIfBinaryIsInstalled(ctx context.Context, t *testing.T, binary string) {
	if !isBinaryInstalled(ctx, binary) {
		t.Skipf("%s is not available for tests, skipping", binary)
	}
}

func isBinaryInstalled(ctx context.Context, binary string) bool {
	// This is almost a re-implementation of IsAvailable... but relying on that presupposes the
	// binary exists to allow the New**** to run (it gathers info from the CLI).
	cmd := exec.CommandContext(ctx, binary, "--help")
	return cmd.Run() == nil
}

// First, we attempt to stop the containers that may be running.
// Next, we start the containers
// After that, we attempt to wait for the containers for up to 20 seconds
// Then we return
// Caller is always expected to call cleanup
func spawnTestContainers(ctx context.Context, feBinary string, names ...string) (func(), error) {
	_ = removeContainers(ctx, feBinary, names...) // best effort
	err := startTestContainers(ctx, feBinary, names...)
	cleanup := func() {
		_ = removeContainers(ctx, feBinary, names...) // best-effort
	}
	if err != nil {
		return cleanup, err
	}
	err = waitForContainers(ctx, feBinary, names...)
	return cleanup, err
}

func startTestContainers(ctx context.Context, feBinary string, names ...string) error {
	var err error
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	image := "docker.io/library/nginx:1.21"
	pullErr := pullImageIfNecessary(ctx, feBinary, image)
	if pullErr != nil {
		return fmt.Errorf("failed to pull image %s: %w", image, pullErr)
	}
	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			cmd := exec.CommandContext(ctx, feBinary, "run", "-d", "--rm", "--name", name, image, "sh", "-c", `echo output stream&&>&2 echo error stream&&sleep 100`)
			output, createErr := cmd.CombinedOutput()
			m.Lock()
			defer m.Unlock()
			if createErr != nil {
				// the frontend exists but is non-functional. This is... not likely to work at all.
				err = multierror.Append(err, errors.Wrap(createErr, string(output)))
			}
		}(name)
	}
	wg.Wait()
	return err
}

// pullImageIfNecessary will only pull the image if it does not exist locally
// This helps us avoid unauthenticated rate limits in tests
func pullImageIfNecessary(ctx context.Context, feBinary string, image string) error {
	cmd := exec.CommandContext(ctx, feBinary, "inspect", "--type=image", image)
	_, inspectErr := cmd.CombinedOutput()
	if inspectErr == nil {
		// If we are able to inspect the image then it must exist locally
		return nil
	}
	cmd = exec.CommandContext(ctx, feBinary, "pull", image)
	_, pullErr := cmd.CombinedOutput()
	if pullErr != nil {
		return fmt.Errorf("failed to pull image %s: %w", image, pullErr)
	}
	return nil
}

func removeContainers(ctx context.Context, feBinary string, names ...string) error {
	var err error
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			removeCmd := exec.CommandContext(ctx, feBinary, "rm", "-f", name)
			_, removeErr := removeCmd.CombinedOutput()
			m.Lock()
			defer m.Unlock()
			if removeErr != nil {
				err = multierror.Append(err, fmt.Errorf("failed to remove container %s", name))
			}
		}(name)
	}
	wg.Wait()
	return err
}

func waitForContainers(ctx context.Context, feBinary string, names ...string) error {
	var err error
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, name := range names {
		const maxAttempts = 100
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			attempts := 0
			for attempts < maxAttempts {
				attempts++
				// docker inspect -f {{.State.Running}} CONTAINERNAME`"=="true"
				cmd := exec.CommandContext(ctx, feBinary, "inspect", "-f", "{{.State.Running}}", name)
				output, inspectErr := cmd.CombinedOutput()
				if inspectErr != nil {
					m.Lock()
					err = multierror.Append(err, inspectErr)
					m.Unlock()
					return
				}
				if strings.Contains(string(output), "true") {
					return
				}
				time.Sleep(time.Millisecond * 200)
			}
			m.Lock()
			defer m.Unlock()
			err = multierror.Append(err, fmt.Errorf("failed to wait for container %s to start", name))
		}(name)
	}
	wg.Wait()
	return err
}

func spawnTestImages(ctx context.Context, feBinary string, refs ...string) (func(), error) {
	var err error
	for _, ref := range refs {
		cmd := exec.CommandContext(ctx, feBinary, "image", "pull", "docker.io/nginx:1.21")
		output, createErr := cmd.CombinedOutput()
		if createErr != nil {
			// the frontend exists but is non-functional. This is... not likely to work at all.
			err = multierror.Append(err, errors.Wrap(createErr, string(output)))
			break
		}

		cmd = exec.CommandContext(ctx, feBinary, "image", "tag", "docker.io/nginx:1.21", ref)
		output, tagErr := cmd.CombinedOutput()
		if tagErr != nil {
			// the frontend exists but is non-functional. This is... not likely to work at all.
			err = multierror.Append(err, errors.Wrap(tagErr, string(output)))
			break
		}
	}

	return func() {
		for _, ref := range refs {
			cmd := exec.CommandContext(ctx, feBinary, "image", "rm", "-f", ref)
			_ = cmd.Run() // Just best effort
		}
	}, err
}

func spawnTestVolumes(ctx context.Context, feBinary string, names ...string) (func(), error) {
	var err error
	for _, name := range names {
		cmd := exec.CommandContext(ctx, feBinary, "volume", "create", name)
		output, createErr := cmd.CombinedOutput()
		if createErr != nil {
			// the frontend exists but is non-functional. This is... not likely to work at all.
			err = multierror.Append(err, errors.Wrap(createErr, string(output)))
		}
	}

	return func() {
		for _, name := range names {
			cmd := exec.CommandContext(ctx, feBinary, "volume", "rm", "-f", name)
			_ = cmd.Run() // Just best effort
		}
	}, nil
}

func testLogger() conslogging.ConsoleLogger {
	var logs strings.Builder
	logger := conslogging.Current(conslogging.NoColor, conslogging.DefaultPadding, conslogging.Info, false)
	return logger.WithWriter(&logs)
}
