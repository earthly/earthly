package ast

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/earthly/earthly/ast/antlrhandler"
	"github.com/earthly/earthly/ast/hint"
	"github.com/earthly/earthly/ast/parser"
	"github.com/earthly/earthly/ast/spec"
	"github.com/pkg/errors"
)

// Parse parses an earthfile into an AST.
func Parse(ctx context.Context, filePath string, enableSourceMap bool) (ef spec.Earthfile, err error) {
	var opts []Opt
	if enableSourceMap {
		opts = append(opts, WithSourceMap())
	}
	return ParseOpts(ctx, FromPath(filePath), opts...)
}

// ParseOpts parses an earthfile into an AST. This is the functional option
// version, which uses option functions to change how a file is parsed.
func ParseOpts(ctx context.Context, from FromOpt, opts ...Opt) (spec.Earthfile, error) {
	defaultPrefs := prefs{
		done: func() {},
	}
	prefs, err := from(defaultPrefs)
	if err != nil {
		return spec.Earthfile{}, errors.Wrap(err, "ast: could not apply FromOpt")
	}
	for _, opt := range opts {
		newPrefs, err := opt(prefs)
		if err != nil {
			return spec.Earthfile{}, errors.Wrap(err, "ast: could not apply options")
		}
		prefs = newPrefs
	}

	defer prefs.done()

	var versionOpts []Opt
	if prefs.enableSourceMap {
		versionOpts = append(versionOpts, WithSourceMap())
	}
	version, err := ParseVersionOpts(FromReader(prefs.reader), versionOpts...)
	if err != nil {
		return spec.Earthfile{}, err
	}

	errorListener := antlrhandler.NewReturnErrorListener()
	errorStrategy := antlrhandler.NewReturnErrorStrategy(parser.GetLexerLiteralNames(), parser.GetLexerSymbolicNames())

	if _, err := prefs.reader.Seek(0, 0); err != nil {
		return spec.Earthfile{}, errors.Wrap(err, "ast: could not seek to beginning of file")
	}
	b, err := io.ReadAll(prefs.reader)
	if err != nil {
		return spec.Earthfile{}, errors.Wrap(err, "ast: could not read Earthfile for parsing")
	}
	stream, tree, err := newEarthfileTree(string(b), errorListener, errorStrategy)
	if err != nil {
		return spec.Earthfile{}, err
	}
	ef, walkErr := walkTree(newListener(ctx, stream, prefs.reader.Name(), prefs.enableSourceMap), tree)
	if len(errorListener.Errs) > 0 {
		firstErr := errorListener.Errs[0].Error()
		if strings.Contains(firstErr, "token recognition error") {
			unknownHint := fmt.Sprintf("'%s' is not a recognized keyword.", extractUnknownToken(firstErr))
			return spec.Earthfile{}, hint.Wrap(errors.New(firstErr), unknownHint)
		}
		errString := []string{fmt.Sprintf("lexer error: %s", prefs.reader.Name())}
		for _, err := range errorListener.Errs {
			errString = append(errString, err.Error())
		}
		return spec.Earthfile{}, errors.Errorf(strings.Join(errString, "\n"))
	}
	if errorStrategy.Err != nil {
		err := errors.Wrapf(
			errorStrategy.Err, "%s:%d:%d '%s'",
			prefs.reader.Name(),
			errorStrategy.RE.GetOffendingToken().GetLine(),
			errorStrategy.RE.GetOffendingToken().GetColumn(),
			errorStrategy.RE.GetOffendingToken().GetText())
		if errorStrategy.Hint != "" {
			err = hint.Wrap(err, errorStrategy.Hint)
		}
		return spec.Earthfile{}, err
	}
	if walkErr != nil {
		return spec.Earthfile{}, walkErr
	}

	ef.Version = version

	if err := validateAst(ef); err != nil {
		return spec.Earthfile{}, err
	}

	return ef, nil
}

func extractUnknownToken(errMsg string) string {
	// Extract the part between single quotes using regex
	re := regexp.MustCompile(`'([^']+)'`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) > 1 {
		return matches[1]
	}
	return "unknown"
}

func walkTree(l *listener, tree parser.IEarthFileContext) (spec.Earthfile, error) {
	antlr.ParseTreeWalkerDefault.Walk(l, tree)
	if err := l.Err(); err != nil {
		return spec.Earthfile{}, err
	}
	return l.Earthfile(), nil
}

func newEarthfileTree(body string, errorListener *antlrhandler.ReturnErrorListener, errorStrategy antlr.ErrorStrategy) (*antlr.CommonTokenStream, parser.IEarthFileContext, error) {
	input := antlr.NewInputStream(body)
	lexer := newLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	if lexer.Err() != nil {
		return nil, nil, lexer.Err()
	}
	p := parser.NewEarthParser(stream)
	p.AddErrorListener(errorListener)
	p.SetErrorHandler(errorStrategy)
	p.BuildParseTrees = true
	return stream, p.EarthFile(), nil
}
