package features

import (
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
	ReferencedSaveOnly         bool `long:"referenced-save-only" description:"only save artifacts that are directly referenced"`
	UseCopyIncludePatterns     bool `long:"use-copy-include-patterns" description:"specify an include pattern to buildkit when performing copies"`
	ForIn                      bool `long:"for-in" description:"allow the use of the FOR command"`
	RequireForceForUnsafeSaves bool `long:"require-force-for-unsafe-saves" description:"require the --force flag when saving to path outside of current path"`
	NoImplicitIgnore           bool `long:"no-implicit-ignore" description:"disable implicit ignore rules to exclude .tmp-earthly-out/, build.earth, Earthfile, .earthignore and .earthlyignore when resolving local context"`
	CheckDuplicateImages       bool `long:"check-duplicate-images" description:"check for duplicate images during output"`
	EarthlyVersionArg          bool `long:"earthly-version-arg" description:"includes EARTHLY_VERSION and EARTHLY_BUILD_SHA ARGs"`
	ExplicitGlobal             bool `long:"explicit-global" description:"require base target args to have explicit settings to be considered global args"`
	UseCacheCommand            bool `long:"use-cache-command" description:"allow use of CACHE command in Earthfiles"`
	UseHostCommand             bool `long:"use-host-command" description:"allow use of HOST command in Earthfiles"`
	ExecAfterParallel          bool `long:"exec-after-parallel" description:"force execution after parallel conversion"`
	UseCopyLink                bool `long:"use-copy-link" description:"use the equivalent of COPY --link for all copy-like operations"`
	ParallelLoad               bool `long:"parallel-load" description:"perform parallel loading of images into WITH DOCKER"`
	NoTarBuildOutput           bool `long:"no-tar-build-output" description:"do not print output when creating a tarball to load into WITH DOCKER"`
	ShellOutAnywhere           bool `long:"shell-out-anywhere" description:"allow shelling-out in the middle of ARGs, or any other command"`
	NewPlatform                bool `long:"new-platform" description:"enable new platform behavior"`
	UseNoManifestList          bool `long:"use-no-manifest-list" description:"enable the SAVE IMAGE --no-manifest-list option"`
	UseChmod                   bool `long:"use-chmod" description:"enable the SAVE IMAGE --no-manifest-list option"`
	UseRegistryForWithDocker   bool `long:"use-registry-for-with-docker" description:"use embedded Docker registry for WITH DOCKER load operations"`
	WaitBlock                  bool `long:"wait-block" description:"enable WITH/END feature, also allows RUN --push mixed with non-push commands"`

	Major int
	Minor int
}

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
	return nil
}

var errUnexpectedArgs = fmt.Errorf("unexpected VERSION arguments; should be VERSION [flags] <major-version>.<minor-version>")

func instrumentVersion(_ string, opt *goflags.Option, s *string) (*string, error) {
	analytics.Count("version-feature-flags", opt.LongName)
	return s, nil // don't modify the flag, just pass it back.
}

// GetFeatures returns a features struct for a particular version
func GetFeatures(version *spec.Version) (*Features, bool, error) {
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

	parsedArgs, err := flagutil.ParseArgsWithValueModifier("VERSION", &ftrs, version.Args, instrumentVersion)
	if err != nil {
		return nil, false, err
	}

	if len(parsedArgs) != 1 {
		return nil, false, errUnexpectedArgs
	}

	majorAndMinor := strings.Split(parsedArgs[0], ".")
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

	// Enable version-specific features.
	if versionAtLeast(ftrs, 0, 5) {
		ftrs.ExecAfterParallel = true
		ftrs.ParallelLoad = true
	}
	if versionAtLeast(ftrs, 0, 6) {
		ftrs.ReferencedSaveOnly = true
		ftrs.UseCopyIncludePatterns = true
		ftrs.ForIn = true
		ftrs.RequireForceForUnsafeSaves = true
		ftrs.NoImplicitIgnore = true
	}
	if versionAtLeast(ftrs, 0, 7) {
		ftrs.ExplicitGlobal = true
		ftrs.CheckDuplicateImages = true
		ftrs.EarthlyVersionArg = true
		ftrs.UseCacheCommand = true
		ftrs.UseHostCommand = true
		ftrs.UseCopyLink = true
		ftrs.NewPlatform = true
		ftrs.NoTarBuildOutput = true
		ftrs.UseNoManifestList = true
		ftrs.UseChmod = true
	}

	return &ftrs, hasVersion, nil
}

// versionAtLeast returns true if the version configured in `ftrs`
// are greater than or equal to the provided major and minor versions.
func versionAtLeast(ftrs Features, majorVersion, minorVersion int) bool {
	return (ftrs.Major > majorVersion) || (ftrs.Major == majorVersion && ftrs.Minor >= minorVersion)
}
