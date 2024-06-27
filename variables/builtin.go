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
	ret.Add(arg.EarthlyTarget, NewStringVariable(target.StringCanonical()))
	ret.Add(arg.EarthlyTargetProject, NewStringVariable(target.ProjectCanonical()))
	targetNoTag := target
	targetNoTag.Tag = ""
	ret.Add(arg.EarthlyTargetProjectNoTag, NewStringVariable(targetNoTag.ProjectCanonical()))
	ret.Add(arg.EarthlyTargetName, NewStringVariable(target.Target))

	setTargetTag(ret, target, gitMeta)

	if platr != nil {
		SetPlatformArgs(ret, platr)
		setUserPlatformArgs(ret, platr)
		if ftrs.NewPlatform {
			setNativePlatformArgs(ret, platr)
		}
	}

	if ftrs.WaitBlock {
		ret.Add(arg.EarthlyPush, NewStringVariable(fmt.Sprintf("%t", push)))
	}

	if ftrs.EarthlyVersionArg {
		ret.Add(arg.EarthlyVersion, NewStringVariable(defaultArgs.EarthlyVersion))
		ret.Add(arg.EarthlyBuildSha, NewStringVariable(defaultArgs.EarthlyBuildSha))
	}

	if ftrs.EarthlyCIArg {
		ret.Add(arg.EarthlyCI, NewStringVariable(fmt.Sprintf("%t", ci)))
	}

	if ftrs.EarthlyLocallyArg {
		SetLocally(ret, false)
	}

	if gitMeta != nil {
		ret.Add(arg.EarthlyGitHash, NewStringVariable(gitMeta.Hash))
		ret.Add(arg.EarthlyGitShortHash, NewStringVariable(gitMeta.ShortHash))
		branch := ""
		if len(gitMeta.Branch) > 0 {
			branch = gitMeta.Branch[0]
		}
		ret.Add(arg.EarthlyGitBranch, NewStringVariable(branch))
		tag := ""
		if len(gitMeta.Tags) > 0 {
			tag = gitMeta.Tags[0]
		}
		ret.Add(arg.EarthlyGitTag, NewStringVariable(tag))
		ret.Add(arg.EarthlyGitOriginURL, NewStringVariable(gitMeta.RemoteURL))
		ret.Add(arg.EarthlyGitOriginURLScrubbed, NewStringVariable(stringutil.ScrubCredentials(gitMeta.RemoteURL)))
		ret.Add(arg.EarthlyGitProjectName, NewStringVariable(getProjectName(gitMeta.RemoteURL)))
		ret.Add(arg.EarthlyGitCommitTimestamp, NewStringVariable(gitMeta.CommitterTimestamp))

		if ftrs.GitCommitAuthorTimestamp {
			ret.Add(arg.EarthlyGitCommitAuthorTimestamp, NewStringVariable(gitMeta.AuthorTimestamp))
		}
		if gitMeta.CommitterTimestamp == "" {
			ret.Add(arg.EarthlySourceDateEpoch, NewStringVariable("0"))
		} else {
			ret.Add(arg.EarthlySourceDateEpoch, NewStringVariable(gitMeta.CommitterTimestamp))
		}
		if ftrs.EarthlyGitAuthorArgs {
			ret.Add(arg.EarthlyGitAuthor, NewStringVariable(gitMeta.AuthorEmail))
			ret.Add(arg.EarthlyGitCoAuthors, NewStringVariable(strings.Join(gitMeta.CoAuthors, " ")))
		}
		if ftrs.GitAuthorEmailNameArgs {
			if gitMeta.AuthorName != "" && gitMeta.AuthorEmail != "" {
				ret.Add(arg.EarthlyGitAuthor, NewStringVariable(fmt.Sprintf("%s <%s>", gitMeta.AuthorName, gitMeta.AuthorEmail)))
			}
			ret.Add(arg.EarthlyGitAuthorEmail, NewStringVariable(gitMeta.AuthorEmail))
			ret.Add(arg.EarthlyGitAuthorName, NewStringVariable(gitMeta.AuthorName))
		}

		if ftrs.GitRefs {
			ret.Add(arg.EarthlyGitRefs, NewStringVariable(strings.Join(gitMeta.Refs, " ")))
		}
	} else {
		// Ensure SOURCE_DATE_EPOCH is always available
		ret.Add(arg.EarthlySourceDateEpoch, NewStringVariable("0"))
	}

	if ftrs.EarthlyCIRunnerArg {
		ret.Add(arg.EarthlyCIRunner, NewStringVariable(strconv.FormatBool(earthlyCIRunner)))
	}
	return ret
}

// SetPlatformArgs sets the platform-specific built-in args to a specific platform.
func SetPlatformArgs(s *Scope, platr *platutil.Resolver) {
	platform := platr.Materialize(platr.Current())
	llbPlatform := platr.ToLLBPlatform(platform)
	s.Add(arg.TargetPlatform, NewStringVariable(platform.String()))
	s.Add(arg.TargetOS, NewStringVariable(llbPlatform.OS))
	s.Add(arg.TargetArch, NewStringVariable(llbPlatform.Architecture))
	s.Add(arg.TargetVariant, NewStringVariable(llbPlatform.Variant))
}

func setUserPlatformArgs(s *Scope, platr *platutil.Resolver) {
	platform := platr.LLBUser()
	s.Add(arg.UserPlatform, NewStringVariable(platforms.Format(platform)))
	s.Add(arg.UserOS, NewStringVariable(platform.OS))
	s.Add(arg.UserArch, NewStringVariable(platform.Architecture))
	s.Add(arg.UserVariant, NewStringVariable(platform.Variant))
}

func setNativePlatformArgs(s *Scope, platr *platutil.Resolver) {
	platform := platr.LLBNative()
	s.Add(arg.NativePlatform, NewStringVariable(platforms.Format(platform)))
	s.Add(arg.NativeOS, NewStringVariable(platform.OS))
	s.Add(arg.NativeArch, NewStringVariable(platform.Architecture))
	s.Add(arg.NativeVariant, NewStringVariable(platform.Variant))
}

// SetLocally sets the locally built-in arg value
func SetLocally(s *Scope, locally bool) {
	s.Add(arg.EarthlyLocally, NewStringVariable(fmt.Sprintf("%v", locally)))
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
		ret.Add(arg.EarthlyTargetTag, NewStringVariable(branch))
		ret.Add(arg.EarthlyTargetTagDocker, NewStringVariable(llbutil.DockerTagSafe(branch)))
		return
	}
	ret.Add(arg.EarthlyTargetTag, NewStringVariable(target.Tag))
	ret.Add(arg.EarthlyTargetTagDocker, NewStringVariable(llbutil.DockerTagSafe(target.Tag)))
}
