package features

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/util/flagutil"
)

// Features is used to denote which features to flip on or off; this is for use in maintaining
// backwards compatibility
type Features struct {
	// DoSaves flags the feature to only save artifacts under directly referenced targets
	// of those referenced by BUILDs
	ReferencedSaveOnly bool `long:"referenced-save-only" description:"only save artifacts that are directly referenced"`
}

var errUnexpectedArgs = fmt.Errorf("unexpected VERSION arguments; should be VERSION [flags] <major-version>.<minor-version>")

// GetFeatures returns a features struct for a particular version
func GetFeatures(version *spec.Version) (*Features, error) {
	var ftrs Features

	if version == nil {
		return &ftrs, nil
	}

	if version.Args == nil {
		return nil, errUnexpectedArgs
	}

	parsedArgs, err := flagutil.ParseArgs("VERSION", &ftrs, version.Args)
	if err != nil {
		return nil, err
	}

	if len(parsedArgs) != 1 {
		return nil, errUnexpectedArgs
	}

	majorAndMinor := strings.Split(parsedArgs[0], ".")
	if len(majorAndMinor) != 2 {
		return nil, errUnexpectedArgs
	}
	major, err := strconv.Atoi(majorAndMinor[0])
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse major version %q", majorAndMinor[0])
	}
	minor, err := strconv.Atoi(majorAndMinor[1])
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse minor version %q", majorAndMinor[1])
	}

	// TODO depending on versions here, flip on/off features.
	_ = major
	_ = minor

	return &ftrs, nil
}
