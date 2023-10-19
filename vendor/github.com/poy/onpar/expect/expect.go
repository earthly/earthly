package expect

import (
	"path"
	"runtime"

	"github.com/poy/onpar/matchers"
)

// ToMatcher is a type that can be passed to (*To).To().
type ToMatcher interface {
	Match(actual any) (resultValue any, err error)
}

// Differ is a type of matcher that will need to diff its expected and
// actual values.
type DiffMatcher interface {
	UseDiffer(matchers.Differ)
}

// T is a type that we can perform assertions with.
type T interface {
	Fatalf(format string, args ...any)
}

// THelper has the method that tells the testing framework that it can declare
// itself a test helper.
type THelper interface {
	Helper()
}

// Opt is an option that can be passed to New to modify Expectations.
type Opt func(To) To

// WithDiffer stores the diff.Differ to be used when displaying diffs between
// actual and expected values.
func WithDiffer(d matchers.Differ) Opt {
	return func(t To) To {
		t.differ = d
		return t
	}
}

// Expectation is provided to make it clear what the expect function does.
type Expectation func(actual any) *To

// New creates a new Expectation
func New(t T, opts ...Opt) Expectation {
	return func(actual any) *To {
		to := To{
			actual: actual,
			t:      t,
		}
		for _, opt := range opts {
			to = opt(to)
		}
		return &to
	}
}

// Expect performs New(t)(actual).
func Expect(t T, actual any) *To {
	return New(t)(actual)
}

// To is a type that stores actual values prior to running them through
// matchers.
type To struct {
	actual    any
	parentErr error

	t      T
	differ matchers.Differ
}

// To takes a matcher and passes it the actual value, failing t's T value
// if the matcher returns an error.
func (t *To) To(matcher matchers.Matcher) {
	if helper, ok := t.t.(THelper); ok {
		helper.Helper()
	}

	if d, ok := matcher.(DiffMatcher); ok {
		d.UseDiffer(t.differ)
	}

	_, err := matcher.Match(t.actual)
	if err != nil {
		_, fileName, lineNumber, _ := runtime.Caller(1)
		t.t.Fatalf("%s\n%s:%d", err.Error(), path.Base(fileName), lineNumber)
	}
}
