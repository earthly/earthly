package ast

import (
	"fmt"
	"strings"

	"github.com/earthly/earthly/ast/spec"
)

type astValidator func(spec.Earthfile) []error

var astValidations = []astValidator{
	noTargetsWithSameName,
	// TODO other checks go here
}

func validateAst(ef spec.Earthfile) error {
	var errors []error

	for _, v := range astValidations {
		if err := v(ef); err != nil {
			errors = append(errors, err...)
		}
	}

	if len(errors) > 0 {
		errorStrings := make([]string, len(errors))
		for i, err := range errors {
			errorStrings[i] = err.Error()
		}

		return fmt.Errorf("%v validation issues.\n- %s", len(errors), strings.Join(errorStrings, "\n- "))
	}

	return nil
}

func noTargetsWithSameName(ef spec.Earthfile) []error {
	var errors []error
	seenTargets := map[string]struct{}{}

	for _, t := range ef.Targets {
		if _, seen := seenTargets[t.Name]; seen {
			errors = append(errors, fmt.Errorf("%s line %v:%v duplicate target \"%s\"", t.SourceLocation.File, t.SourceLocation.StartLine, t.SourceLocation.StartColumn, t.Name))
		}

		seenTargets[t.Name] = struct{}{}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
