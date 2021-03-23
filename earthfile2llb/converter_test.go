package earthfile2llb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type failer struct {
	failTimes   int
	currentFail int
}

func (f *failer) Go() error {
	if f.currentFail >= f.failTimes {
		return nil
	}

	f.currentFail++

	return fmt.Errorf("%v of %v fails happened", f.currentFail, f.failTimes)
}

func TestRetry(t *testing.T) {

	f := &failer{1, 0}

	err := doWithRetries(f.Go, func(err error) bool { return err != nil }, 2)

	assert.NoError(t, err)
}
