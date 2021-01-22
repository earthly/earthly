package llbutil

import "github.com/moby/buildkit/client/llb"

var (
	// DefaultPathEnv is the default PATH to use.
	DefaultPathEnv = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
)

type StatesAdapter func(llb.State) llb.State
