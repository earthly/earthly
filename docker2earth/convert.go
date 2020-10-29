package docker2earth

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/earthly/earthly/fileutils"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/pkg/errors"
)

func getArtifactName(s string) string {
	split := strings.Split(s, "/")
	n := len(split)
	return split[n-1]
}

// Docker2Earth converts an existing Dockerfile in the current directory and writes out an Earthfile in the current directory
// and error is returned if an Earthfile already exists.
func Docker2Earth(dockerfilePath, EarthfilePath string) error {
	if fileutils.FileExists(EarthfilePath) {
		return fmt.Errorf("Earthfile already exists; please delete it if you wish to continue")
	}

	var in io.Reader
	if dockerfilePath == "-" {
		in = bufio.NewReader(os.Stdin)
	} else {
		in2, err := os.Open(dockerfilePath)
		if err != nil {
			return errors.Wrapf(err, "failed to open %q", dockerfilePath)
		}
		defer in2.Close()
		in = in2
	}

	targets := [][]string{
		{
			"# This Earthfile was generated using docker2earth",
			"# the conversion is done on a best-effort basis",
			"# and might not follow best practices, please",
			"# visit http://docs.earthly.dev for Earthfile guides",
		},
	}

	dockerfile, err := parser.Parse(in)
	if err != nil {
		return errors.Wrapf(err, "failed to parse Dockerfile located at %q", dockerfilePath)
	}

	stages, _, err := instructions.Parse(dockerfile.AST)
	if err != nil {
		return errors.Wrapf(err, "failed to parse Dockerfile located at %q", dockerfilePath)
	}

	names := map[string]int{}

	for i, stage := range stages {
		targets = append(targets, []string{
			fmt.Sprintf("FROM %s", stage.BaseName),
		})
		if stage.Name == "" {
			names[fmt.Sprintf("%d", i)] = i
		} else {
			names[stage.Name] = i
		}

		for _, cmd := range stage.Commands {
			l := fmt.Sprintf("%v", cmd)
			if strings.HasPrefix(l, "COPY ") && strings.Contains(l, "--from") {
				parts := strings.Split(l, " ")
				if len(parts) != 4 {
					return fmt.Errorf("failed to parse %q", l)
				}
				kv := strings.Split(parts[1], "=")
				if len(kv) != 2 {
					return fmt.Errorf("failed to parse %q", l)
				}
				fromStageName := kv[1]
				n := names[fromStageName]
				artifactName := getArtifactName(parts[2])
				l = fmt.Sprintf("COPY +subbuild%d/%s %s", n+1, artifactName, parts[3])
				targets[n+1] = append(targets[n+1], fmt.Sprintf("SAVE ARTIFACT %s %s\n", parts[2], artifactName))
			}
			if strings.HasPrefix(l, "ADD ") {
				return fmt.Errorf("earthly does not support ADD, please convert to COPY instead")
			}
			targets[i+1] = append(targets[i+1], l)
		}
	}
	i := len(targets) - 1
	targets[i] = append(targets[i], "SAVE IMAGE myimage:latest")

	var out io.Writer
	if EarthfilePath == "-" {
		out2 := bufio.NewWriter(os.Stdout)
		defer out2.Flush()
		out = out2
	} else {
		out2, err := os.Create(EarthfilePath)
		if err != nil {
			return errors.Wrap(err, "failed to create Earthfile")
		}
		defer out2.Close()
		out = out2
	}

	fmt.Fprintf(out, "\n")

	for i, lines := range targets {
		for j, l := range lines {
			if i == 0 {
				fmt.Fprintf(out, "%s\n", l)
			} else {
				if j == 0 {
					fmt.Fprintf(out, "subbuild%d:\n", i)
				}
				fmt.Fprintf(out, "    %s\n", l)
			}
		}
	}

	fmt.Fprintf(out, "\nbuild:\n    BUILD +subbuild%d\n", i)

	fmt.Fprintf(os.Stderr, "An Earthfile has been generated; to run it use: earth +build; then run with docker run -ti myimage:latest\n")
	return nil
}
