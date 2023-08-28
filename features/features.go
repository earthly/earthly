package features

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	goflags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/util/flagutil"
)

// Features is used to denote which features to flip on or off; this is for use in maintaining
// backwards compatibility
type Features struct {
	// VERSION 0.5
	ExecAfterParallel        bool `long:"exec-after-parallel" description:"force execution after parallel conversion"`
	ParallelLoad             bool `long:"parallel-load" description:"perform parallel loading of images into WITH DOCKER"`
	UseRegistryForWithDocker bool `long:"use-registry-for-with-docker" description:"use embedded Docker registry for WITH DOCKER load operations"`

	// VERSION 0.6
	ForIn                      bool `long:"for-in" description:"allow the use of the FOR command"`
	NoImplicitIgnore           bool `long:"no-implicit-ignore" description:"disable implicit ignore rules to exclude .tmp-earthly-out/, build.earth, Earthfile, .earthignore and .earthlyignore when resolving local context"`
	ReferencedSaveOnly         bool `long:"referenced-save-only" description:"only save artifacts that are directly referenced"`
	RequireForceForUnsafeSaves bool `long:"require-force-for-unsafe-saves" description:"require the --force flag when saving to path outside of current path"`
	UseCopyIncludePatterns     bool `long:"use-copy-include-patterns" description:"specify an include pattern to buildkit when performing copies"`

	// VERSION 0.7
	CheckDuplicateImages     bool `long:"check-duplicate-images" description:"check for duplicate images during output"`
	EarthlyCIArg             bool `long:"ci-arg" description:"include EARTHLY_CI arg"`
	EarthlyGitAuthorArgs     bool `long:"earthly-git-author-args" description:"includes EARTHLY_GIT_AUTHOR and EARTHLY_GIT_CO_AUTHORS ARGs"`
	EarthlyLocallyArg        bool `long:"earthly-locally-arg" description:"includes EARTHLY_LOCALLY ARG"`
	EarthlyVersionArg        bool `long:"earthly-version-arg" description:"includes EARTHLY_VERSION and EARTHLY_BUILD_SHA ARGs"`
	ExplicitGlobal           bool `long:"explicit-global" description:"require base target args to have explicit settings to be considered global args"`
	GitCommitAuthorTimestamp bool `long:"git-commit-author-timestamp" description:"include EARTHLY_GIT_COMMIT_AUTHOR_TIMESTAMP arg"`
	NewPlatform              bool `long:"new-platform" description:"enable new platform behavior"`
	NoTarBuildOutput         bool `long:"no-tar-build-output" description:"do not print output when creating a tarball to load into WITH DOCKER"`
	SaveArtifactKeepOwn      bool `long:"save-artifact-keep-own" description:"always apply the --keep-own flag with SAVE ARTIFACT"`
	ShellOutAnywhere         bool `long:"shell-out-anywhere" description:"allow shelling-out in the middle of ARGs, or any other command"`
	UseCacheCommand          bool `long:"use-cache-command" description:"allow use of CACHE command in Earthfiles"`
	UseChmod                 bool `long:"use-chmod" description:"enable the COPY --chmod option"`
	UseCopyLink              bool `long:"use-copy-link" description:"use the equivalent of COPY --link for all copy-like operations"`
	UseHostCommand           bool `long:"use-host-command" description:"allow use of HOST command in Earthfiles"`
	UseNoManifestList        bool `long:"use-no-manifest-list" description:"enable the SAVE IMAGE --no-manifest-list option"`
	UsePipelines             bool `long:"use-pipelines" description:"enable the PIPELINE and TRIGGER commands"`
	UseProjectSecrets        bool `long:"use-project-secrets" description:"enable project-based secret resolution"`
	WaitBlock                bool `long:"wait-block" description:"enable WITH/END feature, also allows RUN --push mixed with non-push commands"`

	// unreleased
	NoUseRegistryForWithDocker bool `long:"no-use-registry-for-with-docker" description:"disable use-registry-for-with-docker"`
	TryFinally                 bool `long:"try" description:"allow the use of the TRY/FINALLY commands"`
	NoNetwork                  bool `long:"no-network" description:"allow the use of RUN --network=none commands"`
	ArgScopeSet                bool `long:"arg-scope-and-set" description:"enable SET to reassign ARGs and prevent ARGs from being redeclared in the same scope"`
	EarthlyCIRunnerArg         bool `long:"earthly-ci-runner-arg" description:"includes EARTHLY_CI_RUNNER ARG"`
	UseDockerIgnore            bool `long:"use-docker-ignore" description:"fallback to .dockerignore incase .earthlyignore or .earthignore do not exist in a local \"FROM DOCKERFILE\" target"`
	PassArgs                   bool `long:"pass-args" description:"Allow the use of the --pass-arg flag in FROM, BUILD, COPY, WITH DOCKER, and DO commands"`

	Major int
	Minor int
}

type ctxKey struct{}

// Version returns the current version
func (f *Features) Version() string {
	return fmt.Sprintf("%d.%d", f.Major, f.Minor)
}

func parseFlagOverrides(env string) map[string]string {
	env = strings.TrimSpace(env)
	m := map[string]string{}
	if env != "" {
		for _, flag := range strings.Split(env, ",") {
			flagNameAndValue := strings.SplitN(flag, "=", 2)
			var flagValue string
			flagName := strings.TrimSpace(flagNameAndValue[0])
			flagName = strings.TrimPrefix(flagName, "--")
			if len(flagNameAndValue) > 1 {
				flagValue = strings.TrimSpace(flagNameAndValue[1])
			}
			m[flagName] = flagValue
		}
	}
	return m
}

// String returns a string representation of the version and set flags
func (f *Features) String() string {
	if f == nil {
		return "<nil>"
	}

	v := reflect.ValueOf(*f)
	typeOf := v.Type()

	flags := []string{}
	for i := 0; i < typeOf.NumField(); i++ {
		tag := typeOf.Field(i).Tag
		if flagName, ok := tag.Lookup("long"); ok {
			ifaceVal := v.Field(i).Interface()
			if boolVal, ok := ifaceVal.(bool); ok && boolVal {
				flags = append(flags, fmt.Sprintf("--%v", flagName))
			}
		}
	}
	sort.Strings(flags)
	args := []string{"VERSION"}
	if len(flags) > 0 {
		args = append(args, strings.Join(flags, " "))
	}
	args = append(args, fmt.Sprintf("%d.%d", f.Major, f.Minor))
	return strings.Join(args, " ")
}

// ApplyFlagOverrides parses a comma separated list of feature flag overrides (without the -- flag name prefix)
// and sets them in the referenced features.
func ApplyFlagOverrides(ftrs *Features, envOverrides string) error {
	overrides := parseFlagOverrides(envOverrides)

	fieldIndices := map[string]int{}
	typeOf := reflect.ValueOf(*ftrs).Type()
	for i := 0; i < typeOf.NumField(); i++ {
		f := typeOf.Field(i)
		tag := f.Tag
		if flagName, ok := tag.Lookup("long"); ok {
			fieldIndices[flagName] = i
		}
	}

	ftrsStruct := reflect.ValueOf(ftrs).Elem()
	for key := range overrides {
		analytics.Count("override-feature-flags", key)
		i, ok := fieldIndices[key]
		if !ok {
			return fmt.Errorf("unable to set %s: invalid flag", key)
		}
		fv := ftrsStruct.Field(i)
		if fv.IsValid() && fv.CanSet() {
			fv.SetBool(true)
		} else {
			return fmt.Errorf("unable to set %s: field is invalid or cant be set", key)
		}
		ifaceVal := fv.Interface()
		if _, ok := ifaceVal.(bool); ok {
			fv.SetBool(true)
		} else {
			return fmt.Errorf("unable to set %s: only boolean fields are currently supported", key)
		}
	}
	processNegativeFlags(ftrs)
	return nil
}

var errUnexpectedArgs = fmt.Errorf("unexpected VERSION arguments; should be VERSION [flags] <major-version>.<minor-version>")

func instrumentVersion(_ string, opt *goflags.Option, s *string) (*string, error) {
	analytics.Count("version-feature-flags", opt.LongName)
	return s, nil // don't modify the flag, just pass it back.
}

// Get returns a features struct for a particular version
func Get(version *spec.Version) (*Features, bool, error) {
	var ftrs Features
	hasVersion := (version != nil)
	if !hasVersion {
		// If no version is specified, we default to 0.5 (the Earthly version
		// before the VERSION command was introduced).
		version = &spec.Version{
			Args: []string{"0.5"},
		}
	}

	if version.Args == nil {
		return nil, false, errUnexpectedArgs
	}

	parsedArgs, err := flagutil.ParseArgsWithValueModifierAndOptions("VERSION", &ftrs, version.Args, instrumentVersion, goflags.PassDoubleDash|goflags.PassAfterNonOption)
	if err != nil {
		return nil, false, err
	}

	if len(parsedArgs) != 1 {
		return nil, false, errUnexpectedArgs
	}

	versionValueStr := parsedArgs[0]
	majorAndMinor := strings.Split(versionValueStr, ".")
	if len(majorAndMinor) != 2 {
		return nil, false, errUnexpectedArgs
	}
	ftrs.Major, err = strconv.Atoi(majorAndMinor[0])
	if err != nil {
		return nil, false, errors.Wrapf(err, "failed to parse major version %q", majorAndMinor[0])
	}
	ftrs.Minor, err = strconv.Atoi(majorAndMinor[1])
	if err != nil {
		return nil, false, errors.Wrapf(err, "failed to parse minor version %q", majorAndMinor[1])
	}

	if hasVersion {
		analytics.Count("version", versionValueStr)
	} else {
		analytics.Count("version", "missing")
	}

	// Enable version-specific features.
	if versionAtLeast(ftrs, 0, 5) {
		ftrs.ExecAfterParallel = true
		ftrs.ParallelLoad = true
		ftrs.UseRegistryForWithDocker = true
	}
	if versionAtLeast(ftrs, 0, 6) {
		ftrs.ForIn = true
		ftrs.NoImplicitIgnore = true
		ftrs.ReferencedSaveOnly = true
		ftrs.RequireForceForUnsafeSaves = true
		ftrs.UseCopyIncludePatterns = true
	}
	if versionAtLeast(ftrs, 0, 7) {
		ftrs.CheckDuplicateImages = true
		ftrs.EarthlyCIArg = true
		ftrs.EarthlyGitAuthorArgs = true
		ftrs.EarthlyLocallyArg = true
		ftrs.EarthlyVersionArg = true
		ftrs.ExplicitGlobal = true
		ftrs.GitCommitAuthorTimestamp = true
		ftrs.NewPlatform = true
		ftrs.NoTarBuildOutput = true
		ftrs.SaveArtifactKeepOwn = true
		ftrs.ShellOutAnywhere = true
		ftrs.UseCacheCommand = true
		ftrs.UseChmod = true
		ftrs.UseCopyLink = true
		ftrs.UseHostCommand = true
		ftrs.UseNoManifestList = true
		ftrs.UsePipelines = true
		ftrs.UseProjectSecrets = true
		ftrs.WaitBlock = true
	}
	processNegativeFlags(&ftrs)

	if ftrs.ArgScopeSet && !ftrs.ShellOutAnywhere {
		// ArgScopeSet uses new ARG declaration logic that requires
		// ShellOutAnywhere. We're erroring here to ensure that users get that
		// feedback as early as possible.
		return nil, false, errors.New("--arg-scope-and-set requires --shell-out-anywhere")
	}

	return &ftrs, hasVersion, nil
}

// versionAtLeast returns true if the version configured in `ftrs`
// are greater than or equal to the provided major and minor versions.
func versionAtLeast(ftrs Features, majorVersion, minorVersion int) bool {
	return (ftrs.Major > majorVersion) || (ftrs.Major == majorVersion && ftrs.Minor >= minorVersion)
}

func processNegativeFlags(ftrs *Features) {
	if ftrs.NoUseRegistryForWithDocker {
		ftrs.UseRegistryForWithDocker = false
	}
}

// WithContext adds the current *Features into the given context and returns a new context.
// Trying to add the *Features to the context more than once will result in an error.
func (f *Features) WithContext(ctx context.Context) (context.Context, error) {
	if ctx.Value(ctxKey{}) != nil {
		return ctx, errors.New("features is already set")
	}
	return context.WithValue(ctx, ctxKey{}, f), nil
}

// FromContext returns the *Features associated with the ctx.
// If no features is found, nil is returned.
func FromContext(ctx context.Context) *Features {
	if f, ok := ctx.Value(ctxKey{}).(*Features); ok {
		return f
	}
	return nil
}
