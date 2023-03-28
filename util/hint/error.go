package hint

import "fmt"

type hintError struct {
	err  error
	hint string
}

func (e hintError) Error() string {
	return fmt.Sprintf("%v\n    [hint: %v]", e.err, e.hint)
}

func (e hintError) Unwrap() error {
	return e.err
}

// Wrap wraps up an error with a hint, to help display hints to a user about
// what might fix the problem.
func Wrap(err error, hint string) error {
	return hintError{err: err, hint: hint}
}

// Wrapf is like Wrap but supports formatting arguments.
func Wrapf(err error, hintf string, args ...any) error {
	return hintError{err: err, hint: fmt.Sprintf(hintf, args...)}
}
