package buildcontext

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/builder/dockerignore"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/pkg/errors"
)

const earthIgnoreFile = ".earthignore"
const earthlyIgnoreFile = ".earthlyignore"

// ImplicitExcludes is a list of implicit patterns to exclude.
var ImplicitExcludes = []string{
	".tmp-earthly-out/",
	"build.earth",
	"Earthfile",
	earthIgnoreFile,
	earthlyIgnoreFile,
}

func readExcludes(dir string) ([]string, error) {
	var ignoreFile = earthIgnoreFile

	//earthIgnoreFile
	var earthIgnoreFilePath = filepath.Join(dir, earthIgnoreFile)
	earthExists := fileutil.FileExists(earthIgnoreFilePath)

	//earthlyIgnoreFile
	var earthlyIgnoreFilePath = filepath.Join(dir, earthlyIgnoreFile)
	earthlyExists := fileutil.FileExists(earthlyIgnoreFilePath)

	// Check which ones exists and which don't
	if earthExists && earthlyExists {
		// if both exist then throw an error
		return ImplicitExcludes, errors.New("both .earthignore and .earthlyignore exist - please remove one")
	} else if earthExists == earthlyExists {
		// return just ImplicitExcludes if neither of them exist
		return ImplicitExcludes, nil
	} else if earthlyExists {
		ignoreFile = earthlyIgnoreFile
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
