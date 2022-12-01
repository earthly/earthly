package variables

import (
	"fmt"
	"strings"

	"github.com/containerd/containerd/platforms"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/platutil"
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
func BuiltinArgs(target domain.Target, platr *platutil.Resolver, gitMeta *gitutil.GitMetadata, defaultArgs DefaultArgs, ftrs *features.Features, push bool, ci bool) *Scope {
	ret := NewScope()
	ret.AddInactive(arg.EarthlyTarget, target.StringCanonical())
	ret.AddInactive(arg.EarthlyTargetProject, target.ProjectCanonical())
	targetNoTag := target
	targetNoTag.Tag = ""
	ret.AddInactive(arg.EarthlyTargetProjectNoTag, targetNoTag.ProjectCanonical())
	ret.AddInactive(arg.EarthlyTargetName, target.Target)
	ret.AddInactive(arg.EarthlyTargetTag, target.Tag)
	ret.AddInactive(arg.EarthlyTargetTagDocker, llbutil.DockerTagSafe(target.Tag))
	SetPlatformArgs(ret, platr)
	setUserPlatformArgs(ret, platr)
	if ftrs.NewPlatform {
		setNativePlatformArgs(ret, platr)
	}
	if ftrs.WaitBlock {
		ret.AddInactive(arg.EarthlyPush, fmt.Sprintf("%t", push))
	}

	if ftrs.EarthlyVersionArg {
		ret.AddInactive(arg.EarthlyVersion, defaultArgs.EarthlyVersion)
		ret.AddInactive(arg.EarthlyBuildSha, defaultArgs.EarthlyBuildSha)
	}

	if ftrs.EarthlyCIArg {
		ret.AddInactive(arg.EarthlyCI, fmt.Sprintf("%t", ci))
	}

	if ftrs.EarthlyLocallyArg {
		SetLocally(ret, false)
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
		ret.AddInactive(arg.EarthlyGitCommitTimestamp, gitMeta.CommitterTimestamp)

		if ftrs.GitCommitAuthorTimestamp {
			ret.AddInactive(arg.EarthlyGitCommitAuthorTimestamp, gitMeta.AuthorTimestamp)
		}
		if gitMeta.CommitterTimestamp == "" {
			ret.AddInactive(arg.EarthlySourceDateEpoch, "0")
		} else {
			ret.AddInactive(arg.EarthlySourceDateEpoch, gitMeta.CommitterTimestamp)
		}
		if ftrs.EarthlyGitAuthorArgs {
			ret.AddInactive(arg.EarthlyGitAuthor, gitMeta.Author)
			ret.AddInactive(arg.EarthlyGitCoAuthors, strings.Join(gitMeta.CoAuthors, " "))
		}
	} else {
		// Ensure SOURCE_DATE_EPOCH is always available
		ret.AddInactive(arg.EarthlySourceDateEpoch, "0")
	}
	return ret
}

// SetPlatformArgs sets the platform-specific built-in args to a specific platform.
func SetPlatformArgs(s *Scope, platr *platutil.Resolver) {
	platform := platr.Materialize(platr.Current())
	llbPlatform := platr.ToLLBPlatform(platform)
	s.AddInactive(arg.TargetPlatform, platform.String())
	s.AddInactive(arg.TargetOS, llbPlatform.OS)
	s.AddInactive(arg.TargetArch, llbPlatform.Architecture)
	s.AddInactive(arg.TargetVariant, llbPlatform.Variant)
}

func setUserPlatformArgs(s *Scope, platr *platutil.Resolver) {
	platform := platr.LLBUser()
	s.AddInactive(arg.UserPlatform, platforms.Format(platform))
	s.AddInactive(arg.UserOS, platform.OS)
	s.AddInactive(arg.UserArch, platform.Architecture)
	s.AddInactive(arg.UserVariant, platform.Variant)
}

func setNativePlatformArgs(s *Scope, platr *platutil.Resolver) {
	platform := platr.LLBNative()
	s.AddInactive(arg.NativePlatform, platforms.Format(platform))
	s.AddInactive(arg.NativeOS, platform.OS)
	s.AddInactive(arg.NativeArch, platform.Architecture)
	s.AddInactive(arg.NativeVariant, platform.Variant)
}

// SetLocally sets the locally built-in arg value
func SetLocally(s *Scope, locally bool) {
	s.AddInactive(arg.EarthlyLocally, fmt.Sprintf("%v", locally))
}

// getProjectName returns the deprecated PROJECT_NAME value
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
