package buildcontext

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/earthly/earthly/domain"
	"github.com/pkg/errors"
)

// detectBuildFile detects whether to use Earthfile, build.earth or Dockerfile.
func detectBuildFile(target domain.Target, localDir string) (string, error) {
	if target.Target == DockerfileMetaTarget {
		return filepath.Join(localDir, "Dockerfile"), nil
	}
	earthfilePath := filepath.Join(localDir, "Earthfile")
	_, err := os.Stat(earthfilePath)
	if os.IsNotExist(err) {
		buildEarthPath := filepath.Join(localDir, "build.earth")
		_, err := os.Stat(buildEarthPath)
		if os.IsNotExist(err) {
			return "", fmt.Errorf(
				"No Earthfile nor build.earth file found for target %s", target.String())
		} else if err != nil {
			return "", errors.Wrapf(err, "stat file %s", buildEarthPath)
		}
		return buildEarthPath, nil
	} else if err != nil {
		return "", errors.Wrapf(err, "stat file %s", earthfilePath)
	}
	return earthfilePath, nil
}
