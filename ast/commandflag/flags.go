package commandflag

import (
	"time"
)

// NOTE: Any new flags must be accompanied by the introduction of a new `VERSION` feature flag.
// This applies to new features which do **not** break backwards compatibility, which is needed
// to ensure an Earthfile that uses `VERSION 0.7` can be built by **any** of the earthly-v0.7.x binaries.

type IfOpts struct {
	Privileged bool     `long:"privileged" description:"Enable privileged mode"`
	WithSSH    bool     `long:"ssh" description:"Make available the SSH agent of the host"`
	NoCache    bool     `long:"no-cache" description:"Always run this specific item, ignoring cache"`
	Secrets    []string `long:"secret" description:"Make available a secret"`
	Mounts     []string `long:"mount" description:"Mount a file or directory"`
}

type ForOpts struct {
	Privileged bool     `long:"privileged" description:"Enable privileged mode"`
	WithSSH    bool     `long:"ssh" description:"Make available the SSH agent of the host"`
	NoCache    bool     `long:"no-cache" description:"Always run this specific item, ignoring cache"`
	Secrets    []string `long:"secret" description:"Make available a secret"`
	Mounts     []string `long:"mount" description:"Mount a file or directory"`
	Separators string   `long:"sep" description:"The separators to use for tokenizing the output of the IN expression. Defaults to '\n\t '"`
}

type RunOpts struct {
	Push            bool     `long:"push" description:"Execute this command only if the build succeeds and also if earthly is invoked in push mode"`
	Privileged      bool     `long:"privileged" description:"Enable privileged mode"`
	WithEntrypoint  bool     `long:"entrypoint" description:"Include the entrypoint of the image when running the command"`
	WithDocker      bool     `long:"with-docker" description:"Deprecated"`
	WithSSH         bool     `long:"ssh" description:"Make available the SSH agent of the host"`
	WithAWS         bool     `long:"aws" description:"Make any AWS credentials set in the environment available to RUN commands"`
	NoCache         bool     `long:"no-cache" description:"Always run this specific item, ignoring cache"`
	Interactive     bool     `long:"interactive" description:"Run this command with an interactive session, without saving changes"`
	InteractiveKeep bool     `long:"interactive-keep" description:"Run this command with an interactive session, saving changes"`
	Secrets         []string `long:"secret" description:"Make available a secret"`
	Mounts          []string `long:"mount" description:"Mount a file or directory"`
	Network         string   `long:"network" description:"Network to use; currently network=none is only supported"`
}

type FromOpts struct {
	AllowPrivileged bool     `long:"allow-privileged" description:"Allow commands under remote targets to enable privileged mode"`
	PassArgs        bool     `long:"pass-args" description:"Pass arguments to external targets"`
	BuildArgs       []string `long:"build-arg" description:"A build arg override passed on to a referenced Earthly target"`
	Platform        string   `long:"platform" description:"The platform to use"`
}

type FromDockerfileOpts struct {
	AllowPrivileged bool     `long:"allow-privileged" description:"Allow command to assume privileged mode"`
	BuildArgs       []string `long:"build-arg" description:"A build arg override passed on to a referenced Earthly target and also to the Dockerfile build"`
	Platform        string   `long:"platform" description:"The platform to use"`
	Target          string   `long:"target" description:"The Dockerfile target to inherit from"`
	Path            string   `short:"f" description:"The Dockerfile location on the host, relative to the current Earthfile, or as an artifact reference"`
}

type CopyOpts struct {
	From            string   `long:"from" description:"Not supported"`
	IsDirCopy       bool     `long:"dir" description:"Copy entire directories, not just the contents"`
	Chown           string   `long:"chown" description:"Apply a specific group and/or owner to the copied files and directories"`
	Chmod           string   `long:"chmod" description:"Apply a mode to the copied files and directories"`
	KeepTs          bool     `long:"keep-ts" description:"Keep created time file timestamps"`
	KeepOwn         bool     `long:"keep-own" description:"Keep owner info"`
	IfExists        bool     `long:"if-exists" description:"Do not fail if the artifact does not exist"`
	SymlinkNoFollow bool     `long:"symlink-no-follow" description:"Do not follow symlinks"`
	AllowPrivileged bool     `long:"allow-privileged" description:"Allow targets to assume privileged mode"`
	PassArgs        bool     `long:"pass-args" description:"Pass arguments to external targets"`
	Platform        string   `long:"platform" description:"The platform to use"`
	BuildArgs       []string `long:"build-arg" description:"A build arg override passed on to a referenced Earthly target"`
}

type SaveArtifactOpts struct {
	KeepTs          bool `long:"keep-ts" description:"Keep created time file timestamps"`
	KeepOwn         bool `long:"keep-own" description:"Keep owner info"`
	IfExists        bool `long:"if-exists" description:"Do not fail if the artifact does not exist"`
	SymlinkNoFollow bool `long:"symlink-no-follow" description:"Do not follow symlinks"`
	Force           bool `long:"force" description:"Force artifact to be saved, even if it means overwriting files or directories outside of the relative directory"`
}

type SaveImageOpts struct {
	Push           bool     `long:"push" description:"Push the image to the remote registry provided that the build succeeds and also that earthly is invoked in push mode"`
	CacheHint      bool     `long:"cache-hint" description:"Instruct Earthly that the current target should be saved entirely as part of the remote cache"`
	Insecure       bool     `long:"insecure" description:"Use unencrypted connection for the push"`
	NoManifestList bool     `long:"no-manifest-list" description:"Do not include a manifest list (specifying the platform) in the creation of the image"`
	CacheFrom      []string `long:"cache-from" description:"Declare additional cache import as a Docker tag"`
}

type BuildOpts struct {
	Platforms       []string `long:"platform" description:"The platform to use"`
	BuildArgs       []string `long:"build-arg" description:"A build arg override passed on to a referenced Earthly target"`
	AllowPrivileged bool     `long:"allow-privileged" description:"Allow targets to assume privileged mode"`
	PassArgs        bool     `long:"pass-args" description:"Pass arguments to external targets"`
	AutoSkip        bool     `long:"auto-skip" description:"Use auto-skip to bypass the target if nothing has changed"`
}

type GitCloneOpts struct {
	Branch string `long:"branch" description:"The git ref to use when cloning"`
	KeepTs bool   `long:"keep-ts" description:"Keep created time file timestamps"`
}

type HealthCheckOpts struct {
	Interval      time.Duration `long:"interval" description:"The interval between healthchecks" default:"30s"`
	Timeout       time.Duration `long:"timeout" description:"The timeout before the command is considered failed" default:"30s"`
	StartPeriod   time.Duration `long:"start-period" description:"An initialization time period in which failures are not counted towards the maximum number of retries"`
	Retries       int           `long:"retries" description:"The number of retries before a container is considered unhealthy" default:"3"`
	StartInterval time.Duration `long:"start-interval" description:"The time interval between health checks during the start period" default:"5s"`
}

type WithDockerOpts struct {
	ComposeFiles    []string `long:"compose" description:"A compose file used to bring up services from"`
	ComposeServices []string `long:"service" description:"A compose service to bring up"`
	Loads           []string `long:"load" description:"An image produced by Earthly which is loaded as a Docker image"`
	Platform        string   `long:"platform" description:"The platform to use"`
	BuildArgs       []string `long:"build-arg" description:"A build arg override passed on to a referenced Earthly target"`
	Pulls           []string `long:"pull" description:"An image which is pulled and made available in the docker cache"`
	AllowPrivileged bool     `long:"allow-privileged" description:"Allow targets referenced by load to assume privileged mode"`
	PassArgs        bool     `long:"pass-args" description:"Pass arguments to external targets"`
}

type DoOpts struct {
	AllowPrivileged bool `long:"allow-privileged" description:"Allow targets to assume privileged mode"`
	PassArgs        bool `long:"pass-args" description:"Pass arguments to external targets"`
}

type ImportOpts struct {
	AllowPrivileged bool `long:"allow-privileged" description:"Allow targets to assume privileged mode"`
	PassArgs        bool `long:"pass-args" description:"Pass arguments to external targets"`
}

type ArgOpts struct {
	Required bool `long:"required" description:"Require argument to be non-empty"`
	Global   bool `long:"global" description:"Global argument to make available to all other targets"`
}

type ProjectOpts struct{}

type SetOpts struct{}

type LetOpts struct{}

type PipelineOpts struct {
	Push            bool `long:"push" description:"Trigger a build in Cloud CI"`
	NoPipelineCache bool `long:"no-pipeline-cache" description:"Disable pipeline caching"`
}

type CacheOpts struct {
	Sharing string `long:"sharing" description:"The cache sharing mode: locked (default), shared, private"`
	Mode    string `long:"chmod" description:"Apply a mode to the cache folder" default:"0644"`
	ID      string `long:"id" description:"Cache ID, to reuse the same cache across different targets and Earthfiles"`
	Persist bool   `long:"persist" description:"If should persist cache state in image"`
}

// NewForOpts creates and returns a ForOpts with default separators.
func NewForOpts() ForOpts {
	return ForOpts{
		Separators: "\n\t ",
	}
}
