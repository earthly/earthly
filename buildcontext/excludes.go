package buildcontext

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/builder/dockerignore"
	"github.com/pkg/errors"
)

// EarthIgnoreFile is the name of the earth ignore file.
const EarthIgnoreFile = ".earthignore"

// ImplicitExcludes is a list of implicit patterns to exclude.
var ImplicitExcludes = []string{
	".tmp-earth-out/",
	"build.earth",
	EarthIgnoreFile,
}

func readExcludes(dir string) ([]string, error) {
	filePath := filepath.Join(dir, EarthIgnoreFile)
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// No earthignore file present.
			return ImplicitExcludes, nil
		}
		return nil, errors.Wrapf(err, "read %s", filePath)
	}
	excludes, err := dockerignore.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "parse %s", filePath)
	}
	return append(excludes, ImplicitExcludes...), nil
}
