package flag

import (
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cmd/earthly/common"
	"github.com/earthly/earthly/util/containerutil"
)

const (
	// DefaultBuildkitdContainerSuffix is the suffix of the buildkitd container.
	DefaultBuildkitdContainerSuffix = "-buildkitd"
	// DefaultBuildkitdVolumeSuffix is the suffix of the docker volume used for storing the cache.
	DefaultBuildkitdVolumeSuffix = "-cache"

	DefaultEnvFile = ".env"
	EnvFileFlag    = "env-file-path"

	DefaultArgFile = ".arg"
	ArgFileFlag    = "arg-file-path"

	DefaultSecretFile = ".secret"
	SecretFileFlag    = "secret-file-path"
)

// Put flags on Flags instead as there are other things in the CLI that are being called + set
// by the subcommands so I thought it made since to declare them just once there and then
// pass them in
type Global struct {
	DockerfilePath             string
	EnableProfiler             bool
	InstallationName           string
	ConfigPath                 string
	GitUsernameOverride        string
	GitPasswordOverride        string
	GitBranchOverride          string
	ExecStatsSummary           string
	SSHAuthSock                string
	Verbose                    bool
	Debug                      bool
	DisplayExecStats           bool
	CloudHTTPAddr              string
	CloudGRPCAddr              string
	CloudGRPCInsecure          bool
	SatelliteAddress           string
	AuthToken                  string
	AuthJWT                    string
	DisableAnalytics           bool
	FeatureFlagOverrides       string
	EnvFile                    string
	ArgFile                    string
	SecretFile                 string
	NoBuildkitUpdate           bool
	LogstreamDebugFile         string
	LogstreamDebugManifestFile string
	LogstreamAddressOverride   string
	RequestID                  string
	BuildID                    string
	ServerConnTimeout          time.Duration
	BuildkitHost               string
	BuildkitdImage             string
	ContainerName              string
	GitLFSPullInclude          string
	BuildkitdSettings          buildkitd.Settings
	InteractiveDebugging       bool
	BootstrapNoBuildkit        bool
	ConversionParallelism      int
	LocalRegistryHost          string
	ContainerFrontend          containerutil.ContainerFrontend
	SatelliteName              string
	NoSatellite                bool
	ProjectName                string
	OrgName                    string
	CloudName                  string
	EarthlyCIRunner            bool
	ArtifactMode               bool
	ImageMode                  bool
	Pull                       bool
	Push                       bool
	CI                         bool
	UseTickTockBuildkitImage   bool
	Output                     bool
	NoOutput                   bool
	NoCache                    bool
	SkipBuildkit               bool
	AllowPrivileged            bool
	MaxRemoteCache             bool
	SaveInlineCache            bool
	UseInlineCache             bool
	NoFakeDep                  bool
	Strict                     bool
	GlobalWaitEnd              bool
	RemoteCache                string
	LocalSkipDB                string
	DisableRemoteRegistryProxy bool
	NoAutoSkip                 bool
	GithubAnnotations          bool
}

func (global *Global) RootFlags(installName string, bkImage string) []cli.Flag {
	defaultInstallationName := installName
	if defaultInstallationName == "" {
		defaultInstallationName = "earthly"
	}
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "installation-name",
			Value:       defaultInstallationName,
			EnvVars:     []string{"EARTHLY_INSTALLATION_NAME"},
			Usage:       "The earthly installation name to use when naming the buildkit container, the docker volume and the ~/.earthly directory",
			Destination: &global.InstallationName,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "config",
			Value:       "", // the default value will be applied in the "Before" fn, after flag.installationName is set.
			EnvVars:     []string{"EARTHLY_CONFIG"},
			Usage:       "Path to config file",
			Destination: &global.ConfigPath,
		},
		&cli.StringFlag{
			Name:        "ssh-auth-sock",
			Value:       os.Getenv("SSH_AUTH_SOCK"),
			EnvVars:     []string{"EARTHLY_SSH_AUTH_SOCK"},
			Usage:       "The SSH auth socket to use for ssh-agent forwarding",
			Destination: &global.SSHAuthSock,
		},
		&cli.StringFlag{
			Name:        "auth-token",
			EnvVars:     []string{"EARTHLY_TOKEN"},
			Usage:       "Force Earthly account login to authenticate with supplied token",
			Destination: &global.AuthToken,
		},
		&cli.StringFlag{
			Name:        "auth-jwt",
			EnvVars:     []string{"EARTHLY_JWT"},
			Usage:       "Force Earthly account to use supplied JWT token",
			Destination: &global.AuthJWT,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "git-username",
			EnvVars:     []string{"GIT_USERNAME"},
			Usage:       "The git username to use for git HTTPS authentication",
			Destination: &global.GitUsernameOverride,
		},
		&cli.StringFlag{
			Name:        "git-password",
			EnvVars:     []string{"GIT_PASSWORD"},
			Usage:       "The git password to use for git HTTPS authentication",
			Destination: &global.GitPasswordOverride,
		},
		&cli.StringFlag{
			Name:        "git-branch",
			EnvVars:     []string{"EARTHLY_GIT_BRANCH_OVERRIDE"},
			Usage:       "The git branch the build should be considered running in",
			Destination: &global.GitBranchOverride,
			Hidden:      true, // primarily used by CI to pass branch context
		},
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"V"},
			EnvVars:     []string{"EARTHLY_VERBOSE"},
			Usage:       "Enable verbose logging",
			Destination: &global.Verbose,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"D"},
			EnvVars:     []string{"EARTHLY_DEBUG"},
			Usage:       "Enable debug mode. This flag also turns on the debug mode of buildkitd, which may cause it to restart",
			Destination: &global.Debug,
			Hidden:      true, // For development purposes only.
		},
		&cli.BoolFlag{
			Name:        "exec-stats",
			EnvVars:     []string{"EARTHLY_EXEC_STATS"},
			Usage:       "Display container stats (e.g. cpu and memory usage)",
			Destination: &global.DisplayExecStats,
			Hidden:      true, // Experimental
		},
		&cli.StringFlag{
			Name:        "exec-stats-summary",
			EnvVars:     []string{"EARTHLY_EXEC_STATS_SUMMARY"},
			Usage:       "Output summarized container stats (e.g. cpu and memory usage) to the specified file",
			Destination: &global.ExecStatsSummary,
			Hidden:      true, // Experimental
		},
		&cli.BoolFlag{
			Name:        "profiler",
			EnvVars:     []string{"EARTHLY_PROFILER"},
			Usage:       "Enable the profiler",
			Destination: &global.EnableProfiler,
			Hidden:      true, // Dev purposes only.
		},
		&cli.StringFlag{
			Name:    "buildkit-host",
			Value:   "",
			EnvVars: []string{"EARTHLY_BUILDKIT_HOST"},
			Usage: `The URL to use for connecting to a buildkit host
		If empty, earthly will attempt to start a buildkitd instance via docker run`,
			Destination: &global.BuildkitHost,
		},
		&cli.StringFlag{
			Name:        "server",
			Value:       "https://api.earthly.dev",
			EnvVars:     []string{"EARTHLY_SERVER_ADDRESS"},
			Usage:       "API server override for dev purposes",
			Destination: &global.CloudHTTPAddr,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "grpc",
			Value:       "ci.earthly.dev:443",
			EnvVars:     []string{"EARTHLY_GRPC_ADDRESS"},
			Usage:       "gRPC server override for dev purposes",
			Destination: &global.CloudGRPCAddr,
			Hidden:      true, // Internal.
		},
		&cli.BoolFlag{
			Name:        "grpc-insecure",
			EnvVars:     []string{"EARTHLY_GRPC_INSECURE"},
			Usage:       "Makes gRPC connections insecure for dev purposes",
			Destination: &global.CloudGRPCInsecure,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "satellite-address",
			EnvVars:     []string{"EARTHLY_SATELLITE_ADDRESS"},
			Usage:       "Satellite address override for dev purposes",
			Destination: &global.SatelliteAddress,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "request-id",
			EnvVars:     []string{"EARTHLY_REQUEST_ID"},
			Usage:       "Override a request ID to the backend API. Useful for debugging or manually retrying a request.",
			Destination: &global.RequestID,
			Hidden:      true, // Internal
		},
		&cli.BoolFlag{
			Name:        "no-buildkit-update",
			EnvVars:     []string{"EARTHLY_NO_BUILDKIT_UPDATE"},
			Usage:       "Disable the automatic update of buildkitd",
			Destination: &global.NoBuildkitUpdate,
			Hidden:      true, // Internal.
		},
		&cli.BoolFlag{
			EnvVars:     []string{"EARTHLY_DISABLE_ANALYTICS", "DO_NOT_TRACK"},
			Usage:       "Disable collection of analytics",
			Destination: &global.DisableAnalytics,
		},
		&cli.StringFlag{
			Name:        "version-flag-overrides",
			EnvVars:     []string{"EARTHLY_VERSION_FLAG_OVERRIDES"},
			Usage:       "Apply additional flags after each VERSION command across all Earthfiles, multiple flags can be separated by commas",
			Destination: &global.FeatureFlagOverrides,
			Hidden:      true, // used for feature-flipping from ./earthly dev script
		},
		&cli.StringFlag{
			Name:        EnvFileFlag,
			EnvVars:     []string{"EARTHLY_ENV_FILE_PATH"},
			Usage:       "Use values from this file as earthly environment variables; values are no longer used as --build-arg's or --secret's",
			Value:       DefaultEnvFile,
			Destination: &global.EnvFile,
		},
		&cli.StringFlag{
			Name:        ArgFileFlag,
			EnvVars:     []string{"EARTHLY_ARG_FILE_PATH"},
			Usage:       "Use values from this file as earthly buildargs",
			Value:       DefaultArgFile,
			Destination: &global.ArgFile,
		},
		&cli.StringFlag{
			Name:        SecretFileFlag,
			EnvVars:     []string{"EARTHLY_SECRET_FILE_PATH"},
			Usage:       "Use values from this file as earthly secrets",
			Value:       DefaultSecretFile,
			Destination: &global.SecretFile,
		},
		&cli.StringFlag{
			Name:        "logstream-debug-file",
			EnvVars:     []string{"EARTHLY_LOGSTREAM_DEBUG_FILE"},
			Usage:       "Enable log streaming debugging output to a file",
			Destination: &global.LogstreamDebugFile,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "logstream-debug-manifest-file",
			EnvVars:     []string{"EARTHLY_LOGSTREAM_DEBUG_MANIFEST_FILE"},
			Usage:       "Enable log streaming manifest debugging output to a file",
			Destination: &global.LogstreamDebugManifestFile,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "logstream-address",
			EnvVars:     []string{"EARTHLY_LOGSTREAM_ADDRESS"},
			Usage:       "Override the Logstream address",
			Destination: &global.LogstreamAddressOverride,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "build-id",
			EnvVars:     []string{"EARTHLY_BUILD_ID"},
			Usage:       "The build ID to use for identifying the build in Earthly Cloud. If not specified, a random ID will be generated",
			Destination: &global.BuildID,
			Hidden:      true, // Internal.
		},
		&cli.DurationFlag{
			Name:        "server-conn-timeout",
			Usage:       "Earthly API server connection timeout value",
			EnvVars:     []string{"EARTHLY_SERVER_CONN_TIMEOUT"},
			Hidden:      true, // Internal.
			Value:       5 * time.Second,
			Destination: &global.ServerConnTimeout,
		},
		&cli.BoolFlag{
			Name:        "artifact",
			Aliases:     []string{"a"},
			Usage:       "Output specified artifact; a wildcard (*) can be used to output all artifacts",
			Destination: &global.ArtifactMode,
		},
		&cli.BoolFlag{
			Name:        "image",
			Usage:       "Output only docker image of the specified target",
			Destination: &global.ImageMode,
		},
		&cli.BoolFlag{
			Name:        "pull",
			EnvVars:     []string{"EARTHLY_PULL"},
			Usage:       "Force pull any referenced Docker images",
			Destination: &global.Pull,
			Hidden:      true, // Experimental
		},
		&cli.BoolFlag{
			Name:        "push",
			EnvVars:     []string{"EARTHLY_PUSH"},
			Usage:       "Push docker images and execute RUN --push commands",
			Destination: &global.Push,
		},
		&cli.BoolFlag{
			Name:        "ci",
			EnvVars:     []string{"EARTHLY_CI"},
			Usage:       common.Wrap("Execute in CI mode. ", "Implies --no-output --strict"),
			Destination: &global.CI,
		},
		&cli.BoolFlag{
			Name:        "ticktock",
			EnvVars:     []string{"EARTHLY_TICKTOCK"},
			Usage:       "Use earthly's experimental buildkit ticktock codebase",
			Destination: &global.UseTickTockBuildkitImage,
			Hidden:      true, // Experimental
		},
		&cli.BoolFlag{
			Name:        "output",
			EnvVars:     []string{"EARTHLY_OUTPUT"},
			Usage:       "Allow artifacts or images to be output, even when running under --ci mode",
			Destination: &global.Output,
		},
		&cli.BoolFlag{
			Name:        "no-output",
			EnvVars:     []string{"EARTHLY_NO_OUTPUT"},
			Usage:       common.Wrap("Do not output artifacts or images", "(using --push is still allowed)"),
			Destination: &global.NoOutput,
		},
		&cli.BoolFlag{
			Name:        "no-cache",
			EnvVars:     []string{"EARTHLY_NO_CACHE"},
			Usage:       "Do not use cache while building",
			Destination: &global.NoCache,
		},
		&cli.BoolFlag{
			Name:        "auto-skip",
			EnvVars:     []string{"EARTHLY_AUTO_SKIP"},
			Usage:       "Skip buildkit if target has already been built",
			Destination: &global.SkipBuildkit,
		},
		&cli.BoolFlag{
			Name:        "allow-privileged",
			Aliases:     []string{"P"},
			EnvVars:     []string{"EARTHLY_ALLOW_PRIVILEGED"},
			Usage:       "Allow build to use the --privileged flag in RUN commands",
			Destination: &global.AllowPrivileged,
		},
		&cli.BoolFlag{
			Name:        "max-remote-cache",
			EnvVars:     []string{"EARTHLY_MAX_REMOTE_CACHE"},
			Usage:       "Saves all intermediate images too in the remote cache",
			Destination: &global.MaxRemoteCache,
		},
		&cli.BoolFlag{
			Name:        "save-inline-cache",
			EnvVars:     []string{"EARTHLY_SAVE_INLINE_CACHE"},
			Usage:       "Enable cache inlining when pushing images",
			Destination: &global.SaveInlineCache,
		},
		&cli.BoolFlag{
			Name:        "use-inline-cache",
			EnvVars:     []string{"EARTHLY_USE_INLINE_CACHE"},
			Usage:       common.Wrap("Attempt to use any inline cache that may have been previously pushed ", "uses image tags referenced by SAVE IMAGE --push or SAVE IMAGE --cache-from"),
			Destination: &global.UseInlineCache,
		},
		&cli.BoolFlag{
			Name:        "interactive",
			Aliases:     []string{"i"},
			EnvVars:     []string{"EARTHLY_INTERACTIVE"},
			Usage:       "Enable interactive debugging",
			Destination: &global.InteractiveDebugging,
		},
		&cli.BoolFlag{
			Name:        "no-fake-dep",
			EnvVars:     []string{"EARTHLY_NO_FAKE_DEP"},
			Usage:       "Internal feature flag for fake-dep",
			Destination: &global.NoFakeDep,
			Hidden:      true, // Internal.
		},
		&cli.BoolFlag{
			Name:        "strict",
			EnvVars:     []string{"EARTHLY_STRICT"},
			Usage:       "Disallow usage of features that may create unrepeatable builds",
			Destination: &global.Strict,
		},
		&cli.BoolFlag{
			Name:        "global-wait-end",
			EnvVars:     []string{"EARTHLY_GLOBAL_WAIT_END"},
			Usage:       "enables global wait-end code in place of builder code",
			Destination: &global.GlobalWaitEnd,
			Hidden:      true, // used to force code-coverage of future builder.go refactor (once we remove support for 0.6)
		},
		&cli.StringFlag{
			Name:        "git-lfs-pull-include",
			EnvVars:     []string{"EARTHLY_GIT_LFS_PULL_INCLUDE"},
			Usage:       "When referencing a remote target, perform a git lfs pull include prior to running the target. Note that this flag is (hopefully) temporary, see https://github.com/earthly/earthly/issues/2921 for details.",
			Destination: &global.GitLFSPullInclude,
			Hidden:      true, // Experimental
		},
		&cli.BoolFlag{
			Name:        "earthly-ci-runner",
			EnvVars:     []string{"EARTHLY_CI_RUNNER"},
			Usage:       "Internal flag to indicate the build is running within Earthly CI",
			Destination: &global.EarthlyCIRunner,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "auto-skip-db-path",
			EnvVars:     []string{"EARTHLY_AUTO_SKIP_DB_PATH"},
			Usage:       "use a local database instead of the cloud db",
			Destination: &global.LocalSkipDB,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "org",
			EnvVars:     []string{"EARTHLY_ORG"},
			Usage:       common.Wrap("The name of the organization that the satellite belongs to. ", "Required when using --satellite and user is a member of multiple organizations."),
			Required:    false,
			Destination: &global.OrgName,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "project",
			EnvVars:     []string{"EARTHLY_PROJECT"},
			Usage:       "The name of the project that may be used during log streaming.",
			Required:    false,
			Destination: &global.ProjectName,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "satellite",
			Aliases:     []string{"sat"},
			EnvVars:     []string{"EARTHLY_SATELLITE"},
			Usage:       "The name of satellite to use for this build.",
			Required:    false,
			Destination: &global.SatelliteName,
		},
		&cli.BoolFlag{
			Name:        "no-satellite",
			Aliases:     []string{"no-sat"},
			EnvVars:     []string{"EARTHLY_NO_SATELLITE"},
			Usage:       "Disables the use of a selected satellite for this build.",
			Required:    false,
			Destination: &global.NoSatellite,
		},
		&cli.StringFlag{
			Name:        "buildkit-image",
			Value:       bkImage,
			EnvVars:     []string{"EARTHLY_BUILDKIT_IMAGE"},
			Usage:       "The docker image to use for the buildkit daemon",
			Destination: &global.BuildkitdImage,
		},
		&cli.StringFlag{
			Name:        "buildkit-container-name",
			Value:       defaultInstallationName + DefaultBuildkitdContainerSuffix,
			EnvVars:     []string{"EARTHLY_CONTAINER_NAME"},
			Usage:       "The docker container name to use for the buildkit daemon",
			Destination: &global.ContainerName,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "buildkit-volume-name",
			Value:       defaultInstallationName + DefaultBuildkitdVolumeSuffix,
			EnvVars:     []string{"EARTHLY_VOLUME_NAME"},
			Usage:       "The docker volume name to use for the buildkit daemon cache",
			Destination: &global.BuildkitdSettings.VolumeName,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "remote-cache",
			EnvVars:     []string{"EARTHLY_REMOTE_CACHE"},
			Usage:       "A remote docker image tag use as explicit cache and optionally additional attributes to set in the image (Format: \"<image-tag>[,<attr1>=<val1>,<attr2>=<val2>,...]\")",
			Destination: &global.RemoteCache,
		},
		&cli.BoolFlag{
			Name:        "disable-remote-registry-proxy",
			EnvVars:     []string{"EARTHLY_DISABLE_REMOTE_REGISTRY_PROXY"},
			Usage:       "Don't use the Docker registry proxy when transferring images",
			Destination: &global.DisableRemoteRegistryProxy,
			Value:       false,
		},
		&cli.BoolFlag{
			Name:        "no-auto-skip",
			EnvVars:     []string{"EARTHLY_NO_AUTO_SKIP"},
			Usage:       "Disable auto-skip functionality",
			Destination: &global.NoAutoSkip,
			Value:       false,
		},
		&cli.BoolFlag{
			Name:        "github-annotations",
			EnvVars:     []string{"GITHUB_ACTIONS"},
			Usage:       "Enable Git Hub Actions workflow specific output",
			Destination: &global.GithubAnnotations,
			Value:       false,
		},
	}
}
