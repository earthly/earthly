package containerutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/earthly/earthly/conslogging"
	"github.com/hashicorp/go-multierror"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/pkg/errors"
)

type containerInfo struct {
	ID      string    `json:"Id"`
	Name    string    `json:"Name"`
	Created time.Time `json:"Created"`
	State   struct {
		Status string `json:"Status"`
	} `json:"State"`
	NetworkSettings struct {
		Networks map[string]struct {
			IPAddress string `json:"IPAddress"`
		} `json:"Networks"`
		Ports map[string][]struct {
			HostIP   string `json:"HostIP"`
			HostPort string `json:"HostPort"`
		} `json:"Ports"`
	} `json:"NetworkSettings"`
	Config struct {
		Image  string            `json:"Image"`
		Labels map[string]string `json:"Labels"`
	} `json:"Config"`
	Image string `json:"Image"`
}

type shellFrontend struct {
	binaryName              string
	rootless                bool
	runCompatibilityArgs    []string
	globalCompatibilityArgs []string
	likelyPodman            bool

	urls    *FrontendURLs
	Console conslogging.ConsoleLogger
}

func (sf *shellFrontend) IsAvailable(ctx context.Context) bool {
	args := append(sf.globalCompatibilityArgs, "ps")
	cmd := exec.CommandContext(ctx, sf.binaryName, args...)
	err := cmd.Run()
	return err == nil
}

const containerDateFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

func (sf *shellFrontend) ContainerList(ctx context.Context) ([]*ContainerInfo, error) {
	// The custom format below is supported by Docker and Podman.
	args := []string{"ps", "--format", `{{.ID}},{{.Names}},{{.Status}},{{.Image}},{{.CreatedAt}}`}
	output, err := sf.commandContextOutput(ctx, args...)
	if err != nil {
		return nil, err
	}
	return parseContainerList(output.stdout.String())
}

func parseContainerList(output string) ([]*ContainerInfo, error) {
	ret := []*ContainerInfo{}
	// The Docker & Podman JSON output format differs, so we parse the standard output here.
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) != 5 {
			continue
		}
		createdAt, err := time.Parse(containerDateFormat, parts[4])
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse container date")
		}
		ret = append(ret, &ContainerInfo{
			ID:      parts[0],
			Name:    parts[1],
			Status:  parts[2],
			Image:   parts[3],
			Created: createdAt,
		})
	}
	return ret, nil
}

func (sf *shellFrontend) ContainerInfo(ctx context.Context, namesOrIDs ...string) (map[string]*ContainerInfo, error) {
	args := append([]string{"container", "inspect"}, namesOrIDs...)

	// Ignore the error. This is because one or more of the provided names or IDs could be missing.
	// This allows for Info to report that the container itself is missing.
	output, _ := sf.commandContextOutput(ctx, args...)

	infos := map[string]*ContainerInfo{}
	for _, nameOrID := range namesOrIDs {
		// Preinitialize all as missing. It will get overwritten when we encounter a real one from the actual output.
		infos[nameOrID] = &ContainerInfo{
			Name:   nameOrID,
			Status: StatusMissing,
		}
	}

	containers := []containerInfo{}
	err := json.Unmarshal([]byte(output.stdout.String()), &containers)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal container inspect output %s", output.stdout.String())
	}

	for i, container := range containers {
		ipAddresses := map[string]string{}
		for k, v := range container.NetworkSettings.Networks {
			ipAddresses[k] = v.IPAddress
		}

		infos[namesOrIDs[i]] = &ContainerInfo{
			ID:      container.ID,
			Name:    container.Name,
			Created: container.Created,
			Status:  container.State.Status,
			IPs:     ipAddresses,
			Image:   container.Config.Image,
			ImageID: container.Image,
			Labels:  container.Config.Labels,
			Ports:   formatPorts(container),
		}
	}

	return infos, nil
}

func formatPorts(info containerInfo) []string {
	ret := []string{}
	for key, ports := range info.NetworkSettings.Ports {
		for _, port := range ports {
			ret = append(ret, fmt.Sprintf("%s:%s:%s", port.HostIP, port.HostPort, key))
		}
	}
	return ret
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

	baseArgs := append(sf.globalCompatibilityArgs, "logs")
	for _, nameOrID := range namesOrIDs {
		// Don't use the wrapper so we can capture stderr and stdout individually
		args := append(baseArgs, nameOrID)
		cmd := exec.CommandContext(ctx, sf.binaryName, args...)

		var stdout, stderr strings.Builder
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		cmdErr := cmd.Run()
		if cmdErr != nil {
			err = multierror.Append(err, cmdErr)
			continue
		}
		logs[nameOrID] = &ContainerLogs{
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

		args = append(args, "-d") // Run detached, this feels implied by the API
		args = append(args, "--pull", "missing")
		args = append(args, "--name", container.NameOrID)
		args = append(args, container.AdditionalArgs...)
		args = append(args, sf.runCompatibilityArgs...)
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
		// preinitialize all as missing. It will get overwritten when we encounter a real one from the actual output.
		infos[ref] = &ImageInfo{}
	}

	// Anonymous struct to just pick out what we need
	images := []struct {
		ID   string   `json:"Id"`
		Tags []string `json:"RepoTags"`
	}{}
	err := json.Unmarshal([]byte(output.stdout.String()), &images)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse image info")
	}

	for i, image := range images {
		infos[refs[i]] = &ImageInfo{
			ID:   image.ID,
			Tags: image.Tags,
		}
	}

	return infos, nil
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
			err = multierror.Append(err, cmdErr)
		}
	}

	return err
}

type commandContextOutput struct {
	stdout strings.Builder
	stderr strings.Builder
}

func (cco *commandContextOutput) string() string {
	return strings.TrimSpace(cco.stdout.String() + cco.stderr.String())
}

func (sf *shellFrontend) commandContextStrings(args ...string) (string, []string) {
	allArgs := append(sf.globalCompatibilityArgs, args...)
	return sf.binaryName, allArgs
}

func (sf *shellFrontend) commandContextOutput(ctx context.Context, args ...string) (*commandContextOutput, error) {
	output := &commandContextOutput{}
	binary, args := sf.commandContextStrings(args...)
	sf.Console.VerbosePrintf("Running command: %s %s\n", binary, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Env = os.Environ() // Ensure all shellouts are using the current environment, picks up DOCKER_/PODMAN_ env vars when they matter
	cmd.Stdout = &output.stdout
	cmd.Stderr = &output.stderr
	err := cmd.Run()
	if err != nil {
		return output, errors.Wrapf(err, "command failed: %s %s: %s: %s", sf.binaryName, strings.Join(args, " "), err.Error(), output.string())
	}
	return output, nil
}

func (sf *shellFrontend) setupAndValidateAddresses(feType string, cfg *FrontendConfig) (*FrontendURLs, error) {
	calculatedBuildkitHost := cfg.BuildkitHostCLIValue
	if cfg.BuildkitHostCLIValue == "" {
		if cfg.BuildkitHostFileValue != "" {
			calculatedBuildkitHost = cfg.BuildkitHostFileValue
		} else {
			var err error
			calculatedBuildkitHost, err = DefaultAddressForSetting(feType, cfg.InstallationName, cfg.DefaultPort)
			if err != nil {
				return nil, errors.Wrap(err, "could not validate default address")
			}

		}
	}

	bkURL, err := parseAndValidateURL(calculatedBuildkitHost)
	if err != nil {
		return nil, err
	}

	lrURL := &url.URL{}
	if IsLocal(calculatedBuildkitHost) && cfg.LocalRegistryHostFileValue != "" {
		// Local registry only matters when local, and specified.
		lrURL, err = parseAndValidateURL(cfg.LocalRegistryHostFileValue)
		if err != nil {
			return nil, err
		}
		if !IsLocal(cfg.LocalRegistryHostFileValue) && bkURL.Hostname() != lrURL.Hostname() {
			cfg.Console.Warnf("Buildkit and local registry URLs are pointed at different hosts (%s vs. %s)", bkURL.Hostname(), lrURL.Hostname())
		}
	} else {
		if cfg.LocalRegistryHostFileValue != "" {
			cfg.Console.VerbosePrintf("Local registry host is specified while using remote buildkit. Local registry will not be used.")
		}
	}

	return &FrontendURLs{
		BuildkitHost:      bkURL,
		LocalRegistryHost: lrURL,
	}, nil
}

// DefaultAddressForSetting returns an address (signifying the desired/default transport) for a given frontend specified by setting.
func DefaultAddressForSetting(setting string, installationName string, defaultPort int) (string, error) {
	switch setting {
	case FrontendDockerShell:
		return fmt.Sprintf(DockerAddressFmt, installationName), nil

	case FrontendPodmanShell:
		return fmt.Sprintf(TCPAddressFmt, defaultPort), nil // Right now, podman only works over TCP. There are weird errors when trying to use the provided helper from buildkit.

	case FrontendStub:
		return fmt.Sprintf(DockerAddressFmt, installationName), nil // Maintain old behavior
	}

	return "", fmt.Errorf("no default buildkit address for %s", setting)
}

func parseAndValidateURL(addr string) (*url.URL, error) {
	parsed, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", addr, errURLParseFailure)
	}

	if parsed.Scheme != "tcp" && parsed.Scheme != "docker-container" && parsed.Scheme != "podman-container" {
		return nil, fmt.Errorf("%s is not a valid scheme. Only tcp or docker-container is allowed at this time: %w", parsed.Scheme, errURLValidationFailure)
	}

	if parsed.Port() == "" && parsed.Scheme == "tcp" {
		return nil, fmt.Errorf("%s does not contain a port number: %w", addr, errURLValidationFailure)
	}

	return parsed, nil
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
		parsed.Scheme == "docker-container" || // Accomodate feature flagging during transition. This will have omitted TLS?
		parsed.Scheme == "podman-container"
}
