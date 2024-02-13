package docker2earthly

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/earthly/earthly/util/fileutil"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/pkg/errors"
)

// Ideally this would point to "the current version" rather than being hard-coded, but the single
// "source of truth" (in ast/validator) isn't currently exported.
const earthlyCurrentVersion = "0.7"

func getArtifactName(s string) string {
	split := strings.Split(s, "/")
	n := len(split)
	return split[n-1]
}

// Docker2Earthly converts an existing Dockerfile in the current directory and writes out an Earthfile in the current directory
// and error is returned if an Earthfile already exists.
func Docker2Earthly(dockerfilePath, earthfilePath, imageTag string) error {
	if exists, _ := fileutil.FileExists(earthfilePath); exists {
		return errors.Errorf("earthfile already exists; please delete it if you wish to continue")
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
			fmt.Sprintf("VERSION %s\n", earthlyCurrentVersion),
			"# This Earthfile was generated using docker2earthly",
			"# the conversion is done on a best-effort basis",
			"# and might not follow best practices, please",
			"# visit https://docs.earthly.dev for Earthfile guides",
		},
	}

	dockerfile, err := parser.Parse(in)
	if err != nil {
		return errors.Wrapf(err, "failed to parse Dockerfile located at %q", dockerfilePath)
	}

	stages, initialArgs, err := instructions.Parse(dockerfile.AST)
	if err != nil {
		return errors.Wrapf(err, "failed to parse Dockerfile located at %q", dockerfilePath)
	}

	names := map[string]int{}

	for i, stage := range stages {
		fromCmd := []string{fmt.Sprintf("FROM %s", stage.BaseName)}
		// These args are in scope *only* for the very first FROM
		if i == 0 && len(initialArgs) > 0 {
			var fromArgs []string
			for _, arg := range initialArgs {
				fromArgs = append(fromArgs, arg.String())
			}
			fromCmd = append(fromArgs, fromCmd...)
		}
		targets = append(targets, fromCmd)

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
					return errors.Errorf("failed to parse %q", l)
				}
				kv := strings.Split(parts[1], "=")
				if len(kv) != 2 {
					return errors.Errorf("failed to parse %q", l)
				}
				fromStageName := kv[1]
				n := names[fromStageName]
				artifactName := getArtifactName(parts[2])
				l = fmt.Sprintf("COPY +subbuild%d/%s %s", n+1, artifactName, parts[3])
				targets[n+1] = append(targets[n+1], fmt.Sprintf("SAVE ARTIFACT %s %s\n", parts[2], artifactName))
			}
			if strings.HasPrefix(l, "ADD ") {
				return errors.Errorf("earthly does not support ADD, please convert to COPY instead")
			}
			targets[i+1] = append(targets[i+1], l)
		}
	}
	i := len(targets) - 1
	targets[i] = append(targets[i], fmt.Sprintf("SAVE IMAGE %s", imageTag))

	var out io.Writer
	if earthfilePath == "-" {
		out2 := bufio.NewWriter(os.Stdout)
		defer out2.Flush()
		out = out2
	} else {
		out2, err := os.Create(earthfilePath)
		if err != nil {
			return errors.Wrapf(err, "failed to create Earthfile under %q", earthfilePath)
		}
		defer out2.Close()
		out = out2
	}

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
	return nil
}

var earthfileTemplate = `
VERSION --use-docker-ignore {{.Version}}
# This Earthfile was generated using {{.CommandName}} command
docker:
	{{- range .BuildArgs}}
	ARG {{. -}}
	{{- end}}
	FROM DOCKERFILE \
	{{- range .BuildArgs}}
	--build-arg {{.}}=${{.}} \
	{{- end}}
   	{{- if .Target}}
	--target {{.Target}} \
	{{- end}}
	-f {{.Dockerfile}} \
	{{.BuildContext}}
	{{- if .ImageTags}}
	SAVE IMAGE --push
	{{- range .ImageTags}} {{. -}}{{- end}}
	{{- end}}

build:
	BUILD{{- range .Platforms}} --platform {{ . -}}{{- end}} +docker
`

type earthfileTemplateArgs struct {
	Version      string
	CommandName  string
	BuildArgs    []string
	Target       string
	Dockerfile   string
	BuildContext string
	ImageTags    []string
	Platforms    []string
}

// GenerateEarthfile returns an Earthfile content string which contains a target to build a docker image using FROM DOCKERFILE
func GenerateEarthfile(buildContextPath string, dockerfilePath string, imageTags []string, buildArgs []string, platforms []string, target string) (string, error) {
	t, err := template.New("earthfile").Parse(earthfileTemplate)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse Earthfile template")
	}
	buf := &bytes.Buffer{}

	if !filepath.IsAbs(dockerfilePath) {
		dockerfilePath, err = filepath.Abs(filepath.Join(buildContextPath, dockerfilePath))
		if err != nil {
			return "", errors.Wrapf(err, "failed to get get absolute path for dockerfile")
		}
	}

	err = t.Execute(buf, &earthfileTemplateArgs{
		Version:      earthlyCurrentVersion,
		CommandName:  "docker-build",
		BuildArgs:    buildArgs,
		Target:       target,
		Dockerfile:   dockerfilePath,
		BuildContext: buildContextPath,
		ImageTags:    imageTags,
		Platforms:    platforms,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to create Earthfile content from template")
	}

	return buf.String(), nil
}
