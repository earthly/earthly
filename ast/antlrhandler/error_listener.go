package antlrhandler

import (
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/earthly/earthly/ast/parser"
	"github.com/pkg/errors"
)

// humanName makes an attempt to translate the name of a symbol into something
// humans can read clearly.
func humanName(symbol string) string {
	switch symbol {
	case "NL":
		return `\n`
	case "EQUALS":
		return "="
	default:
		return symbol
	}
}

// ReturnErrorListener allows for the errors to be collected and returned after parsing.
type ReturnErrorListener struct {
	*antlr.DefaultErrorListener
	Errs []error
}

// NewReturnErrorListener returns a new ReturnErrorListener.
func NewReturnErrorListener() *ReturnErrorListener {
	return &ReturnErrorListener{}
}

// SyntaxError implements ErrorListener SyntaxError.
func (rel *ReturnErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	p, ok := recognizer.(*antlr.BaseParser)
	if !ok {
		rel.Errs = append(rel.Errs, errors.Errorf("syntax error: line %d:%d: %v", line, column, msg))
		return
	}

	// The line/column arguments seem to be passed in as the start of the
	// statement that failed to parse. But it seems like we get closer to the
	// real problem using GetCurrentToken() and its location.
	currTok := p.GetCurrentToken()
	tokLine := currTok.GetLine()
	tokCol := currTok.GetColumn()
	currLit := currTok.GetText()

	hintErr := hintError{
		err: errors.Errorf("syntax error: line %d:%d: unexpected '%v': %s", tokLine, tokCol, humanName(currLit), msg),
	}

	expected := p.GetExpectedTokens().StringVerbose(p.LiteralNames, p.SymbolicNames, false)
	hintErr.hints = []string{fmt.Sprintf("I got lost looking for '%v'", humanName(expected))}

	stream := p.GetInputStream()
	currIdx := stream.Index()
	switch e.(type) {
	case *antlr.NoViableAltException:
		// TODO: this error doesn't give us much option to give good hints.
		// Usually, when there's "no viable alternative", antlr rolls back to a
		// previous token, so the offendingSymbol is misleading.
		//
		// What has been tried:
		//
		// - walk forward using stream.Seek(idx) until p.GetCurrentToken() does
		//   not match p.GetExpectedTokens(). Ideally, this gives us the first
		//   token that doesn't match the expected token set. Unfortunately,
		//   p.GetExpectedTokens() returns a *antlr.IntervalSet which ... has
		//   zero useful exported methods or fields.

		// Until we can figure that out, the "I got lost looking for..." message
		// is pretty likely to be misleading.
		hintErr.hints[0] = "I couldn't find a pattern that completes the current statement - check your quote pairs, paren pairs, and newlines"
	default:
	}

	// Just to prevent duplicates, since we seem to run into them sometimes
	hintSet := map[string]struct{}{
		hintErr.hints[0]: {},
	}

	for idx := currIdx; idx >= 0; idx-- {
		stream.Seek(idx)
		tok := p.GetCurrentToken()
		if tok.GetTokenType() != parser.EarthLexerAtom {
			// The Atom type is our catch-all, and is the most likely candidate
			// for consuming something the user intended as a keyword. Other
			// tokens would probably provide misleading hints.
			continue
		}
		if currLine, currCol := tok.GetLine(), tok.GetColumn(); currLine < line || (currLine == line && currCol < column) {
			break
		}
		tokLit := tok.GetText()
		for _, lit := range p.LiteralNames {
			lit = strings.Trim(lit, "'")
			if lit == "" {
				continue
			}
			if tokLit == lit {
				msg := fmt.Sprintf("I parsed '%v' as a word, but it looks like it should be a keyword - is it on the wrong line?", lit)
				if _, ok := hintSet[msg]; ok {
					break
				}
				hintSet[msg] = struct{}{}
				hintErr.hints = append(hintErr.hints, msg)
				break
			}
		}
	}
	rel.Errs = append(rel.Errs, hintErr)
}
