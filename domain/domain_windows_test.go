//go:build windows
// +build windows

package domain

import (
	"testing"
)

var targetTests = []struct {
	in  string
	out Target
}{
	{"+target", Target{Target: "target", LocalPath: "."}},
	{"+another-target", Target{Target: "another-target", LocalPath: "."}},
	{`.\rel\win\dir+target`, Target{Target: "target", LocalPath: `.\rel\win\dir`}},
	{`./rel/win/dir+target`, Target{Target: "target", LocalPath: `./rel/win/dir`}},
	{`C:\abs\win\dir+target`, Target{Target: "target", LocalPath: `C:\abs\win\dir`}},
	{`.\rel\space here\dir+target`, Target{Target: "target", LocalPath: `.\rel\space here\dir`}},
	{`.\rel\fwd/slash\dir+target`, Target{Target: "target", LocalPath: `.\rel\fwd/slash\dir`}},
	{"github.com/foo/bar+target", Target{Target: "target", GitURL: "github.com/foo/bar"}},
	{"github.com/foo/bar:tag+target", Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag"}},
	{"github.com/foo/bar:tag/with/slash+target", Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag/with/slash"}},
	{"import+target", Target{Target: "target", ImportRef: "import"}},
	{"github.com/foo/bar/dir-with-\\+-in+target", Target{Target: "target", GitURL: "github.com/foo/bar/dir-with-+-in"}},
	{"github.com/foo/bar:tag-with-\\+-in+target", Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag-with-+-in"}},
}

var targetNegativeTests = []string{
	"+COMMAND", "./something+COMMAND", "nope", "abc+cde+efg", "+target/artifact",
}

func TestTargetParserWin(t *testing.T) {
	for _, tt := range targetTests {
		t.Run(tt.in, func(t *testing.T) {
			out, err := ParseTarget(tt.in)
			NoError(t, err, "parse target failed")
			Equal(t, tt.out, out)
		})
	}
}

func TestTargetParserNegativeWin(t *testing.T) {
	for _, tt := range targetNegativeTests {
		t.Run(tt, func(t *testing.T) {
			_, err := ParseTarget(tt)
			Error(t, err, "parse target should have failed")
		})
	}
}

func TestTargetToStringWin(t *testing.T) {
	for _, tt := range targetTests {
		t.Run(tt.in, func(t *testing.T) {
			str := tt.out.String()
			Equal(t, tt.in, str)
		})
	}
}

var artifactTests = []struct {
	in  string
	out Artifact
}{
	{"+target/artifact", Artifact{Target: Target{Target: "target", LocalPath: "."}, Artifact: "/artifact"}},
	{"+another-target/another-artifact", Artifact{Target: Target{Target: "another-target", LocalPath: "."}, Artifact: "/another-artifact"}},
	{"+another-target/deep/artifact", Artifact{Target: Target{Target: "another-target", LocalPath: "."}, Artifact: "/deep/artifact"}},
	{"+another-target/deep/artifact/with/*", Artifact{Target: Target{Target: "another-target", LocalPath: "."}, Artifact: "/deep/artifact/with/*"}},
	{`.\rel\win\dir+target/artifact`, Artifact{Target: Target{Target: "target", LocalPath: `.\rel\win\dir`}, Artifact: "/artifact"}},
	{`./rel/win/dir+target/artifact`, Artifact{Target: Target{Target: "target", LocalPath: `./rel/win/dir`}, Artifact: "/artifact"}},
	{`.\rel\space here\dir+target/artifact`, Artifact{Target: Target{Target: "target", LocalPath: `.\rel\space here\dir`}, Artifact: "/artifact"}},
	{`.\rel\fwd/slash\dir+target/artifact`, Artifact{Target: Target{Target: "target", LocalPath: `.\rel\fwd/slash\dir`}, Artifact: "/artifact"}},
	{`C:\abs\win\dir+target/artifact`, Artifact{Target: Target{Target: "target", LocalPath: `C:\abs\win\dir`}, Artifact: "/artifact"}},
	{"github.com/foo/bar+target/artifact", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar"}, Artifact: "/artifact"}},
	{"github.com/foo/bar:tag+target/artifact", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag"}, Artifact: "/artifact"}},
	{"github.com/foo/bar:tag/with/slash+target/artifact", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag/with/slash"}, Artifact: "/artifact"}},
	{"github.com/foo/bar/dir-with-\\+-in+target/artifact", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar/dir-with-+-in"}, Artifact: "/artifact"}},
	{"github.com/foo/bar:tag-with-\\+-in+target/artifact", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag-with-+-in"}, Artifact: "/artifact"}},
	{"github.com/foo/bar/dir-with-\\+-in+target/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar/dir-with-+-in"}, Artifact: "/artifact-with-+/in/it"}},
	{"github.com/foo/bar:tag-with-\\+-in+target/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag-with-+-in"}, Artifact: "/artifact-with-+/in/it"}},
}

var artifactNegativeTests = []string{
	"+COMMAND/art", "./something+COMMAND/art", "nope/art", "abc+cde+efg/art", "+just-target",
}

func TestArtifactParserWin(t *testing.T) {
	for _, tt := range artifactTests {
		t.Run(tt.in, func(t *testing.T) {
			out, err := ParseArtifact(tt.in)
			NoError(t, err, "parse artifact failed")
			Equal(t, tt.out, out)
		})
	}
}

func TestArtifactParserNegativeWin(t *testing.T) {
	for _, tt := range artifactNegativeTests {
		t.Run(tt, func(t *testing.T) {
			_, err := ParseArtifact(tt)
			Error(t, err, "parse artifact should have failed")
		})
	}
}

func TestArtifactToStringWin(t *testing.T) {
	for _, tt := range artifactTests {
		t.Run(tt.in, func(t *testing.T) {
			str := tt.out.String()
			Equal(t, tt.in, str)
		})
	}
}

var commandTests = []struct {
	in  string
	out Command
}{
	{"+COMMAND", Command{Command: "COMMAND", LocalPath: "."}},
	{"+ANOTHER_COMMAND", Command{Command: "ANOTHER_COMMAND", LocalPath: "."}},
	{`.\rel\win\dir+COMMAND`, Command{Command: "COMMAND", LocalPath: `.\rel\win\dir`}},
	{`./rel/win/dir+COMMAND`, Command{Command: "COMMAND", LocalPath: `./rel/win/dir`}},
	{`.\rel\space here\dir+COMMAND`, Command{Command: "COMMAND", LocalPath: `.\rel\space here\dir`}},
	{`.\rel\fwd/slash\dir+COMMAND`, Command{Command: "COMMAND", LocalPath: `.\rel\fwd/slash\dir`}},
	{`C:\abs\win\dir+COMMAND`, Command{Command: "COMMAND", LocalPath: `C:\abs\win\dir`}},
	{"github.com/foo/bar+COMMAND", Command{Command: "COMMAND", GitURL: "github.com/foo/bar"}},
	{"github.com/foo/bar:tag+COMMAND", Command{Command: "COMMAND", GitURL: "github.com/foo/bar", Tag: "tag"}},
	{"github.com/foo/bar:tag/with/slash+COMMAND", Command{Command: "COMMAND", GitURL: "github.com/foo/bar", Tag: "tag/with/slash"}},
	{"import+COMMAND", Command{Command: "COMMAND", ImportRef: "import"}},
	{"github.com/foo/bar/dir-with-\\+-in+COMMAND", Command{Command: "COMMAND", GitURL: "github.com/foo/bar/dir-with-+-in"}},
	{"github.com/foo/bar:tag-with-\\+-in+COMMAND", Command{Command: "COMMAND", GitURL: "github.com/foo/bar", Tag: "tag-with-+-in"}},
}

var commandNegativeTests = []string{
	"+target", "./something+target", "nope", "NOPE", "ABC+DEF+EFG", "+COMMAND/artifact",
}

func TestCommandParserWin(t *testing.T) {
	for _, tt := range commandTests {
		t.Run(tt.in, func(t *testing.T) {
			out, err := ParseCommand(tt.in)
			NoError(t, err, "parse target failed")
			Equal(t, tt.out, out)
		})
	}
}

func TestCommandParserNegativeWin(t *testing.T) {
	for _, tt := range commandNegativeTests {
		t.Run(tt, func(t *testing.T) {
			_, err := ParseCommand(tt)
			Error(t, err, "parse command should have failed")
		})
	}
}

func TestCommandToStringWin(t *testing.T) {
	for _, tt := range commandTests {
		t.Run(tt.in, func(t *testing.T) {
			str := tt.out.String()
			Equal(t, tt.in, str)
		})
	}
}
