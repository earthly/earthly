package version

// We use this package to export ldflags main vars to other packages
var (
	Version string
	GitSha  string
	BuiltBy string
)
