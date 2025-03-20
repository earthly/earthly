package gitutil

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/domain"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

var (
	// ErrNoGitBinary is an error returned when no git binary is found.
	ErrNoGitBinary = errors.New("No git binary found")
	// ErrNotAGitDir is an error returned when a given directory is not a git dir.
	ErrNotAGitDir = errors.New("Not a git directory")
	// ErrCouldNotDetectRemote is an error returned when git remote could not be detected or parsed.
	ErrCouldNotDetectRemote = errors.New("Could not auto-detect or parse Git remote URL")
	// ErrCouldNotDetectGitHash is an error returned when git hash could not be detected.
	ErrCouldNotDetectGitHash = errors.New("Could not auto-detect or parse Git hash")
	// ErrCouldNotDetectGitShortHash is an error returned when git short hash could not be detected.
	ErrCouldNotDetectGitShortHash = errors.New("Could not auto-detect or parse Git short hash")
	// ErrCouldNotDetectGitBranch is an error returned when git branch could not be detected.
	ErrCouldNotDetectGitBranch = errors.New("Could not auto-detect or parse Git branch")
	// ErrCouldNotDetectGitTags is an error returned when git tags could not be detected.
	ErrCouldNotDetectGitTags = errors.New("Could not auto-detect or parse Git tags")
	// ErrCouldNotDetectGitRefs is an error returned when git refs could not be detected.
	ErrCouldNotDetectGitRefs = errors.New("Could not auto-detect or parse Git refs")
)

// GitMetadata is a collection of git information about a certain directory.
type GitMetadata struct {
	BaseDir              string
	RelDir               string
	RemoteURL            string
	GitURL               string
	Hash                 string
	ShortHash            string
	BranchOverrideTagArg bool
	Branch               []string
	Tags                 []string
	CommitterTimestamp   string
	AuthorTimestamp      string
	AuthorEmail          string
	AuthorName           string
	CoAuthors            []string
	Refs                 []string
	Message              string
	FullMessage          string
}

// Metadata performs git metadata detection on the provided directory.
func Metadata(ctx context.Context, dir, gitBranchOverride string) (*GitMetadata, error) {
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
	var retErr error
	remoteURL, err := detectGitRemoteURL(ctx, dir)
	if err != nil {
		retErr = err
		// Keep going.
	}
	var gitURL string
	if remoteURL != "" {
		gitURL, err = ParseGitRemoteURL(remoteURL)
		if err != nil {
			return nil, err
		}
	}
	hash, err := detectGitHash(ctx, dir)
	if err != nil {
		retErr = err
		// Keep going.
	}
	shortHash, err := detectGitShortHash(ctx, dir)
	if err != nil {
		retErr = err
		// Keep going.
	}
	branch, err := detectGitBranch(ctx, dir, gitBranchOverride)
	if err != nil {
		retErr = err
		// Keep going.
	}
	tags, err := detectGitTags(ctx, dir)
	if err != nil {
		retErr = err
		// Keep going.
	}
	committerTimestamp, err := detectGitTimestamp(ctx, dir, committer)
	if err != nil {
		retErr = err
		// Keep going.
	}
	authorTimestamp, err := detectGitTimestamp(ctx, dir, author)
	if err != nil {
		retErr = err
		// Keep going.
	}
	authorEmail, err := detectGitAuthor(ctx, dir, "%ae")
	if err != nil {
		retErr = err
		// Keep going.
	}
	authorName, err := detectGitAuthor(ctx, dir, "%an")
	if err != nil {
		retErr = err
		// Keep going.
	}
	coAuthors, err := detectGitCoAuthors(ctx, dir)
	if err != nil {
		retErr = err
		// Keep going.
	}
	refs, err := detectGitRefs(ctx, dir)
	if err != nil {
		retErr = err
		// Keep going.
	}
	message, err := detectGitMessage(ctx, dir, oneline)
	if err != nil {
		retErr = err
		// Keep going.
	}
	fullMessage, err := detectGitMessage(ctx, dir, full)
	if err != nil {
		retErr = err
		// Keep going.
	}

	relDir, isRel, err := gitRelDir(baseDir, dir)
	if err != nil {
		return nil, errors.Wrapf(err, "get rel dir for %s when base git path is %s", dir, baseDir)
	}
	if !isRel {
		return nil, errors.New("unexpected non-relative path within git dir")
	}

	return &GitMetadata{
		BaseDir:              filepath.ToSlash(baseDir),
		RelDir:               filepath.ToSlash(relDir),
		RemoteURL:            remoteURL,
		GitURL:               gitURL,
		Hash:                 hash,
		ShortHash:            shortHash,
		BranchOverrideTagArg: gitBranchOverride != "",
		Branch:               branch,
		Tags:                 tags,
		CommitterTimestamp:   committerTimestamp,
		AuthorTimestamp:      authorTimestamp,
		AuthorEmail:          authorEmail,
		AuthorName:           authorName,
		CoAuthors:            coAuthors,
		Refs:                 refs,
		Message:              message,
		FullMessage:          fullMessage,
	}, retErr
}

// Clone returns a copy of the GitMetadata object.
func (gm *GitMetadata) Clone() *GitMetadata {
	return &GitMetadata{
		BaseDir:              gm.BaseDir,
		RelDir:               gm.RelDir,
		RemoteURL:            gm.RemoteURL,
		GitURL:               gm.GitURL,
		Hash:                 gm.Hash,
		ShortHash:            gm.ShortHash,
		BranchOverrideTagArg: gm.BranchOverrideTagArg,
		Branch:               gm.Branch,
		Tags:                 gm.Tags,
		CommitterTimestamp:   gm.CommitterTimestamp,
		AuthorTimestamp:      gm.AuthorTimestamp,
		AuthorEmail:          gm.AuthorEmail,
		AuthorName:           gm.AuthorName,
		CoAuthors:            gm.CoAuthors,
		Refs:                 gm.Refs,
		Message:              gm.Message,
		FullMessage:          gm.FullMessage,
	}
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

// ParseGitRemoteURL converts a gitURL like user@host.com:path/to.git or https://host.com/path/to.git to host.com/path/to
func ParseGitRemoteURL(gitURL string) (string, error) {
	s := gitURL

	// remove transport
	parts := strings.SplitN(gitURL, "://", 2)
	if len(parts) == 2 {
		s = parts[1]
	}

	// remove user
	parts = strings.SplitN(s, "@", 2)
	if len(parts) == 2 {
		s = parts[1]
	}

	s = strings.Replace(s, ":", "/", 1)
	s = strings.TrimSuffix(s, ".git")
	return s, nil
}

func detectGitRemoteURL(ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "config", "--get", "remote.origin.url")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(
			ErrCouldNotDetectRemote, "returned error %s: %s", err.Error(), string(out))
	}
	outStr := string(out)
	if outStr == "" {
		return "", errors.Wrapf(ErrCouldNotDetectRemote, "no remote origin url output")
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
	// cmd.Output produces a path with forward slash, but on Windows, we should preserve backslash paths.
	// E.g. This would convert `C:/my/path` to `C:\my\path`, but only when the platform is Windows.
	outStr = strings.Join(strings.Split(outStr, "/"), string(filepath.Separator))
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
		return "", errors.Wrapf(ErrCouldNotDetectGitHash, "returned error %s: %s", err.Error(), string(out))
	}
	outStr := string(out)
	if outStr == "" {
		return "", errors.Wrapf(ErrCouldNotDetectGitHash, "no remote origin url output")
	}
	return strings.SplitN(outStr, "\n", 2)[0], nil
}

func detectGitShortHash(ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--short=8", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(ErrCouldNotDetectGitShortHash, "returned error %s: %s", err.Error(), string(out))
	}
	outStr := string(out)
	if outStr == "" {
		return "", errors.Wrapf(ErrCouldNotDetectGitShortHash, "no remote origin url output")
	}
	return strings.SplitN(outStr, "\n", 2)[0], nil
}

func detectGitBranch(ctx context.Context, dir, gitBranchOverride string) ([]string, error) {
	if gitBranchOverride != "" {
		return []string{gitBranchOverride}, nil
	}
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrapf(ErrCouldNotDetectGitBranch, "returned error %s: %s", err.Error(), string(out))
	}
	outStr := string(out)
	if outStr != "" {
		return strings.Split(outStr, "\n"), nil
	}
	return nil, nil
}

func detectGitTags(ctx context.Context, dir string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "describe", "--exact-match", "--tags")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrapf(ErrCouldNotDetectGitTags, "returned error %s: %s", err.Error(), string(out))
	}
	outStr := string(out)
	if outStr != "" {
		return strings.Split(outStr, "\n"), nil
	}
	return nil, nil
}

func detectGitRefs(ctx context.Context, dir string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "for-each-ref", "--points-at", "HEAD", "--format", "'%(refname:lstrip=-1)'")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrapf(ErrCouldNotDetectGitRefs, "returned error %s: %s", err.Error(), string(out))
	}
	outStr := string(out)
	if outStr != "" {
		refs := []string{}
		for _, ref := range strings.Split(outStr, "\n") {
			ref = strings.Trim(ref, "'\"")
			if ref != "" && ref != "HEAD" && !slices.Contains(refs, ref) {
				refs = append(refs, ref)
			}
		}
		return refs, nil
	}
	return nil, nil
}

type gitTimestampType int

const (
	author gitTimestampType = iota
	committer
)

type gitCommitMessageType int

const (
	oneline gitCommitMessageType = iota
	full
)

func detectGitTimestamp(ctx context.Context, dir string, tsType gitTimestampType) (string, error) {
	var format string
	switch tsType {
	case author:
		format = "%at"
	case committer:
		format = "%ct"
	}
	cmd := exec.CommandContext(ctx, "git", "log", "-1", "--format="+format)
	cmd.Dir = dir
	cmd.Stderr = nil // force capture of stderr on errors
	out, err := cmd.Output()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok && strings.Contains(string(exitError.Stderr), "does not have any commits yet") {
			return "", nil
		}
		return "", errors.Wrap(err, "detect git timestamp")
	}
	outStr := string(out)
	if outStr == "" {
		return "", nil
	}
	return strings.SplitN(outStr, "\n", 2)[0], nil
}

func detectGitAuthor(ctx context.Context, dir string, format string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "log", "-1", fmt.Sprintf("--format=%s", format))
	cmd.Dir = dir
	cmd.Stderr = nil // force capture of stderr on errors
	out, err := cmd.Output()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok && strings.Contains(string(exitError.Stderr), "does not have any commits yet") {
			return "", nil
		}
		return "", errors.Wrap(err, "detect git author")
	}
	outStr := string(out)
	if outStr == "" {
		return "", nil
	}
	return strings.SplitN(outStr, "\n", 2)[0], nil
}

func detectGitMessage(ctx context.Context, dir string, formatType gitCommitMessageType) (string, error) {
	var format string
	switch formatType {
	case oneline:
		format = "%s"
	case full:
		format = "%B"
	}
	cmd := exec.CommandContext(ctx, "git", "log", "-1", fmt.Sprintf("--format=%s", format))
	cmd.Dir = dir
	cmd.Stderr = nil // force capture of stderr on errors
	out, err := cmd.Output()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok && strings.Contains(string(exitError.Stderr), "does not have any commits yet") {
			return "", nil
		}
		return "", errors.Wrap(err, "detect git comment")
	}
	outStr := string(out)
	if outStr == "" {
		return "", nil
	}
	if formatType == oneline {
		return strings.SplitN(outStr, "\n", 2)[0], nil
	}
	return outStr, nil
}

// ConfigEmail returns the user's currently configured (global) email address
func ConfigEmail(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "config", "--get", "user.email")
	cmd.Stderr = nil // force capture of stderr on errors
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "detect git global email")
	}
	return strings.TrimSpace(string(out)), nil
}

func detectGitCoAuthors(ctx context.Context, dir string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "log", "-1", "--format=%b")
	cmd.Dir = dir
	cmd.Stderr = nil // force capture of stderr on errors
	out, err := cmd.Output()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok && strings.Contains(string(exitError.Stderr), "does not have any commits yet") {
			return nil, nil
		}
		if out != nil && strings.Contains(string(out), "does not have any commits yet") {
			return nil, nil
		}
		return nil, errors.Wrap(err, "detect git co-authors")
	}
	return ParseCoAuthorsFromBody(string(out)), nil
}

// ParseCoAuthorsFromBody returns a list of coauthor emails from a git body
func ParseCoAuthorsFromBody(body string) []string {
	coAuthors := []string{}
	coAuthorsSeen := map[string]struct{}{}
	for _, s := range strings.Split(body, "\n") {
		s = strings.TrimSpace(s)
		splits := strings.Split(s, " ")
		n := len(splits)
		if n > 2 {
			if splits[0] == "Co-authored-by:" {
				email := splits[n-1]
				n = len(email)
				if n > 2 {
					if email[0] == '<' && email[n-1] == '>' {
						email = email[1:(n - 1)]
						_, seen := coAuthorsSeen[email]
						if !seen {
							coAuthors = append(coAuthors, email)
							coAuthorsSeen[email] = struct{}{}
						}
					}
				}
			}
		}
	}
	return coAuthors
}

// gitRelDir returns the relative path from git root (where .git directory locates in the project)
// This function validates the input data (basePath, path) as well.
func gitRelDir(basePath string, path string) (string, bool, error) {
	if !filepath.IsAbs(basePath) {
		return "", false, errors.Errorf("git base path %s is not absolute", basePath)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", false, errors.Wrapf(err, "get abs path for %s", path)
	}
	absPath2, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return "", false, errors.Wrapf(err, "eval symlinks for %s", absPath)
	}

	basePathParts := strings.Split(basePath, string(filepath.Separator))
	pathParts := strings.Split(absPath2, string(filepath.Separator))

	// `basePath` must be the part of `path`
	// So it's length split by filepath.Separator must be shorter than `part`.
	if len(pathParts) < len(basePathParts) {
		return "", false, nil
	}

	a, err := os.Stat(basePath)
	if err != nil {
		return "", false, errors.Wrapf(err, "stat for %s", basePath)
	}

	// Checking VolumeName determines if we have a fully-qualified Windows path like `C:\my\dir`.
	if filepath.VolumeName(basePath) == "" {
		// `pathParts` here has lost the root filepath.Separator since it was built by strings.Split beforehand.
		// So putting heading separator is required to make it absolute again.
		pathParts[0] = string(filepath.Separator) + pathParts[0]
	} else {
		// In Window's style absolute paths, we need append the file-separator after the first element.
		// e.g. We want: `C:\some\dir`, not `\C:some\dir`
		pathParts[0] = pathParts[0] + string(filepath.Separator)
	}
	b, err := os.Stat(filepath.Join(pathParts[:len(basePathParts)]...))
	if err != nil {
		return "", false, errors.Wrapf(err, "stat for %v", pathParts)
	}
	// Here checks if `path` is included in `basePath` in filesystem agnostic way.
	// Case-sensitivity difference (like HFS+ in OSX) is also covered by os.SameFile.
	if !os.SameFile(a, b) {
		return "", false, nil
	}

	// Now we are sure that inclusion of `basePath` in `path` is OK.
	// Finally, here extracts the relative path from `basePath` to return.
	relPath := filepath.Join(pathParts[len(basePathParts):]...)
	if relPath == "" {
		return ".", true, nil
	}
	return relPath, true, nil
}

// ReferenceWithGitMeta applies git metadata to the target naming.
func ReferenceWithGitMeta(ref domain.Reference, gitMeta *GitMetadata) domain.Reference {
	if gitMeta == nil || gitMeta.GitURL == "" {
		return ref
	}
	gitURL := gitMeta.GitURL
	if gitMeta.RelDir != "" {
		gitURL = path.Join(gitURL, gitMeta.RelDir)
	}
	tag := ref.GetTag()
	localPath := ref.GetLocalPath()
	name := ref.GetName()
	importRef := ref.GetImportRef()

	if tag == "" {
		if len(gitMeta.Tags) > 0 {
			tag = gitMeta.Tags[0]
		} else if len(gitMeta.Branch) > 0 {
			tag = gitMeta.Branch[0]
		} else {
			tag = gitMeta.Hash
		}
	}

	switch ref.(type) {
	case domain.Target:
		return domain.Target{
			GitURL:    gitURL,
			Tag:       tag,
			LocalPath: localPath,
			ImportRef: importRef,
			Target:    name,
		}
	case domain.Command:
		return domain.Command{
			GitURL:    gitURL,
			Tag:       tag,
			LocalPath: localPath,
			ImportRef: importRef,
			Command:   name,
		}
	default:
		panic("not supported for this type")
	}
}
