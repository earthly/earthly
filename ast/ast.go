package ast

import (
	"context"
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/earthly/earthly/ast/antlrhandler"
	"github.com/earthly/earthly/ast/parser"
	"github.com/earthly/earthly/ast/spec"
	"github.com/pkg/errors"
)

// Parse parses an earthfile into an AST.
func Parse(ctx context.Context, filePath string, enableSourceMap bool) (ef spec.Earthfile, err error) {
	// Convert.
	errorListener := antlrhandler.NewReturnErrorListener()
	errorStrategy := antlrhandler.NewReturnErrorStrategy()
	tree, err := newEarthfileTree(filePath, errorListener, errorStrategy)
	if err != nil {
		return spec.Earthfile{}, err
	}
	ef, walkErr := walkTree(newListener(ctx, filePath, enableSourceMap), tree)
	if len(errorListener.Errs) > 0 {
		var errString []string
		for _, err := range errorListener.Errs {
			errString = append(errString, err.Error())
		}
		return spec.Earthfile{}, fmt.Errorf(strings.Join(errString, "\n"))
	}
	if errorStrategy.Err != nil {
		var errString []string
		errString = append(errString,
			fmt.Sprintf(
				"syntax error: line %d:%d",
				errorStrategy.RE.GetOffendingToken().GetLine(),
				errorStrategy.RE.GetOffendingToken().GetColumn()))
		errString = append(errString,
			fmt.Sprintf("Details: %s", errorStrategy.RE.GetMessage()))
		return spec.Earthfile{}, errors.Wrapf(errorStrategy.Err, "%s", strings.Join(errString, "\n"))
	}
	if walkErr != nil {
		return spec.Earthfile{}, walkErr
	}
	return ef, nil
}

func walkTree(l *listener, tree parser.IEarthFileContext) (spec.Earthfile, error) {
	antlr.ParseTreeWalkerDefault.Walk(l, tree)
	err := l.Err()
	if err != nil {
		return spec.Earthfile{}, errors.Wrap(err, "parse")
	}
	return l.Earthfile(), nil
}

func newEarthfileTree(filename string, errorListener antlr.ErrorListener, errorStrategy antlr.ErrorStrategy) (parser.IEarthFileContext, error) {
	input, err := antlr.NewFileStream(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "new file stream %s", filename)
	}
	lexer := newLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewEarthParser(stream)
	p.AddErrorListener(errorListener)
	p.SetErrorHandler(errorStrategy)
	p.BuildParseTrees = true
	return p.EarthFile(), nil
}
