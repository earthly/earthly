package antlrhandler

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// ReturnErrorStrategy allows for the error to be returned after parsing.
type ReturnErrorStrategy struct {
	*antlr.DefaultErrorStrategy
	Err        error
	ErrContext antlr.ParserRuleContext
	RE         antlr.RecognitionException
}

var _ antlr.ErrorStrategy = &ReturnErrorStrategy{}

// NewReturnErrorStrategy returns a new instance of ReturnErrorStrategy.
func NewReturnErrorStrategy() *ReturnErrorStrategy {
	res := new(ReturnErrorStrategy)
	res.DefaultErrorStrategy = antlr.NewDefaultErrorStrategy()
	return res
}

// Recover implements ErrorStrategy Recover.
func (res *ReturnErrorStrategy) Recover(recognizer antlr.Parser, e antlr.RecognitionException) {
	if res.Err == nil {
		res.RE = e
		res.Err = fmt.Errorf("parse error")
		res.ErrContext = recognizer.GetParserRuleContext()
	}
	context := recognizer.GetParserRuleContext()
	for context != nil {
		context.SetException(e)
		var ok bool
		context, ok = context.GetParent().(antlr.ParserRuleContext)
		if !ok {
			break
		}
	}
}

// RecoverInline implements ErrorStrategy RecoverInline.
func (res *ReturnErrorStrategy) RecoverInline(recognizer antlr.Parser) antlr.Token {
	res.Recover(recognizer, antlr.NewInputMisMatchException(recognizer))
	return recognizer.GetCurrentToken()
}

// Sync implements ErrorStrategy Sync.
func (res *ReturnErrorStrategy) Sync(recognizer antlr.Parser) {
}
