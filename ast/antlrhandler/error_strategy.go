package antlrhandler

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"
)

// ReturnErrorStrategy allows for the error to be returned after parsing.
type ReturnErrorStrategy struct {
	*antlr.DefaultErrorStrategy
	Err        error
	Hint       string
	ErrContext antlr.ParserRuleContext
	RE         antlr.RecognitionException

	litNames, symbNames []string
}

var _ antlr.ErrorStrategy = &ReturnErrorStrategy{}

// NewReturnErrorStrategy returns a new instance of ReturnErrorStrategy.
func NewReturnErrorStrategy(litNames, symbNames []string) *ReturnErrorStrategy {
	res := &ReturnErrorStrategy{
		litNames:  litNames,
		symbNames: symbNames,
	}
	res.DefaultErrorStrategy = antlr.NewDefaultErrorStrategy()
	return res
}

// Recover implements ErrorStrategy Recover.
func (res *ReturnErrorStrategy) Recover(recognizer antlr.Parser, e antlr.RecognitionException) {
	if res.Err == nil {
		res.RE = e
		res.Err = errors.Errorf("invalid syntax")
		res.ErrContext = recognizer.GetParserRuleContext()
		expected := recognizer.GetExpectedTokens().StringVerbose(res.litNames, res.symbNames, false)
		res.Hint = fmt.Sprintf("I got lost looking for '%v'", humanName(expected))
		switch expected {
		case "EQUALS":
			res.Hint += " - did you define a key/value pair without a value?"
		}
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
