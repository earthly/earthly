package containerutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/pkg/errors"
)

type shellFrontend struct {
	binaryName        string
	rootless          bool
	compatibilityArgs []string
}

func (sf *shellFrontend) IsAvaliable(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, sf.binaryName, "ps")
	err := cmd.Run()
	return err == nil
}

func (sf *shellFrontend) ContainerInfo(ctx context.Context, namesOrIDs ...string) (map[string]*ContainerInfo, error) {
	args := append([]string{"container", "inspect"}, namesOrIDs...)

	// Ignore the error. This is because one or more of the provided names or IDs could be missing.
	// This allows for Info to report that the container itself is missing.
	output, _ := sf.commandContextOutput(ctx, args...)

	infos := map[string]*ContainerInfo{}
	for _, nameOrID := range namesOrIDs {
		// Pre-initialize all as missing. It will get overwritten when we encounter a real one from the actual output.
		infos[nameOrID] = &ContainerInfo{
			Name:   nameOrID,
			Status: StatusMissing,
		}
	}

	// Anonymous struct to just pick out what we need
	containers := []struct {
		ID    string `json:"Id"`
		Name  string `json:"Name"`
		State struct {
			Status string `json:"Status"`
		} `json:"State"`
		NetworkSettings struct {
			Networks map[string]struct {
				IPAddress string `json:"IPAddress"`
			} `json:"Networks"`
		} `json:"NetworkSettings"`
		Config struct {
			Image  string            `json:"Image"`
			Labels map[string]string `json:"Labels"`
		} `json:"Config"`
		Image string `json:"Image"`
	}{}
	json.Unmarshal([]byte(output.stdout.String()), &containers)

	for i, container := range containers {
		ipAddresses := map[string]string{}
		for k, v := range container.NetworkSettings.Networks {
			ipAddresses[k] = v.IPAddress
		}

		infos[namesOrIDs[i]] = &ContainerInfo{
			ID:      container.ID,
			Name:    container.Name,
			Status:  container.State.Status,
			IPs:     ipAddresses,
			Image:   container.Config.Image,
			ImageID: container.Image,
			Labels:  container.Config.Labels,
		}
	}

	return infos, nil
}

func (sf *shellFrontend) ContainerRemove(ctx context.Context, force bool, namesOrIDs ...string) error {
	args := []string{"rm"}

	if force {
		args = append(args, "-f")
	}

	args = append(args, namesOrIDs...)

	_, err := sf.commandContextOutput(ctx, args...)
	return err
}

func (sf *shellFrontend) ContainerStop(ctx context.Context, timeoutSec uint, namesOrIDs ...string) error {
	args := append([]string{"stop", "-t", strconv.FormatUint(uint64(timeoutSec), 10)}, namesOrIDs...)

	_, err := sf.commandContextOutput(ctx, args...)
	return err
}

func (sf *shellFrontend) ContainerLogs(ctx context.Context, namesOrIDs ...string) (map[string]*ContainerLogs, error) {
	logs := map[string]*ContainerLogs{}
	var err error

	for _, nameOrId := range namesOrIDs {
		// Don't use the wrapper so we can capture stderr and stdout individually
		cmd := exec.CommandContext(ctx, sf.binaryName, "logs", nameOrId)

		var stdout, stderr strings.Builder
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		cmdErr := cmd.Run()
		if cmdErr != nil {
			err = multierror.Append(err, cmdErr)
			continue
		}
		logs[nameOrId] = &ContainerLogs{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
		}
	}

	return logs, err
}

func (sf *shellFrontend) ContainerRun(ctx context.Context, containers ...ContainerRun) error {
	var err error
	for _, container := range containers {
		args := []string{"run"}

		if container.Privileged {
			args = append(args, "--privileged")
		}

		for k, v := range container.Envs {
			env := fmt.Sprintf("%s=%s", k, v)
			args = append(args, "--env", env)
		}

		for k, v := range container.Labels {
			label := fmt.Sprintf("%s=%s", k, v)
			args = append(args, "--label", label)
		}

		for _, mnt := range container.Mounts {
			mount := fmt.Sprintf("type=%s,source=%s,dst=%s", mnt.Type, mnt.Source, mnt.Dest)
			// Older podmans do not support "readonly" as an option for the mount, but all CLIs so far support "ro"
			// Also some older podmans interpret the presence of the "ro" flag existing at all as meaning readonly
			if mnt.ReadOnly {
				mount = fmt.Sprintf("%s,ro=%t", mount, mnt.ReadOnly)
			}
			args = append(args, "--mount", mount)
		}

		for _, prt := range container.Ports {
			hostPort := strconv.FormatInt(int64(prt.HostPort), 10)
			if prt.HostPort <= 0 {
				// Docker allows 0 as a port for autoassign. Podman does not.
				// Both honor omission to allow a random open host port.
				hostPort = ""
			}

			port := fmt.Sprintf("%s:%v:%v", prt.IP, hostPort, prt.ContainerPort)

			if prt.Protocol != "" {
				// Unspecified protocol means we dont specify a protocol either.
				port = fmt.Sprintf("%s/%s", port, prt.Protocol)
			}

			args = append(args, "--publish", port)
		}

		if sf.supportsPlatformArg(ctx) {
			args = append(args, "--platform", getPlatform())
		}

		args = append(args, "-d") // Run detached, this feels implied by the API
		args = append(args, "--name", container.NameOrID)
		args = append(args, container.AdditionalArgs...)
		args = append(args, sf.compatibilityArgs...)
		args = append(args, container.ImageRef)
		args = append(args, container.ContainerArgs...)

		_, cmdErr := sf.commandContextOutput(ctx, args...)
		if cmdErr != nil {
			err = multierror.Append(err, cmdErr)
		}
	}

	return err
}

func (sf *shellFrontend) ImageInfo(ctx context.Context, refs ...string) (map[string]*ImageInfo, error) {
	args := append([]string{"image", "inspect"}, refs...)

	// Ignore the error. This is because one or more of the provided refs could be missing.
	// This allows for Info to report that the image itself is missing.
	output, _ := sf.commandContextOutput(ctx, args...)

	infos := map[string]*ImageInfo{}
	for _, ref := range refs {
		// Pre-initialize all as missing. It will get overwritten when we encounter a real one from the actual output.
		infos[ref] = &ImageInfo{}
	}

	// Anonymous struct to just pick out what we need
	images := []struct {
		ID   string   `json:"Id"`
		Tags []string `json:"RepoTags"`
	}{}
	json.Unmarshal([]byte(output.stdout.String()), &images)

	for i, image := range images {
		infos[refs[i]] = &ImageInfo{
			ID:   image.ID,
			Tags: image.Tags,
		}
	}

	return infos, nil
}

func (sf *shellFrontend) ImagePull(ctx context.Context, refs ...string) error {
	var err error
	for _, ref := range refs {
		_, cmdErr := sf.commandContextOutput(ctx, "pull", ref)
		if cmdErr != nil {
			err = multierror.Append(err, cmdErr)
		}
	}

	return err
}

func (sf *shellFrontend) ImageRemove(ctx context.Context, force bool, refs ...string) error {
	args := []string{"image", "rm"}
	if force {
		args = append(args, "-f")
	}

	args = append(args, refs...)

	_, err := sf.commandContextOutput(ctx, args...)
	return err
}

func (sf *shellFrontend) ImageTag(ctx context.Context, tags ...ImageTag) error {
	var err error
	for _, tag := range tags {
		_, cmdErr := sf.commandContextOutput(ctx, "tag", tag.SourceRef, tag.TargetRef)
		if cmdErr != nil {
			multierror.Append(err)
		}
	}

	return err
}

func (sf *shellFrontend) ImageLoad(ctx context.Context, images ...io.Reader) error {
	var err error
	for _, image := range images {
		// Do not use the wrapper to allow the image to come in on stdin
		cmd := exec.CommandContext(ctx, sf.binaryName, "load")
		cmd.Stdin = image
		output, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil {
			err = multierror.Append(err, errors.Wrapf(cmdErr, "image load failed: %s", string(output)))
		}
	}

	return err
}

type commmandContextOutput struct {
	stdout strings.Builder
	stderr strings.Builder
}

func (cco *commmandContextOutput) string() string {
	return strings.TrimSpace(cco.stdout.String() + cco.stderr.String())
}

func (sf *shellFrontend) commandContextOutput(ctx context.Context, args ...string) (*commmandContextOutput, error) {
	output := &commmandContextOutput{}

	cmd := exec.CommandContext(ctx, sf.binaryName, args...)
	cmd.Env = os.Environ() // Ensure all shellouts are using the current environment, picks up DOCKER_/PODMAN_ env vars when they matter
	cmd.Stdout = &output.stdout
	cmd.Stderr = &output.stderr

	err := cmd.Run()
	if err != nil {
		return output, errors.Wrapf(err, "command failed: %s %s: %s: %s", sf.binaryName, strings.Join(args, " "), err.Error(), output.string())
	}

	return output, nil
}

func (sf *shellFrontend) supportsPlatformArg(ctx context.Context) bool {
	// We can't run scratch, but the error is different depending on whether
	// --platform is supported or not. This is faster than attempting to run
	// an actual image which may require downloading.
	output, _ := sf.commandContextOutput(ctx, "run", "--rm", "--platform", getPlatform(), "scratch")
	return strings.Contains(output.string(), "Unable to find image")
}
