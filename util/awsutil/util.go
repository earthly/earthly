package awsutil

import (
	"fmt"
	"os"

	"github.com/earthly/earthly/util/fileutil"
	"github.com/pkg/errors"
)

const (
	secretKeyEnv = "AWS_SECRET_ACCESS_KEY"
	accessKeyEnv = "AWS_ACCESS_KEY_ID"
)

// CredsAvailable determines if AWS credentials are available in the environment
// or in the standard ~/.aws location.
func CredsAvailable() (bool, error) {

	homeDir, _ := fileutil.HomeDir()
	credsPath := fmt.Sprintf("%s/.aws/credentials", homeDir)

	credsFileExists := true
	_, err := os.Stat(credsPath)
	switch {
	case errors.Is(err, os.ErrExist):
		credsFileExists = false
	case err != nil:
		return false, errors.Wrapf(err, "failed to stat %s", credsPath)
	}

	secretNames := []string{
		secretKeyEnv,
		accessKeyEnv,
	}

	allEnvsSet := true

	for _, name := range secretNames {
		if _, ok := os.LookupEnv(name); !ok {
			allEnvsSet = false
			break
		}
	}

	return credsFileExists || allEnvsSet, nil
}
