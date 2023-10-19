package ast

import (
	"fmt"
	"strings"

	"github.com/earthly/earthly/ast/spec"
	"github.com/pkg/errors"
)

// List of valid Earthfile versions.
// At some point we might want to break out Earthfile versioning
// into it's own package with some helper functions that are
// consumable from other packages.
var validEarthfileVersions = []string{
	"0.0", // Meant only for testing/debugging. Disables all feature flags.
	"0.6",
	"0.7",
	"0.8",
}

var errUnexpectedVersionArgs = fmt.Errorf("unexpected VERSION arguments; should be VERSION [flags] <major-version>.<minor-version>")

type astValidator func(spec.Earthfile) []error

var astValidations = []astValidator{
	noTargetsWithSameName,
	noTargetsWithKeywords,
	validVersion,
	// TODO other checks go here
}

func validateAst(ef spec.Earthfile) error {
	var errs []error

	for _, v := range astValidations {
		if err := v(ef); err != nil {
			errs = append(errs, err...)
		}
	}

	if len(errs) > 0 {
		errorStrings := make([]string, len(errs))
		for i, err := range errs {
			errorStrings[i] = err.Error()
		}

		return errors.Errorf("%v validation issues.\n- %s", len(errs), strings.Join(errorStrings, "\n- "))
	}

	return nil
}

func getValidVersionsFormatted() string {
	if validEarthfileVersions[0] != "0.0" {
		panic("validEarthfileVersions should start with 0.0")
	}
	var sb strings.Builder
	latestIndex := len(validEarthfileVersions) - 1
	for i := 1; i < latestIndex; i++ {
		sb.WriteString(validEarthfileVersions[i] + ", ")
	}
	sb.WriteString("or " + validEarthfileVersions[latestIndex])
	return sb.String()
}

func validVersion(ef spec.Earthfile) []error {
	var errs []error

	// VERSION is not required in Earthfile for now
	if ef.Version == nil {
		return nil
	}

	// if VERSION is specified, it's invalid to have no args
	if len(ef.Version.Args) == 0 {
		errs = append(errs, errUnexpectedVersionArgs)
		return errs
	}

	// version is always last in VERSION command
	earthFileVersion := ef.Version.Args[len(ef.Version.Args)-1]

	isVersionValid := false
	for _, version := range validEarthfileVersions {
		if version == earthFileVersion {
			isVersionValid = true
			break
		}
	}

	if !isVersionValid {
		errs = append(errs, errors.Errorf("Earthfile version is invalid, supported versions are %v", getValidVersionsFormatted()))
	}

	return errs
}

func noTargetsWithSameName(ef spec.Earthfile) []error {
	var errs []error
	seenTargets := map[string]struct{}{}

	for _, t := range ef.Targets {
		if _, seen := seenTargets[t.Name]; seen {
			errs = append(errs, errors.Errorf("%s line %v:%v duplicate target \"%s\"", t.SourceLocation.File, t.SourceLocation.StartLine, t.SourceLocation.StartColumn, t.Name))
		}

		seenTargets[t.Name] = struct{}{}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func noTargetsWithKeywords(ef spec.Earthfile) []error {
	var errs []error

	for _, t := range ef.Targets {
		if t.Name == "base" {
			errs = append(errs, errors.Errorf("%s line %v:%v invalid target \"%s\": %s is a reserved target name", t.SourceLocation.File, t.SourceLocation.StartLine, t.SourceLocation.StartColumn, t.Name, t.Name))
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
