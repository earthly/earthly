package hint

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var internal = errors.New("internal")

func TestWrapf(t *testing.T) {
	t.Run("without args", func(t *testing.T) {
		res := Wrapf(internal, "some hint")
		assert.Equal(t, &Error{
			err:   internal,
			hints: []string{"some hint"},
		}, res)
	})
	t.Run("with args", func(t *testing.T) {
		res := Wrapf(internal, "some hint with arg %s", "my-arg")
		assert.Equal(t, &Error{
			err:   internal,
			hints: []string{"some hint with arg my-arg"},
		}, res)
	})
}

func TestWrap(t *testing.T) {
	t.Run("with one hint", func(t *testing.T) {
		res := Wrap(internal, "some hint")
		assert.Equal(t, &Error{
			err:   internal,
			hints: []string{"some hint"},
		}, res)
	})
	t.Run("with multiple hints", func(t *testing.T) {
		res := Wrap(internal, "some hint", "another hint")
		assert.Equal(t, &Error{
			err:   internal,
			hints: []string{"some hint", "another hint"},
		}, res)
	})
}

func TestReceivers(t *testing.T) {
	err := Wrap(internal, "some hint", "another hint")

	t.Run("test Error", func(t *testing.T) {
		assert.Equal(t, "internal:Hint: some hint\nanother hint\n", err.Error())
	})

	t.Run("test Message", func(t *testing.T) {
		assert.Equal(t, "internal", err.(*Error).Message())
	})

	t.Run("test Hint", func(t *testing.T) {
		assert.Equal(t, "internal", err.(*Error).Message())
		assert.Equal(t, "some hint\nanother hint\n", err.(*Error).Hint())
	})
}

func TestFromError(t *testing.T) {

	tests := map[string]struct {
		err                 error
		expectedErr         *Error
		expectedIsHintError bool
	}{
		"err is nil": {},
		"err is not a hint err (but close)": {
			err: errors.New("some error: Hint 123"),
		},
		"err is a hint error": {
			err:                 Wrap(internal, "some hint"),
			expectedErr:         Wrap(internal, "some hint\n").(*Error),
			expectedIsHintError: true,
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res, isHintErr := FromError(tc.err)
			assert.Equal(t, tc.expectedErr, res)
			assert.Equal(t, tc.expectedIsHintError, isHintErr)
		})
	}
}
