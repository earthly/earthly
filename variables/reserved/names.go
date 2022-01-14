package reserved

const (
	EarthlyBuildSha             = "EARTHLY_BUILD_SHA"
	EarthlyGitBranch            = "EARTHLY_GIT_BRANCH"
	EarthlyGitCommitTimestamp   = "EARTHLY_GIT_COMMIT_TIMESTAMP"
	EarthlyGitHash              = "EARTHLY_GIT_HASH"
	EarthlyGitOriginURL         = "EARTHLY_GIT_ORIGIN_URL"
	EarthlyGitOriginURLScrubbed = "EARTHLY_GIT_ORIGIN_URL_SCRUBBED"
	EarthlyGitProjectName       = "EARTHLY_GIT_PROJECT_NAME"
	EarthlyGitShortHash         = "EARTHLY_GIT_SHORT_HASH"
	EarthlyGitTag               = "EARTHLY_GIT_TAG"
	EarthlySourceDateEpoch      = "EARTHLY_SOURCE_DATE_EPOCH"
	EarthlyTarget               = "EARTHLY_TARGET"
	EarthlyTargetName           = "EARTHLY_TARGET_NAME"
	EarthlyTargetProject        = "EARTHLY_TARGET_PROJECT"
	EarthlyTargetProjectNoTag   = "EARTHLY_TARGET_PROJECT_NO_TAG"
	EarthlyTargetTag            = "EARTHLY_TARGET_TAG"
	EarthlyTargetTagDocker      = "EARTHLY_TARGET_TAG_DOCKER"
	EarthlyVersion              = "EARTHLY_VERSION"
	TargetArch                  = "TARGETARCH"
	TargetOS                    = "TARGETOS"
	TargetPlatform              = "TARGETPLATFORM"
	TargetVariant               = "TARGETVARIANT"
	UserArch                    = "USERARCH"
	UserOS                      = "USEROS"
	UserPlatform                = "USERPLATFORM"
	UserVariant                 = "USERVARIANT"
)

var args map[string]struct{}

func init() {
	args = map[string]struct{}{
		EarthlyBuildSha:             struct{}{},
		EarthlyGitBranch:            struct{}{},
		EarthlyGitCommitTimestamp:   struct{}{},
		EarthlyGitHash:              struct{}{},
		EarthlyGitOriginURL:         struct{}{},
		EarthlyGitOriginURLScrubbed: struct{}{},
		EarthlyGitProjectName:       struct{}{},
		EarthlyGitShortHash:         struct{}{},
		EarthlyGitTag:               struct{}{},
		EarthlySourceDateEpoch:      struct{}{},
		EarthlyTarget:               struct{}{},
		EarthlyTargetName:           struct{}{},
		EarthlyTargetProject:        struct{}{},
		EarthlyTargetProjectNoTag:   struct{}{},
		EarthlyTargetTag:            struct{}{},
		EarthlyTargetTagDocker:      struct{}{},
		EarthlyVersion:              struct{}{},
		TargetArch:                  struct{}{},
		TargetOS:                    struct{}{},
		TargetPlatform:              struct{}{},
		TargetVariant:               struct{}{},
		UserArch:                    struct{}{},
		UserOS:                      struct{}{},
		UserPlatform:                struct{}{},
		UserVariant:                 struct{}{},
	}
}

// IsBuiltIn returns true if s is the name of a builtin arg
func IsBuiltIn(s string) bool {
	_, exists := args[s]
	return exists
}
