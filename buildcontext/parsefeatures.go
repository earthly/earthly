package buildcontext

import (
	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/features"
)

type buildFile struct {
	path string
	ftrs *features.Features
}

func parseFeatures(buildFilePath string, featureFlagOverrides string, projectRef string, console conslogging.ConsoleLogger) (*features.Features, error) {
	version, err := ast.ParseVersion(buildFilePath, false)
	if err != nil {
		return nil, err
	}

	ftrs, hasVersion, err := features.GetFeatures(version)
	if err != nil {
		return nil, err
	}
	if !hasVersion {
		console.Warnf(
			"Warning: No version specified in %s/Earthfile. Implying VERSION 0.5, which is not the latest available. Please note that in the future, the VERSION command will be required for all Earthfiles.\n", projectRef)
	}

	err = features.ApplyFlagOverrides(ftrs, featureFlagOverrides)
	if err != nil {
		return nil, err
	}

	return ftrs, nil
}
