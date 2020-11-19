package buildcontext

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/llbutil"
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

const (
	defaultGitImage = "alpine/git:v2.24.1"
)

type gitResolver struct {
	cleanCollection *cleanup.Collection

	projectCache map[string]*resolvedGitProject
	gitLookup    *GitLookup
}

type resolvedGitProject struct {
	// gitMetaAndEarthfileRef is the ref containing the git metadata and build files.
	gitMetaAndEarthfileRef gwclient.Reference
	// hash is the git hash.
	hash string
	// branches is the git branches.
	branches []string
	// tags is the git tags
	tags []string
	// state is the state holding the git files.
	state llb.State
}

func (gr *gitResolver) resolveEarthProject(ctx context.Context, gwClient gwclient.Client, target domain.Target) (*Data, error) {
	if !target.IsRemote() {
		return nil, fmt.Errorf("Unexpected local target %s", target.String())
	}
	rgp, gitURL, subDir, err := gr.resolveGitProject(ctx, gwClient, target)
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
			rgp.state, []string{subDir}, buildContext, "./", false, false, "",
			llb.WithCustomNamef("[internal] COPY git context %s", target.String()))
	}

	earthfileTmpDir, err := ioutil.TempDir("/tmp", "earthly-git")
	if err != nil {
		return nil, errors.Wrap(err, "create temp dir for Earthfile")
	}
	gr.cleanCollection.Add(func() error {
		return os.RemoveAll(earthfileTmpDir)
	})
	buildFile, err := detectBuildFileInRef(ctx, target, rgp.gitMetaAndEarthfileRef, subDir)
	if err != nil {
		return nil, err
	}
	buildFileBytes, err := rgp.gitMetaAndEarthfileRef.ReadFile(ctx, gwclient.ReadRequest{
		Filename: buildFile,
	})
	if err != nil {
		return nil, errors.Wrap(err, "read build file")
	}
	localBuildFilePath := filepath.Join(earthfileTmpDir, path.Base(buildFile))
	err = ioutil.WriteFile(localBuildFilePath, buildFileBytes, 0700)
	if err != nil {
		return nil, errors.Wrapf(err, "write build file to tmp dir at %s", localBuildFilePath)
	}

	// TODO: Apply excludes / .earthignore.
	return &Data{
		BuildFilePath: localBuildFilePath,
		BuildContext:  buildContext,
		GitMetadata: &GitMetadata{
			BaseDir:   "",
			RelDir:    subDir,
			RemoteURL: gitURL,
			Hash:      rgp.hash,
			Branch:    rgp.branches,
			Tags:      rgp.tags,
		},
	}, nil
}

func (gr *gitResolver) resolveGitProject(ctx context.Context, gwClient gwclient.Client, target domain.Target) (rgp *resolvedGitProject, gitURL string, subDir string, finalErr error) {
	ref := target.Tag

	var err error
	var keyScan string
	gitURL, subDir, keyScan, err = gr.gitLookup.GetCloneURL(target.GitURL)
	if err != nil {
		return nil, "", "", errors.Wrap(err, "failed to get url for cloning")
	}

	// Check the cache first.
	cacheKey := fmt.Sprintf("%s#%s", gitURL, ref)
	data, found := gr.projectCache[cacheKey]
	if found {
		return data, gitURL, subDir, nil
	}
	// Not cached.

	// Copy all Earthfile, build.earth and Dockerfile files.
	gitOpts := []llb.GitOption{
		llb.WithCustomNamef("[internal] GIT CLONE %s", gitURL),
		llb.KeepGitDir(),
	}
	if keyScan != "" {
		gitOpts = append(gitOpts, llb.KnownSSHHosts(keyScan))
	}
	gitState := llb.Git(gitURL, ref, gitOpts...)
	copyOpts := []llb.RunOption{
		llb.Args([]string{
			"find",
			"-type", "f",
			"(", "-name", "build.earth", "-o", "-name", "Earthfile", "-o", "-name", "Dockerfile", ")",
			"-exec", "cp", "--parents", "{}", "/dest", ";",
		}),
		llb.Dir("/git-src"),
		llb.ReadonlyRootFS(),
		llb.AddMount("/git-src", gitState, llb.Readonly),
		llb.WithCustomNamef("[internal] COPY GIT CLONE %s Metadata", target.ProjectCanonical()),
	}
	opImg := llb.Image(
		defaultGitImage, llb.MarkImageInternal, llb.ResolveModePreferLocal,
		llb.Platform(llbutil.TargetPlatform))
	copyOp := opImg.Run(copyOpts...)
	earthfileState := copyOp.AddMount("/dest", llb.Scratch().Platform(llbutil.TargetPlatform))

	// Get git hash.
	gitHashOpts := []llb.RunOption{
		llb.Args([]string{
			"/bin/sh", "-c",
			"git rev-parse HEAD >/dest/git-hash ; " +
				"git rev-parse --abbrev-ref HEAD >/dest/git-branch  || touch /dest/git-branch ; " +
				"git describe --exact-match --tags >/dest/git-tags || touch /dest/git-tags",
		}),
		llb.Dir("/git-src"),
		llb.ReadonlyRootFS(),
		llb.AddMount("/git-src", gitState, llb.Readonly),
		llb.WithCustomNamef("[internal] GET GIT META %s", target.ProjectCanonical()),
	}
	gitHashOp := opImg.Run(gitHashOpts...)
	gitMetaAndEarthfileState := gitHashOp.AddMount("/dest", earthfileState)

	gitMetaAndEarthfileRef, err := llbutil.StateToRef(ctx, gwClient, gitMetaAndEarthfileState)
	if err != nil {
		return nil, "", "", errors.Wrap(err, "state to ref git meta")
	}
	gitHashBytes, err := gitMetaAndEarthfileRef.ReadFile(ctx, gwclient.ReadRequest{
		Filename: "git-hash",
	})
	if err != nil {
		return nil, "", "", errors.Wrap(err, "read git-hash")
	}
	gitBranchBytes, err := gitMetaAndEarthfileRef.ReadFile(ctx, gwclient.ReadRequest{
		Filename: "git-branch",
	})
	if err != nil {
		return nil, "", "", errors.Wrap(err, "read git-branch")
	}
	gitTagsBytes, err := gitMetaAndEarthfileRef.ReadFile(ctx, gwclient.ReadRequest{
		Filename: "git-tags",
	})
	if err != nil {
		return nil, "", "", errors.Wrap(err, "read git-tags")
	}

	gitHash := strings.SplitN(string(gitHashBytes), "\n", 2)[0]
	gitBranches := strings.SplitN(string(gitBranchBytes), "\n", 2)
	var gitBranches2 []string
	for _, gitBranch := range gitBranches {
		if gitBranch != "" {
			gitBranches2 = append(gitBranches2, gitBranch)
		}
	}
	gitTags := strings.SplitN(string(gitTagsBytes), "\n", 2)
	var gitTags2 []string
	for _, gitTag := range gitTags {
		if gitTag != "" && gitTag != "HEAD" {
			gitTags2 = append(gitTags2, gitTag)
		}
	}

	gitOpts = []llb.GitOption{
		llb.WithCustomNamef("[context %s] git context %s", gitURL, target.StringCanonical()),
	}
	if keyScan != "" {
		gitOpts = append(gitOpts, llb.KnownSSHHosts(keyScan))
	}

	// Add to cache.
	resolved := &resolvedGitProject{
		gitMetaAndEarthfileRef: gitMetaAndEarthfileRef,
		hash:                   gitHash,
		branches:               gitBranches2,
		tags:                   gitTags2,
		state: llb.Git(
			gitURL,
			gitHash,
			gitOpts...,
		),
	}
	gr.projectCache[cacheKey] = resolved
	cacheKey2 := fmt.Sprintf("%s#%s", gitURL, gitHash)
	gr.projectCache[cacheKey2] = resolved
	if len(gitBranches2) > 0 {
		cacheKey3 := fmt.Sprintf("%s#%s", gitURL, gitBranches2[0])
		gr.projectCache[cacheKey3] = resolved
	}
	if len(gitTags2) > 0 {
		cacheKey4 := fmt.Sprintf("%s#%s", gitURL, gitTags2[0])
		gr.projectCache[cacheKey4] = resolved
	}
	return resolved, gitURL, subDir, nil
}
