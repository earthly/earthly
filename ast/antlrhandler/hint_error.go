package antlrhandler

import (
	"fmt"
	"strings"
)

type hintError struct {
	err   error
	hints []string
}

func WithHints(err error, hints ...string) error {
	if err == nil {
		return nil
	}
	return hintError{err: err, hints: hints}
}

func (e hintError) Error() string {
	if len(e.hints) == 0 {
		return e.err.Error()
	}
	return fmt.Sprintf(`%v

Hints:
  - %v`, e.err, strings.Join(e.hints, "\n  - "))
}

func (e hintError) Unwrap() error {
	return e.err
}
