package buildcontext

import (
	"os"
	"path/filepath"

	"github.com/earthly/earthly/util/fileutil"

	"github.com/moby/buildkit/frontend/dockerfile/dockerignore"
	"github.com/pkg/errors"
)

const earthIgnoreFile = ".earthignore"
const earthlyIgnoreFile = ".earthlyignore"

var errDuplicateIgnoreFile = errors.New("both .earthignore and .earthlyignore exist - please remove one")

// ImplicitExcludes is a list of implicit patterns to exclude.
var ImplicitExcludes = []string{
	".tmp-earthly-out/",
	"build.earth",
	"Earthfile",
	earthIgnoreFile,
	earthlyIgnoreFile,
}

func readExcludes(dir string, noImplicitIgnore bool) ([]string, error) {
	var ignoreFile = earthIgnoreFile

	//earthIgnoreFile
	var earthIgnoreFilePath = filepath.Join(dir, earthIgnoreFile)
	earthExists, err := fileutil.FileExists(earthIgnoreFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to check if %s exists", earthIgnoreFilePath)
	}

	//earthlyIgnoreFile
	var earthlyIgnoreFilePath = filepath.Join(dir, earthlyIgnoreFile)
	earthlyExists, err := fileutil.FileExists(earthlyIgnoreFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to check if %s exists", earthlyIgnoreFilePath)
	}

	defaultExcludes := ImplicitExcludes
	if noImplicitIgnore {
		defaultExcludes = []string{}
	}

	// Check which ones exists and which don't
	if earthExists && earthlyExists {
		// if both exist then throw an error
		return defaultExcludes, errDuplicateIgnoreFile
	} else if earthExists == earthlyExists {
		// return just ImplicitExcludes if neither of them exist
		return defaultExcludes, nil
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
	return append(excludes, defaultExcludes...), nil
}
