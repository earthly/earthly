package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/earthly/earthly/ast/spec"
	"github.com/pkg/errors"
)

// List of valid Earthfile versions.
// At some point we might want to break out Earthfile versioning
// into it's own package with some helper functions that are
// consumable from other packages.
var validEarthfileVersions = []string{
	"0.0",
	"0.5",
}

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

func validVersion(ef spec.Earthfile) []error {
	var errs []error

	// VERSION is not required in Earthfile for now
	if ef.Version == nil {
		return nil
	}

	major, minor, err := parseSemver(ef.Version)
	if err != nil {
		errs = append(errs, err)
	}

	majorStr := strconv.Itoa(major)
	minorStr := strconv.Itoa(minor)
	earthFileVersion := fmt.Sprintf("%s.%s", majorStr, minorStr)

	isVersionValid := false
	for _, version := range validEarthfileVersions {
		if earthFileVersion == version {
			isVersionValid = true
			break
		}
	}

	if !isVersionValid {
		errs = append(errs, errors.Errorf("Earthfile version is invalid, supported versions are %v", validEarthfileVersions))
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
