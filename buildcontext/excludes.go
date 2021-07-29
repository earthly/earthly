package buildcontext

import (
	"github.com/earthly/earthly/util/fileutil"
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
	var ignoreFile = EarthIgnoreFile

	// if non (doesNotExist) errors appear then we need to return them.
	//EarthIgnoreFile
	var earthIgnoreFilePath = filepath.Join(dir, EarthIgnoreFile)
	earthExists := fileutil.FileExists(earthIgnoreFilePath)

	//EarthlyIgnoreFile
	var earthlyIgnoreFilePath = filepath.Join(dir, EarthlyIgnoreFile)
	earthlyExists := fileutil.FileExists(earthlyIgnoreFilePath)

	// Check which ones exists and which don't
	if earthExists && earthlyExists {
		// if both exist then throw an error
		return ImplicitExcludes, errors.New("both .earthignore and .earthlyignore exist - please remove one")
	} else if earthExists == earthlyExists {
		// return just ImplicitExcludes if neither of them exist
		return ImplicitExcludes, nil
	} else if earthlyExists {
		ignoreFile = EarthlyIgnoreFile
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
