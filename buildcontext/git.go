package buildcontext

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/llbutil/llbfactory"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/stringutil"
	"github.com/earthly/earthly/util/syncutil/synccache"

	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

const (
	defaultGitImage = "alpine/git:v2.30.1"
)

type gitResolver struct {
	cleanCollection *cleanup.Collection

	projectCache   *synccache.SyncCache // "gitURL#gitRef" -> *resolvedGitProject
	buildFileCache *synccache.SyncCache // project ref -> local path
	gitLookup      *GitLookup
}

type resolvedGitProject struct {
	// hash is the git hash.
	hash string
	// branches is the git branches.
	branches []string
	// tags is the git tags
	tags []string
	// state is the state holding the git files.
	state pllb.State
}

func (gr *gitResolver) resolveEarthProject(ctx context.Context, gwClient gwclient.Client, ref domain.Reference) (*Data, error) {
	if !ref.IsRemote() {
		return nil, errors.Errorf("unexpected local reference %s", ref.String())
	}
	rgp, gitURL, subDir, err := gr.resolveGitProject(ctx, gwClient, ref)
	if err != nil {
		return nil, err
	}

	var buildContextFactory llbfactory.Factory
	if _, isTarget := ref.(domain.Target); isTarget {
		// Restrict the resulting build context to the right subdir.
		if subDir == "." {
			// Optimization.
			buildContextFactory = llbfactory.PreconstructedState(rgp.state)
		} else {
			buildContextFactory = llbfactory.PreconstructedState(llbutil.CopyOp(
				rgp.state, []string{subDir}, llbutil.ScratchWithPlatform(), "./", false, false, false, "root:root", false, false,
				llb.WithCustomNamef("[internal] COPY git context %s", ref.String())))
		}
	} else {
		// Commands don't come with a build context.
	}

	key := ref.ProjectCanonical()
	if strings.HasPrefix(ref.GetName(), DockerfileMetaTarget) {
		// Different key for dockerfiles to include the dockerfile name itself.
		key = ref.StringCanonical()
	}
	localBuildFilePathValue, err := gr.buildFileCache.Do(ctx, key, func(ctx context.Context, _ interface{}) (interface{}, error) {
		earthfileTmpDir, err := ioutil.TempDir(os.TempDir(), "earthly-git")
		if err != nil {
			return nil, errors.Wrap(err, "create temp dir for Earthfile")
		}
		gr.cleanCollection.Add(func() error {
			return os.RemoveAll(earthfileTmpDir)
		})
		gitState, err := llbutil.StateToRef(ctx, gwClient, rgp.state, nil, nil)
		if err != nil {
			return nil, errors.Wrap(err, "state to ref git meta")
		}
		buildFile, err := detectBuildFileInRef(ctx, ref, gitState, subDir)
		if err != nil {
			return nil, err
		}
		buildFileBytes, err := gitState.ReadFile(ctx, gwclient.ReadRequest{
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
		return localBuildFilePath, nil
	})
	if err != nil {
		return nil, err
	}
	localBuildFilePath := localBuildFilePathValue.(string)
	// TODO: Apply excludes / .earthignore.
	return &Data{
		BuildFilePath:       localBuildFilePath,
		BuildContextFactory: buildContextFactory,
		GitMetadata: &gitutil.GitMetadata{
			BaseDir:   "",
			RelDir:    subDir,
			RemoteURL: gitURL,
			Hash:      rgp.hash,
			Branch:    rgp.branches,
			Tags:      rgp.tags,
		},
	}, nil
}

func (gr *gitResolver) resolveGitProject(ctx context.Context, gwClient gwclient.Client, ref domain.Reference) (rgp *resolvedGitProject, gitURL string, subDir string, finalErr error) {
	gitRef := ref.GetTag()

	var err error
	var keyScan string
	gitURL, subDir, keyScan, err = gr.gitLookup.GetCloneURL(ref.GetGitURL())
	if err != nil {
		return nil, "", "", errors.Wrap(err, "failed to get url for cloning")
	}
	analytics.Count("gitResolver.resolveEarthProject", analytics.RepoHashFromCloneURL(gitURL))

	// Check the cache first.
	cacheKey := fmt.Sprintf("%s#%s", gitURL, gitRef)
	rgpValue, err := gr.projectCache.Do(ctx, cacheKey, func(ctx context.Context, k interface{}) (interface{}, error) {
		// Copy all Earthfile, build.earth and Dockerfile files.
		gitOpts := []llb.GitOption{
			llb.WithCustomNamef("[internal] GIT CLONE %s", stringutil.ScrubCredentials(gitURL)),
			llb.KeepGitDir(),
		}
		if keyScan != "" {
			gitOpts = append(gitOpts, llb.KnownSSHHosts(keyScan))
		}
		gitState := llb.Git(gitURL, gitRef, gitOpts...)
		opImg := pllb.Image(
			defaultGitImage, llb.MarkImageInternal, llb.ResolveModePreferLocal,
			llb.Platform(llbutil.DefaultPlatform()))

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
			llb.WithCustomNamef("[internal] GET GIT META %s", ref.ProjectCanonical()),
		}
		gitHashOp := opImg.Run(gitHashOpts...)
		gitMetaState := gitHashOp.AddMount("/dest", llbutil.ScratchWithPlatform())

		gitMetaRef, err := llbutil.StateToRef(ctx, gwClient, gitMetaState, nil, nil)
		if err != nil {
			return nil, errors.Wrap(err, "state to ref git meta")
		}
		gitHashBytes, err := gitMetaRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-hash",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-hash")
		}
		gitBranchBytes, err := gitMetaRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-branch",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-branch")
		}
		gitTagsBytes, err := gitMetaRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-tags",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-tags")
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
			llb.WithCustomNamef("[context %s] git context %s", gitURL, ref.StringCanonical()),
			llb.KeepGitDir(),
		}
		if keyScan != "" {
			gitOpts = append(gitOpts, llb.KnownSSHHosts(keyScan))
		}

		rgp := &resolvedGitProject{
			hash:     gitHash,
			branches: gitBranches2,
			tags:     gitTags2,
			state: pllb.Git(
				gitURL,
				gitHash,
				gitOpts...,
			),
		}
		go func() {
			// Add cache entries for the branch and for the tag (if any).
			if len(gitBranches2) > 0 {
				cacheKey3 := fmt.Sprintf("%s#%s", gitURL, gitBranches2[0])
				_ = gr.projectCache.Add(ctx, cacheKey3, rgp, nil)
			}
			if len(gitTags2) > 0 {
				cacheKey4 := fmt.Sprintf("%s#%s", gitURL, gitTags2[0])
				_ = gr.projectCache.Add(ctx, cacheKey4, rgp, nil)
			}
		}()
		return rgp, nil
	})
	if err != nil {
		return nil, "", "", err
	}
	rgp = rgpValue.(*resolvedGitProject)
	return rgp, gitURL, subDir, nil
}
