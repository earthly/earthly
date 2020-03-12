package earthfile2llb

import (
	"context"
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/pkg/errors"
	"github.com/vladaionescu/earthly/buildcontext"
	"github.com/vladaionescu/earthly/cleanup"
	"github.com/vladaionescu/earthly/domain"
	"github.com/vladaionescu/earthly/earthfile2llb/parser"
	"github.com/vladaionescu/earthly/earthfile2llb/variables"
	"github.com/vladaionescu/earthly/logging"
)

// Earthfile2LLB parses a earthfile and executes the statements for a given target.
func Earthfile2LLB(ctx context.Context, target domain.Target, resolver *buildcontext.Resolver, dockerBuilderFun DockerBuilderFun, cleanCollection *cleanup.Collection, visitedStates map[string][]*SingleTargetStates, buildArgs map[string]variables.Variable) (mts *MultiTargetStates, err error) {
	if visitedStates == nil {
		visitedStates = make(map[string][]*SingleTargetStates)
	}
	// Check if we have previously converted this target, with the same
	// build args.
	targetStr := target.String()
	for _, sts := range visitedStates[targetStr] {
		same := true
		for _, bai := range sts.TargetInput.BuildArgs {
			if sts.Ongoing && !bai.IsConstant {
				return nil, fmt.Errorf(
					"Use of recursive targets with variable build args is not supported: %s", targetStr)
			}
			variable, found := buildArgs[bai.Name]
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
	tree, err := newEarthfileTree(bc.EarthfilePath)
	if err != nil {
		return nil, err
	}
	converter, err := NewConverter(
		targetCtx, bc.Target, resolver, dockerBuilderFun, cleanCollection, bc,
		visitedStates, buildArgs)
	if err != nil {
		return nil, err
	}
	err = walkTree(newListener(targetCtx, converter, target.Target), tree)
	if err != nil {
		return nil, err
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
	if l.err != nil {
		return errors.Wrap(l.err, "parse error")
	}
	return nil
}

// ParseDebug parses a earthfile and prints debug information about it.
func ParseDebug(filename string) error {
	tree, err := newEarthfileTree(filename)
	if err != nil {
		return errors.Wrap(err, "new earthfile tree")
	}
	antlr.ParseTreeWalkerDefault.Walk(newDebugListener(), tree)
	return nil
}

func newEarthfileTree(filename string) (parser.IEarthFileContext, error) {
	input, err := antlr.NewFileStream(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "new file stream %s", filename)
	}
	lexer := newLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewEarthParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.SetErrorHandler(antlr.NewBailErrorStrategy())
	p.BuildParseTrees = true
	return p.EarthFile(), nil
}
