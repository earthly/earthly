package containerutil_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/containerutil"
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
			assert.NoError(t, err)
			assert.NotNil(t, fe)
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
			assert.NoError(t, err)

			scheme := fe.Scheme()
			assert.Equal(t, tC.scheme, scheme)
		})
	}
}

func TestFrontendIsAvaliable(t *testing.T) {
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
			assert.NoError(t, err)

			avaliable := fe.IsAvaliable(ctx)
			assert.True(t, avaliable)
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
			assert.NoError(t, err)

			info, err := fe.Information(ctx)
			assert.NoError(t, err)
			assert.NotNil(t, info)
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
			assert.NoError(t, err)
			defer cleanup()

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			assert.NoError(t, err)

			getInfos := append(testContainers, "missing")
			info, err := fe.ContainerInfo(ctx, getInfos...)
			assert.NoError(t, err)
			assert.NotNil(t, info)

			assert.Len(t, info, 3)

			assert.Equal(t, getInfos[0], info[getInfos[0]].Name)
			assert.Equal(t, "docker.io/library/nginx:1.21", info[getInfos[0]].Image)

			assert.Equal(t, getInfos[1], info[getInfos[1]].Name)
			assert.Equal(t, "docker.io/library/nginx:1.21", info[getInfos[1]].Image)

			assert.Equal(t, getInfos[2], info[getInfos[2]].Name)
			assert.Equal(t, containerutil.StatusMissing, info[getInfos[2]].Status)
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
			assert.NoError(t, err)
			defer cleanup()

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			assert.NoError(t, err)

			info, err := fe.ContainerInfo(ctx, testContainers...)
			assert.NoError(t, err)
			assert.Len(t, info, 2)

			err = fe.ContainerRemove(ctx, true, testContainers...)
			assert.NoError(t, err)

			info, err = fe.ContainerInfo(ctx, testContainers...)
			assert.NoError(t, err)
			assert.Equal(t, containerutil.StatusMissing, info[testContainers[0]].Status)
			assert.Equal(t, containerutil.StatusMissing, info[testContainers[1]].Status)
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
			assert.NoError(t, err)
			defer cleanup()

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			assert.NoError(t, err)

			info, err := fe.ContainerInfo(ctx, testContainers...)
			assert.NoError(t, err)
			assert.Len(t, info, 2)

			err = fe.ContainerStop(ctx, 0, testContainers...)
			assert.NoError(t, err)

			_, err = fe.ContainerInfo(ctx, testContainers...)
			assert.NoError(t, err)
			assert.Len(t, info, 2)
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
			assert.NoError(t, err)
			defer cleanup()

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			assert.NoError(t, err)

			logs, err := fe.ContainerLogs(ctx, testContainers...)
			assert.NoError(t, err)
			assert.Len(t, logs, 2)

			assert.Empty(t, logs[testContainers[0]].Stdout)
			assert.NotEmpty(t, logs[testContainers[0]].Stderr)

			assert.Empty(t, logs[testContainers[1]].Stdout)
			assert.NotEmpty(t, logs[testContainers[1]].Stderr)
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
			assert.NoError(t, err)

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
					cmd.Run() // Just best effort

					cmd = exec.CommandContext(ctx, tC.binary, "volume", "rm", "-f", fmt.Sprintf("vol-%s", name))
					cmd.Run()
				}
			}()

			info, err := fe.ContainerInfo(ctx, testContainers...)
			assert.NoError(t, err)
			assert.Equal(t, containerutil.StatusMissing, info[testContainers[0]].Status)
			assert.Equal(t, containerutil.StatusMissing, info[testContainers[1]].Status)

			err = fe.ContainerRun(ctx, runs...)
			assert.NoError(t, err)

			info, err = fe.ContainerInfo(ctx, testContainers...)
			assert.NoError(t, err)
			assert.Equal(t, containerutil.StatusRunning, info[testContainers[0]].Status)
			assert.Equal(t, containerutil.StatusRunning, info[testContainers[1]].Status)
		})
	}
}

func TestFrontendImagePull(t *testing.T) {
	testCases := []struct {
		binary  string
		newFunc func(context.Context, *containerutil.FrontendConfig) (containerutil.ContainerFrontend, error)
		refList []string
	}{
		{"docker", containerutil.NewDockerShellFrontend, []string{"nginx:1.21", "alpine:3.15"}},
		{"podman", containerutil.NewPodmanShellFrontend, []string{"docker.io/nginx:1.21", "docker.io/alpine:3.15"}}, // Podman prefers... and exports fully-qualified image tags
	}
	for _, tC := range testCases {
		t.Run(tC.binary, func(t *testing.T) {
			ctx := context.Background()
			onlyIfBinaryIsInstalled(ctx, t, tC.binary)

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{
				LocalRegistryHostFileValue: "tcp://some-host:5309", // podman pull needs some potentially valid address to check against, otherwise panic
				Console:                    testLogger(),
			})
			assert.NoError(t, err)

			err = fe.ImagePull(ctx, tC.refList...)
			assert.NoError(t, err)

			defer func() {
				for _, ref := range tC.refList {
					cmd := exec.CommandContext(ctx, "docker", "image", "rm", "-f", ref)
					cmd.Run()
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
			assert.NoError(t, err)
			defer cleanup()

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			assert.NoError(t, err)

			info, err := fe.ImageInfo(ctx, tC.refList...)
			assert.NoError(t, err)

			assert.Len(t, info, 2)

			assert.Contains(t, info[tC.refList[0]].Tags, tC.refList[0])
			assert.Contains(t, info[tC.refList[1]].Tags, tC.refList[1])
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
			assert.NoError(t, err)
			defer cleanup()

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			assert.NoError(t, err)

			info, err := fe.ImageInfo(ctx, refList...)
			assert.NoError(t, err)
			assert.Len(t, info, 2)

			err = fe.ImageRemove(ctx, true, refList...)
			assert.NoError(t, err)

			info, err = fe.ImageInfo(ctx, refList...)
			assert.NoError(t, err)
			assert.Empty(t, info[refList[0]].ID)
			assert.Empty(t, info[refList[1]].ID)
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
			assert.NoError(t, err)
			defer cleanup()

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			assert.NoError(t, err)

			info, err := fe.ImageInfo(ctx, ref)
			assert.NoError(t, err)

			imageID := info[ref].ID
			tags := []containerutil.ImageTag{}
			for _, tagName := range tC.tagList {
				tags = append(tags, containerutil.ImageTag{
					SourceRef: imageID,
					TargetRef: tagName,
				})
			}

			err = fe.ImageTag(ctx, tags...)
			assert.NoError(t, err)

			info, err = fe.ImageInfo(ctx, tC.tagList...)
			assert.NoError(t, err)

			assert.Contains(t, info[tC.tagList[0]].Tags, tC.tagList[0])
			assert.Contains(t, info[tC.tagList[1]].Tags, tC.tagList[1])
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
			assert.NoError(t, err)

			imgBuffer := &bytes.Buffer{}
			cmd := exec.CommandContext(ctx, tC.binary, "image", "save", tC.ref)
			cmd.Stdout = bufio.NewWriter(imgBuffer)
			err = cmd.Run()
			assert.NoError(t, err)

			cleanup()

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			assert.NoError(t, err)

			err = fe.ImageLoad(ctx, bufio.NewReader(imgBuffer))
			assert.NoError(t, err)

			defer func() {
				cmd := exec.CommandContext(ctx, tC.binary, "image", "rm", "-f", tC.ref)
				cmd.Run()
			}()

			info, err := fe.ImageInfo(ctx, tC.ref)
			assert.NoError(t, err)
			assert.Contains(t, info[tC.ref].Tags, tC.ref)
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
			assert.NoError(t, err)

			data, err := os.ReadFile("./testdata/hybrid.tar")
			assert.NoError(t, err)
			reader := bytes.NewReader(data)

			err = fe.ImageLoad(ctx, reader)
			assert.NoError(t, err)

			defer func() {
				cmd := exec.CommandContext(ctx, tC.binary, "image", "rm", "-f", tC.ref)
				cmd.Run()
			}()

			info, err := fe.ImageInfo(ctx, tC.ref)
			assert.NoError(t, err)
			assert.Contains(t, info[tC.ref].Tags, tC.ref)
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
			assert.NoError(t, err)
			defer cleanup()

			fe, err := tC.newFunc(ctx, &containerutil.FrontendConfig{Console: testLogger()})
			assert.NoError(t, err)

			info, err := fe.VolumeInfo(ctx, volList...)
			assert.NoError(t, err)
			assert.Len(t, info, 2)
		})
	}
}

func onlyIfBinaryIsInstalled(ctx context.Context, t *testing.T, binary string) {
	if !isBinaryInstalled(ctx, binary) {
		t.Skipf("%s is not avaliable for tests, skipping", binary)
	}
}

func isBinaryInstalled(ctx context.Context, binary string) bool {
	// This is almost a re-implementation of IsAvaliable... but relying on that presupposes the
	// binary exists to allow the New**** to run (it gathers info from the CLI).
	cmd := exec.CommandContext(ctx, binary, "--help")
	return cmd.Run() == nil
}

func spawnTestContainers(ctx context.Context, feBinary string, names ...string) (func(), error) {
	var err error
	for _, name := range names {
		cmd := exec.CommandContext(ctx, feBinary, "run", "-d", "--name", name, "docker.io/library/nginx:1.21", `-text="test"`)
		output, createErr := cmd.CombinedOutput()
		if err != nil {
			// the frontend exists but is non-functional. This is... not likely to work at all.
			multierror.Append(err, errors.Wrap(createErr, string(output)))
		}
	}

	return func() {
		for _, name := range names {
			cmd := exec.CommandContext(ctx, feBinary, "rm", "-f", name)
			cmd.Run() // Just best effort
		}
	}, nil
}

func spawnTestImages(ctx context.Context, feBinary string, refs ...string) (func(), error) {
	var err error
	for _, ref := range refs {
		cmd := exec.CommandContext(ctx, feBinary, "image", "pull", "docker.io/nginx:1.21")
		output, createErr := cmd.CombinedOutput()
		if err != nil {
			// the frontend exists but is non-functional. This is... not likely to work at all.
			multierror.Append(err, errors.Wrap(createErr, string(output)))
			break
		}

		cmd = exec.CommandContext(ctx, feBinary, "image", "tag", "docker.io/nginx:1.21", ref)
		output, tagErr := cmd.CombinedOutput()
		if err != nil {
			// the frontend exists but is non-functional. This is... not likely to work at all.
			multierror.Append(err, errors.Wrap(tagErr, string(output)))
			break
		}
	}

	return func() {
		for _, ref := range refs {
			cmd := exec.CommandContext(ctx, feBinary, "image", "rm", "-f", ref)
			cmd.Run() // Just best effort
		}
	}, nil
}

func spawnTestVolumes(ctx context.Context, feBinary string, names ...string) (func(), error) {
	var err error
	for _, name := range names {
		cmd := exec.CommandContext(ctx, feBinary, "volume", "create", name)
		output, createErr := cmd.CombinedOutput()
		if err != nil {
			// the frontend exists but is non-functional. This is... not likely to work at all.
			multierror.Append(err, errors.Wrap(createErr, string(output)))
		}
	}

	return func() {
		for _, name := range names {
			cmd := exec.CommandContext(ctx, feBinary, "volume", "rm", "-f", name)
			cmd.Run() // Just best effort
		}
	}, nil
}

func testLogger() conslogging.ConsoleLogger {
	var logs strings.Builder
	logger := conslogging.Current(conslogging.NoColor, conslogging.DefaultPadding, false)
	return logger.WithWriter(&logs)
}
