package params

import (
	"fmt"
	"github.com/pkg/errors"
)

type Error struct {
	msg   string
	cause error
}

func Errorf(format string, args ...any) error {
	return &Error{
		msg: fmt.Sprintf(format, args...),
	}
}

func Wrapf(err error, format string, args ...any) error {
	return &Error{
		msg:   fmt.Sprintf(format, args...),
		cause: err,
	}
}

func (e *Error) Error() string {
	if e.cause != nil {
		fmt.Errorf("%s: %w", e.msg, e.cause)
	}
	return e.msg
}

func (e *Error) Cause() error {
	return errors.Cause(e.cause)
}

func (e *Error) Is(err error) bool {
	_, ok := err.(*Error)
	return ok
}

func (e *Error) ParentError() string {
	return e.msg
}
