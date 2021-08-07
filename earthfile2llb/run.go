package earthfile2llb

import "github.com/moby/buildkit/client/llb"

// ConvertRunOpts represents a set of options needed for the RUN command.
type ConvertRunOpts struct {
	CommandName     string
	Args            []string
	Locally         bool
	Mounts          []string
	Secrets         []string
	WithEntrypoint  bool
	WithShell       bool
	Privileged      bool
	Push            bool
	Transient       bool
	WithSSH         bool
	NoCache         bool
	Interactive     bool
	InteractiveKeep bool

	// Internal.
	shellWrap    shellWrapFun
	extraRunOpts []llb.RunOption
}
