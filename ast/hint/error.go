package hint

import (
	"fmt"
	"strings"
)

// TODO: once WrapWithDisplay & WrapfWithDisplay are out of use,
// remove displayHints field and let earthly/cmd/earthly/app/run.go use the Hint() value instead of printing it
// as part of the Error() function

type Error struct {
	err   error
	hints []string
	// Indicate whether to display the hint as part of the Error() function
	// or expose the hints via Hint().
	// Once the deprecated functions are removed, this flag can be removed as well.
	displayHints bool
}

func (e Error) Error() string {
	if !e.displayHints {
		return e.err.Error()
	}
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

func (e Error) Hint() string {
	if e.displayHints {
		// if true we can leave this empty so that it is not displayed twice in cmd/earthly/app/run.go
		return ""
	}
	return strings.Join(e.hints, "\n")
}

// Wrap wraps up an error with hints, to help display hints to a user about what
// might fix the problem.
func Wrap(err error, firstHint string, extraHints ...string) error {
	return Error{err: err, hints: append([]string{firstHint}, extraHints...)}
}

// WrapWithDisplay wraps up an error with hints, to help display hints to a user about what
// might fix the problem.
// Deprecated: use Wrap instead after verifying hint is properly displayed
func WrapWithDisplay(err error, firstHint string, extraHints ...string) error {
	return Error{err: err, displayHints: true, hints: append([]string{firstHint}, extraHints...)}
}

// Wrapf wraps an error with a single hint with formatting arguments.
func Wrapf(err error, hintf string, args ...any) error {
	return Wrap(err, fmt.Sprintf(hintf, args...))
}

// WrapfWithDisplay wraps an error with a single hint with formatting arguments.
// Deprecated: use Wrapf instead after verifying hint is properly displayed
func WrapfWithDisplay(err error, hintf string, args ...any) error {
	return Wrap(err, fmt.Sprintf(hintf, args...))
}
