package containerutil

import (
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// ContainerInfo contains things we may care about from inspect output for a given container.
type ContainerInfo struct {
	ID       string
	Name     string
	Platform string
	Created  time.Time
	Status   string
	IPs      map[string]string
	Ports    []string
	Image    string
	ImageID  string
	Labels   map[string]string
}

const (
	// StatusMissing signifies that a container is not present.
	StatusMissing = "missing"

	// StatusCreated signifies that a container has been created, but not started.
	StatusCreated = "created"

	// StatusRestarting signifies that a container has started, stopped, and is currently restarting.
	StatusRestarting = "restarting"

	// StatusRunning signifies that a container is currently running.
	StatusRunning = "running"

	// StatusRemoving signifies that a container has exited and is currently being removed.
	StatusRemoving = "removing"

	// StatusPaused means a container has been suspended.
	StatusPaused = "paused"

	// StatusExited means that a container was running and has been stopped, but not removed.
	StatusExited = "exited"

	// StatusDead means that a container was killed for some reason and has not yet been restarted.
	StatusDead = "dead"
)

// ContainerLogs contains the stdout and stderr logs of a given container.
type ContainerLogs struct {
	Stdout string
	Stderr string
}

// FrontendInfo contains the client and server information for a frontend.
type FrontendInfo struct {
	ClientVersion    string
	ClientAPIVersion string
	ClientPlatform   string

	ServerVersion    string
	ServerAPIVersion string
	ServerPlatform   string
	ServerAddress    string
}

// ImageInfo contains information about a given image ref, including all relevant tags.
type ImageInfo struct {
	ID           string
	OS           string
	Architecture string
	Tags         []string
}

// VolumeInfo contains information about a given volume, including its name, where its mounted from, and the size of the volume.
type VolumeInfo struct {
	Name       string
	Mountpoint string
	SizeBytes  uint64
}

// ImageTag contains a source and target ref, used for tagging an image. It means that the SourceRef is tagged as the value in TargetRef.
type ImageTag struct {
	SourceRef string
	TargetRef string
}

// MountType constrains the kinds of mounts the Frontend API needs to support. Current valid values are bind and volume.
type MountType string

const (
	// MountBind is the bind MountType
	MountBind = MountType("bind")

	// MountVolume is the volume MountType
	MountVolume = MountType("volume")
)

// Mount contains the needed data to construct a mount for a container in a given frontend.
type Mount struct {
	Type     MountType
	Source   string
	Dest     string
	ReadOnly bool
}

// MountOpt is a list of Mounts to perform
type MountOpt []Mount

// ProtocolType constrains the kinds of protocols the frontend API needs to support. Current valid values are tcp and udp.
type ProtocolType string

const (
	// ProtocolTCP is the TCP protocol type
	ProtocolTCP = ProtocolType("tcp")

	// ProtocolUDP is the UDP protocol type
	ProtocolUDP = ProtocolType("udp")
)

// Port contains the needed data to publish a port for a given container in a given frontend.
type Port struct {
	IP            string
	HostPort      int
	ContainerPort int
	Protocol      ProtocolType
}

// PortOpt is a list of Ports to publish
type PortOpt []Port

// EnvMap is a map of environment variable names (key) to values. Values must be strings.
type EnvMap map[string]string

// LabelMap is a map of label names (key) to values. Values must be strings.
type LabelMap map[string]string

// ContainerRun contains the information needed to create and run a container.
type ContainerRun struct {
	NameOrID      string
	ImageRef      string
	Privileged    bool
	Envs          EnvMap
	Labels        LabelMap
	Mounts        MountOpt
	Ports         PortOpt
	ContainerArgs []string

	// We would like to shift to the non-shell providers. However, we do provide an option for supplying
	// additional arguments to the CLI when starting buildkit. While this allowed great flexibility, we
	// also do not know what or how it is being used. This gives us the option to support those users until
	// we decide to pull the plug. This argument is ignored by non-shell providers.
	AdditionalArgs []string
}

const (
	// FrontendAuto is automatic frontend detection.
	FrontendAuto = "auto"

	// FrontendDocker forces usage of the (future, currently unimplemented) docker API for container operations.
	FrontendDocker = "docker"

	// FrontendDockerShell forces usage of the docker binary for container operations.
	FrontendDockerShell = "docker-shell"

	// FrontendPodman forces usage of the (future, currently unimplemented) podman API for container operations.
	FrontendPodman = "podman"

	// FrontendPodmanShell forces usage of the podman binary for container operations.
	FrontendPodmanShell = "podman-shell"

	// FrontendStub is for when there is no valid provider but attempting to run anyways is desired; like integration tests, or the earthly/earthly image when NO_DOCKER is set.
	FrontendStub = "stub"
)

// CurrentFrontend contains information about the current container frontend. Useful when the frontend was configured using FrontendAuto for deciding transport.
type CurrentFrontend struct {
	*FrontendURLs
	Setting string
	Binary  string
	Type    string
}

const (
	// FrontendTypeShell signifies that a given frontend is using an external binary via shell.
	FrontendTypeShell = "shell"

	// FrontendTypeAPI signifies that a given frontend is using an API, typically through some external daemon.
	FrontendTypeAPI = "api"
)

const (
	// TCPAddressFmt is the address at which the daemon is available when using TCP.
	TCPAddressFmt = "tcp://127.0.0.1:%d"

	// DockerAddressFmt is the address at which the daemon is available when using a Docker Container directly
	DockerAddressFmt = "docker-container://%s-buildkitd"

	// PodmanAddressFmt is the address at which the daemon is available when using a Podman Container directly.
	// Currently unused due to image export issues
	PodmanAddressFmt = "podman-container://%s-buildkitd"

	// SatelliteAddress is the remote address when using a Satellite to execute your builds remotely.
	SatelliteAddress = "tcp://satellite.earthly.dev:8372"
)

var (
	errURLParseFailure      = errors.New("Invalid URL")
	errURLValidationFailure = errors.New("URL did not pass validation")
)

// FrontendURLs is a struct containing the relevant URLs to contact a container, if it is a buildkit container.
type FrontendURLs struct {
	// TODO this could probably be elsewhere still; but because we need at least one of these (LocalRegistryHost) to better
	//   handle the additional TLS argument needed for podman to pull images. Its at least _better_ than what we had before
	//   without a more major refactor.

	BuildkitHost      *url.URL
	LocalRegistryHost *url.URL
}
