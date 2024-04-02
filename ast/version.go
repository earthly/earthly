package ast

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/earthly/earthly/ast/spec"

	"github.com/pkg/errors"
)

// ParseVersion reads the VERSION command for an Earthfile and returns spec.Version
func ParseVersion(filePath string, enableSourceMap bool) (*spec.Version, error) {
	var opts []Opt
	if enableSourceMap {
		opts = append(opts, WithSourceMap())
	}
	return ParseVersionOpts(FromPath(filePath), opts...)
}

// ParseVersionOpts reads the VERSION command for an Earthfile and returns a
// spec.Version. This is the functional option version, which uses options to
// change how a file is parsed.
func ParseVersionOpts(fromOpt FromOpt, opts ...Opt) (*spec.Version, error) {
	defaultPrefs := prefs{
		done: func() {},
	}
	prefs, err := fromOpt(defaultPrefs)
	if err != nil {
		return nil, errors.Wrap(err, "ast: could not apply ParseVersion from opt")
	}

	for _, opt := range opts {
		newPrefs, err := opt(prefs)
		if err != nil {
			return nil, errors.Wrap(err, "ast: could not apply ParseVersion opts")
		}
		prefs = newPrefs
	}
	file := prefs.reader
	defer prefs.done()

	var version spec.Version

	foundVersion := false

	scanner := bufio.NewScanner(file)
	i := 0
	var startLine int
	var endLine int
	args := []string{}

outer:
	for scanner.Scan() {
		i++
		l := scanner.Text()
		lineWidth := len(l)
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		if l[0] == '#' {
			continue
		}
		fields := strings.Fields(l)

		if len(fields) == 0 {
			continue
		}

		foundComment := false
		trailingLine := false
		for _, f := range fields {
			if foundComment {
				continue // ignore rest of line
			}
			if trailingLine {
				if strings.HasPrefix(f, "#") {
					foundComment = true
					continue
				}
				// found something other than a '#' after a '\'
				// e.g. VERSION    \    UNEXPECTED
				return nil, fmt.Errorf("malformed trailing line on %s:%d", file.Name(), i)
			}
			if f == "VERSION" && !foundVersion {
				foundVersion = true
				startLine = i
				continue
			}
			if f == `\` {
				trailingLine = true
				continue
			}

			if strings.HasPrefix(f, "#") {
				foundComment = true
				continue
			}

			if !foundVersion {
				// received a keyword other than VERSION
				break outer
			}

			args = append(args, f)
		}

		if trailingLine {
			continue
		}
		endLine = i

		version.Args = args

		if prefs.enableSourceMap {
			version.SourceLocation = &spec.SourceLocation{
				File:        file.Name(),
				StartLine:   startLine,
				StartColumn: 0,
				EndLine:     endLine,
				EndColumn:   lineWidth,
			}
		}

		return &version, nil
	}

	// No version was found
	return nil, nil
}
