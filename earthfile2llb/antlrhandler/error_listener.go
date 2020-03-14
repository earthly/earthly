package antlrhandler

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
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
	rel.Errs = append(rel.Errs, fmt.Errorf("Syntax error: line %d:%d %s", line, column, msg))
}
