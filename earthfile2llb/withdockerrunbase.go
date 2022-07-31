package earthfile2llb

import (
	"context"
	"fmt"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/platutil"
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

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
	runOpts := []llb.RunOption{
		llb.AddMount(
			dockerAutoInstallScriptPath, llb.Scratch(), llb.HostBind(), llb.SourcePath(dockerAutoInstallScriptPath)),
		llb.Args(args),
		llb.WithCustomNamef("%sWITH DOCKER (install deps)", w.c.vertexPrefix(false, false, false)),
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
	runOpts := []llb.RunOption{
		llb.AddMount(
			dockerdWrapperPath, llb.Scratch(), llb.HostBind(), llb.SourcePath(dockerdWrapperPath)),
		llb.Args(args),
		llb.WithCustomNamef("%sWITH DOCKER (docker-compose config)", w.c.vertexPrefix(false, false, false)),
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
