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
	"github.com/earthly/earthly/gitutil"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/stringutil"
	"github.com/earthly/earthly/syncutil/synccache"

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
	// gitMetaAndEarthfileState is the state containing the git metadata and build files.
	gitMetaAndEarthfileState llb.State
	// hash is the git hash.
	hash string
	// branches is the git branches.
	branches []string
	// tags is the git tags
	tags []string
	// state is the state holding the git files.
	state llb.State
}

func (gr *gitResolver) resolveEarthProject(ctx context.Context, gwClient gwclient.Client, ref domain.Reference) (*Data, error) {
	if !ref.IsRemote() {
		return nil, fmt.Errorf("unexpected local reference %s", ref.String())
	}
	rgp, gitURL, subDir, err := gr.resolveGitProject(ctx, gwClient, ref)
	if err != nil {
		return nil, err
	}

	var buildContext llb.State
	if _, isTarget := ref.(domain.Target); isTarget {
		// Restrict the resulting build context to the right subdir.
		if subDir == "." {
			// Optimization.
			buildContext = rgp.state
		} else {
			buildContext = llbutil.ScratchWithPlatform()
			buildContext = llbutil.CopyOp(
				rgp.state, []string{subDir}, buildContext, "./", false, false, false, "root:root", false,
				llb.WithCustomNamef("[internal] COPY git context %s", ref.String()))
		}
	} else {
		// Commands don't come with a build context.
	}

	key := ref.ProjectCanonical()
	if ref.GetName() == DockerfileMetaTarget {
		// Different key for dockerfiles.
		key = key + "@" + DockerfileMetaTarget
	}
	localBuildFilePathValue, err := gr.buildFileCache.Do(key, func(_ interface{}) (interface{}, error) {
		earthfileTmpDir, err := ioutil.TempDir(os.TempDir(), "earthly-git")
		if err != nil {
			return nil, errors.Wrap(err, "create temp dir for Earthfile")
		}
		gr.cleanCollection.Add(func() error {
			return os.RemoveAll(earthfileTmpDir)
		})
		gitMetaAndEarthfileRef, err := llbutil.StateToRef(ctx, gwClient, rgp.gitMetaAndEarthfileState, nil, nil)
		if err != nil {
			return nil, errors.Wrap(err, "state to ref git meta")
		}
		buildFile, err := detectBuildFileInRef(ctx, ref, gitMetaAndEarthfileRef, subDir)
		if err != nil {
			return nil, err
		}
		buildFileBytes, err := gitMetaAndEarthfileRef.ReadFile(ctx, gwclient.ReadRequest{
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
		BuildFilePath: localBuildFilePath,
		BuildContext:  buildContext,
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

	// Check the cache first.
	cacheKey := fmt.Sprintf("%s#%s", gitURL, gitRef)
	rgpValue, err := gr.projectCache.Do(cacheKey, func(k interface{}) (interface{}, error) {
		// Copy all Earthfile, build.earth and Dockerfile files.
		gitOpts := []llb.GitOption{
			llb.WithCustomNamef("[internal] GIT CLONE %s", stringutil.ScrubCredentials(gitURL)),
			llb.KeepGitDir(),
		}
		if keyScan != "" {
			gitOpts = append(gitOpts, llb.KnownSSHHosts(keyScan))
		}
		gitState := llb.Git(gitURL, gitRef, gitOpts...)
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
			llb.WithCustomNamef("[internal] COPY GIT CLONE %s Metadata", ref.ProjectCanonical()),
		}
		opImg := llb.Image(
			defaultGitImage, llb.MarkImageInternal, llb.ResolveModePreferLocal,
			llb.Platform(llbutil.DefaultPlatform()))
		copyOp := opImg.Run(copyOpts...)
		earthfileState := copyOp.AddMount("/dest", llbutil.ScratchWithPlatform())

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
		gitMetaAndEarthfileState := gitHashOp.AddMount("/dest", earthfileState)

		gitMetaAndEarthfileRef, err := llbutil.StateToRef(ctx, gwClient, gitMetaAndEarthfileState, nil, nil)
		if err != nil {
			return nil, errors.Wrap(err, "state to ref git meta")
		}
		gitHashBytes, err := gitMetaAndEarthfileRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-hash",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-hash")
		}
		gitBranchBytes, err := gitMetaAndEarthfileRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-branch",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-branch")
		}
		gitTagsBytes, err := gitMetaAndEarthfileRef.ReadFile(ctx, gwclient.ReadRequest{
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

		return &resolvedGitProject{
			gitMetaAndEarthfileState: gitMetaAndEarthfileState,
			hash:                     gitHash,
			branches:                 gitBranches2,
			tags:                     gitTags2,
			state: llb.Git(
				gitURL,
				gitHash,
				gitOpts...,
			),
		}, nil
	})
	if err != nil {
		return nil, "", "", err
	}
	rgp = rgpValue.(*resolvedGitProject)
	return rgp, gitURL, subDir, nil
}
