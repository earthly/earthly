package earthfile2llb

import (
	"fmt"

	"github.com/earthly/earthly/ast/spec"
	"github.com/pkg/errors"
)

var _ error = &InterpreterError{}

// InterpreterError is an error of the interpreter, which contains optional references to the original
// source code location.
type InterpreterError struct {
	SourceLocation *spec.SourceLocation
	text           string
	cause          error
}

// Errorf creates a new interpreter error.
func Errorf(sl *spec.SourceLocation, format string, args ...interface{}) *InterpreterError {
	return &InterpreterError{
		SourceLocation: sl,
		text:           fmt.Sprintf(format, args...),
	}
}

// WrapError wraps another error into a new interpreter error.
func WrapError(cause error, sl *spec.SourceLocation, format string, args ...interface{}) *InterpreterError {
	return &InterpreterError{
		cause:          cause,
		SourceLocation: sl,
		text:           fmt.Sprintf(format, args...),
	}
}

func (ie InterpreterError) Error() string {
	var err error
	if ie.cause != nil {
		err = errors.Wrap(ie.cause, ie.text)
	} else {
		err = errors.New(ie.text)
	}
	if ie.SourceLocation == nil {
		return err.Error()
	}
	return fmt.Sprintf(
		"%s line %d:%d %s",
		ie.SourceLocation.File, ie.SourceLocation.StartLine, ie.SourceLocation.StartColumn,
		err.Error())
}

// Unwrap returns the cause of the error (if any).
func (ie InterpreterError) Unwrap() error {
	return ie.cause
}

// GetInterpreterError finds the first InterpreterError in the wrap chain and returns it.
func GetInterpreterError(err error) (*InterpreterError, bool) {
	if err == nil {
		return nil, false
	}
	ie, ok := err.(*InterpreterError)
	if ok {
		return ie, true
	}
	unwrapped := errors.Unwrap(err)
	if unwrapped != nil {
		return GetInterpreterError(unwrapped)
	}
	return nil, false
}
