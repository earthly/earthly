package variables

import (
	"fmt"
	"strconv"
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
func BuiltinArgs(target domain.Target, platr *platutil.Resolver, gitMeta *gitutil.GitMetadata, defaultArgs DefaultArgs, ftrs *features.Features, push bool, ci bool, earthlyCIRunner bool) *Scope {
	ret := NewScope()
	ret.Add(arg.EarthlyTarget, target.StringCanonical())
	ret.Add(arg.EarthlyTargetProject, target.ProjectCanonical())
	targetNoTag := target
	targetNoTag.Tag = ""
	ret.Add(arg.EarthlyTargetProjectNoTag, targetNoTag.ProjectCanonical())
	ret.Add(arg.EarthlyTargetName, target.Target)

	setTargetTag(ret, target, gitMeta)

	if platr != nil {
		SetPlatformArgs(ret, platr)
		setUserPlatformArgs(ret, platr)
		if ftrs.NewPlatform {
			setNativePlatformArgs(ret, platr)
		}
	}

	if ftrs.WaitBlock {
		ret.Add(arg.EarthlyPush, fmt.Sprintf("%t", push))
	}

	if ftrs.EarthlyVersionArg {
		ret.Add(arg.EarthlyVersion, defaultArgs.EarthlyVersion)
		ret.Add(arg.EarthlyBuildSha, defaultArgs.EarthlyBuildSha)
	}

	if ftrs.EarthlyCIArg {
		ret.Add(arg.EarthlyCI, fmt.Sprintf("%t", ci))
	}

	if ftrs.EarthlyLocallyArg {
		SetLocally(ret, false)
	}

	if gitMeta != nil {
		ret.Add(arg.EarthlyGitHash, gitMeta.Hash)
		ret.Add(arg.EarthlyGitShortHash, gitMeta.ShortHash)
		branch := ""
		if len(gitMeta.Branch) > 0 {
			branch = gitMeta.Branch[0]
		}
		ret.Add(arg.EarthlyGitBranch, branch)
		tag := ""
		if len(gitMeta.Tags) > 0 {
			tag = gitMeta.Tags[0]
		}
		ret.Add(arg.EarthlyGitTag, tag)
		ret.Add(arg.EarthlyGitOriginURL, gitMeta.RemoteURL)
		ret.Add(arg.EarthlyGitOriginURLScrubbed, stringutil.ScrubCredentials(gitMeta.RemoteURL))
		ret.Add(arg.EarthlyGitProjectName, getProjectName(gitMeta.RemoteURL))
		ret.Add(arg.EarthlyGitCommitTimestamp, gitMeta.CommitterTimestamp)

		if ftrs.GitCommitAuthorTimestamp {
			ret.Add(arg.EarthlyGitCommitAuthorTimestamp, gitMeta.AuthorTimestamp)
		}
		if gitMeta.CommitterTimestamp == "" {
			ret.Add(arg.EarthlySourceDateEpoch, "0")
		} else {
			ret.Add(arg.EarthlySourceDateEpoch, gitMeta.CommitterTimestamp)
		}
		if ftrs.EarthlyGitAuthorArgs {
			ret.Add(arg.EarthlyGitAuthor, gitMeta.AuthorEmail)
			ret.Add(arg.EarthlyGitCoAuthors, strings.Join(gitMeta.CoAuthors, " "))
		}
		if ftrs.GitAuthorEmailNameArgs {
			if gitMeta.AuthorName != "" && gitMeta.AuthorEmail != "" {
				ret.Add(arg.EarthlyGitAuthor, fmt.Sprintf("%s <%s>", gitMeta.AuthorName, gitMeta.AuthorEmail))
			}
			ret.Add(arg.EarthlyGitAuthorEmail, gitMeta.AuthorEmail)
			ret.Add(arg.EarthlyGitAuthorName, gitMeta.AuthorName)
		}

		if ftrs.GitRefs {
			ret.Add(arg.EarthlyGitRefs, strings.Join(gitMeta.Refs, " "))
		}

		if ftrs.GitMessages {
			ret.Add(arg.EarthlyGitMessage, gitMeta.Message)
			ret.Add(arg.EarthlyGitMessageFull, gitMeta.FullMessage)
		}
	} else {
		// Ensure SOURCE_DATE_EPOCH is always available
		ret.Add(arg.EarthlySourceDateEpoch, "0")
	}

	if ftrs.EarthlyCIRunnerArg {
		ret.Add(arg.EarthlyCIRunner, strconv.FormatBool(earthlyCIRunner))
	}
	return ret
}

// SetPlatformArgs sets the platform-specific built-in args to a specific platform.
func SetPlatformArgs(s *Scope, platr *platutil.Resolver) {
	platform := platr.Materialize(platr.Current())
	llbPlatform := platr.ToLLBPlatform(platform)
	s.Add(arg.TargetPlatform, platform.String())
	s.Add(arg.TargetOS, llbPlatform.OS)
	s.Add(arg.TargetArch, llbPlatform.Architecture)
	s.Add(arg.TargetVariant, llbPlatform.Variant)
}

func setUserPlatformArgs(s *Scope, platr *platutil.Resolver) {
	platform := platr.LLBUser()
	s.Add(arg.UserPlatform, platforms.Format(platform))
	s.Add(arg.UserOS, platform.OS)
	s.Add(arg.UserArch, platform.Architecture)
	s.Add(arg.UserVariant, platform.Variant)
}

func setNativePlatformArgs(s *Scope, platr *platutil.Resolver) {
	platform := platr.LLBNative()
	s.Add(arg.NativePlatform, platforms.Format(platform))
	s.Add(arg.NativeOS, platform.OS)
	s.Add(arg.NativeArch, platform.Architecture)
	s.Add(arg.NativeVariant, platform.Variant)
}

// SetLocally sets the locally built-in arg value
func SetLocally(s *Scope, locally bool) {
	s.Add(arg.EarthlyLocally, fmt.Sprintf("%v", locally))
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

func setTargetTag(ret *Scope, target domain.Target, gitMeta *gitutil.GitMetadata) {
	// We prefer branch for these tags if the build is triggered from an action on a branch (pr / push)
	// https://github.com/earthly/cloud-issues/issues/11#issuecomment-1467308267
	if gitMeta != nil && gitMeta.BranchOverrideTagArg && len(gitMeta.Branch) > 0 {
		branch := gitMeta.Branch[0]
		ret.Add(arg.EarthlyTargetTag, branch)
		ret.Add(arg.EarthlyTargetTagDocker, llbutil.DockerTagSafe(branch))
		return
	}
	ret.Add(arg.EarthlyTargetTag, target.Tag)
	ret.Add(arg.EarthlyTargetTagDocker, llbutil.DockerTagSafe(target.Tag))
}
