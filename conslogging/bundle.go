package conslogging

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type BundleBuilder struct {
	RootPath   string
	Entrypoint string

	logmap  map[string]*TargetLogger
	started time.Time
}

type TargetLogger struct {
	writer  *strings.Builder
	result  string
	status  string
	started time.Time
}

func (tl *TargetLogger) Write(p []byte) (n int, err error) {
	return tl.writer.Write(p)
}

func NewBundleBuilder(rootPath, entrypoint string) *BundleBuilder {
	return &BundleBuilder{
		RootPath:   rootPath,
		Entrypoint: entrypoint,
		logmap:     map[string]*TargetLogger{},
		started:    time.Now(),
	}
}

func (bb *BundleBuilder) PrefixResult(prefix, result string) {
	if builder, ok := bb.logmap[prefix]; ok {
		builder.result = result
	}
}

func (bb *BundleBuilder) PrefixStatus(prefix, status string) {
	if builder, ok := bb.logmap[prefix]; ok {
		builder.status = status
	}
}

func (bb *BundleBuilder) PrefixWriter(prefix string) io.Writer {
	if builder, ok := bb.logmap[prefix]; ok {
		return builder
	}

	writer := &TargetLogger{
		writer:  &strings.Builder{},
		status:  StatusWaiting,
		result:  "",
		started: time.Now(),
	}
	bb.logmap[prefix] = writer
	return writer
}

func (bb *BundleBuilder) WriteToDisk() error {
	fmt.Println(bb.RootPath)

	targetPath := path.Join(bb.RootPath, "target")
	err := os.MkdirAll(targetPath, 0700)
	if err != nil {
		return errors.Wrapf(err, "Failed to write targets directory for bundle at %s", bb.RootPath)
	}

	manifest := &Manifest{
		Version:    1,
		Duration:   int(time.Since(bb.started).Seconds()),
		Status:     "complete",
		Result:     "success",
		CreatedAt:  time.Now().In(time.UTC),
		Entrypoint: bb.Entrypoint,
		Targets:    make([]TargetManifest, 0),
	}

	for prefix, lines := range bb.logmap {
		if lines.writer.Len() <= 0 {
			// Do not write empty logs, if the prefix didn't write anything
			continue
		}

		trimmed := strings.TrimSpace(prefix)
		escaped := url.PathEscape(trimmed)
		logPath := path.Join(targetPath, escaped)

		command, summary := bb.GetCommandAndSummary(prefix, lines.writer)

		manifest.Targets = append(manifest.Targets, TargetManifest{
			Name:     prefix,
			Status:   lines.status,
			Result:   lines.result,
			Duration: int(time.Since(lines.started).Seconds()),
			Size:     lines.writer.Len(),
			Command:  command,
			Summary:  summary,
		})

		if lines.result != ResultSuccess {
			manifest.Result = lines.result
		}

		if lines.status != StatusComplete {
			manifest.Status = lines.status
		}

		tgtErr := ioutil.WriteFile(logPath, []byte(lines.writer.String()), 0600)
		if err != nil {
			err = multierror.Append(err, tgtErr)
		}
	}
	if err != nil {
		return errors.Wrap(err, "errors while writing targets for log bundle")
	}

	manifestJSON, _ := json.Marshal(&manifest)
	err = ioutil.WriteFile(path.Join(bb.RootPath, "manifest"), manifestJSON, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to serialize bundle manifest")
	}

	permissionsJSON, _ := json.Marshal(&Permissions{
		Version: 1,
		Users:   make([]uint64, 0),
		Orgs:    make([]uint64, 0),
	})
	err = ioutil.WriteFile(path.Join(bb.RootPath, "permissions"), permissionsJSON, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to serialize bundle permissions placeholder")
	}

	return nil
}

// Nobody expects ANSI in the command/summary.
// So, even if we don't inject color we should strip it since a tool inside could have done an ANSI too. *SIGH*
const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func (bb *BundleBuilder) GetCommandAndSummary(prefix string, builder *strings.Builder) (string, string) {
	rawText := builder.String()
	text := re.ReplaceAllString(rawText, "")

	prettyPrefix := prettyPrefix2(DefaultPadding, prefix)

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
	command := TruncateString(remainder, 120) // The line up to a newline or 120 chars

	// regex to get the last line, (ab)use groups to get the line without the prefix. Truncate it like command.
	regexStr2 := fmt.Sprintf(`%s \| (.*)\n?$`, regexp.QuoteMeta(prettyPrefix))
	r2 := regexp.MustCompile(regexStr2)
	matches2 := r2.FindAllStringSubmatch(remainder, -1)
	if len(matches2) == 0 {
		return command, ""
	}
	summary := TruncateString(matches2[len(matches2)-1][1], 120)

	return command, summary
}

// TruncateString truncates a string honoring weird glyphs that are larger than one byte... like Japanese.
func TruncateString(str string, length int) string {
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

// --- Copied types below

const (
	StatusWaiting    = "waiting"
	StatusInProgress = "in_progress"
	StatusComplete   = "complete"
	StatusCancelled  = "cancelled"

	ResultSuccess   = "success"
	ResultFailure   = "failure"
	ResultCancelled = "cancelled"
)

type Manifest struct {
	Version    int              `json:"version"`
	Duration   int              `json:"duration"`
	Status     string           `json:"status"`
	Result     string           `json:"result"`
	CreatedAt  time.Time        `json:"created_at"`
	Targets    []TargetManifest `json:"targets"`
	Entrypoint string           `json:"entrypoint"`
}

type TargetManifest struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Result   string `json:"result"`
	Duration int    `json:"duration"`
	Size     int    `json:"size"`
	Command  string `json:"command,omitempty"`
	Summary  string `json:"summary,omitempty"`
}

type Permissions struct {
	Version int      `json:"version"`
	Users   []uint64 `json:"users"`
	Orgs    []uint64 `json:"orgs"`
}
