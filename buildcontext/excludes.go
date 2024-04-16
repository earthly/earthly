package buildcontext

import (
	"os"
	"path/filepath"

	"github.com/earthly/earthly/util/fileutil"
	"github.com/moby/patternmatcher/ignorefile"
	"github.com/pkg/errors"
)

const earthIgnoreFile = ".earthignore"
const earthlyIgnoreFile = ".earthlyignore"
const dockerIgnoreFile = ".dockerignore"

var errDuplicateIgnoreFile = errors.New("both .earthignore and .earthlyignore exist - please remove one")

// ImplicitExcludes is a list of implicit patterns to exclude.
var ImplicitExcludes = []string{
	".tmp-earthly-out/",
	"build.earth",
	"Earthfile",
	earthIgnoreFile,
	earthlyIgnoreFile,
}

func readExcludes(dir string, noImplicitIgnore bool, useDockerIgnore bool) ([]string, error) {
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

	//dockerIgnoreFile
	var dockerIgnoreFilePath = filepath.Join(dir, dockerIgnoreFile)
	dockerExists := false
	if useDockerIgnore {
		dockerExists, err = fileutil.FileExists(dockerIgnoreFilePath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to check if %s exists", dockerIgnoreFilePath)
		}
	}

	defaultExcludes := ImplicitExcludes
	if noImplicitIgnore {
		defaultExcludes = []string{}
	}

	// Check which ones exists and which don't
	if earthExists && earthlyExists {
		// if both exist then throw an error
		return defaultExcludes, errDuplicateIgnoreFile
	}
	if earthExists == earthlyExists {
		if !dockerExists {
			// return just ImplicitExcludes if neither of them exist
			return defaultExcludes, nil
		}
		ignoreFile = dockerIgnoreFile
	} else if earthlyExists {
		ignoreFile = earthlyIgnoreFile
	}

	filePath := filepath.Join(dir, ignoreFile)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "read %s", filePath)
	}
	defer f.Close()
	excludes, err := ignorefile.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "parse %s", filePath)
	}
	return append(excludes, defaultExcludes...), nil
}
