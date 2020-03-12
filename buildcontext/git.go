package buildcontext

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
	"github.com/vladaionescu/earthly/conslogging"
	"github.com/vladaionescu/earthly/domain"
	"github.com/vladaionescu/earthly/llbutil"
	"github.com/vladaionescu/earthly/llbutil/llbgit"
	"golang.org/x/sync/errgroup"
)

const (
	defaultGitImage = "alpine/git:v2.24.1"
)

type gitResolver struct {
	bkClient *client.Client
	console  conslogging.ConsoleLogger

	projectCache map[string]resolvedGitProject
}

type resolvedGitProject struct {
	// localGitDir is where the git dir exists locally (only build.earth files).
	localGitDir string
	// hash is the git hash.
	hash string
	// state is the state holding the git files.
	state llb.State
}

func (gr *gitResolver) resolveEarthProject(ctx context.Context, target domain.Target) (*Data, error) {
	if !target.IsRemote() {
		return nil, fmt.Errorf("Unexpected local target %s", target.String())
	}
	rgp, gitURL, subDir, err := gr.resolveGitProject(ctx, target)
	if err != nil {
		return nil, err
	}

	// Restrict the resulting build context to the right subdir.
	var buildContext llb.State
	if subDir == "." {
		// Optimization.
		buildContext = rgp.state
	} else {
		buildContext = llb.Scratch().Platform(llbutil.TargetPlatform)
		buildContext = llbutil.CopyOp(
			rgp.state, []string{subDir}, buildContext, "./", false,
			llb.WithCustomNamef("[internal] COPY git context %s", target.String()))
	}

	// TODO: Apply excludes / .earthignore.
	return &Data{
		EarthfilePath: filepath.Join(
			rgp.localGitDir, filepath.FromSlash(subDir), "build.earth"),
		BuildContext: buildContext,
		GitMetadata: &GitMetadata{
			BaseDir:    "",
			RelDir:     subDir,
			RemoteURL:  gitURL,
			GitVendor:  target.Registry,
			GitProject: target.ProjectPath,
			Hash:       rgp.hash,
			Branch:     nil,
			Tags:       nil,
		},
	}, nil
}

func (gr *gitResolver) resolveGitProject(ctx context.Context, target domain.Target) (rgp *resolvedGitProject, gitURL string, subDir string, finalErr error) {
	if target.Registry != "github.com" {
		return nil, "", "", fmt.Errorf("Registry %s not supported", target.Registry)
	}
	projectPathParts := strings.Split(target.ProjectPath, "/")
	if len(projectPathParts) < 2 {
		return nil, "", "", fmt.Errorf("Invalid github project path %s", target.ProjectPath)
	}
	githubUsername := projectPathParts[0]
	githubProject := projectPathParts[1]
	subDir = strings.Join(projectPathParts[2:], "/")
	gitURL = fmt.Sprintf("git@github.com:%s/%s.git", githubUsername, githubProject)
	ref := target.Tag
	if ref == "" {
		ref = "HEAD"
	}

	// Check the cache first.
	cacheKey := fmt.Sprintf("%s#%s", gitURL, ref)
	data, found := gr.projectCache[cacheKey]
	if found {
		return &data, gitURL, subDir, nil
	}
	// Not cached.

	// Copy all build.earth files.
	gitState := llbgit.Git(gitURL, ref, llb.KeepGitDir())
	copyOpts := []llb.RunOption{
		llb.Args([]string{
			"find",
			"-type", "f",
			"-name", "build.earth",
			"-exec", "cp", "--parents", "{}", "/dest", ";",
		}),
		llb.Dir("/git-src"),
		llb.ReadonlyRootFS(),
		llb.AddMount("/git-src", gitState, llb.Readonly),
		llb.WithCustomNamef("COPY GIT CLONE %s/build.earth", target.ProjectCanonical()),
	}
	opImg := llb.Image(
		defaultGitImage, llb.MarkImageInternal,
		llb.Platform(llbutil.TargetPlatform))
	copyOp := opImg.Run(copyOpts...)
	earthfileState := copyOp.AddMount("/dest", llb.Scratch().Platform(llbutil.TargetPlatform))

	// Get git hash.
	gitHashOpts := []llb.RunOption{
		llb.Args([]string{"/bin/sh", "-c", "git rev-parse HEAD >/dest/git-hash"}),
		llb.Dir("/git-src"),
		llb.ReadonlyRootFS(),
		llb.AddMount("/git-src", gitState, llb.Readonly),
		llb.WithCustomNamef("GET GIT HASH %s", target.ProjectCanonical()),
	}
	gitHashOp := opImg.Run(gitHashOpts...)
	hashAndEarthfileState := gitHashOp.AddMount("/dest", earthfileState)

	// Build.
	dt, err := hashAndEarthfileState.Marshal(llb.LinuxAmd64)
	if err != nil {
		return nil, "", "", errors.Wrapf(err, "get build.earth from %s", target.ProjectCanonical())
	}
	earthfileTmpDir, err := ioutil.TempDir("/tmp", "earthly-git")
	if err != nil {
		return nil, "", "", errors.Wrap(err, "create temp dir for build.earth")
	}
	defer func() {
		if finalErr != nil {
			os.RemoveAll(earthfileTmpDir)
		}
	}()
	solveOpt := newSolveOptGit(earthfileTmpDir)
	ch := make(chan *client.SolveStatus)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		_, err = gr.bkClient.Solve(ctx, dt, *solveOpt, ch)
		return err
	})
	eg.Go(func() error {
		for {
			select {
			case ss, ok := <-ch:
				if !ok {
					return nil
				}
				for _, vertex := range ss.Vertexes {
					if vertex.Error != "" {
						gr.console.Printf("ERROR: %s\n", vertex.Error)
					}
				}
			case <-ctx.Done():
				return nil
			}
		}
	})
	err = eg.Wait()
	if err != nil {
		return nil, "", "", err
	}

	// Use built files.
	gitHashFile, err := os.Open(filepath.Join(earthfileTmpDir, "git-hash"))
	if err != nil {
		return nil, "", "", errors.Wrap(err, "open git hash file after solve")
	}
	gitHashBytes, err := ioutil.ReadAll(gitHashFile)
	if err != nil {
		return nil, "", "", errors.Wrap(err, "read git hash after solve")
	}
	gitHash := strings.SplitN(string(gitHashBytes), "\n", 2)[0]

	// Add to cache.
	resolved := resolvedGitProject{
		localGitDir: earthfileTmpDir,
		hash:        gitHash,
		state: llbgit.Git(
			gitURL,
			gitHash,
			llb.WithCustomNamef("[context] git context %s", target.StringCanonical()),
		),
	}
	gr.projectCache[cacheKey] = resolved
	return &resolved, gitURL, subDir, nil
}

func (gr *gitResolver) close() error {
	var lastErr error
	for _, rgp := range gr.projectCache {
		lastErr = os.RemoveAll(rgp.localGitDir)
	}
	if lastErr != nil {
		return errors.Wrap(lastErr, "remove temp git dir")
	}
	return nil
}

func newSolveOptGit(outDir string) *client.SolveOpt {
	return &client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type:      client.ExporterLocal,
				OutputDir: outDir,
			},
		},
		LocalDirs: make(map[string]string),
	}
}
