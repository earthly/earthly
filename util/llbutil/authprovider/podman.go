package authprovider

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/pkg/errors"
)

const (
	dockerDockerhubKey = "https://index.docker.io/v1/"
	podmanDockerhubKey = "docker.io"
	podmanAuthFile     = "auth.json"
)

// OS contains methods that are similar to the os package functions. It is
// provided so that os-level functions may be mocked.
type OS interface {
	Open(string) (io.ReadCloser, error)
	Getenv(string) string
}

type defaultOS struct{}

func (o defaultOS) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (o defaultOS) Getenv(name string) string {
	return os.Getenv(name)
}

type podmanCfg struct {
	os OS
}

// PodmanOpt is an option which may be used when constructing a podman auth
// provider.
type PodmanOpt func(podmanCfg) podmanCfg

// WithOS returns an option that provides custom OS-level functions for the
// podman auth provider to use.
func WithOS(o OS) PodmanOpt {
	return func(c podmanCfg) podmanCfg {
		c.os = o
		return c
	}
}

func NewPodman(stderr io.Writer, opts ...PodmanOpt) session.Attachable {
	conf := podmanCfg{
		os: defaultOS{},
	}
	for _, o := range opts {
		conf = o(conf)
	}
	if authfile := conf.os.Getenv("REGISTRY_AUTH_FILE"); authfile != "" {
		cfg, err := podmanAuth(conf.os, authfile)
		if err != nil {
			fmt.Fprintf(stderr, "WARNING: Error loading config file: %v\n", err)
			return authprovider.NewDockerAuthProvider(cfg)
		}
		syncDockerKey(cfg)
		return authprovider.NewDockerAuthProvider(cfg)
	}

	xdgRuntime := conf.os.Getenv("XDG_RUNTIME_DIR")
	if xdgRuntime == "" {
		idCmd := exec.Command("id", "-u")
		out, err := idCmd.CombinedOutput()
		if err != nil {
			return authprovider.NewDockerAuthProvider(config.LoadDefaultConfigFile(stderr))
		}

		id := strings.TrimSpace(string(out))
		// TODO: figure out how podman finds this path - on first pass we
		// couldn't find a good location for it.
		path := filepath.Join("/run", "containers", id, "auth.json")
		cfg, err := podmanAuth(conf.os, path)
		if errors.Is(err, fs.ErrNotExist) {
			return authprovider.NewDockerAuthProvider(config.LoadDefaultConfigFile(stderr))
		}
		if err != nil {
			fmt.Fprintf(stderr, "WARNING: Error loading config file: %v\n", err)
			return authprovider.NewDockerAuthProvider(cfg)
		}
		syncDockerKey(cfg)
		return authprovider.NewDockerAuthProvider(cfg)
	}

	path := filepath.Join(xdgRuntime, "containers", podmanAuthFile)
	cfg, err := podmanAuth(conf.os, path)
	if errors.Is(err, fs.ErrNotExist) {
		return authprovider.NewDockerAuthProvider(config.LoadDefaultConfigFile(stderr))
	}
	if err != nil {
		fmt.Fprintf(stderr, "WARNING: Error loading config file: %v\n", err)
		return authprovider.NewDockerAuthProvider(cfg)
	}
	syncDockerKey(cfg)
	return authprovider.NewDockerAuthProvider(cfg)
}

func podmanAuth(o OS, path string) (*configfile.ConfigFile, error) {
	f, err := o.Open(path)
	cfg := configfile.New(path)
	if err != nil {
		return cfg, errors.Wrap(err, path)
	}
	defer f.Close()

	if err := cfg.LoadFromReader(f); err != nil {
		return cfg, errors.Wrap(err, path)
	}
	return cfg, nil
}

func syncDockerKey(cfg *configfile.ConfigFile) {
	if _, ok := cfg.AuthConfigs[dockerDockerhubKey]; ok {
		return
	}
	v, ok := cfg.AuthConfigs[podmanDockerhubKey]
	if !ok {
		return
	}
	cfg.AuthConfigs[dockerDockerhubKey] = v
}
