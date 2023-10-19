package authprovider

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/pkg/errors"

	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/util/bklog"
)

const (
	// PodmanConfigFileName is the name of config file
	PodmanConfigFileName = "auth.json"

	// XDGSubPath is the subpath in the XDG_RUNTIME_DIR directory that contains the podman config file
	XDGSubPath = "containers"
)

// NewPodmanAuthProvider provides an attachable that loads the Podman config instead of the hardcoded docker ones. It is
// More or less a copy of the relevant pieces for loading Docker config files used by the DockerAuthProvider, just using
// Podmans default config paths & names.
func NewPodmanAuthProvider(stderr io.Writer) session.Attachable {
	xdgRuntime, ok := os.LookupEnv("XDG_RUNTIME_DIR")
	if !ok {
		// Podman uses docker's default settings location when the XDG_RUNTIME_DIR is missing.
		// See here for more details: https://docs.podman.io/en/latest/markdown/podman-login.1.html
		bklog.G(context.TODO()).Debugf("WARNING: XDG_RUNTIME_DIR is not set, trying Docker config")
		cfg := config.LoadDefaultConfigFile(stderr)
		return NewDockerAuthProvider(cfg, nil)
	}

	filename := filepath.Join(xdgRuntime, XDGSubPath, PodmanConfigFileName)
	configFile := configfile.New(filename)

	if file, err := os.Open(filename); err == nil {
		defer file.Close()

		err = configFile.LoadFromReader(file)
		if err != nil {
			err = errors.Wrap(err, filename)
			fmt.Fprintf(stderr, "WARNING: Error loading config file: %v\n", err)
		}

		return NewDockerAuthProvider(configFile, nil)
	} else if !os.IsNotExist(err) {
		// if file is there but we can't stat it for any reason other
		// than it doesn't exist then stop
		fmt.Fprintf(stderr, "WARNING: Error loading config file: %v\n", err)
		return NewDockerAuthProvider(configFile, nil)
	}

	// No Podman config file, just use the docker one then to pick up docker credentials.
	// Podman uses docker's default settings location when the XDG path does not exist.
	// See here for more details: https://docs.podman.io/en/latest/markdown/podman-login.1.html
	bklog.G(context.TODO()).Debugf("WARNING: %s did not exist, trying Docker config", filename)
	cfg := config.LoadDefaultConfigFile(stderr)
	return NewDockerAuthProvider(cfg, nil)
}
