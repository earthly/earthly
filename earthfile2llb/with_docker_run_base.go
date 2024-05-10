package earthfile2llb

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/containerd/containerd/platforms"
	debuggercommon "github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/platutil"
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	dockerdWrapperPath          = "/var/earthly/dockerd-wrapper.sh"
	dockerAutoInstallScriptPath = "/var/earthly/docker-auto-install.sh"
	composeConfigFile           = "compose-config.yml"
	suggestedDINDImage          = "earthly/dind:alpine-3.19-docker-25.0.5-r0"
)

// DockerLoadOpt holds parameters for WITH DOCKER --load parameter.
type DockerLoadOpt struct {
	Target          string
	ImageName       string
	Platform        platutil.Platform
	BuildArgs       []string
	AllowPrivileged bool
	PassArgs        bool
}

// DockerPullOpt holds parameters for the WITH DOCKER --pull parameter.
type DockerPullOpt struct {
	ImageName string
	Platform  platutil.Platform
}

// WithDockerOpt holds parameters for WITH DOCKER run.
type WithDockerOpt struct {
	Mounts                []string
	Secrets               []string
	WithShell             bool
	WithEntrypoint        bool
	WithSSH               bool
	NoCache               bool
	Interactive           bool
	interactiveKeep       bool
	Pulls                 []DockerPullOpt
	Loads                 []DockerLoadOpt
	ComposeFiles          []string
	ComposeServices       []string
	TryCatchSaveArtifacts []debuggercommon.SaveFilesSettings
	extraRunOpts          []llb.RunOption
	CacheID               string
}

type withDockerRunBase struct {
	c *Converter
}

func (w *withDockerRunBase) installDeps(ctx context.Context, opt WithDockerOpt) error {
	params := composeParams(opt)
	args := []string{
		"/bin/sh", "-c",
		fmt.Sprintf(
			"%s %s",
			strings.Join(params, " "),
			dockerAutoInstallScriptPath),
	}
	prefix, _, err := w.c.newVertexMeta(ctx, false, false, false, opt.Secrets)
	if err != nil {
		return err
	}
	runOpts := []llb.RunOption{
		llb.AddMount(
			dockerAutoInstallScriptPath, llb.Scratch(), llb.HostBind(), llb.SourcePath(dockerAutoInstallScriptPath)),
		llb.Args(args),
		llb.WithCustomNamef("%sWITH DOCKER (install deps)", prefix),
	}
	w.c.mts.Final.MainState = w.c.mts.Final.MainState.Run(runOpts...).Root()
	return nil
}

func (w *withDockerRunBase) getComposePulls(ctx context.Context, opt WithDockerOpt) ([]DockerPullOpt, error) {
	if len(opt.ComposeFiles) == 0 {
		// Quick way out. Compose not used.
		return nil, nil
	}
	// Get compose images from compose config.
	composeConfigDt, err := w.getComposeConfig(ctx, opt)
	if err != nil {
		return nil, err
	}
	type composeService struct {
		Image    string `yaml:"image"`
		Platform string `yaml:"platform"`
	}
	type composeData struct {
		Services map[string]composeService `yaml:"services"`
	}
	var config composeData
	err = yaml.Unmarshal(composeConfigDt, &config)
	if err != nil {
		return nil, errors.Wrapf(err, "parse compose config for %v", opt.ComposeFiles)
	}

	// Collect relevant images from the compose config.
	composeServicesSet := make(map[string]bool)
	for _, composeService := range opt.ComposeServices {
		composeServicesSet[composeService] = true
	}
	var pulls []DockerPullOpt
	for serviceName, serviceInfo := range config.Services {
		if serviceInfo.Image == "" {
			// Image not specified in yaml.
			continue
		}
		platform := w.c.platr.Current()
		if serviceInfo.Platform != "" {
			p, err := platforms.Parse(serviceInfo.Platform)
			if err != nil {
				return nil, errors.Wrapf(
					err, "parse platform for image %s: %s", serviceInfo.Image, serviceInfo.Platform)
			}
			platform = platutil.FromLLBPlatform(p)
		}
		if len(opt.ComposeServices) > 0 {
			if composeServicesSet[serviceName] {
				pulls = append(pulls, DockerPullOpt{
					ImageName: serviceInfo.Image,
					Platform:  platform,
				})
			}
		} else {
			// No services specified. Special case: collect all.
			pulls = append(pulls, DockerPullOpt{
				ImageName: serviceInfo.Image,
				Platform:  platform,
			})
		}
	}
	return pulls, nil
}

func (w *withDockerRunBase) getComposeConfig(ctx context.Context, opt WithDockerOpt) ([]byte, error) {
	// Add the right run to fetch the docker compose config.
	params := composeParams(opt)
	args := []string{
		"/bin/sh", "-c",
		fmt.Sprintf(
			"%s %s get-compose-config",
			strings.Join(params, " "),
			dockerdWrapperPath),
	}
	prefix, _, err := w.c.newVertexMeta(ctx, false, false, false, opt.Secrets)
	if err != nil {
		return nil, err
	}
	runOpts := []llb.RunOption{
		llb.AddMount(
			dockerdWrapperPath, llb.Scratch(), llb.HostBind(), llb.SourcePath(dockerdWrapperPath)),
		llb.Args(args),
		llb.WithCustomNamef("%sWITH DOCKER (docker-compose config)", prefix),
	}
	state := w.c.mts.Final.MainState.Run(runOpts...).Root()
	ref, err := llbutil.StateToRef(
		ctx, w.c.opt.GwClient, state, w.c.opt.NoCache,
		w.c.platr, w.c.opt.CacheImports.AsSlice())
	if err != nil {
		return nil, errors.Wrap(err, "state to ref compose config")
	}
	composeConfigDt, err := ref.ReadFile(ctx, gwclient.ReadRequest{
		Filename: fmt.Sprintf("/tmp/earthly/%s", composeConfigFile),
	})
	if err != nil {
		return nil, errors.Wrap(err, "read compose config file")
	}
	return composeConfigDt, nil
}

func makeWithDockerdWrapFun(dindID string, tarPaths []string, imgsWithDigests []string, opt WithDockerOpt) shellWrapFun {
	cacheDataRoot := strings.HasPrefix(dindID, "cache_")
	dockerRoot := path.Join("/var/earthly/dind", dindID)
	params := []string{
		fmt.Sprintf("EARTHLY_DOCKERD_DATA_ROOT=\"%s\"", dockerRoot),
		fmt.Sprintf("EARTHLY_DOCKERD_CACHE_DATA=\"%v\"", cacheDataRoot),
		fmt.Sprintf("EARTHLY_DOCKER_LOAD_FILES=\"%s\"", strings.Join(tarPaths, " ")),
		// This is not actually used, but it is needed in order to bust the cache
		// in case an image is updated.
		fmt.Sprintf("EARTHLY_IMAGES_WITH_DIGESTS=\"%s\"", strings.Join(imgsWithDigests, " ")),
	}
	params = append(params, composeParams(opt)...)
	return func(args []string, envVars []string, isWithShell, withDebugger, forceDebugger bool) []string {
		envVars2 := append(params, envVars...)
		return []string{
			"/bin/sh", "-c",
			strWithEnvVarsAndDocker(args, envVars2, isWithShell, withDebugger, forceDebugger, true, false, "", ""),
		}
	}
}

func composeParams(opt WithDockerOpt) []string {
	return []string{
		fmt.Sprintf("EARTHLY_START_COMPOSE=\"%t\"", (len(opt.ComposeFiles) > 0)),
		fmt.Sprintf("EARTHLY_COMPOSE_FILES=\"%s\"", strings.Join(opt.ComposeFiles, " ")),
		fmt.Sprintf("EARTHLY_COMPOSE_SERVICES=\"%s\"", strings.Join(opt.ComposeServices, " ")),
		// fmt.Sprintf("EARTHLY_DEBUG=\"true\""),
	}
}

func platformIncompatMsg(platr *platutil.Resolver) string {
	currentPlatStr := platr.Materialize(platr.Current()).String()
	nativePlatStr := platr.Materialize(platutil.NativePlatform).String()
	return "running WITH DOCKER as a non-native CPU architecture. This is not supported.\n" +
		fmt.Sprintf("Current platform: %s\n", currentPlatStr) +
		fmt.Sprintf("Native platform of the worker: %s\n", nativePlatStr) +
		fmt.Sprintf("Try using\n\n\tFROM --platform=native %s\n\ninstead.\n", suggestedDINDImage) +
		"You may still --load and --pull images of a different platform.\n"
}
