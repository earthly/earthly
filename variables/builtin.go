package variables

import (
	"github.com/containerd/containerd/platforms"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/gitutil"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/stringutil"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// BuiltinArgs returns a scope containing the builtin args.
func BuiltinArgs(target domain.Target, platform specs.Platform, gitMeta *gitutil.GitMetadata) *Scope {
	ret := NewScope()
	ret.AddInactive("EARTHLY_TARGET", Var{Type: StringType, Value: target.StringCanonical()})
	ret.AddInactive("EARTHLY_TARGET_PROJECT", Var{Type: StringType, Value: target.ProjectCanonical()})
	ret.AddInactive("EARTHLY_TARGET_NAME", Var{Type: StringType, Value: target.Target})
	ret.AddInactive("EARTHLY_TARGET_TAG", Var{Type: StringType, Value: target.Tag})
	ret.AddInactive("EARTHLY_TARGET_TAG_DOCKER", Var{Type: StringType, Value: llbutil.DockerTagSafe(target.Tag)})
	SetPlatformArgs(ret, platform)

	if gitMeta != nil {
		ret.AddInactive("EARTHLY_GIT_HASH", Var{Type: StringType, Value: gitMeta.Hash})
		branch := ""
		if len(gitMeta.Branch) > 0 {
			branch = gitMeta.Branch[0]
		}
		ret.AddInactive("EARTHLY_GIT_BRANCH", Var{Type: StringType, Value: branch})
		tag := ""
		if len(gitMeta.Tags) > 0 {
			tag = gitMeta.Tags[0]
		}
		ret.AddInactive("EARTHLY_GIT_TAG", Var{Type: StringType, Value: tag})
		ret.AddInactive("EARTHLY_GIT_ORIGIN_URL", Var{Type: StringType, Value: gitMeta.RemoteURL})
		ret.AddInactive("EARTHLY_GIT_ORIGIN_URL_SCRUBBED", Var{Type: StringType, Value: stringutil.ScrubCredentials(gitMeta.RemoteURL)})
		ret.AddInactive("EARTHLY_GIT_PROJECT_NAME", Var{Type: StringType, Value: getProjectName(gitMeta.RemoteURL)})
	}
	return ret
}

// SetPlatformArgs sets the platform-specific built-in args to a specific platform.
func SetPlatformArgs(s *Scope, platform specs.Platform) {
	s.AddInactive("TARGETPLATFORM", Var{Type: StringType, Value: platforms.Format(platform)})
	s.AddInactive("TARGETOS", Var{Type: StringType, Value: platform.OS})
	s.AddInactive("TARGETARCH", Var{Type: StringType, Value: platform.Architecture})
	s.AddInactive("TARGETVARIANT", Var{Type: StringType, Value: platform.Variant})
}
