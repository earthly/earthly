package main

import (
	"fmt"
	"strings"

	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/util/platutil"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (app *earthlyApp) actionDocumentTarget(cliCtx *cli.Context) error {
	app.commandName = "docTarget"

	if cliCtx.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	}

	var tgtPath string
	if cliCtx.NArg() > 0 {
		tgtPath = cliCtx.Args().Get(0)
		switch tgtPath[0] {
		case '.', '/', '+':
		default:
			return errors.New("remote-paths are not currently supported - documentation targets must start with one of ['.', '/', '+']")
		}
	}

	singleTgt := true
	if !strings.ContainsRune(tgtPath, '+') {
		tgtPath += "+base"
		singleTgt = false
	}

	target, err := domain.ParseTarget(tgtPath)
	if err != nil {
		return errors.Errorf("unable to parse target [%v]", tgtPath)
	}

	gitLookup := buildcontext.NewGitLookup(app.console, app.sshAuthSock)
	resolver := buildcontext.NewResolver(nil, gitLookup, app.console, "")
	platr := platutil.NewResolver(platutil.GetUserPlatform())
	var gwClient gwclient.Client
	bc, err := resolver.Resolve(cliCtx.Context, gwClient, platr, target)
	if err != nil {
		return errors.Wrap(err, "failed to resolve target")
	}

	const docsIndent = "  "

	if singleTgt {
		tgt, err := findTarget(bc.Earthfile, target.Target)
		if err != nil {
			return errors.Wrap(err, "failed to look up target")
		}
		return app.documentSingleTarget(cliCtx, "", docsIndent, tgt, true)
	}

	tgts := bc.Earthfile.Targets
	fmt.Println("TARGETS:")
	const tgtIndent = docsIndent
	for _, tgt := range tgts {
		_ = app.documentSingleTarget(cliCtx, tgtIndent, docsIndent, tgt, false)
	}

	return nil
}

func docString(body string, names ...string) (string, error) {
	firstWordEnd := strings.IndexRune(body, ' ')
	if firstWordEnd == -1 {
		return "", errors.Errorf("failed to parse first word of documentation comments")
	}
	firstWord := body[:firstWordEnd]
	for _, n := range names {
		if firstWord == n {
			return body, nil
		}
	}
	return "", errors.Errorf("no doc comment found [hint: a comment was found but the first word was not one of (%s)]", strings.Join(names, ", "))
}

type docSection struct {
	identifier string
	body       string
}

func docSectionsOutput(currIndent, scopeIndent, title string, sections ...docSection) string {
	if len(sections) == 0 {
		return ""
	}
	out := indent(currIndent, title+":") + "\n"
	currIndent += scopeIndent
	for _, section := range sections {
		out += indent(currIndent, section.identifier) + "\n"
		if section.body == "" {
			continue
		}
		indented := indent(currIndent+scopeIndent, section.body)
		out += strings.Trim(indented, "\n") + "\n"
	}
	return out
}

type blockIO struct {
	// TODO: globals
	requiredArgs   []docSection
	optionalArgs   []docSection
	artifacts      []docSection
	localArtifacts []docSection
	images         []docSection
}

func (io blockIO) options() string {
	var options []string
	for _, arg := range io.requiredArgs {
		options = append(options, arg.identifier)
	}
	for _, arg := range io.optionalArgs {
		options = append(options, fmt.Sprintf("[%s]", arg.identifier))
	}
	return strings.Join(options, " ")
}

func (io blockIO) help(indent, scopeIndent string) string {
	return docSectionsOutput(indent, scopeIndent, "REQUIRED ARGS", io.requiredArgs...) +
		docSectionsOutput(indent, scopeIndent, "OPTIONAL ARGS", io.optionalArgs...) +
		docSectionsOutput(indent, scopeIndent, "ARTIFACTS", io.artifacts...) +
		docSectionsOutput(indent, scopeIndent, "LOCAL ARTIFACTS", io.localArtifacts...) +
		docSectionsOutput(indent, scopeIndent, "IMAGES", io.images...)
}

func parseDocSections(cliCtx *cli.Context, cmds spec.Block) (*blockIO, error) {
	var io blockIO
	for _, rb := range cmds {
		if rb.Command == nil {
			continue
		}
		cmd := *rb.Command
		identifiers, err := earthfile2llb.Name(cliCtx.Context, cmd)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse name(s) of command type %v", cmd.Name)
		}
		switch cmd.Name {
		case "ARG":
			var dflt string
			if len(identifiers) == 2 {
				dflt = identifiers[1]
				identifiers = identifiers[:1]
			}
			if len(identifiers) != 1 {
				return nil, errors.Errorf("ARG should have exactly 1 identifier after consuming the default; got %v", len(identifiers))
			}
			docs, _ := docString(cmd.Docs, identifiers...)
			argDoc := docSection{
				identifier: "--" + identifiers[0],
				body:       docs,
			}
			if dflt != "" {
				argDoc.identifier += "=" + dflt
			}
			if isRequired(cmd) {
				io.requiredArgs = append(io.requiredArgs, argDoc)
				continue
			}
			io.optionalArgs = append(io.optionalArgs, argDoc)
		case "SAVE ARTIFACT":
			docs, _ := docString(cmd.Docs, identifiers...)
			artDoc := docSection{
				body: docs,
			}
			if len(identifiers) == 1 {
				artDoc.identifier = identifiers[0]
				io.artifacts = append(io.artifacts, artDoc)
				continue
			}
			artDoc.identifier = fmt.Sprintf("%s -> %s", identifiers[0], identifiers[1])
			io.localArtifacts = append(io.localArtifacts, artDoc)
		case "SAVE IMAGE":
			if len(identifiers) == 0 {
				continue
			}
			docs, _ := docString(cmd.Docs, identifiers...)
			io.images = append(io.images, docSection{
				identifier: strings.Join(identifiers, ", "),
				body:       docs,
			})
		}
	}
	return &io, nil
}

func (app *earthlyApp) documentSingleTarget(cliCtx *cli.Context, currIndent, scopeIndent string, tgt spec.Target, includeBlockDocs bool) error {
	if tgt.Docs == "" {
		return errors.Errorf("no doc comment found [hint: add a comment starting with the word '%s' on the line immediately above this target]", tgt.Name)
	}

	docs, err := docString(tgt.Docs, tgt.Name)
	if err != nil {
		return err
	}

	blockIO, err := parseDocSections(cliCtx, tgt.Recipe)
	if err != nil {
		return errors.Wrapf(err, "failed to parse body of recipe '%v'", tgt.Name)
	}

	usage := indent(currIndent, "+"+tgt.Name)
	options := blockIO.options()
	if options != "" {
		usage += " " + options
	}
	fmt.Println(usage)
	docIndent := currIndent + scopeIndent + scopeIndent
	indented := indent(docIndent, docs)
	fmt.Println(strings.Trim(indented, "\n"))

	if !includeBlockDocs {
		return nil
	}

	fmt.Println(blockIO.help(currIndent+scopeIndent, scopeIndent))
	return nil
}

func indent(indent, s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		if l == "" {
			continue
		}
		lines[i] = indent + l
	}
	return strings.Join(lines, "\n")
}

func isRequired(cmd spec.Command) bool {
	for _, arg := range cmd.Args {
		if arg == "--required" {
			return true
		}
	}
	return false
}

func findTarget(ef spec.Earthfile, name string) (spec.Target, error) {
	for _, tgt := range ef.Targets {
		if tgt.Name == name {
			return tgt, nil
		}
	}
	return spec.Target{}, errors.Errorf("could not find target named [%v]", name)
}
