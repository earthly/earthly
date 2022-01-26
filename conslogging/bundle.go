package conslogging

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
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
	var err error

	fmt.Println(bb.RootPath)

	manifest := &Manifest{
		Version:    1,
		Duration:   int(time.Since(bb.started).Seconds()),
		Status:     "complete",
		Result:     "success",
		CreatedAt:  time.Now(),
		Entrypoint: bb.Entrypoint,
		Targets:    make([]TargetManifest, 0),
	}

	for prefix, lines := range bb.logmap {
		trimmed := strings.TrimSpace(prefix)
		escaped := url.PathEscape(trimmed)
		logPath := path.Join(bb.RootPath, escaped)

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

		tgtErr := ioutil.WriteFile(logPath, []byte(lines.writer.String()), 0666)
		if err != nil {
			err = multierror.Append(err, tgtErr)
		}
	}

	manifestJSON, _ := json.Marshal(&manifest)
	err = ioutil.WriteFile(path.Join(bb.RootPath, "manifest.json"), manifestJSON, 0666)

	return err
}

// Nobody expects ANSI in the command/summary.
// So, even if we don't inject color we should strip it since a tool inside could have done an ANSI too.
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

	// find the next newline, and go there

	// If there is any left, take the last line, up to 120 characters
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
