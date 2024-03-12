package buildcontext

import (
	"fmt"
	"strings"

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

	ftrs, hasVersion, err := features.Get(version)
	if err != nil {
		return nil, err
	}
	if !hasVersion {
		return nil, fmt.Errorf("No version specified in %s/Earthfile", projectRef)
	}

	warningStrs, err := ftrs.Adjust()
	if err != nil {
		return nil, err
	}

	console.Printf("Note that the feature flag %s is not necessary for version %s, since the feature was enabled by default. You may remove this flag from your Earthfile.", strings.Join(warningStrs, ", "), ftrs.Version())

	err = features.ApplyFlagOverrides(ftrs, featureFlagOverrides)
	if err != nil {
		return nil, err
	}

	return ftrs, nil
}
