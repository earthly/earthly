package ast

import (
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/earthly/earthly/ast/parser"
)

const (
	indentChannel  = antlr.LexerDefaultTokenChannel
	newlineChannel = antlr.LexerDefaultTokenChannel
)

// lexer is a lexer for an earthly file, which also emits indentation
// and dedentation tokens.
type lexer struct {
	*parser.EarthLexer

	prevIndentLevel  int
	indentLevel      int
	afterNewLine     bool
	afterLineComment bool

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
	defer func() {
		// Recovering here prevents a panic due to erroneous indentation. The
		// tokens will still fail to match our parser grammar, so we have no
		// need to add an error of our own.
		_ = recover()
	}()
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

	if l.afterLineComment {
		if peek.GetTokenType() == parser.EarthParserNL {
			// Comments on their own line consume one trailing newline so that
			// we can recognize comment blocks separately from other comment
			// blocks. We still need to handle the indentation level for the
			// newline, though.
			l.processIndentation(peek)
			peek = l.EarthLexer.NextToken()
		}
		l.afterLineComment = false
	}

	ret := peek
	if peek.GetTokenType() == parser.EarthParserEOF {
		// Add a NL before EOF. It simplifies the logic a lot if we know
		// that all lines have been completed.
		l.tokenQueue = append(l.tokenQueue, l.makeNL(peek))
		if l.debug {
			fmt.Printf("NL ")
		}
		// This ensures that any necessary DEDENT tokens are added before the
		// EOF.
		l.afterNewLine = true
		l.indentLevel = 0

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

type seeker interface {
	Index() int
	Seek(int)
}

// handleCommentIndentLevel checks whether or not a comment may need to trigger
// an INDENT or DEDENT. The indent level will be set during this function if
// necessary. This happens mainly in the following scenarios:
//
//   - When a comment is unindented after a recipe body and is followed by a new
//     recipe.
//
//   - When a comment is indented as the first line of a recipe body and is
//     followed by a non-comment token.
//
// In these scenarios, the comment may be documentation and needs to trigger the
// INDENT/DEDENT _before_ the comment in the token sequence.
func (l *lexer) handleCommentIndentLevel(seeker seeker, comment antlr.Token) bool {
	idx := seeker.Index()
	defer seeker.Seek(idx)

	text := comment.GetText()
	// TODO: with whitespace on its own channel, we can probably remove the
	// whitespace from the COMMENT token and read these as WS tokens instead.
	indented := strings.HasPrefix(text, " ") || strings.HasPrefix(text, "\t")

	next := l.EarthLexer.NextToken()
	for ; next.GetTokenType() == parser.EarthLexerCOMMENT || next.GetTokenType() == parser.EarthLexerNL; next = l.EarthLexer.NextToken() {
		if next.GetTokenType() != parser.EarthLexerCOMMENT {
			continue
		}
		text := next.GetText()
		alsoIndented := strings.HasPrefix(text, " ") || strings.HasPrefix(text, "\t")
		if indented != alsoIndented {
			return false
		}
	}
	switch next.GetTokenType() {
	case parser.EarthLexerWS:
		if indented {
			l.indentLevel = 1
			return true
		}
		return false
	default:
		if !indented {
			l.indentLevel = 0
			return true
		}
		return false
	}
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
	case parser.EarthLexerCOMMENT:
		if !l.afterNewLine {
			return
		}

		l.afterLineComment = true

		if l.handleCommentIndentLevel(l.GetInputStream(), peek) {
			l.handleIndentLevel(peek)
		}
	default:
		l.handleIndentLevel(peek)
	}
}

func (l *lexer) handleIndentLevel(peek antlr.Token) {
	if !l.afterNewLine {
		return
	}
	l.afterNewLine = false

	if l.prevIndentLevel == l.indentLevel {
		return
	}
	prevIndent := l.prevIndentLevel
	l.prevIndentLevel = l.indentLevel

	if prevIndent < l.indentLevel {
		l.tokenQueue = append(l.tokenQueue, l.makeIndent())
		if l.debug {
			fmt.Printf("INDENT ")
		}
		return
	}

	l.tokenQueue = append(l.tokenQueue, l.makeDedent())
	if l.debug {
		fmt.Printf("DEDENT ")
	}
	switch peek.GetTokenType() {
	case parser.EarthLexerTarget, parser.EarthLexerUserCommand:
	default:
		l.popRecipeMode()
	}
}

func (l *lexer) makeIndent() antlr.Token {
	return l.GetTokenFactory().Create(
		l.GetTokenSourceCharStreamPair(), parser.EarthLexerINDENT, "",
		indentChannel, l.wsStart, l.wsStop, l.wsLine, l.wsColumn)
}

func (l *lexer) makeDedent() antlr.Token {
	return l.GetTokenFactory().Create(
		l.GetTokenSourceCharStreamPair(), parser.EarthLexerDEDENT, "",
		indentChannel, l.wsStart, l.wsStop, l.wsLine, l.wsColumn)
}

func (l *lexer) makeNL(peek antlr.Token) antlr.Token {
	return l.GetTokenFactory().Create(
		l.GetTokenSourceCharStreamPair(), parser.EarthLexerNL, "",
		newlineChannel, peek.GetStart(), peek.GetStop(),
		peek.GetLine(), peek.GetColumn())
}
