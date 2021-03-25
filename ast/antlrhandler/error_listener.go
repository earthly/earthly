package antlrhandler

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"
)

// ReturnErrorListener allows for the errors to be collected and returned after parsing.
type ReturnErrorListener struct {
	*antlr.DefaultErrorListener
	Errs []error
}

// NewReturnErrorListener returns a new ReturnErrorListener.
func NewReturnErrorListener() *ReturnErrorListener {
	return new(ReturnErrorListener)
}

// SyntaxError implements ErrorListener SyntaxError.
func (rel *ReturnErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	rel.Errs = append(rel.Errs, errors.Errorf("syntax error: line %d:%d %s", line, column, msg))
}
