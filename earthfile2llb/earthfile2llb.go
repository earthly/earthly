package earthfile2llb

import (
	"context"
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb/antlrhandler"
	"github.com/earthly/earthly/earthfile2llb/parser"
	"github.com/earthly/earthly/earthfile2llb/variables"
	"github.com/earthly/earthly/logging"
	"github.com/pkg/errors"
)

// Earthfile2LLB parses a earthfile and executes the statements for a given target.
func Earthfile2LLB(ctx context.Context, target domain.Target, resolver *buildcontext.Resolver, dockerBuilderFun DockerBuilderFun, cleanCollection *cleanup.Collection, visitedStates map[string][]*SingleTargetStates, varCollection *variables.Collection) (mts *MultiTargetStates, err error) {
	if visitedStates == nil {
		visitedStates = make(map[string][]*SingleTargetStates)
	}
	// Check if we have previously converted this target, with the same build args.
	targetStr := target.String()
	for _, sts := range visitedStates[targetStr] {
		same := true
		for _, bai := range sts.TargetInput.BuildArgs {
			if sts.Ongoing && !bai.IsConstant {
				return nil, fmt.Errorf(
					"Use of recursive targets with variable build args is not supported: %s", targetStr)
			}
			variable, _, found := varCollection.Get(bai.Name)
			if found {
				if !variable.BuildArgInput(bai.Name, bai.DefaultValue).Equals(bai) {
					same = false
					break
				}
			} else {
				if !bai.IsDefaultValue() {
					same = false
					break
				}
			}
		}
		if same {
			if sts.Ongoing {
				return nil, fmt.Errorf(
					"Infinite recursion detected for target %s", targetStr)
			}
			// Use the already built states.
			return &MultiTargetStates{
				FinalStates:   sts,
				VisitedStates: visitedStates,
			}, nil
		}
	}
	// Resolve build context.
	bc, err := resolver.Resolve(ctx, target)
	if err != nil {
		return nil, errors.Wrapf(err, "resolve build context for target %s", target.String())
	}
	// Convert.
	targetCtx := logging.With(ctx, "target", target)
	errorListener := antlrhandler.NewReturnErrorListener()
	errorStrategy := antlrhandler.NewReturnErrorStrategy()
	tree, err := newEarthfileTree(bc.BuildFilePath, errorListener, errorStrategy)
	if err != nil {
		return nil, err
	}
	converter, err := NewConverter(
		targetCtx, bc.Target, resolver, dockerBuilderFun, cleanCollection, bc,
		visitedStates, varCollection)
	if err != nil {
		return nil, err
	}
	walkErr := walkTree(newListener(targetCtx, converter, target.Target), tree)
	if len(errorListener.Errs) > 0 {
		var errString []string
		for _, err := range errorListener.Errs {
			errString = append(errString, err.Error())
		}
		return nil, fmt.Errorf(strings.Join(errString, "\n"))
	}
	if errorStrategy.Err != nil {
		var errString []string
		errString = append(errString,
			fmt.Sprintf(
				"Syntax error: line %d:%d when parsing %s",
				errorStrategy.RE.GetOffendingToken().GetLine(),
				errorStrategy.RE.GetOffendingToken().GetColumn(),
				errorStrategy.ErrContext.GetText()))
		errString = append(errString,
			fmt.Sprintf("Details: %s", errorStrategy.RE.GetMessage()))
		return nil, errors.Wrapf(errorStrategy.Err, "%s", strings.Join(errString, "\n"))
	}
	if walkErr != nil {
		return nil, walkErr
	}
	return converter.FinalizeStates(), nil
}

func walkTree(l *listener, tree parser.IEarthFileContext) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("parser failure: %v", r)
		}
	}()
	antlr.ParseTreeWalkerDefault.Walk(l, tree)
	err = l.Err()
	if err != nil {
		return errors.Wrap(err, "parse")
	}
	return nil
}

// ParseDebug parses a earthfile and prints debug information about it.
func ParseDebug(filename string) error {
	tree, err := newEarthfileTree(
		filename, antlr.NewConsoleErrorListener(), antlr.NewBailErrorStrategy())
	if err != nil {
		return errors.Wrap(err, "new earthfile tree")
	}
	antlr.ParseTreeWalkerDefault.Walk(newDebugListener(), tree)
	return nil
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
