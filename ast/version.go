package ast

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/earthly/earthly/ast/spec"

	"github.com/pkg/errors"
)

func parseVersion(filePath string, enableSourceMap bool) (*spec.Version, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open %q", filePath)
	}
	defer file.Close()

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
				return nil, fmt.Errorf("malformed trailing line on %s:%d", filePath, i)
			}
			if f == "VERSION" && !foundVersion {
				foundVersion = true
				startLine = i
				continue
			}
			if f == "\\" {
				trailingLine = true
				continue
			}

			if strings.HasPrefix(f, "#") {
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

		if enableSourceMap {
			version.SourceLocation = &spec.SourceLocation{
				File:        filePath,
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
