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
	version, err := ast.MustParseVersion(buildFilePath, false)
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

	warningStrs, err := ftrs.ProcessFlags()
	if err != nil {
		return nil, err
	}

	if len(warningStrs) > 0 {
		console.Printf("NOTE: The %s feature is enabled by default under VERSION %s, and can be safely removed from the VERSION command", strings.Join(warningStrs, ", "), ftrs.Version())
	}

	err = features.ApplyFlagOverrides(ftrs, featureFlagOverrides)
	if err != nil {
		return nil, err
	}

	return ftrs, nil
}
