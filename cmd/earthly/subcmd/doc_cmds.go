package subcmd

import (
	"fmt"
	"strings"

	"github.com/earthly/earthly/ast/hint"
	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/buildcontext"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/util/platutil"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type Doc struct {
	cli CLI

	docShowLong bool
}

func NewDoc(cli CLI) *Doc {
	return &Doc{
		cli: cli,
	}
}

func (a *Doc) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "doc",
			Usage:       "Document targets from an Earthfile",
			UsageText:   "earthly [options] doc [<project-ref>[+<target-ref>]]",
			Description: "Document targets from an Earthfile by reading in line comments.",
			Action:      a.action,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "long",
					Aliases:     []string{"l"},
					Usage:       "Show full details for all target inputs and outputs",
					Destination: &a.docShowLong,
				},
			},
		},
	}
}

func (a *Doc) action(cliCtx *cli.Context) error {
	a.cli.SetCommandName("docTarget")

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
		return errors.Errorf("unable to parse target %q", tgtPath)
	}

	gitLookup := buildcontext.NewGitLookup(a.cli.Console(), a.cli.Flags().SSHAuthSock)
	resolver := buildcontext.NewResolver(nil, gitLookup, a.cli.Console(), "", a.cli.Flags().GitBranchOverride, "", 0, "")
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
		return a.documentSingleTarget(cliCtx, "", docsIndent, bc.Features, bc.Earthfile.BaseRecipe, tgt, true)
	}

	tgts := bc.Earthfile.Targets
	fmt.Println("TARGETS:")
	const tgtIndent = docsIndent
	for _, tgt := range tgts {
		_ = a.documentSingleTarget(cliCtx, tgtIndent, docsIndent, bc.Features, bc.Earthfile.BaseRecipe, tgt, a.docShowLong)
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
	return "", hint.Wrapf(errors.New("no doc comment found"), "a comment was found but the first word was not one of (%s)", strings.Join(names, ", "))
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

func addArg(cliCtx *cli.Context, io *blockIO, ft *features.Features, stmt spec.Statement, isBase bool, onlyGlobal bool) error {
	if stmt.Command == nil {
		return nil
	}
	cmd := *stmt.Command
	if cmd.Name != "ARG" {
		return nil
	}
	ident, dflt, isRequired, isGlobal, err := earthfile2llb.ArgName(cliCtx.Context, cmd, isBase, ft.ExplicitGlobal)
	if err != nil {
		return errors.Wrap(err, "failed to parse ARG statement")
	}
	if onlyGlobal && !isGlobal {
		return nil
	}
	docs, _ := docString(cmd.Docs, ident)
	doc := docSection{
		identifier: "--" + ident,
		body:       docs,
	}
	if dflt != nil {
		doc.identifier += "=" + *dflt
	}
	if isRequired {
		io.requiredArgs = append(io.requiredArgs, doc)
		return nil
	}
	io.optionalArgs = append(io.optionalArgs, doc)
	return nil
}

func parseDocSections(cliCtx *cli.Context, ft *features.Features, baseRcp, cmds spec.Block) (*blockIO, error) {
	var io blockIO
	for _, base := range baseRcp {
		err := addArg(cliCtx, &io, ft, base, true, true)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse global ARG in base recipe")
		}
	}
	for _, rb := range cmds {
		if rb.Command == nil {
			continue
		}
		cmd := *rb.Command
		switch cmd.Name {
		case "ARG":
			err := addArg(cliCtx, &io, ft, rb, false, false)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse non-global ARG")
			}
		case "SAVE ARTIFACT":
			name, localName, err := earthfile2llb.ArtifactName(cliCtx.Context, cmd)
			if err != nil {
				return nil, errors.Wrap(err, "could not parse SAVE ARTIFACT name")
			}
			idents := []string{name}
			if localName != nil {
				idents = append(idents, *localName)
			}
			docs, _ := docString(cmd.Docs, idents...)
			artDoc := docSection{
				identifier: name,
				body:       docs,
			}
			if localName != nil {
				artDoc.identifier += " -> " + *localName
				io.localArtifacts = append(io.localArtifacts, artDoc)
				continue
			}
			io.artifacts = append(io.artifacts, artDoc)
		case "SAVE IMAGE":
			identifiers, err := earthfile2llb.ImageNames(cliCtx.Context, cmd)
			if err != nil {
				return nil, errors.Wrap(err, "could not parse SAVE IMAGE name(s)")
			}
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

func (a *Doc) documentSingleTarget(cliCtx *cli.Context, currIndent, scopeIndent string, ft *features.Features, baseRcp spec.Block, tgt spec.Target, includeBlockDocs bool) error {
	if tgt.Docs == "" {
		return hint.Wrapf(errors.New("no doc comment found"), "add a comment starting with the word '%s' on the line immediately above this target", tgt.Name)
	}

	docs, err := docString(tgt.Docs, tgt.Name)
	if err != nil {
		return err
	}

	blockIO, err := parseDocSections(cliCtx, ft, baseRcp, tgt.Recipe)
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

func findTarget(ef spec.Earthfile, name string) (spec.Target, error) {
	for _, tgt := range ef.Targets {
		if tgt.Name == name {
			return tgt, nil
		}
	}
	return spec.Target{}, errors.Errorf("could not find target named %q", name)
}
