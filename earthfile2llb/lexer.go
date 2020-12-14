package earthfile2llb

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/earthly/earthly/earthfile2llb/parser"
)

// lexer is a lexer for an earthly file, which also emits indentation
// and dedentation tokens.
type lexer struct {
	*parser.EarthLexer
	prevIndentLevel                              int
	indentLevel                                  int
	afterNewLine                                 bool
	tokenQueue                                   []antlr.Token
	wsChannel, wsStart, wsStop, wsLine, wsColumn int
}

func newLexer(input antlr.CharStream) antlr.Lexer {
	l := new(lexer)
	l.EarthLexer = parser.NewEarthLexer(input)
	return l
}

func (l *lexer) NextToken() antlr.Token {
	peek := l.EarthLexer.NextToken()
	ret := peek
	switch peek.GetTokenType() {
	case parser.EarthLexerWS:
		if l.afterNewLine {
			l.indentLevel++
		}
		l.wsChannel, l.wsStart, l.wsStop, l.wsLine, l.wsColumn =
			peek.GetChannel(), peek.GetStart(), peek.GetStop(), peek.GetLine(), peek.GetColumn()
	case parser.EarthLexerNL:
		l.indentLevel = 0
		l.afterNewLine = true
	default:
		if l.afterNewLine {
			if l.prevIndentLevel < l.indentLevel {
				l.tokenQueue = append(l.tokenQueue, l.GetTokenFactory().Create(
					l.GetTokenSourceCharStreamPair(), parser.EarthLexerINDENT, "",
					l.wsChannel, l.wsStart, l.wsStop, l.wsLine, l.wsColumn))
			} else if l.prevIndentLevel > l.indentLevel {
				l.tokenQueue = append(l.tokenQueue, l.GetTokenFactory().Create(
					l.GetTokenSourceCharStreamPair(), parser.EarthLexerDEDENT, "",
					l.wsChannel, l.wsStart, l.wsStop, l.wsLine, l.wsColumn))
				l.PopMode() // Pop RECIPE mode.
			}
		}
		l.prevIndentLevel = l.indentLevel
		l.afterNewLine = false
	}
	if len(l.tokenQueue) > 0 {
		l.tokenQueue = append(l.tokenQueue, peek)
		ret = l.tokenQueue[0]
		l.tokenQueue = l.tokenQueue[1:]
	}
	return ret
}
