package matchers

// Differ is a type that can show a detailed difference between an
// actual and an expected value.
type Differ interface {
	Diff(actual, expected any) string
}
