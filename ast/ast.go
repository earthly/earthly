package ast

import (
	"context"
	"fmt"
	"strings"

	"github.com/earthly/earthly/ast/antlrhandler"
	"github.com/earthly/earthly/ast/parser"
	"github.com/earthly/earthly/ast/spec"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"
)

// Parse parses an earthfile into an AST.
func Parse(ctx context.Context, filePath string, enableSourceMap bool) (ef spec.Earthfile, err error) {
	version, err := ParseVersion(filePath, enableSourceMap)
	if err != nil {
		return spec.Earthfile{}, err
	}

	// Convert.
	errorListener := antlrhandler.NewReturnErrorListener()
	errorStrategy := antlrhandler.NewReturnErrorStrategy()
	tree, err := newEarthfileTree(filePath, errorListener, errorStrategy)
	if err != nil {
		return spec.Earthfile{}, err
	}
	ef, walkErr := walkTree(newListener(ctx, filePath, enableSourceMap), tree)
	if len(errorListener.Errs) > 0 {
		errString := []string{fmt.Sprintf("lexer error: %s", filePath)}
		for _, err := range errorListener.Errs {
			errString = append(errString, err.Error())
		}
		return spec.Earthfile{}, errors.Errorf(strings.Join(errString, "\n"))
	}
	if errorStrategy.Err != nil {
		return spec.Earthfile{}, errors.Wrapf(
			errorStrategy.Err, "%s line %d:%d '%s'",
			filePath,
			errorStrategy.RE.GetOffendingToken().GetLine(),
			errorStrategy.RE.GetOffendingToken().GetColumn(),
			errorStrategy.RE.GetOffendingToken().GetText())
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

func walkTree(l *listener, tree parser.IEarthFileContext) (spec.Earthfile, error) {
	antlr.ParseTreeWalkerDefault.Walk(l, tree)
	err := l.Err()
	if err != nil {
		return spec.Earthfile{}, err
	}
	return l.Earthfile(), nil
}

func newEarthfileTree(filename string, errorListener *antlrhandler.ReturnErrorListener, errorStrategy antlr.ErrorStrategy) (parser.IEarthFileContext, error) {
	input, err := antlr.NewFileStream(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "new file stream %s", filename)
	}
	lexer := newLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	if lexer.Err() != nil {
		return nil, lexer.Err()
	}
	p := parser.NewEarthParser(stream)
	p.AddErrorListener(errorListener)
	p.SetErrorHandler(errorStrategy)
	p.BuildParseTrees = true
	return p.EarthFile(), nil
}
