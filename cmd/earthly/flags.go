package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/util/containerutil"
)

func (app *earthlyApp) rootFlags() []cli.Flag {
	defaultInstallationName := DefaultInstallationName
	if defaultInstallationName == "" {
		defaultInstallationName = "earthly"
	}
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Value:       defaultConfigPath(defaultInstallationName),
			EnvVars:     []string{"EARTHLY_CONFIG"},
			Usage:       "Path to config file",
			Destination: &app.configPath,
		},
		&cli.StringFlag{
			Name:        "ssh-auth-sock",
			Value:       os.Getenv("SSH_AUTH_SOCK"),
			EnvVars:     []string{"EARTHLY_SSH_AUTH_SOCK"},
			Usage:       wrap("The SSH auth socket to use for ssh-agent forwarding", ""),
			Destination: &app.sshAuthSock,
		},
		&cli.StringFlag{
			Name:        "auth-token",
			EnvVars:     []string{"EARTHLY_TOKEN"},
			Usage:       "Force Earthly account login to authenticate with supplied token",
			Destination: &app.authToken,
		},
		&cli.StringFlag{
			Name:        "auth-jwt",
			EnvVars:     []string{"EARTHLY_JWT"},
			Usage:       "Force Earthly account to use supplied JWT token",
			Destination: &app.authJWT,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "git-username",
			EnvVars:     []string{"GIT_USERNAME"},
			Usage:       "The git username to use for git HTTPS authentication",
			Destination: &app.gitUsernameOverride,
		},
		&cli.StringFlag{
			Name:        "git-password",
			EnvVars:     []string{"GIT_PASSWORD"},
			Usage:       "The git password to use for git HTTPS authentication",
			Destination: &app.gitPasswordOverride,
		},
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"V"},
			EnvVars:     []string{"EARTHLY_VERBOSE"},
			Usage:       "Enable verbose logging",
			Destination: &app.verbose,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"D"},
			EnvVars:     []string{"EARTHLY_DEBUG"},
			Usage:       "Enable debug mode. This flag also turns on the debug mode of buildkitd, which may cause it to restart",
			Destination: &app.debug,
			Hidden:      true, // For development purposes only.
		},
		&cli.BoolFlag{
			Name:        "profiler",
			EnvVars:     []string{"EARTHLY_PROFILER"},
			Usage:       "Enable the profiler",
			Destination: &app.enableProfiler,
			Hidden:      true, // Dev purposes only.
		},
		&cli.StringFlag{
			Name:        "buildkit-host",
			Value:       "",
			EnvVars:     []string{"EARTHLY_BUILDKIT_HOST"},
			Usage:       wrap("The URL to use for connecting to a buildkit host. ", "If empty, earthly will attempt to start a buildkitd instance via docker run"),
			Destination: &app.buildkitHost,
		},
		&cli.StringFlag{
			Name:        "server",
			Value:       "https://api.earthly.dev",
			EnvVars:     []string{"EARTHLY_SERVER_ADDRESS"},
			Usage:       "API server override for dev purposes",
			Destination: &app.cloudHTTPAddr,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "grpc",
			Value:       "ci.earthly.dev:443",
			EnvVars:     []string{"EARTHLY_GRPC_ADDRESS"},
			Usage:       "gRPC server override for dev purposes",
			Destination: &app.cloudGRPCAddr,
			Hidden:      true, // Internal.
		},
		&cli.BoolFlag{
			Name:        "grpc-insecure",
			EnvVars:     []string{"EARTHLY_GRPC_INSECURE"},
			Usage:       "Makes gRPC connections insecure for dev purposes",
			Destination: &app.cloudGRPCInsecure,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "satellite-address",
			Value:       containerutil.SatelliteAddress,
			EnvVars:     []string{"EARTHLY_SATELLITE_ADDRESS"},
			Usage:       "Satellite address override for dev purposes",
			Destination: &app.satelliteAddress,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "request-id",
			EnvVars:     []string{"EARTHLY_REQUEST_ID"},
			Usage:       "Override a request ID to the backend API. Useful for debugging or manually retrying a request.",
			Destination: &app.requestID,
			Hidden:      true, // Internal
		},
		&cli.BoolFlag{
			Name:        "no-buildkit-update",
			EnvVars:     []string{"EARTHLY_NO_BUILDKIT_UPDATE"},
			Usage:       "Disable the automatic update of buildkitd",
			Destination: &app.noBuildkitUpdate,
			Hidden:      true, // Internal.
		},
		&cli.BoolFlag{
			EnvVars:     []string{"EARTHLY_DISABLE_ANALYTICS", "DO_NOT_TRACK"},
			Usage:       "Disable collection of analytics",
			Destination: &app.disableAnalytics,
		},
		&cli.StringFlag{
			Name:        "version-flag-overrides",
			EnvVars:     []string{"EARTHLY_VERSION_FLAG_OVERRIDES"},
			Usage:       "Apply additional flags after each VERSION command across all Earthfiles, multiple flags can be separated by commas",
			Destination: &app.featureFlagOverrides,
			Hidden:      true, // used for feature-flipping from ./earthly dev script
		},
		&cli.StringFlag{
			Name:        envFileFlag,
			EnvVars:     []string{"EARTHLY_ENV_FILE"},
			Usage:       "Use values from this file as earthly environment variables, buildargs, or secrets",
			Value:       defaultEnvFile,
			Destination: &app.envFile,
		},
		&cli.BoolFlag{
			Name:        "logstream",
			EnvVars:     []string{"EARTHLY_LOGSTREAM"},
			Usage:       "Enable log streaming only locally",
			Destination: &app.logstream,
			Hidden:      true, // Internal.
		},
		&cli.BoolFlag{
			Name:        "logstream-upload",
			EnvVars:     []string{"EARTHLY_LOGSTREAM_UPLOAD"},
			Usage:       "Enable log stream uploading",
			Destination: &app.logstreamUpload,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "logstream-debug-file",
			EnvVars:     []string{"EARTHLY_LOGSTREAM_DEBUG_FILE"},
			Usage:       "Enable log streaming debugging output to a file",
			Destination: &app.logstreamDebugFile,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "logstream-debug-manifest-file",
			EnvVars:     []string{"EARTHLY_LOGSTREAM_DEBUG_MANIFEST_FILE"},
			Usage:       "Enable log streaming manifest debugging output to a file",
			Destination: &app.logstreamDebugManifestFile,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "build-id",
			EnvVars:     []string{"EARTHLY_BUILD_ID"},
			Usage:       "The build ID to use for identifying the build in Earthly Cloud. If not specified, a random ID will be generated",
			Destination: &app.buildID,
			Hidden:      true, // Internal.
		},
	}
}

func (app *earthlyApp) buildFlags() []cli.Flag {
	defaultInstallationName := DefaultInstallationName
	if defaultInstallationName == "" {
		defaultInstallationName = "earthly"
	}
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "platform",
			EnvVars: []string{"EARTHLY_PLATFORMS"},
			Usage:   "Specify the target platform to build for",
			Value:   &app.platformsStr,
		},
		&cli.StringSliceFlag{
			Name:    "build-arg",
			EnvVars: []string{"EARTHLY_BUILD_ARGS"},
			Usage:   "A build arg override, specified as <key>=[<value>]",
			Value:   &app.buildArgs,
			Hidden:  true, // Deprecated
		},
		&cli.StringSliceFlag{
			Name:    "secret",
			Aliases: []string{"s"},
			EnvVars: []string{"EARTHLY_SECRETS"},
			Usage:   "A secret override, specified as <key>=[<value>]",
			Value:   &app.secrets,
		},
		&cli.StringSliceFlag{
			Name:    "secret-file",
			EnvVars: []string{"EARTHLY_SECRET_FILES"},
			Usage:   "A secret override, specified as <key>=<path>",
			Value:   &app.secretFiles,
		},
		&cli.BoolFlag{
			Name:        "artifact",
			Aliases:     []string{"a"},
			Usage:       "Output specified artifact; a wildcard (*) can be used to output all artifacts",
			Destination: &app.artifactMode,
		},
		&cli.BoolFlag{
			Name:        "image",
			Usage:       "Output only docker image of the specified target",
			Destination: &app.imageMode,
		},
		&cli.BoolFlag{
			Name:        "pull",
			EnvVars:     []string{"EARTHLY_PULL"},
			Usage:       "Force pull any referenced Docker images",
			Destination: &app.pull,
			Hidden:      true, // Experimental
		},
		&cli.BoolFlag{
			Name:        "push",
			EnvVars:     []string{"EARTHLY_PUSH"},
			Usage:       "Push docker images and execute RUN --push commands",
			Destination: &app.push,
		},
		&cli.BoolFlag{
			Name:        "ci",
			EnvVars:     []string{"EARTHLY_CI"},
			Usage:       wrap("Execute in CI mode. ", "Implies --no-output --strict"),
			Destination: &app.ci,
		},
		&cli.BoolFlag{
			Name:        "output",
			EnvVars:     []string{"EARTHLY_OUTPUT"},
			Usage:       "Allow artifacts or images to be output, even when running under --ci mode",
			Destination: &app.output,
		},
		&cli.BoolFlag{
			Name:        "no-output",
			EnvVars:     []string{"EARTHLY_NO_OUTPUT"},
			Usage:       wrap("Do not output artifacts or images", "(using --push is still allowed)"),
			Destination: &app.noOutput,
		},
		&cli.BoolFlag{
			Name:        "no-cache",
			EnvVars:     []string{"EARTHLY_NO_CACHE"},
			Usage:       "Do not use cache while building",
			Destination: &app.noCache,
		},
		&cli.BoolFlag{
			Name:        "allow-privileged",
			Aliases:     []string{"P"},
			EnvVars:     []string{"EARTHLY_ALLOW_PRIVILEGED"},
			Usage:       "Allow build to use the --privileged flag in RUN commands",
			Destination: &app.allowPrivileged,
		},
		&cli.StringFlag{
			Name:        "org",
			EnvVars:     []string{"EARTHLY_ORG"},
			Usage:       wrap("The name of the organization the satellite belongs to. ", "Required when using --satellite and user is a member of multiple organizations."),
			Required:    false,
			Destination: &app.orgName,
		},
		&cli.StringFlag{
			Name:        "satellite",
			Aliases:     []string{"sat"},
			EnvVars:     []string{"EARTHLY_SATELLITE"},
			Usage:       "The name of satellite to use for this build.",
			Required:    false,
			Destination: &app.satelliteName,
		},
		&cli.BoolFlag{
			Name:        "no-satellite",
			Aliases:     []string{"no-sat"},
			EnvVars:     []string{"EARTHLY_NO_SATELLITE"},
			Usage:       "Disables the use of a selected satellite for this build.",
			Required:    false,
			Destination: &app.noSatellite,
		},
		&cli.StringFlag{
			Name:        "tlscert",
			Value:       "./certs/earthly_cert.pem",
			EnvVars:     []string{"EARTHLY_TLS_CERT"},
			Usage:       wrap("The path to the client TLS cert", "If relative, will be interpreted as relative to the ~/.earthly folder."),
			Destination: &app.certPath,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "tlskey",
			Value:       "./certs/earthly_key.pem",
			EnvVars:     []string{"EARTHLY_TLS_KEY"},
			Usage:       wrap("The path to the client TLS key.", "If relative, will be interpreted as relative to the ~/.earthly folder."),
			Destination: &app.keyPath,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "buildkit-image",
			Value:       DefaultBuildkitdImage,
			EnvVars:     []string{"EARTHLY_BUILDKIT_IMAGE"},
			Usage:       "The docker image to use for the buildkit daemon",
			Destination: &app.buildkitdImage,
		},
		&cli.StringFlag{
			Name:        "buildkit-container-name",
			Value:       defaultInstallationName + DefaultBuildkitdContainerSuffix,
			EnvVars:     []string{"EARTHLY_CONTAINER_NAME"},
			Usage:       "The docker container name to use for the buildkit daemon",
			Destination: &app.containerName,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "buildkit-volume-name",
			Value:       defaultInstallationName + DefaultBuildkitdVolumeSuffix,
			EnvVars:     []string{"EARTHLY_VOLUME_NAME"},
			Usage:       "The docker volume name to use for the buildkit daemon cache",
			Destination: &app.buildkitdSettings.VolumeName,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "installation-name",
			Value:       defaultInstallationName,
			EnvVars:     []string{"EARTHLY_INSTALLATION_NAME"},
			Usage:       "The earthly installation name to use when naming the buildkit container, the docker volume and the ~/.earthly directory",
			Destination: &app.installationName,
		},
		&cli.StringSliceFlag{
			Name:    "cache-from",
			EnvVars: []string{"EARTHLY_CACHE_FROM"},
			Usage:   "Remote docker image tags to use as readonly explicit cache (experimental)",
			Value:   &app.cacheFrom,
			Hidden:  true, // Experimental
		},
		&cli.StringFlag{
			Name:        "remote-cache",
			EnvVars:     []string{"EARTHLY_REMOTE_CACHE"},
			Usage:       "A remote docker image tag use as explicit cache",
			Destination: &app.remoteCache,
		},
		&cli.BoolFlag{
			Name:        "max-remote-cache",
			EnvVars:     []string{"EARTHLY_MAX_REMOTE_CACHE"},
			Usage:       "Saves all intermediate images too in the remote cache",
			Destination: &app.maxRemoteCache,
		},
		&cli.BoolFlag{
			Name:        "save-inline-cache",
			EnvVars:     []string{"EARTHLY_SAVE_INLINE_CACHE"},
			Usage:       "Enable cache inlining when pushing images",
			Destination: &app.saveInlineCache,
		},
		&cli.BoolFlag{
			Name:        "use-inline-cache",
			EnvVars:     []string{"EARTHLY_USE_INLINE_CACHE"},
			Usage:       wrap("Attempt to use any inline cache that may have been previously pushed ", "uses image tags referenced by SAVE IMAGE --push or SAVE IMAGE --cache-from"),
			Destination: &app.useInlineCache,
		},
		&cli.BoolFlag{
			Name:        "interactive",
			Aliases:     []string{"i"},
			EnvVars:     []string{"EARTHLY_INTERACTIVE"},
			Usage:       "Enable interactive debugging",
			Destination: &app.interactiveDebugging,
		},
		&cli.BoolFlag{
			Name:        "no-fake-dep",
			EnvVars:     []string{"EARTHLY_NO_FAKE_DEP"},
			Usage:       "Internal feature flag for fake-dep",
			Destination: &app.noFakeDep,
			Hidden:      true, // Internal.
		},
		&cli.BoolFlag{
			Name:        "strict",
			EnvVars:     []string{"EARTHLY_STRICT"},
			Usage:       "Disallow usage of features that may create unrepeatable builds",
			Destination: &app.strict,
		},
		&cli.BoolFlag{
			Name:        "global-wait-end",
			EnvVars:     []string{"EARTHLY_GLOBAL_WAIT_END"},
			Usage:       "enables global wait-end code in place of builder code",
			Destination: &app.globalWaitEnd,
			Hidden:      true, // used to force code-coverage of future builder.go refactor (once we remove support for 0.6)
		},
	}
}

func (app *earthlyApp) hiddenBuildFlags() []cli.Flag {
	_, isAutocomplete := os.LookupEnv("COMP_LINE")
	flags := app.buildFlags()
	if isAutocomplete {
		// Don't hide the build flags for autocomplete.
		return flags
	}
	for _, flag := range flags {
		switch f := flag.(type) {
		case *cli.StringSliceFlag:
			f.Hidden = true
		case *cli.StringFlag:
			f.Hidden = true
		case *cli.BoolFlag:
			f.Hidden = true
		case *cli.IntFlag:
			f.Hidden = true
		}
	}
	return flags
}
