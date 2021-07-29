package buildcontext

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/builder/dockerignore"
	"github.com/pkg/errors"
)

// EarthIgnoreFile is the name of the earthly ignore file.
const EarthIgnoreFile = ".earthignore"
const EarthlyIgnoreFile = ".earthlyignore"

// ImplicitExcludes is a list of implicit patterns to exclude.
var ImplicitExcludes = []string{
	".tmp-earthly-out/",
	"build.earth",
	"Earthfile",
	EarthIgnoreFile,
	EarthlyIgnoreFile,
}

func readExcludes(dir string) ([]string, error) {
	ignoreFile := EarthIgnoreFile

	// if non (doesNotExist) errors appear then we need to return them.
	earthExists, err := ignoreFileExists(dir, EarthIgnoreFile)
	if err != nil {
		return ImplicitExcludes, err
	}
	earthlyExists, err := ignoreFileExists(dir, EarthlyIgnoreFile)
	if err != nil {
		return ImplicitExcludes, err
	}

	// Check which ones exists and which don't
	if earthExists && earthlyExists {
		// if both exist then throw an error
		return ImplicitExcludes, errors.New("both .earthignore and .earthlyignore exist - please remove one")
	} else if earthExists == earthlyExists {
		// return just ImplicitExcludes if neither of them exist
		return ImplicitExcludes, nil
	} else if earthlyExists {
		ignoreFile := EarthlyIgnoreFile
	}

	filePath := filepath.Join(dir, ignoreFile)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "read %s", filePath)
	}
	excludes, err := dockerignore.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "parse %s", filePath)
	}
	return append(excludes, ImplicitExcludes...), nil
}

// cleanest way I could imagine to iterate and check if both of them exist - If the file is a bash/executable file then it could return an error other than (os.IsNotExist(err))
func ignoreFileExists(dir, file string) (exists bool, err error) {
	filePath := filepath.Join(dir, file)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
