package llbgit

import (
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/util/apicaps"
)

// Git is the same as llb.Git, but with bugs fixed (git@github.com works).
func Git(remote, ref string, opts ...llb.GitOption) llb.State {
	url := ""

	for _, prefix := range []string{
		"http://", "https://", "git://",
	} {
		if strings.HasPrefix(remote, prefix) {
			url = strings.Split(remote, "#")[0]
			remote = strings.TrimPrefix(remote, prefix)
			break
		}
	}

	id := remote

	if ref != "" {
		id += "#" + ref
	}

	gi := &llb.GitInfo{}
	// gi := &llb.GitInfo{
	// 	AuthHeaderSecret: "GIT_AUTH_HEADER",
	// 	AuthTokenSecret:  "GIT_AUTH_TOKEN",
	// }
	for _, o := range opts {
		o.SetGitOption(gi)
	}
	attrs := map[string]string{}
	if gi.KeepGitDir {
		attrs[pb.AttrKeepGitDir] = "true"
		addCap(&gi.Constraints, pb.CapSourceGitKeepDir)
	}
	if url != "" {
		attrs[pb.AttrFullRemoteURL] = url
		addCap(&gi.Constraints, pb.CapSourceGitFullURL)
	}
	// if gi.AuthTokenSecret != "" {
	// 	attrs[pb.AttrAuthTokenSecret] = gi.AuthTokenSecret
	// 	addCap(&gi.Constraints, pb.CapSourceGitHttpAuth)
	// }
	// if gi.AuthHeaderSecret != "" {
	// 	attrs[pb.AttrAuthHeaderSecret] = gi.AuthHeaderSecret
	// 	addCap(&gi.Constraints, pb.CapSourceGitHttpAuth)
	// }

	addCap(&gi.Constraints, pb.CapSourceGit)

	source := llb.NewSource("git://"+id, attrs, gi.Constraints)
	return llb.NewState(source.Output())
}

func addCap(c *llb.Constraints, id apicaps.CapID) {
	if c.Metadata.Caps == nil {
		c.Metadata.Caps = make(map[apicaps.CapID]bool)
	}
	c.Metadata.Caps[id] = true
}
