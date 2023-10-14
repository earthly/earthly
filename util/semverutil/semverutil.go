package semverutil

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

// Version is a semantic version number.
type Version struct {
	Major int
	Minor int
	Patch int
	Tail  string
}

// Parse parses a semantic version number.
func Parse(s string) (Version, error) {
	s = strings.TrimPrefix(s, "v")
	var v Version
	n, err := fmt.Sscanf(s, "%d.%d.%d%s", &v.Major, &v.Minor, &v.Patch, &v.Tail)
	if err == io.EOF && n == 3 { // no tail case
		return v, nil
	}
	if err != nil {
		return Version{}, errors.Wrap(err, "parsing semantic version")
	}
	return v, nil
}

func Equal(s string, ver string) bool {
	if s == ver {
		return true
	}
	return s == "v"+ver
}

// String returns the string representation of the version.
func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d%s", v.Major, v.Minor, v.Patch, v.Tail)
}

// IsCompatible returns true if the two versions are compatible.
func IsCompatible(a, b Version) bool {
	return a.Major == b.Major && a.Minor == b.Minor
}
