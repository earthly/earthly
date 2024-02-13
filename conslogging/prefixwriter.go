package conslogging

import "io"

// PrefixWriter is a writer that can take a prefix.
type PrefixWriter interface {
	io.Writer
	// WithPrefix returns a new PrefixWriter with the given prefix.
	WithPrefix(prefix string) PrefixWriter
}
