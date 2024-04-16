package conslogging

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/earthly/earthly/cleanup"
	"github.com/pkg/errors"
)

const fullLog = "_full"

type targetLogger struct {
	writer  *strings.Builder
	prefix  string
	result  string
	status  string
	started time.Time
}

// BundleBuilder builds log bundles for local storage or upload to a logging server
type BundleBuilder struct {
	entrypoint string
	started    time.Time
	cleanup    *cleanup.Collection

	mu            sync.Mutex
	logsForTarget map[string]*targetLogger
}

// Write implements io.Writer as a passthrough to the underlying strings.Builder for convenience.
func (tl *targetLogger) Write(p []byte) (n int, err error) {
	return tl.writer.Write(p)
}

// NewBundleBuilder makes a new BundleBuilder, that will write logs to the targeted root directory,
// and specify the entrypoint in the resulting manifest.
func NewBundleBuilder(entrypoint string, cleanup *cleanup.Collection) *BundleBuilder {
	return &BundleBuilder{
		entrypoint:    entrypoint,
		logsForTarget: map[string]*targetLogger{},
		started:       time.Now(),
		cleanup:       cleanup,
	}
}

// PrefixResult sets the prefix(aka target) result as it should appear in the manifest for that specific target.
func (bb *BundleBuilder) PrefixResult(prefix, result string) {
	bb.mu.Lock()
	defer bb.mu.Unlock()
	if builder, ok := bb.logsForTarget[prefix]; ok {
		builder.result = result
	}
}

// PrefixStatus sets the prefix(aka target) result as it should appear in the manifest for that specific target.
func (bb *BundleBuilder) PrefixStatus(prefix, status string) {
	bb.mu.Lock()
	defer bb.mu.Unlock()
	if builder, ok := bb.logsForTarget[prefix]; ok {
		builder.status = status
	}
}

// PrefixWriter gets an io.Writer for a given prefix(aka target). If its a prefix we have not seen before,
// then generate a new writer to accommodate it.
func (bb *BundleBuilder) PrefixWriter(prefix string) io.Writer {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	if builder, ok := bb.logsForTarget[prefix]; ok {
		return builder
	}

	writer := &targetLogger{
		writer:  &strings.Builder{},
		status:  StatusWaiting,
		result:  ResultCancelled,
		started: time.Now(),
		prefix:  prefix,
	}
	bb.logsForTarget[prefix] = writer
	return writer
}

// WriteToDisk aggregates all the data in the numerous prefix writers, and generates an Earthly log bundle.
// These bundles include a manifest generated from the aggregation of the prefixes (targets).
func (bb *BundleBuilder) WriteToDisk() (string, error) {
	// Build file and io.Writer for saving log data
	file, err := os.CreateTemp("", "earthly-log*.tar.gz")
	if err != nil {
		return "", errors.Wrapf(err, "could not create tarball")
	}
	defer file.Close()
	bb.cleanup.Add(func() error {
		return os.Remove(file.Name())
	})

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Make a copy so that we keep the lock for as little time as possible.
	bb.mu.Lock()
	logsForTarget := make(map[string]*targetLogger, len(bb.logsForTarget))
	for k, v := range bb.logsForTarget {
		logsForTarget[k] = v
	}
	bb.mu.Unlock()

	// Convert targets to manifest representations, get tar headers for data
	targetData := make([]TargetManifest, 0)
	for _, lines := range logsForTarget {
		mt, err := lines.toManifestTarget()
		if err != nil {
			// Something was wrong with this targets logs (0 length, or blacklisted name...). Ignore it.
			continue
		}

		targetData = append(targetData, mt)

		trimmed := strings.TrimSpace(lines.prefix)
		escaped := url.QueryEscape(trimmed)

		err = tarWriter.WriteHeader(&tar.Header{
			Name:       fmt.Sprintf("target/%s", escaped),
			Size:       int64(lines.writer.Len()),
			Mode:       0600,
			ChangeTime: time.Now(),
		})
		if err != nil {
			return "", errors.Wrapf(err, "could not write target header")
		}
		_, err = tarWriter.Write([]byte(lines.writer.String()))
		if err != nil {
			return "", errors.Wrapf(err, "could not write target data")
		}
	}

	// build manifest and permissions
	mani := bb.buildManifest(targetData)
	manifestJSON, _ := json.Marshal(mani)
	err = tarWriter.WriteHeader(&tar.Header{
		Name:       "manifest",
		Size:       int64(len(manifestJSON)),
		Mode:       0600,
		ChangeTime: time.Now(),
	})
	if err != nil {
		return "", errors.Wrapf(err, "could not write manifest header")
	}
	_, err = tarWriter.Write(manifestJSON)
	if err != nil {
		return "", errors.Wrapf(err, "could not write manifest")
	}

	perm := bb.buildPermissions()
	permissionsJSON, _ := json.Marshal(perm)
	err = tarWriter.WriteHeader(&tar.Header{
		Name:       "permissions",
		Size:       int64(len(permissionsJSON)),
		Mode:       0600,
		ChangeTime: time.Now(),
	})
	if err != nil {
		return "", errors.Wrapf(err, "could not write permissions header")
	}
	_, err = tarWriter.Write(permissionsJSON)
	if err != nil {
		return "", errors.Wrapf(err, "could not write permissions")
	}

	return file.Name(), nil
}

func (bb *BundleBuilder) buildManifest(targetManifests []TargetManifest) *Manifest {
	manifest := &Manifest{
		Version:    1,
		Duration:   int(time.Since(bb.started).Milliseconds()),
		Status:     StatusComplete,
		Result:     ResultSuccess,
		CreatedAt:  time.Now().In(time.UTC),
		Entrypoint: bb.entrypoint,
		Targets:    targetManifests,
	}

	for _, tm := range targetManifests {
		if tm.Name == fullLog {
			// Full Log reserved name should not determine whole build status.
			// Really, we could go back through after determining whole build status to set _full result & status to the
			// values for the whole build; but it doesn't (yet) affect or mean anything to us. So leave it as is.
			continue
		}

		if tm.Result != ResultSuccess {
			manifest.Result = tm.Result
		}

		if tm.Status != StatusComplete {
			manifest.Status = tm.Status
		}
	}

	return manifest
}

func (bb *BundleBuilder) buildPermissions() *Permissions {
	return &Permissions{
		Version: 1,
		Users:   []string{"*"},
		Orgs:    []string{"*"},
	}
}

func (tl *targetLogger) toManifestTarget() (TargetManifest, error) {
	if tl.writer.Len() <= 0 {
		// Do not write empty logs, if the prefix didn't write anything
		return TargetManifest{}, errors.New("0 length target")
	}

	if tl.prefix == "ongoing" || tl.prefix == "buildkitd" {
		// The ongoing & buildkitd init messages end up in here too. Since they are not updates from a vertex, we will
		// never mark them as complete. Additionally, its not useful to have in the output. Ignore it here.
		return TargetManifest{}, fmt.Errorf("blacklisted target name %s", tl.prefix)
	}

	command, summary := tl.getCommandAndSummary()

	manifestTarget := TargetManifest{
		Name:     tl.prefix,
		Status:   tl.status,
		Result:   tl.result,
		Duration: int(time.Since(tl.started).Milliseconds()),
		Size:     tl.writer.Len(),
		Command:  command,
		Summary:  summary,
	}

	return manifestTarget, nil
}

// Nobody expects ANSI in the command/summary.
// So, even if we don't inject color we should strip it since a tool inside could have done an ANSI too. *SIGH*
const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func (tl *targetLogger) getCommandAndSummary() (string, string) {
	rawText := tl.writer.String()
	text := re.ReplaceAllString(rawText, "")

	prettyPrefix := prettyPrefix(DefaultPadding, tl.prefix)

	// regex to find command lines in the output.
	regexStr := fmt.Sprintf(`(?m)^%s \| (\*cached\* |\*local\* | )*--> `, regexp.QuoteMeta(prettyPrefix))
	r := regexp.MustCompile(regexStr)
	matches := r.FindAllStringIndex(text, -1)
	if len(matches) == 0 {
		return "", ""
	}

	// Take the last match, and use the first line up to 120 characters, or first newlines... whichever comes first.
	lastMatch := matches[len(matches)-1]
	remainder := text[lastMatch[1]:]          // The rest of the log from end of the last match
	command := truncateString(remainder, 120) // The line up to a newline or 120 chars

	// regex to get the last line, (ab)use groups to get the line without the prefix. Truncate it like command.
	regexStr2 := fmt.Sprintf(`%s \| (.*)\n?$`, regexp.QuoteMeta(prettyPrefix))
	r2 := regexp.MustCompile(regexStr2)
	matches2 := r2.FindAllStringSubmatch(remainder, -1)
	if len(matches2) == 0 {
		return command, ""
	}
	summary := truncateString(matches2[len(matches2)-1][1], 120)

	return command, summary
}

func truncateString(str string, length int) string {
	// This weird truncation is needed to support multi-byte characters, which the slice notation does not account for.
	if length <= 0 {
		return ""
	}

	truncated := ""
	count := 0
	for _, char := range str {
		if char == '\n' {
			break
		}
		truncated += string(char)
		count++
		if count >= length {
			break
		}
	}
	return truncated
}
