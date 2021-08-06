package earthfile2llb

import "github.com/moby/buildkit/client/llb"

type convertRunOpts struct {
	CommandName     string
	Args            []string
	Mounts          []string
	Secrets         []string
	Privileged      bool
	WithEntrypoint  bool
	WithDocker      bool
	IsWithShell     bool
	Push            bool
	WithSSH         bool
	NoCache         bool
	Interactive     bool
	InteractiveKeep bool
}

type convertInternalRunOpts struct {
	CommandName     string
	Args            []string
	Mounts          []string
	Secrets         []string
	IsWithShell     bool
	ShellWrap       shellWrapFun
	Privileged      bool
	Push            bool
	Transient       bool
	WithSSH         bool
	NoCache         bool
	Interactive     bool
	InteractiveKeep bool
	extraRunOpts    []llb.RunOption

	// TODO: Unify
	Locally bool
}
