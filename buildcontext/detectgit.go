package buildcontext

import (
	"context"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/vladaionescu/earthly/domain"
)

var (
	// ErrNoGitBinary is an error returned when no git binary is found.
	ErrNoGitBinary = fmt.Errorf("No git binary found")
	// ErrNotAGitDir is an error returned when a given directory is not a git dir.
	ErrNotAGitDir = fmt.Errorf("Not a git directory")
)

// GitMetadata is a collection of git information about a certain directory.
type GitMetadata struct {
	BaseDir    string
	RelDir     string
	RemoteURL  string
	GitVendor  string
	GitProject string
	Hash       string
	Branch     []string
	Tags       []string
}

// Metadata performs git metadata detection on the provided directory.
func Metadata(ctx context.Context, dir string) (*GitMetadata, error) {
	err := detectGitBinary(ctx)
	if err != nil {
		return nil, err
	}
	err = detectIsGitDir(ctx, dir)
	if err != nil {
		return nil, err
	}
	baseDir, err := detectGitBaseDir(ctx, dir)
	if err != nil {
		return nil, err
	}
	remoteURL, err := detectGitRemoteURL(ctx, dir)
	if err != nil {
		return nil, err
	}
	var vendor, project string
	if remoteURL != "" {
		vendor, project, err = parseGitRemoteURL(remoteURL)
		if err != nil {
			return nil, err
		}
	}
	hash, err := detectGitHash(ctx, dir)
	if err != nil {
		return nil, err
	}
	branch, err := detectGitBranch(ctx, dir)
	if err != nil {
		return nil, err
	}
	tags, err := detectGitTags(ctx, dir)
	if err != nil {
		// Most likely no tags. Keep going.
		tags = nil
	}

	relDir, isRel, err := gitRelDir(baseDir, dir)
	if err != nil {
		return nil, errors.Wrapf(err, "get rel dir for %s when base git path is %s", dir, baseDir)
	}
	if !isRel {
		return nil, errors.New("Unexpected non-relative path within git dir")
	}

	return &GitMetadata{
		BaseDir:    filepath.ToSlash(baseDir),
		RelDir:     filepath.ToSlash(relDir),
		GitVendor:  vendor,
		GitProject: project,
		Hash:       hash,
		Branch:     branch,
		Tags:       tags,
	}, nil
}

// Clone returns a copy of the GitMetadata object.
func (gm *GitMetadata) Clone() *GitMetadata {
	return &GitMetadata{
		BaseDir:    gm.BaseDir,
		RelDir:     gm.RelDir,
		GitVendor:  gm.GitVendor,
		GitProject: gm.GitProject,
		Hash:       gm.Hash,
		Branch:     gm.Branch,
		Tags:       gm.Tags,
	}
}

func detectGitBinary(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", "which git")
	_, err := cmd.Output()
	if err != nil {
		return ErrNoGitBinary
	}
	return nil
}

func detectIsGitDir(ctx context.Context, dir string) error {
	cmd := exec.CommandContext(ctx, "git", "status")
	cmd.Dir = dir
	_, err := cmd.Output()
	if err != nil {
		return ErrNotAGitDir
	}
	return nil
}

func parseGitRemoteURL(url string) (string, string, error) {
	if strings.HasPrefix(url, "git@github.com:") {
		ret := strings.TrimPrefix(url, "git@github.com:")
		ret = strings.TrimSuffix(ret, ".git")
		return "github.com", ret, nil
	} else if strings.HasPrefix(url, "https://") {
		ret := strings.TrimPrefix(url, "https://")
		parts := strings.SplitN(ret, "/", 2)
		if len(parts) < 2 {
			return "", "", fmt.Errorf("Unable to parse git URL %s", url)
		}
		return parts[0], strings.TrimSuffix(parts[1], ".git"), nil
	} else {
		return "", "", fmt.Errorf("Unable to parse git URL %s", url)
	}
}

func detectGitRemoteURL(ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "config", "--get", "remote.origin.url")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "detect git directory")
	}
	outStr := string(out)
	if outStr == "" {
		return "", errors.New("No output returned for git remote url")
	}
	return strings.SplitN(outStr, "\n", 2)[0], nil
}

func detectGitBaseDir(ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "detect git directory")
	}
	outStr := string(out)
	if outStr == "" {
		return "", errors.New("No output returned for git base dir")
	}
	return strings.SplitN(outStr, "\n", 2)[0], nil
}

func detectGitHash(ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "detect git hash")
	}
	outStr := string(out)
	if outStr == "" {
		return "", errors.New("No output returned for git hash")
	}
	return strings.SplitN(outStr, "\n", 2)[0], nil
}

func detectGitBranch(ctx context.Context, dir string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "detect git current branch")
	}
	outStr := string(out)
	if outStr != "" {
		return strings.Split(outStr, "\n"), nil
	}
	return nil, nil
}

func detectGitTags(ctx context.Context, dir string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "describe", "--tags")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "detect git current tags")
	}
	outStr := string(out)
	if outStr != "" {
		return strings.Split(outStr, "\n"), nil
	}
	return nil, nil
}

func gitRelDir(basePath string, path string) (string, bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", false, errors.Wrapf(err, "get abs path for %s", path)
	}
	if !filepath.IsAbs(basePath) {
		return "", false, fmt.Errorf("Git base path %s is not absolute", basePath)
	}
	basePathSlash := filepath.ToSlash(basePath)
	pathSlash := filepath.ToSlash(absPath)
	basePathParts := strings.Split(basePathSlash, "/")
	pathParts := strings.Split(pathSlash, "/")
	if len(pathParts) < len(basePathParts) {
		return "", false, nil
	}
	for index := range basePathParts {
		if basePathParts[index] != pathParts[index] {
			return "", false, nil
		}
	}
	relPath := strings.Join(pathParts[len(basePathParts):], "/")
	if relPath == "" {
		return ".", true, nil
	}
	return filepath.FromSlash(relPath), true, nil
}

// TargetWithGitMeta applies git metadata to the target naming.
func TargetWithGitMeta(target domain.Target, gitMeta *GitMetadata) domain.Target {
	if gitMeta == nil {
		return target
	}
	targetRet := target
	targetRet.Registry = gitMeta.GitVendor
	targetRet.ProjectPath = path.Join(gitMeta.GitProject, gitMeta.RelDir)
	if len(gitMeta.Branch) > 0 {
		targetRet.Tag = gitMeta.Branch[0]
	}
	return targetRet
}
