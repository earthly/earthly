package buildcontext

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// detectEarthfile detects whether to use Earthfile or build.earth.
func detectEarthfile(targetName string, localDir string) (string, error) {
	earthfilePath := filepath.Join(localDir, "Earthfile")
	buildEarthPath := filepath.Join(localDir, "build.earth")
	_, err := os.Stat(earthfilePath)
	if os.IsNotExist(err) {
		_, err := os.Stat(buildEarthPath)
		if os.IsNotExist(err) {
			return "", fmt.Errorf(
				"No Earthfile nor build.earth file found for target %s", targetName)
		} else if err != nil {
			return "", errors.Wrapf(err, "stat file %s", buildEarthPath)
		}
		return buildEarthPath, nil
	} else if err != nil {
		return "", errors.Wrapf(err, "stat file %s", earthfilePath)
	}
	return earthfilePath, nil
}
