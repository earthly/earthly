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

// ConvertOpt holds conversion parameters needed for conversion.
type ConvertOpt struct {
	// Resolver is the build context resolver.
	Resolver *buildcontext.Resolver
	// DockerBuilderFun is a fun that can be used to execute an image build. This
	// is used as part of operations like DOCKER LOAD and DOCKER PULL, where
	// a tar image is needed in the middle of a build.
	DockerBuilderFun DockerBuilderFun
	// ArtifactBuilderFun is a fun that can be used to execute build of an artifact.
	// This is used as part of operations like FROM DOCKERFILE +.../..., where
	// a generated Dockerfile is needed in the middle of a build.
	ArtifactBuilderFun ArtifactBuilderFun
	// CleanCollection is a collection of cleanup functions.
	CleanCollection *cleanup.Collection
	// VisitedStates is a collection of target states which have been converted to LLB.
	// This is used for deduplication and infinite cycle detection.
	VisitedStates map[string][]*SingleTargetStates
	// VarCollection is a collection of build args used for overriding args in the build.
	VarCollection *variables.Collection
}

// DockerBuilderFun is a function able to build a target into a docker tar file.
type DockerBuilderFun = func(ctx context.Context, mts *MultiTargetStates, dockerTag string, outFile string) error

// ArtifactBuilderFun is a function able to build an artifact and output it locally.
type ArtifactBuilderFun = func(ctx context.Context, mts *MultiTargetStates, artifact domain.Artifact, outFile string) error

// Earthfile2LLB parses a earthfile and executes the statements for a given target.
func Earthfile2LLB(ctx context.Context, target domain.Target, opt ConvertOpt) (mts *MultiTargetStates, err error) {
	if opt.VisitedStates == nil {
		opt.VisitedStates = make(map[string][]*SingleTargetStates)
	}
	// Check if we have previously converted this target, with the same build args.
	targetStr := target.String()
	for _, sts := range opt.VisitedStates[targetStr] {
		same := true
		for _, bai := range sts.TargetInput.BuildArgs {
			if sts.Ongoing && !bai.IsConstant {
				return nil, fmt.Errorf(
					"Use of recursive targets with variable build args is not supported: %s", targetStr)
			}
			variable, _, found := opt.VarCollection.Get(bai.Name)
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
				VisitedStates: opt.VisitedStates,
			}, nil
		}
	}
	// Resolve build context.
	bc, err := opt.Resolver.Resolve(ctx, target)
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
	converter, err := NewConverter(targetCtx, bc.Target, bc, opt)
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

// GetTargets returns a list of targets from an Earthfile
func GetTargets(filename string) ([]string, error) {
	tree, err := newEarthfileTree(
		filename, antlr.NewConsoleErrorListener(), antlr.NewBailErrorStrategy())
	if err != nil {
		return nil, errors.Wrap(err, "new earthfile tree")
	}
	tc := &targetCollector{}
	antlr.ParseTreeWalkerDefault.Walk(tc, tree)
	return tc.targets, nil
}

type targetCollector struct {
	*parser.BaseEarthParserListener
	targets []string
}

func (l *targetCollector) EnterTarget(ctx *parser.TargetContext) {
	l.targets = append(l.targets, strings.TrimSuffix(ctx.TargetHeader().GetText(), ":"))
}
