package buildcontext

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/domain"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

// ErrNotExist is the struct indicating that file does not exist.
type ErrEarthfileNotExist struct {
	Target string
}

// Error is function required by error interface.
func (err ErrEarthfileNotExist) Error() string {
	return "No Earthfile nor build.earth file found for target " + err.Target
}

// detectBuildFile detects whether to use Earthfile, build.earth or Dockerfile.
func detectBuildFile(ref domain.Reference, localDir string) (string, error) {
	if strings.HasPrefix(ref.GetName(), DockerfileMetaTarget) {
		return filepath.Join(localDir, strings.TrimPrefix(ref.GetName(), DockerfileMetaTarget)), nil
	}
	earthfilePath := filepath.Join(localDir, "Earthfile")
	_, err := os.Stat(earthfilePath)
	if os.IsNotExist(err) {
		buildEarthPath := filepath.Join(localDir, "build.earth")
		_, err := os.Stat(buildEarthPath)
		if os.IsNotExist(err) {
			return "", ErrEarthfileNotExist{Target: ref.String()}
		} else if err != nil {
			return "", errors.Wrapf(err, "stat file %s", buildEarthPath)
		}
		return buildEarthPath, nil
	} else if err != nil {
		return "", errors.Wrapf(err, "stat file %s", earthfilePath)
	}
	return earthfilePath, nil
}

func detectBuildFileInRef(ctx context.Context, earthlyRef domain.Reference, ref gwclient.Reference, subDir string) (string, error) {
	if strings.HasPrefix(earthlyRef.GetName(), DockerfileMetaTarget) {
		return filepath.Join(subDir, strings.TrimPrefix(earthlyRef.GetName(), DockerfileMetaTarget)), nil
	}
	earthfilePath := path.Join(subDir, "Earthfile")
	exists, err := fileExists(ctx, ref, earthfilePath)
	if err != nil {
		return "", err
	}
	if exists {
		return earthfilePath, nil
	}
	buildEarthPath := path.Join(subDir, "build.earth")
	exists, err = fileExists(ctx, ref, buildEarthPath)
	if err != nil {
		return "", err
	}
	if exists {
		return buildEarthPath, nil
	}
	return "", ErrEarthfileNotExist{Target: earthlyRef.String()}
}

func fileExists(ctx context.Context, ref gwclient.Reference, fpath string) (bool, error) {
	dir, file := path.Split(fpath)
	fstats, err := ref.ReadDir(ctx, gwclient.ReadDirRequest{
		Path:           dir,
		IncludePattern: file,
	})
	if err != nil {
		return false, errors.Wrapf(err, "cannot read dir %s", dir)
	}
	for _, fstat := range fstats {
		name := path.Base(fstat.GetPath())
		if name == file && !fstat.IsDir() {
			return true, nil
		}
	}
	return false, nil
}
