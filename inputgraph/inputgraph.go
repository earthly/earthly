package inputgraph

import (
	"context"
	"fmt"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
)

func Load(ctx context.Context, target domain.Target, conslog conslogging.ConsoleLogger) {
	if target.IsRemote() {
		panic("remote not supported")
	}
	fmt.Printf("target is %s\n", target.String())
	resolver := buildcontext.NewResolver(nil, nil, conslog, "", "")
	data, err := resolver.Resolve(ctx, nil, nil, target)
	if err != nil {
		panic(err)
	}

	converter := &StubConverter{
		conslog: conslog,
		target:  target,
	}
	interpreter := earthfile2llb.NewInterpreter(converter, target, false, false, conslog, nil)
	err = interpreter.Run(ctx, data.Earthfile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("interpreter.Run finished without an error\n")
}
