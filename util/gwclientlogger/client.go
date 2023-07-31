package gwclientlogger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	digest "github.com/opencontainers/go-digest"
)

type verboseClient struct {
	c gwclient.Client
}

// New returns a new gateway client that logs all calls to the wrapped client
func New(c gwclient.Client) gwclient.Client {
	return &verboseClient{
		c: c,
	}
}

// Solve wraps gwclient.Solve
func (vc *verboseClient) Solve(ctx context.Context, req gwclient.SolveRequest) (*gwclient.Result, error) {
	reqStr, _ := json.MarshalIndent(req, "", "\t")
	res, err := vc.c.Solve(ctx, req)
	resStr, _ := json.MarshalIndent(res, "", "\t")
	fmt.Printf("Solve req=%s res=%s; err=%v\n", reqStr, resStr, err)
	return res, err
}

// Export wraps gwclient.Export
func (vc *verboseClient) Export(ctx context.Context, req gwclient.ExportRequest) error {
	reqStr, _ := json.MarshalIndent(req, "", "\t")
	err := vc.c.Export(ctx, req)
	fmt.Printf("Export req=%s; err=%v\n", reqStr, err)
	return err
}

// ResolveImageConfig wraps gwclient.ResolveImageConfig
func (vc *verboseClient) ResolveImageConfig(ctx context.Context, ref string, opt llb.ResolveImageConfigOpt) (string, digest.Digest, []byte, error) {
	s, _ := json.MarshalIndent(opt, "", "\t")
	fmt.Printf("ResolveImageConfig %s %s\n", ref, s)
	return vc.c.ResolveImageConfig(ctx, ref, opt)
}

// BuildOpts wraps gwclient.BuildOpts
func (vc *verboseClient) BuildOpts() gwclient.BuildOpts {
	opts := vc.c.BuildOpts()
	fmt.Printf("BuildOpts res=%v\n", opts)
	return opts
}

// Inputs wraps gwclient.Inputs
func (vc *verboseClient) Inputs(ctx context.Context) (map[string]llb.State, error) {
	inputs, err := vc.c.Inputs(ctx)
	fmt.Printf("Inputs=%v err=%v\n", inputs, err)
	return inputs, err
}

// NewContainer wraps gwclient.NewContainer
func (vc *verboseClient) NewContainer(ctx context.Context, req gwclient.NewContainerRequest) (gwclient.Container, error) {
	s, _ := json.MarshalIndent(req, "", "\t")
	container, err := vc.c.NewContainer(ctx, req)
	fmt.Printf("NewContainer req=%s container=%v err=%v\n", s, container, err)
	return container, err
}

// Warn wraps gwclient.Warn
func (vc *verboseClient) Warn(ctx context.Context, dgst digest.Digest, msg string, warnOpts gwclient.WarnOpts) error {
	return vc.c.Warn(ctx, dgst, msg, warnOpts)
}
