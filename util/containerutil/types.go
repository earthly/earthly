package containerutil

// ContainerInfo is a representation of things we may care about from inspect output for a given container.
type ContainerInfo struct {
	ID      string
	Name    string
	Status  string
	IPs     map[string]string
	Image   string
	ImageID string
	Labels  map[string]string
}

const (
	StatusMissing    = "missing"
	StatusCreated    = "created"
	StatusRestarting = "restarting"
	StatusRunning    = "running"
	StatusRemoving   = "removing"
	StatusPaused     = "paused"
	StatusExited     = "exited"
	StatusDead       = "dead"
)

type ContainerLogs struct {
	Stdout string
	Stderr string
}

type FrontendInfo struct {
	ClientVersion    string
	ClientAPIVersion string
	ClientPlatform   string

	ServerVersion    string
	ServerAPIVersion string
	ServerPlatform   string
	ServerAddress    string
}

type ImageInfo struct {
	ID   string
	Tags []string
}

type VolumeInfo struct {
	Name       string
	Mountpoint string
	Size       uint64
}

type ImageTag struct {
	SourceRef string
	TargetRef string
}

type MountType string

const (
	MountBind   = MountType("bind")
	MountVolume = MountType("volume")
)

type Mount struct {
	Type     MountType
	Source   string
	Dest     string
	ReadOnly bool
}
type MountOpt []Mount

type ProtocolType string

const (
	ProtocolTCP = ProtocolType("tcp")
	ProtocolUDP = ProtocolType("udp")
)

type Port struct {
	IP            string
	HostPort      int
	ContainerPort int
	Protocol      ProtocolType
}
type PortOpt []Port

type EnvMap map[string]string

type LabelMap map[string]string
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
	FrontendAutomatic   = "automatic"
	FrontendDocker      = "docker"
	FrontendDockerShell = "docker-shell"
	FrontendPodman      = "podman"
	FrontendPodmanShell = "podman-shell"
)
