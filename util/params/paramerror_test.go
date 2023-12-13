package params

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var internal = errors.New("internal")

func TestErrorf(t *testing.T) {
	t.Run("without args", func(t *testing.T) {
		res := Errorf("some error")
		assert.Equal(t, &Error{
			msg: "some error",
		}, res)
	})
	t.Run("with args", func(t *testing.T) {
		res := Errorf("some error %s", "myarg")
		assert.Equal(t, &Error{
			msg: "some error myarg",
		}, res)
	})
}

func TestWrapf(t *testing.T) {
	t.Run("without args", func(t *testing.T) {
		res := Wrapf(internal, "some error")
		assert.Equal(t, &Error{
			msg:   "some error",
			cause: internal,
		}, res)
	})
	t.Run("with args", func(t *testing.T) {
		res := Wrapf(internal, "some error %s", "myarg")
		assert.Equal(t, &Error{
			msg:   "some error myarg",
			cause: internal,
		}, res)
	})
}

func TestError(t *testing.T) {
	t.Run("without cause", func(t *testing.T) {
		res := Errorf("some error").Error()
		assert.Equal(t, "some error", res)
	})
	t.Run("with cause", func(t *testing.T) {
		res := Wrapf(internal, "some error").Error()
		assert.Equal(t, "some error: internal", res)
	})
}

func TestCause(t *testing.T) {
	pe, _ := Wrapf(internal, "some error").(*Error)
	res := pe.Cause()
	assert.Equal(t, errors.Cause(internal), res)
}

func TestIs(t *testing.T) {
	t.Run("non param error", func(t *testing.T) {
		pe, _ := Errorf("some error").(*Error)
		res := pe.Is(internal)
		assert.False(t, res)
	})
	t.Run("param error", func(t *testing.T) {
		pe, _ := Errorf("some error").(*Error)
		res := pe.Is(pe)
		assert.True(t, res)
	})
}

func TestParentError(t *testing.T) {
	pe, _ := Wrapf(internal, "some error").(*Error)
	res := pe.ParentError()
	assert.Equal(t, "some error", res)
}
