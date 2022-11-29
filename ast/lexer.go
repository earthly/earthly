package ast

import (
	"fmt"
	"strings"

	"github.com/earthly/earthly/ast/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// lexer is a lexer for an earthly file, which also emits indentation
// and dedentation tokens.
type lexer struct {
	*parser.EarthLexer

	prevIndentLevel int
	indentLevel     int
	afterNewLine    bool

	tokenQueue                                   []antlr.Token
	wsChannel, wsStart, wsStop, wsLine, wsColumn int

	err error

	debug bool
}

func newLexer(input antlr.CharStream) *lexer {
	return &lexer{
		EarthLexer: parser.NewEarthLexer(input),
		// Uncomment to print tokens to stdout.
		// debug: true,
	}
}

func (l *lexer) Err() error {
	return l.err
}

func (l *lexer) getMode() int {
	// TODO: Is there a better way to get this? There's no API for getting
	//       the current mode.
	l.PushMode(0)
	return l.PopMode()
}

// popRecipeMode removes the recipe mode stack frame and then places back anything that existed
// on top of it.
func (l *lexer) popRecipeMode() {
	m := l.getMode()
	if m == parser.EarthLexerRECIPE {
		// Special case: nothing above.
		l.PopMode()
		return
	}
	above := []int{m}
	for {
		m = l.PopMode()
		if m == parser.EarthLexerRECIPE {
			l.PopMode()
			break
		}
		above = append(above, m)
	}
	for i := len(above) - 1; i >= 0; i-- {
		l.PushMode(above[i])
	}
}

func (l *lexer) stackString() string {
	m := l.getMode()
	stack := []int{m}
	for m != 0 {
		m = l.PopMode()
		stack = append(stack, m)
	}
	var str []string
	for i := len(stack) - 1; i >= 0; i-- {
		if i != len(stack)-1 {
			l.PushMode(stack[i])
		}
		str = append(str, parser.GetLexerModeNames()[stack[i]])
	}
	return strings.Join(str, "/")
}

func (l *lexer) NextToken() antlr.Token {
	modeBefore := l.getMode()
	peek := l.EarthLexer.NextToken()
	ret := peek
	if peek.GetTokenType() == parser.EarthParserEOF {
		// Add a NL before EOF. It simplifies the logic a lot if we know
		// that all lines have been completed.
		l.tokenQueue = append(l.tokenQueue, l.makeNL(peek))
		if l.debug {
			fmt.Printf("NL ")
		}
		// Force the default mode.
		l.PushMode(0)
		modeBefore = 0
	}
	switch modeBefore {
	case 0, parser.EarthLexerRECIPE:
		l.processIndentation(peek)
	default:
		// Don't process indentation for any mode other than DEFAULT_MODE and RECIPE.
		l.indentLevel = 0
		l.afterNewLine = true
	}

	if l.debug {
		if modeBefore >= 0 {
			fmt.Printf("%s", parser.GetLexerModeNames()[modeBefore])
		}
		mode := l.getMode()
		if mode >= 0 && peek.GetTokenType() > 0 {
			if mode != modeBefore {
				fmt.Printf(">>%s", parser.GetLexerModeNames()[mode])
			}
			fmt.Printf("-%d(%s) ", l.indentLevel, parser.GetLexerSymbolicNames()[peek.GetTokenType()])
		}
		if peek.GetTokenType() == parser.EarthLexerNL {
			fmt.Printf("\n")
		}
	}

	if len(l.tokenQueue) > 0 {
		l.tokenQueue = append(l.tokenQueue, peek)
		ret = l.tokenQueue[0]
		l.tokenQueue = l.tokenQueue[1:]
	}

	return ret
}

func (l *lexer) processIndentation(peek antlr.Token) {
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
				l.tokenQueue = append(l.tokenQueue, l.makeIndent())
				if l.debug {
					fmt.Printf("INDENT ")
				}
			} else if l.prevIndentLevel > l.indentLevel {
				l.tokenQueue = append(l.tokenQueue, l.makeDedent())
				if l.debug {
					fmt.Printf("DEDENT ")
				}
				if peek.GetTokenType() != parser.EarthLexerTarget && peek.GetTokenType() != parser.EarthLexerUserCommand {
					l.popRecipeMode()
				}
			}
		}
		l.prevIndentLevel = l.indentLevel
		l.afterNewLine = false
	}
}

func (l *lexer) makeIndent() antlr.Token {
	return l.GetTokenFactory().Create(
		l.GetTokenSourceCharStreamPair(), parser.EarthLexerINDENT, "",
		l.wsChannel, l.wsStart, l.wsStop, l.wsLine, l.wsColumn)
}

func (l *lexer) makeDedent() antlr.Token {
	return l.GetTokenFactory().Create(
		l.GetTokenSourceCharStreamPair(), parser.EarthLexerDEDENT, "",
		l.wsChannel, l.wsStart, l.wsStop, l.wsLine, l.wsColumn)
}

func (l *lexer) makeNL(peek antlr.Token) antlr.Token {
	return l.GetTokenFactory().Create(
		l.GetTokenSourceCharStreamPair(), parser.EarthLexerNL, "",
		peek.GetChannel(), peek.GetStart(), peek.GetStop(),
		peek.GetLine(), peek.GetColumn())
}
