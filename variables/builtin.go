package variables

import (
	"strings"

	"github.com/containerd/containerd/platforms"
	specs "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/stringutil"
	arg "github.com/earthly/earthly/variables/reserved"
)

// DefaultArgs contains additional builtin ARG values which need
// to be passed in from outside of the scope of this package.
type DefaultArgs struct {
	EarthlyVersion  string
	EarthlyBuildSha string
}

// BuiltinArgs returns a scope containing the builtin args.
func BuiltinArgs(target domain.Target, platform llbutil.Platform, nativePlatform specs.Platform, gitMeta *gitutil.GitMetadata, defaultArgs DefaultArgs, ftrs *features.Features) *Scope {
	ret := NewScope()
	ret.AddInactive(arg.EarthlyTarget, target.StringCanonical())
	ret.AddInactive(arg.EarthlyTargetProject, target.ProjectCanonical())
	targetNoTag := target
	targetNoTag.Tag = ""
	ret.AddInactive(arg.EarthlyTargetProjectNoTag, targetNoTag.ProjectCanonical())
	ret.AddInactive(arg.EarthlyTargetName, target.Target)
	ret.AddInactive(arg.EarthlyTargetTag, target.Tag)
	ret.AddInactive(arg.EarthlyTargetTagDocker, llbutil.DockerTagSafe(target.Tag))
	SetPlatformArgs(ret, platform, nativePlatform)
	setUserPlatformArgs(ret)
	if ftrs.NewPlatform {
		setNativePlatformArgs(ret, nativePlatform)
	}

	if ftrs != nil && ftrs.EarthlyVersionArg {
		ret.AddInactive(arg.EarthlyVersion, defaultArgs.EarthlyVersion)
		ret.AddInactive(arg.EarthlyBuildSha, defaultArgs.EarthlyBuildSha)
	}

	if gitMeta != nil {
		ret.AddInactive(arg.EarthlyGitHash, gitMeta.Hash)
		ret.AddInactive(arg.EarthlyGitShortHash, gitMeta.ShortHash)
		branch := ""
		if len(gitMeta.Branch) > 0 {
			branch = gitMeta.Branch[0]
		}
		ret.AddInactive(arg.EarthlyGitBranch, branch)
		tag := ""
		if len(gitMeta.Tags) > 0 {
			tag = gitMeta.Tags[0]
		}
		ret.AddInactive(arg.EarthlyGitTag, tag)
		ret.AddInactive(arg.EarthlyGitOriginURL, gitMeta.RemoteURL)
		ret.AddInactive(arg.EarthlyGitOriginURLScrubbed, stringutil.ScrubCredentials(gitMeta.RemoteURL))
		ret.AddInactive(arg.EarthlyGitProjectName, getProjectName(gitMeta.RemoteURL))
		ret.AddInactive(arg.EarthlyGitCommitTimestamp, gitMeta.Timestamp)

		if gitMeta.Timestamp == "" {
			ret.AddInactive(arg.EarthlySourceDateEpoch, "0")
		} else {
			ret.AddInactive(arg.EarthlySourceDateEpoch, gitMeta.Timestamp)
		}
	} else {
		// Ensure SOURCE_DATE_EPOCH is always available
		ret.AddInactive(arg.EarthlySourceDateEpoch, "0")
	}
	return ret
}

// SetPlatformArgs sets the platform-specific built-in args to a specific platform.
func SetPlatformArgs(s *Scope, platform llbutil.Platform, nativePlatform specs.Platform) {
	platform = platform.Resolve(nativePlatform)
	llbPlatform := platform.ToLLBPlatform(nativePlatform)
	s.AddInactive(arg.TargetPlatform, platform.String())
	s.AddInactive(arg.TargetOS, llbPlatform.OS)
	s.AddInactive(arg.TargetArch, llbPlatform.Architecture)
	s.AddInactive(arg.TargetVariant, llbPlatform.Variant)
}

func setUserPlatformArgs(s *Scope) {
	platform := platforms.DefaultSpec()
	s.AddInactive(arg.UserPlatform, platforms.Format(platform))
	s.AddInactive(arg.UserOS, platform.OS)
	s.AddInactive(arg.UserArch, platform.Architecture)
	s.AddInactive(arg.UserVariant, platform.Variant)
}

func setNativePlatformArgs(s *Scope, np specs.Platform) {
	s.AddInactive(arg.NativePlatform, platforms.Format(np))
	s.AddInactive(arg.NativeOS, np.OS)
	s.AddInactive(arg.NativeArch, np.Architecture)
	s.AddInactive(arg.NativeVariant, np.Variant)
}

// getProjectName returns the depricated PROJECT_NAME value
func getProjectName(s string) string {
	protocol := "unknown"
	parts := strings.SplitN(s, "://", 2)
	if len(parts) > 1 {
		protocol = parts[0]
		s = parts[1]
	}
	parts = strings.SplitN(s, "@", 2)
	if len(parts) > 1 {
		s = parts[1]
	}
	if protocol == "unknown" {
		s = strings.Replace(s, ":", "/", 1)
	}
	s = strings.TrimSuffix(s, ".git")
	parts = strings.SplitN(s, "/", 2)
	if len(parts) > 1 {
		s = parts[1]
	}
	return s
}
