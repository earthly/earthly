package docker2earth

import (
	"fmt"
	"os"
	"strconv"
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
func Docker2Earth() error {
	if fileutils.FileExists("Earthfile") {
		return fmt.Errorf("Earthfile already exists; please delete it if you wish to continue")
	}

	in, err := os.Open("Dockerfile")
	if err != nil {
		return errors.Wrap(err, "failed to open ./Dockerfile")
	}
	defer in.Close()

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
		return errors.Wrap(err, "failed to parse Dockerfile")
	}

	stages, _, err := instructions.Parse(dockerfile.AST)
	if err != nil {
		return errors.Wrap(err, "failed to parse Dockerfile")
	}

	//shlex := shell.NewLex(dockerfile.EscapeToken)

	//for _, cmd := range metaArgs {
	//	for _, metaArg := range cmd.Args {
	//		if metaArg.Value != nil {
	//			*metaArg.Value, _ = shlex.ProcessWordWithMap(*metaArg.Value, metaArgsToMap(optMetaArgs))
	//		}
	//		optMetaArgs = append(optMetaArgs, setKVValue(metaArg, opt.BuildArgs))
	//	}
	//}

	for i, stage := range stages {
		targets = append(targets, []string{
			fmt.Sprintf("FROM %s", stage.BaseName),
		})

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
				n, err := strconv.Atoi(kv[1])
				if err != nil {
					return fmt.Errorf("failed to parse %q", l)
				}
				artifactName := getArtifactName(parts[2])
				_ = n
				_ = artifactName
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

	out, err := os.Create("Earthfile")
	if err != nil {
		return errors.Wrap(err, "failed to create Earthfile")
	}
	defer out.Close()

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

	fmt.Printf("An Earthfile has been generated; to run it use: earth +build; then run with docker run -ti myimage:latest\n")
	return nil
}
