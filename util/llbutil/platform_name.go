package llbutil

import (
	"fmt"

	"github.com/earthly/earthly/util/platutil"

	"github.com/docker/distribution/reference"
	"github.com/pkg/errors"
)

// PlatformSpecificImageName returns the PlatformSpecificImageName
func PlatformSpecificImageName(imgName string, platform platutil.Platform) (string, error) {
	platformStr := platform.String()
	if platformStr == "" {
		platformStr = "native"
	}
	r, err := reference.ParseNormalizedNamed(imgName)
	if err != nil {
		return "", errors.Wrapf(err, "parse %s", imgName)
	}
	taggedR, ok := reference.TagNameOnly(r).(reference.Tagged)
	if !ok {
		return "", errors.Wrapf(err, "not tagged %s", reference.TagNameOnly(r).String())
	}
	platformTag := DockerTagSafe(fmt.Sprintf("%s_%s", taggedR.Tag(), platformStr))
	r2, err := reference.WithTag(r, platformTag)
	if err != nil {
		return "", errors.Wrapf(err, "with tag %s - %s", r.String(), platformTag)
	}
	return reference.FamiliarString(r2), nil
}
