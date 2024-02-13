package builder

// BuildError contains an BuildError and log
type BuildError struct {
	err error
	log string
}

// NewBuildError creates a new BuildError with the additional output log of the command that failed
func NewBuildError(err error, vertexLog string) error {
	if vertexLog == "" {
		return err
	}
	return &BuildError{
		err: err,
		log: vertexLog,
	}
}

// BuildError formats the BuildError as a string, omitting the vertex log
func (e *BuildError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.err.Error()
}

// Unwrap returns the wrapped error
func (e *BuildError) Unwrap() error {
	return e.err
}

// VertexLog returns the vertex log associated with the error
func (e *BuildError) VertexLog() string {
	return e.log
}
