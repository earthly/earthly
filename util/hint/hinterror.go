package hint

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/earthly/earthly/util/stringutil"
)

// note this regex should be updated in case the error format changes in Error
var regex = regexp.MustCompile(`(?P<error>.+?):Hint: (?P<hint>(?s).+)`)

// Error is an error that includes hints to be displayed after the error.
type Error struct {
	err   error
	hints []string
}

// Error returns the error string
func (e *Error) Error() string {
	return fmt.Sprintf(`%v:Hint: %v`, e.err, e.Hint())
}

func (e *Error) Message() string {
	return e.err.Error()
}

// Hint returns all hints in a single string separated by a new line
func (e *Error) Hint() string {
	if len(e.hints) == 0 {
		return ""
	}
	res := strings.Join(e.hints, "\n")
	if !strings.HasSuffix(res, "\n") {
		res = fmt.Sprintf("%s\n", res)
	}
	return res
}

// Wrap wraps up an error with hints, to help display hints to a user about what
// might fix the problem.
func Wrap(err error, firstHint string, extraHints ...string) error {
	return &Error{err: err, hints: append([]string{firstHint}, extraHints...)}
}

// Wrapf wraps an error with a single hint with formatting arguments.
func Wrapf(err error, hintf string, args ...any) error {
	return Wrap(err, fmt.Sprintf(hintf, args...))
}

// FromError attempts to parse the given error's string to an *hint.Error
func FromError(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	matches, _ := stringutil.NamedGroupMatches(err.Error(), regex)
	if len(matches) != 2 && len(matches) != 2 {
		return nil, false
	}
	for k := range matches {
		if len(matches[k]) != 1 {
			return nil, false
		}
	}
	errMsg := matches["error"][0]
	hint := matches["hint"][0]

	return Wrap(
		errors.New(errMsg),
		hint,
	).(*Error), true
}
