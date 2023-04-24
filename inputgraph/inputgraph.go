package inputgraph

import (
	"context"
	"fmt"
	"strings"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
)

func Load(ctx context.Context, target domain.Target, conslog conslogging.ConsoleLogger) {
	if target.IsRemote() {
		panic("remote not supported")
	}
	resolver := buildcontext.NewResolver(nil, nil, conslog, "", "")
	data, err := resolver.Resolve(ctx, nil, nil, target)
	if err != nil {
		panic(err)
	}

	//deps := map[string][]string{}

	for _, tgt := range data.Earthfile.Targets {
		fmt.Printf("%s\n", tgt.Name)
		for _, stmt := range tgt.Recipe {
			if stmt.Command != nil {
				cmd := stmt.Command
				fmt.Printf("here with %+v\n", cmd)
				switch cmd.Name {
				case "FROM":
					fmt.Printf("TODO hash FROM %s\n", cmd.Args)
				case "RUN":
					fmt.Printf("TODO hash RUN %s\n", cmd.Args)
				case "COPY":
					fmt.Printf("TODO parse %s\n", cmd.Args)
				case "BUILD":
					fmt.Printf("TODO recurse into BUILD %s\n", cmd.Args)
				case "SAVE ARTIFACT", "ENTRYPOINT", "SAVE IMAGE":
					break // TODO
				default:
					panic(fmt.Sprintf("not supported %q", cmd.Name))
				}
			} else {
				panic("not supported")
			}
		}
	}

	clean := func(s string) string {
		s = strings.TrimPrefix(s, "+")
		s = strings.Replace(s, ":", "", -1)
		s = strings.Replace(s, ".", "", -1)
		s = strings.Replace(s, "-", "", -1)
		return s
	}

	fmt.Printf("----------------------------------------\n\n")
	fmt.Printf("digraph G {\n")
	for _, tgt := range data.Earthfile.Targets {
		for _, stmt := range tgt.Recipe {
			if stmt.Command != nil {
				cmd := stmt.Command
				switch cmd.Name {
				case "FROM", "BUILD":
					fmt.Printf("  %s -> %s;\n", clean(tgt.Name), clean(cmd.Args[0]))
				case "RUN", "COPY", "SAVE ARTIFACT", "ENTRYPOINT", "SAVE IMAGE":
					break
				default:
					panic(fmt.Sprintf("not supported %q", cmd.Name))
				}
			} else {
				panic("not supported")
			}
		}
	}
	fmt.Printf("}\n")

	//fmt.Printf("got %+v\n", data)
}
