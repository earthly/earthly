package variables

import (
	"fmt"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/gitutil"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/states/dedup"
	"github.com/earthly/earthly/stringutil"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// BuiltinArgs returns a scope containing the builtin args.
func BuiltinArgs(target domain.Target, platform specs.Platform, gitMeta *gitutil.GitMetadata) *Scope {
	ret := NewScope()
	ret.AddInactive("EARTHLY_TARGET", target.StringCanonical())
	ret.AddInactive("EARTHLY_TARGET_PROJECT", target.ProjectCanonical())
	ret.AddInactive("EARTHLY_TARGET_NAME", target.Target)
	ret.AddInactive("EARTHLY_TARGET_TAG", target.Tag)
	ret.AddInactive("EARTHLY_TARGET_TAG_DOCKER", llbutil.DockerTagSafe(target.Tag))
	SetPlatformArgs(ret, platform)

	if gitMeta != nil {
		ret.AddInactive("EARTHLY_GIT_HASH", gitMeta.Hash)
		branch := ""
		if len(gitMeta.Branch) > 0 {
			branch = gitMeta.Branch[0]
		}
		ret.AddInactive("EARTHLY_GIT_BRANCH", branch)
		tag := ""
		if len(gitMeta.Tags) > 0 {
			tag = gitMeta.Tags[0]
		}
		ret.AddInactive("EARTHLY_GIT_TAG", tag)
		ret.AddInactive("EARTHLY_GIT_ORIGIN_URL", gitMeta.RemoteURL)
		ret.AddInactive("EARTHLY_GIT_ORIGIN_URL_SCRUBBED", stringutil.ScrubCredentials(gitMeta.RemoteURL))
		ret.AddInactive("EARTHLY_GIT_PROJECT_NAME", getProjectName(gitMeta.RemoteURL))
	}
	// Note: Please update targetinput.go BuiltinVariables if adding more builtin variables.
	for _, key := range ret.SortedAny() {
		if !dedup.BuiltinVariables[key] {
			panic(fmt.Sprintf("you forgot to add %s to the map of BuiltinVariables", key))
		}
	}
	return ret
}

// SetPlatformArgs sets the platform-specific built-in args to a specific platform.
func SetPlatformArgs(s *Scope, platform specs.Platform) {
	s.AddInactive("TARGETPLATFORM", platforms.Format(platform))
	s.AddInactive("TARGETOS", platform.OS)
	s.AddInactive("TARGETARCH", platform.Architecture)
	s.AddInactive("TARGETVARIANT", platform.Variant)
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
