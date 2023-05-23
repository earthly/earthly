package buildcontext

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/llbutil/llbfactory"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/stringutil"
	"github.com/earthly/earthly/util/syncutil/synccache"
	"github.com/earthly/earthly/util/vertexmeta"

	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

const (
	defaultGitImage = "alpine/git:v2.30.1"
)

type gitResolver struct {
	cleanCollection   *cleanup.Collection
	gitBranchOverride string
	lfsInclude        string
	logLevel          llb.GitLogLevel
	gitImage          string
	projectCache      *synccache.SyncCache // "gitURL#gitRef" -> *resolvedGitProject
	buildFileCache    *synccache.SyncCache // project ref -> local path
	gitLookup         *GitLookup
	console           conslogging.ConsoleLogger
}

type resolvedGitProject struct {
	// hash is the git hash.
	hash string
	// shortHash is the short git hash.
	shortHash string
	// branches is the git branches.
	branches []string
	// tags is the git tags.
	tags []string
	// committerTs is the git committer timestamp.
	committerTs string
	// authorTs is the git author timestamp.
	authorTs  string
	author    string
	coAuthors []string
	// state is the state holding the git files.
	state pllb.State
}

func (gr *gitResolver) resolveEarthProject(ctx context.Context, gwClient gwclient.Client, platr *platutil.Resolver, ref domain.Reference, featureFlagOverrides string) (*Data, error) {
	if !ref.IsRemote() {
		return nil, errors.Errorf("unexpected local reference %s", ref.String())
	}
	rgp, gitURL, subDir, err := gr.resolveGitProject(ctx, gwClient, platr, ref)
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
			vm := &vertexmeta.VertexMeta{
				TargetName: ref.String(),
				Internal:   true,
			}
			copyState, err := llbutil.CopyOp(ctx,
				rgp.state, []string{subDir}, platr.Scratch(), "./", false, false, false, "root:root", nil, false, false, false,
				llb.WithCustomNamef("%sCOPY git context %s", vm.ToVertexPrefix(), ref.String()))
			if err != nil {
				return nil, errors.Wrap(err, "copyOp failed in resolveEarthProject")
			}
			buildContextFactory = llbfactory.PreconstructedState(copyState)
		}
	}
	// Else not needed: Commands don't come with a build context.

	key := ref.ProjectCanonical()
	isDockerfile := strings.HasPrefix(ref.GetName(), DockerfileMetaTarget)
	if isDockerfile {
		// Different key for dockerfiles to include the dockerfile name itself.
		key = ref.StringCanonical()
	}
	localBuildFileValue, err := gr.buildFileCache.Do(ctx, key, func(ctx context.Context, _ interface{}) (interface{}, error) {
		earthfileTmpDir, err := os.MkdirTemp(os.TempDir(), "earthly-git")
		if err != nil {
			return nil, errors.Wrap(err, "create temp dir for Earthfile")
		}
		gr.cleanCollection.Add(func() error {
			return os.RemoveAll(earthfileTmpDir)
		})
		gitState, err := llbutil.StateToRef(
			ctx, gwClient, rgp.state, false,
			platr.SubResolver(platutil.NativePlatform), nil)
		if err != nil {
			return nil, errors.Wrap(err, "state to ref git meta")
		}
		bf, err := detectBuildFileInRef(ctx, ref, gitState, subDir)
		if err != nil {
			return nil, err
		}
		bfBytes, err := gitState.ReadFile(ctx, gwclient.ReadRequest{
			Filename: bf,
		})
		if err != nil {
			return nil, errors.Wrap(err, "read build file")
		}
		localBuildFilePath := filepath.Join(earthfileTmpDir, path.Base(bf))
		err = os.WriteFile(localBuildFilePath, bfBytes, 0700)
		if err != nil {
			return nil, errors.Wrapf(err, "write build file to tmp dir at %s", localBuildFilePath)
		}
		var ftrs *features.Features
		if isDockerfile {
			ftrs = new(features.Features)
		} else {
			ftrs, err = parseFeatures(localBuildFilePath, featureFlagOverrides, ref.ProjectCanonical(), gr.console)
			if err != nil {
				return nil, err
			}
		}
		return &buildFile{
			path: localBuildFilePath,
			ftrs: ftrs,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	localBuildFile := localBuildFileValue.(*buildFile)

	// TODO: Apply excludes / .earthignore.
	return &Data{
		BuildFilePath:       localBuildFile.path,
		BuildContextFactory: buildContextFactory,
		GitMetadata: &gitutil.GitMetadata{
			BaseDir:              "",
			RelDir:               subDir,
			RemoteURL:            gitURL,
			Hash:                 rgp.hash,
			ShortHash:            rgp.shortHash,
			BranchOverrideTagArg: gr.gitBranchOverride != "",
			Branch:               rgp.branches,
			Tags:                 rgp.tags,
			CommitterTimestamp:   rgp.committerTs,
			AuthorTimestamp:      rgp.authorTs,
			Author:               rgp.author,
			CoAuthors:            rgp.coAuthors,
		},
		Features: localBuildFile.ftrs,
	}, nil
}

func (gr *gitResolver) resolveGitProject(ctx context.Context, gwClient gwclient.Client, platr *platutil.Resolver, ref domain.Reference) (rgp *resolvedGitProject, gitURL string, subDir string, finalErr error) {
	gitRef := ref.GetTag()

	var err error
	var keyScans []string
	gitURL, subDir, keyScans, err = gr.gitLookup.GetCloneURL(ref.GetGitURL())
	if err != nil {
		return nil, "", "", errors.Wrap(err, "failed to get url for cloning")
	}
	analytics.Count("gitResolver.resolveEarthProject", "")

	// Check the cache first.
	cacheKey := fmt.Sprintf("%s#%s", gitURL, gitRef)
	rgpValue, err := gr.projectCache.Do(ctx, cacheKey, func(ctx context.Context, k interface{}) (interface{}, error) {
		// Copy all Earthfile, build.earth and Dockerfile files.
		vm := &vertexmeta.VertexMeta{
			TargetName: cacheKey,
			Internal:   true,
		}
		gitOpts := []llb.GitOption{
			llb.WithCustomNamef("%sGIT CLONE %s", vm.ToVertexPrefix(), stringutil.ScrubCredentials(gitURL)),
			llb.KeepGitDir(),
			llb.LogLevel(gr.logLevel),
		}
		if len(keyScans) > 0 {
			gitOpts = append(gitOpts, llb.KnownSSHHosts(strings.Join(keyScans, "\n")))
		}
		if gr.lfsInclude != "" {
			// TODO this should eventually be infered by the contents of a COPY command, which means the call to resolveGitProject will need to be lazy-evaluated
			// However this makes it really difficult for an Earthfile which first has an ARG EARTHLY_GIT_HASH, then a RUN, then a COPY
			gitOpts = append(gitOpts, llb.LFSInclude(gr.lfsInclude))
		}

		gitState := llb.Git(gitURL, gitRef, gitOpts...)
		gitImage := gr.gitImage
		if gitImage == "" {
			gitImage = defaultGitImage
		}
		opImg := pllb.Image(
			gitImage, llb.MarkImageInternal, llb.ResolveModePreferLocal,
			llb.Platform(platr.LLBNative()))

		// Get git hash.
		gitHashOpts := []llb.RunOption{
			llb.Args([]string{
				"/bin/sh", "-c",
				"git rev-parse HEAD >/dest/git-hash ; " +
					"git rev-parse --short=8 HEAD >/dest/git-short-hash ; " +
					"git rev-parse --abbrev-ref HEAD >/dest/git-branch  || touch /dest/git-branch ; " +
					"ls .git/refs/heads/ | head -n 1 >/dest/git-default-branch  || touch /dest/git-default-branch ; " +
					"git describe --exact-match --tags >/dest/git-tags || touch /dest/git-tags ; " +
					"git log -1 --format=%ct >/dest/git-committer-ts || touch /dest/git-committer-ts ; " +
					"git log -1 --format=%at >/dest/git-author-ts || touch /dest/git-author-ts ; " +
					"git log -1 --format=%ae >/dest/git-author || touch /dest/git-author ; " +
					"git log -1 --format=%b >/dest/git-body || touch /dest/git-body ; " +
					"",
			}),
			llb.Dir("/git-src"),
			llb.ReadonlyRootFS(),
			llb.AddMount("/git-src", gitState, llb.Readonly),
			llb.WithCustomNamef("%sGET GIT META %s", vm.ToVertexPrefix(), ref.ProjectCanonical()),
		}
		gitHashOp := opImg.Run(gitHashOpts...)
		gitMetaState := gitHashOp.AddMount("/dest", platr.Scratch())

		noCache := false // TODO figure out if we want to propagate --no-cache here
		gitMetaRef, err := llbutil.StateToRef(
			ctx, gwClient, gitMetaState, noCache,
			platr.SubResolver(platutil.NativePlatform), nil)
		if err != nil {
			return nil, errors.Wrap(err, "state to ref git meta")
		}
		gitHashBytes, err := gitMetaRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-hash",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-hash")
		}
		gitShortHashBytes, err := gitMetaRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-short-hash",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-short-hash")
		}
		gitBranch, err := gr.readGitBranch(ctx, gitMetaRef)
		if err != nil {
			return nil, errors.Wrap(err, "read git-branch")
		}
		gitDefaultBranchBytes, err := gitMetaRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-default-branch",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-default-branch")
		}
		gitTagsBytes, err := gitMetaRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-tags",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-tags")
		}
		gitCommitterTsBytes, err := gitMetaRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-committer-ts",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-committer-ts")
		}
		gitAuthorTsBytes, err := gitMetaRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-author-ts",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-author-ts")
		}
		gitAuthorBytes, err := gitMetaRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-author",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-author")
		}
		gitBodyBytes, err := gitMetaRef.ReadFile(ctx, gwclient.ReadRequest{
			Filename: "git-body",
		})
		if err != nil {
			return nil, errors.Wrap(err, "read git-body")
		}

		gitHash := strings.SplitN(string(gitHashBytes), "\n", 2)[0]
		gitShortHash := strings.SplitN(string(gitShortHashBytes), "\n", 2)[0]
		gitBranches := strings.SplitN(gitBranch, "\n", 2)
		gitAuthor := strings.SplitN(string(gitAuthorBytes), "\n", 2)[0]
		gitCoAuthors := gitutil.ParseCoAuthorsFromBody(string(gitBodyBytes))
		var gitBranches2 []string
		for _, gitBranch := range gitBranches {
			if gitBranch != "" && gitBranch != "HEAD" {
				gitBranches2 = append(gitBranches2, gitBranch)
			}
		}
		if len(gitBranches2) == 0 {
			// fallback case for when git rev-parse --abbrev-ref fails
			if gitRef != "" {
				// use the reference name (if given); but only if it is not the git sha
				if !strings.HasPrefix(gitRef, gitShortHash) {
					gitBranches2 = []string{gitRef}
				}
			} else {
				gitBranches2 = []string{strings.SplitN(string(gitDefaultBranchBytes), "\n", 2)[0]}
			}

		}
		gitTags := strings.SplitN(string(gitTagsBytes), "\n", 2)
		var gitTags2 []string
		for _, gitTag := range gitTags {
			if gitTag != "" && gitTag != "HEAD" {
				gitTags2 = append(gitTags2, gitTag)
			}
		}
		gitCommiterTs := strings.SplitN(string(gitCommitterTsBytes), "\n", 2)[0]
		gitAuthorTs := strings.SplitN(string(gitAuthorTsBytes), "\n", 2)[0]

		gitOpts = []llb.GitOption{
			llb.WithCustomNamef("[context %s] git context %s", stringutil.ScrubCredentials(gitURL), ref.StringCanonical()),
			llb.KeepGitDir(),
		}
		if len(keyScans) > 0 {
			gitOpts = append(gitOpts, llb.KnownSSHHosts(strings.Join(keyScans, "\n")))
		}
		if gr.lfsInclude != "" {
			gitOpts = append(gitOpts, llb.LFSInclude(gr.lfsInclude))
		}

		rgp := &resolvedGitProject{
			hash:        gitHash,
			shortHash:   gitShortHash,
			branches:    gitBranches2,
			tags:        gitTags2,
			committerTs: gitCommiterTs,
			authorTs:    gitAuthorTs,
			author:      gitAuthor,
			coAuthors:   gitCoAuthors,
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

func (gr *gitResolver) readGitBranch(ctx context.Context, gitMetaRef gwclient.Reference) (string, error) {
	if gr.gitBranchOverride != "" {
		return gr.gitBranchOverride, nil
	}
	gitBranchBytes, err := gitMetaRef.ReadFile(ctx, gwclient.ReadRequest{
		Filename: "git-branch",
	})
	if err != nil {
		return "", err
	}
	return string(gitBranchBytes), nil
}
