package hint

import (
	"fmt"
	"strings"
)

type hintError struct {
	err   error
	hints []string
}

func (e hintError) Error() string {
	switch len(e.hints) {
	case 0:
		return e.err.Error()
	case 1:
		return fmt.Sprintf(`%v

  Hint: %v
`, e.err, e.hints[0])
	default:
		return fmt.Sprintf(`%v

  Hints:
  - %v
`, e.err, strings.Join(e.hints, "\n  - "))
	}
}

func (e hintError) Unwrap() error {
	return e.err
}

// Wrap wraps up an error with hints, to help display hints to a user about what
// might fix the problem.
func Wrap(err error, firstHint string, extraHints ...string) error {
	return hintError{err: err, hints: append([]string{firstHint}, extraHints...)}
}

// Wrapf wraps an error with a single hint with formatting arguments.
func Wrapf(err error, hintf string, args ...any) error {
	return Wrap(err, fmt.Sprintf(hintf, args...))
}
