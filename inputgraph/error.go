package inputgraph

import (
	"fmt"
	"strings"

	"github.com/earthly/earthly/ast/spec"
	"github.com/pkg/errors"
)

type Error struct {
	srcLoc *spec.SourceLocation
	msg    string
	err    error
}

func (e *Error) Error() string {
	parts := []string{}
	if e.msg != "" {
		parts = append(parts, e.msg)
	}
	if e.err != nil {
		parts = append(parts, e.err.Error())
	}
	return strings.Join(parts, ": ")
}

func FormatError(err error) string {
	e := &Error{}
	if ok := errors.As(err, &e); ok {
		return fmt.Sprintf("%s line %d:%d: %s", e.srcLoc.File, e.srcLoc.StartLine, e.srcLoc.StartColumn, err)
	}
	return e.Error()
}

func newError(srcLoc *spec.SourceLocation, format string, args ...any) error {
	return &Error{
		srcLoc: srcLoc,
		msg:    fmt.Sprintf(format, args...),
	}
}

func wrapError(err error, srcLoc *spec.SourceLocation, format string, args ...any) error {
	e := &Error{
		srcLoc: srcLoc,
		err:    err,
	}
	if format != "" {
		e.msg = fmt.Sprintf(format, args...)
	}
	return e
}

func addErrorSrc(err error, srcLoc *spec.SourceLocation) error {
	return wrapError(err, srcLoc, "")
}
