package earthfile2llb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type failer struct {
	failTimes int
}

func (f *failer) Go(tries int) error {
	if tries >= f.failTimes {
		return nil
	}

	return fmt.Errorf("%v of %v fails happened", tries, f.failTimes)
}

func TestRetry(t *testing.T) {

	f := &failer{9999}

	err := doWithRetries(f.Go, func(err error) bool { return err != nil }, 2)

	assert.NoError(t, err)
}
