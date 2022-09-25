package containerutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type nerdctlShellFrontend struct {
	*shellFrontend
}

// NewNerdctlShellFrontend constructs a new Frontend using the nerdctl binary installed on the host.
// It also ensures that the binary is functional for our needs and collects compatibility information.
func NewNerdctlShellFrontend(ctx context.Context, cfg *FrontendConfig) (ContainerFrontend, error) {
	fe := &nerdctlShellFrontend{
		shellFrontend: &shellFrontend{
			binaryName:              "nerdctl",
			runCompatibilityArgs:    make([]string, 0),
			globalCompatibilityArgs: make([]string, 0),
		},
	}

	output, err := fe.commandContextOutput(ctx, "info", "--format={{.SecurityOptions}}")
	if err != nil {
		return nil, err
	}
	fe.rootless = strings.Contains(output.string(), "rootless")

	fe.urls, err = fe.setupAndValidateAddresses(FrontendNerdctlShell, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate buildkit URLs")
	}

	return fe, nil
}

func (nsf *nerdctlShellFrontend) Scheme() string {
	return "nerdctl-container"
}

func (nsf *nerdctlShellFrontend) Config() *CurrentFrontend {
	return &CurrentFrontend{
		Setting:      FrontendNerdctlShell,
		Binary:       nsf.binaryName,
		Type:         FrontendTypeShell,
		FrontendURLs: nsf.urls,
	}
}

func (nsf *nerdctlShellFrontend) Information(ctx context.Context) (*FrontendInfo, error) {
	output, err := nsf.commandContextOutput(ctx, "version", "--format={{json .}}")
	if err != nil {
		return nil, err
	}

	type info struct {
		Client struct {
			Version   string
			GitCommit string
			OS        string
			Arch      string
		}
		Server struct {
			Components []struct {
				Name    string
				Version string
				Details struct {
					GitCommit string
				}
			}
		}
	}

	allInfo := info{}
	json.Unmarshal([]byte(output.string()), &allInfo)

	var serverVersion, serverAPIVersion string
	for _, component := range allInfo.Server.Components {
		if component.Name == " containerd" {
			serverVersion = component.Version
			serverAPIVersion = component.Details.GitCommit
		}
	}

	// Note that nerdctl does not support remote containerd: https://github.com/containerd/nerdctl/issues/473
	return &FrontendInfo{
		ClientVersion:    allInfo.Client.Version,
		ClientAPIVersion: allInfo.Client.GitCommit,
		ClientPlatform:   fmt.Sprintf("%s/%s", allInfo.Client.OS, allInfo.Client.Arch),
		ServerVersion:    serverVersion,
		ServerAPIVersion: serverAPIVersion,
		ServerPlatform:   fmt.Sprintf("%s/%s", allInfo.Client.OS, allInfo.Client.Arch),
	}, nil
}

func (nsf *nerdctlShellFrontend) ContainerInfo(ctx context.Context, namesOrIDs ...string) (map[string]*ContainerInfo, error) {
	results, err := nsf.shellFrontend.ContainerInfo(ctx, namesOrIDs...)
	if err != nil {
		return nil, err
	}

	for _, v := range results {
		// nerdctl puts an image tag where podman & docker put an image sha. We will unflip this here, and shell out to
		// the image command to get what we can. Note that the Id is (probably) not the same as Docker/Podman - it uses
		// a repo ID. Regular sha ids are not taggable here, but the special nerdctl one is... so we report that one.
		if v.Status == StatusMissing {
			continue
		}

		v.Image = v.ImageID
		v.ImageID = ""

		id, imageErr := nsf.commandContextOutput(ctx, "image", "inspect", `--format="{{json .RepoDigests}}"`, v.Image)
		if imageErr != nil {
			err = multierror.Append(err, imageErr)
			continue
		}

		// We assume the first repo digest is the correct one with no basis in reality, but it seems to work
		repoDigests := make([]string, 0)
		jsonErr := json.Unmarshal([]byte(id.stdout.String()), &repoDigests)
		if jsonErr != nil {
			err = multierror.Append(err, jsonErr)
			continue
		}

		digest, digestErr := nsf.getRepoDigest(repoDigests)
		if digestErr != nil {
			err = multierror.Append(err, errors.Wrapf(digestErr, "%s has no digests", v.Image))
			continue
		}

		v.ImageID = digest
	}

	return results, nil
}

func (nsf *nerdctlShellFrontend) ImageInfo(ctx context.Context, refs ...string) (map[string]*ImageInfo, error) {
	args := append([]string{"image", "inspect"}, refs...)

	// Ignore the error. This is because one or more of the provided refs could be missing.
	// This allows for Info to report that the image itself is missing.
	output, _ := nsf.commandContextOutput(ctx, args...)

	infos := map[string]*ImageInfo{}
	for _, ref := range refs {
		// preinitialize all as missing. It will get overwritten when we encounter a real one from the actual output.
		infos[ref] = &ImageInfo{}
	}

	// Anonymous struct to just pick out what we need
	images := []struct {
		RepoDigests []string `json:"RepoDigests"`
		Tags        []string `json:"RepoTags"`
	}{}
	json.Unmarshal([]byte(output.stdout.String()), &images)

	var err error
	for i, image := range images {
		digest, digestErr := nsf.getRepoDigest(image.RepoDigests)
		if digestErr != nil {
			err = multierror.Append(err, errors.Wrapf(digestErr, "%s has no digests", refs[i]))
			continue
		}

		infos[refs[i]] = &ImageInfo{
			ID:   digest,
			Tags: image.Tags,
		}
	}
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (nsf *nerdctlShellFrontend) ImagePull(ctx context.Context, refs ...string) error {
	var err error
	for _, ref := range refs {
		_, cmdErr := nsf.commandContextOutput(ctx, "pull", ref)
		if cmdErr != nil {
			err = multierror.Append(err, cmdErr)
		}
	}

	return err
}

func (nsf *nerdctlShellFrontend) ImageLoad(ctx context.Context, images ...io.Reader) error {
	var err error
	args := append(nsf.globalCompatibilityArgs, "load")
	for _, image := range images {
		// Do not use the wrapper to allow the image to come in on stdin
		cmd := exec.CommandContext(ctx, nsf.binaryName, args...)
		cmd.Stdin = image
		output, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil {
			err = multierror.Append(err, errors.Wrapf(cmdErr, "image load failed: %s", string(output)))
		}
	}

	return err
}

func (nsf *nerdctlShellFrontend) ImageLoadFromFileCommand(filename string) string {
	binary, args := nsf.commandContextStrings("load")

	all := []string{binary}
	all = append(all, args...)

	return fmt.Sprintf("cat %s | %s", shellescape.Quote(filename), strings.Join(all, " "))
}

func (nsf *nerdctlShellFrontend) VolumeInfo(ctx context.Context, volumeNames ...string) (map[string]*VolumeInfo, error) {
	// Ignore the error. This is because one or more of the provided names could be missing.
	// This allows for Info to report that the volume itself is missing.
	output, _ := nsf.commandContextOutput(ctx, "volume", "ls", "--format={{json  .}}")

	results := map[string]*VolumeInfo{}
	for _, name := range volumeNames {
		// Preinitialize all as missing. It will get overwritten when we encounter a real one from the actual output.
		results[name] = &VolumeInfo{Name: name}
	}

	// simple struct to just pick out what we need
	type volume struct {
		Name string `json:"Name"`
		// Size is not yet supported here. It does not appear to be anywhere in the public APIs?
		Mountpoint string `json:"Mountpoint"`
	}
	volumeInfos := make([]volume, 0)

	var err error
	for _, line := range strings.Split(strings.TrimSpace(output.stdout.String()), "\n") {
		v := volume{}
		jsonErr := json.Unmarshal([]byte(line), &v)
		if jsonErr != nil {
			err = multierror.Append(err, errors.Wrapf(jsonErr, "failed to decode docker volume info: '%v'", line))
			continue
		}
	}
	if err != nil {
		return nil, err
	}

	for _, name := range volumeNames {
		for _, volumeInfo := range volumeInfos {
			if name == volumeInfo.Name {
				results[name] = &VolumeInfo{
					Name:       volumeInfo.Name,
					Mountpoint: volumeInfo.Mountpoint,
				}
				break
			}
		}
	}

	return results, nil
}

func (nsf *nerdctlShellFrontend) getRepoDigest(digests []string) (string, error) {
	if len(digests) == 0 {
		return "", errors.New("no repo digests")
	}
	return strings.TrimLeftFunc(digests[0], func(r rune) bool {
		return r != '@'
	})[1:], nil // Trim off the remaining '@'
}
