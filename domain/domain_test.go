//go:build !windows
// +build !windows

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
	{"./a/local/dir+target", Target{Target: "target", LocalPath: "./a/local/dir"}},
	{"/abs/local/dir+target", Target{Target: "target", LocalPath: "/abs/local/dir"}},
	{"/abs/space here/dir+target", Target{Target: "target", LocalPath: "/abs/space here/dir"}},
	{`/abs/back\slash/dir+target`, Target{Target: "target", LocalPath: `/abs/back\slash/dir`}},
	{"../rel/local/dir+target", Target{Target: "target", LocalPath: "../rel/local/dir"}},
	{"github.com/foo/bar+target", Target{Target: "target", GitURL: "github.com/foo/bar"}},
	{"github.com/foo/bar:tag+target", Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag"}},
	{"github.com/foo/bar:tag/with/slash+target", Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag/with/slash"}},
	{"import+target", Target{Target: "target", ImportRef: "import"}},
	// \+
	{"./a/local/dir-with-\\+-in-it+target", Target{Target: "target", LocalPath: "./a/local/dir-with-+-in-it"}},
	{"/abs/local/dir-with-\\+-in+target", Target{Target: "target", LocalPath: "/abs/local/dir-with-+-in"}},
	{"../rel/local/dir-with-\\+-in+target", Target{Target: "target", LocalPath: "../rel/local/dir-with-+-in"}},
	{"github.com/foo/bar/dir-with-\\+-in+target", Target{Target: "target", GitURL: "github.com/foo/bar/dir-with-+-in"}},
	{"github.com/foo/bar:tag-with-\\+-in+target", Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag-with-+-in"}},
}

var targetNegativeTests = []string{
	"+COMMAND", "./something+COMMAND", "nope", "abc+cde+efg", "+target/artifact",
}

func TestTargetParser(t *testing.T) {
	for _, tt := range targetTests {
		t.Run(tt.in, func(t *testing.T) {
			out, err := ParseTarget(tt.in)
			NoError(t, err, "parse target failed")
			Equal(t, tt.out, out)
		})
	}
}

func TestTargetParserNegative(t *testing.T) {
	for _, tt := range targetNegativeTests {
		t.Run(tt, func(t *testing.T) {
			_, err := ParseTarget(tt)
			Error(t, err, "parse target should have failed")
		})
	}
}

func TestTargetToString(t *testing.T) {
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
	{"./a/local/dir+target/artifact", Artifact{Target: Target{Target: "target", LocalPath: "./a/local/dir"}, Artifact: "/artifact"}},
	{"/abs/space here/dir+target/artifact", Artifact{Target: Target{Target: "target", LocalPath: "/abs/space here/dir"}, Artifact: "/artifact"}},
	{`/abs/back\slash/dir+target/artifact`, Artifact{Target: Target{Target: "target", LocalPath: `/abs/back\slash/dir`}, Artifact: "/artifact"}},
	{"github.com/foo/bar+target/artifact", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar"}, Artifact: "/artifact"}},
	{"github.com/foo/bar:tag+target/artifact", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag"}, Artifact: "/artifact"}},
	{"github.com/foo/bar:tag/with/slash+target/artifact", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag/with/slash"}, Artifact: "/artifact"}},
	{"import+target/artifact", Artifact{Target: Target{Target: "target", ImportRef: "import"}, Artifact: "/artifact"}},
	// \+ in target
	{"./a/local/dir-with-\\+-in-it+target/artifact", Artifact{Target: Target{Target: "target", LocalPath: "./a/local/dir-with-+-in-it"}, Artifact: "/artifact"}},
	{"/abs/local/dir-with-\\+-in+target/artifact", Artifact{Target: Target{Target: "target", LocalPath: "/abs/local/dir-with-+-in"}, Artifact: "/artifact"}},
	{"../rel/local/dir-with-\\+-in+target/artifact", Artifact{Target: Target{Target: "target", LocalPath: "../rel/local/dir-with-+-in"}, Artifact: "/artifact"}},
	{"github.com/foo/bar/dir-with-\\+-in+target/artifact", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar/dir-with-+-in"}, Artifact: "/artifact"}},
	{"github.com/foo/bar:tag-with-\\+-in+target/artifact", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag-with-+-in"}, Artifact: "/artifact"}},
	// \+ in artifact
	{"+target/artifact-with-\\+", Artifact{Target: Target{Target: "target", LocalPath: "."}, Artifact: "/artifact-with-+"}},
	{"+another-target/deep/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "another-target", LocalPath: "."}, Artifact: "/deep/artifact-with-+/in/it"}},
	{"+another-target/deep/artifact/with-\\+/and/*", Artifact{Target: Target{Target: "another-target", LocalPath: "."}, Artifact: "/deep/artifact/with-+/and/*"}},
	// \+ in target and artifact
	{"./a/local/dir-with-\\+-in-it+target/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "target", LocalPath: "./a/local/dir-with-+-in-it"}, Artifact: "/artifact-with-+/in/it"}},
	{"/abs/local/dir-with-\\+-in+target/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "target", LocalPath: "/abs/local/dir-with-+-in"}, Artifact: "/artifact-with-+/in/it"}},
	{"../rel/local/dir-with-\\+-in+target/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "target", LocalPath: "../rel/local/dir-with-+-in"}, Artifact: "/artifact-with-+/in/it"}},
	{"github.com/foo/bar/dir-with-\\+-in+target/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar/dir-with-+-in"}, Artifact: "/artifact-with-+/in/it"}},
	{"github.com/foo/bar:tag-with-\\+-in+target/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "target", GitURL: "github.com/foo/bar", Tag: "tag-with-+-in"}, Artifact: "/artifact-with-+/in/it"}},
}

var artifactNegativeTests = []string{
	"+COMMAND/art", "./something+COMMAND/art", "nope/art", "abc+cde+efg/art", "+just-target",
}

func TestArtifactParser(t *testing.T) {
	for _, tt := range artifactTests {
		t.Run(tt.in, func(t *testing.T) {
			out, err := ParseArtifact(tt.in)
			NoError(t, err, "parse artifact failed")
			Equal(t, tt.out, out)
		})
	}
}

func TestArtifactParserNegative(t *testing.T) {
	for _, tt := range artifactNegativeTests {
		t.Run(tt, func(t *testing.T) {
			_, err := ParseArtifact(tt)
			Error(t, err, "parse artifact should have failed")
		})
	}
}

func TestArtifactToString(t *testing.T) {
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
	{"./a/local/dir+COMMAND", Command{Command: "COMMAND", LocalPath: "./a/local/dir"}},
	{"/abs/local/dir+COMMAND", Command{Command: "COMMAND", LocalPath: "/abs/local/dir"}},
	{"../rel/local/dir+COMMAND", Command{Command: "COMMAND", LocalPath: "../rel/local/dir"}},
	{"/abs/space here/dir+COMMAND", Command{Command: "COMMAND", LocalPath: "/abs/space here/dir"}},
	{`/abs/back\slash/dir+COMMAND`, Command{Command: "COMMAND", LocalPath: `/abs/back\slash/dir`}},
	{"github.com/foo/bar+COMMAND", Command{Command: "COMMAND", GitURL: "github.com/foo/bar"}},
	{"github.com/foo/bar:tag+COMMAND", Command{Command: "COMMAND", GitURL: "github.com/foo/bar", Tag: "tag"}},
	{"github.com/foo/bar:tag/with/slash+COMMAND", Command{Command: "COMMAND", GitURL: "github.com/foo/bar", Tag: "tag/with/slash"}},
	{"import+COMMAND", Command{Command: "COMMAND", ImportRef: "import"}},
	// \+
	{"./a/local/dir-with-\\+-in-it+COMMAND", Command{Command: "COMMAND", LocalPath: "./a/local/dir-with-+-in-it"}},
	{"/abs/local/dir-with-\\+-in+COMMAND", Command{Command: "COMMAND", LocalPath: "/abs/local/dir-with-+-in"}},
	{"../rel/local/dir-with-\\+-in+COMMAND", Command{Command: "COMMAND", LocalPath: "../rel/local/dir-with-+-in"}},
	{"github.com/foo/bar/dir-with-\\+-in+COMMAND", Command{Command: "COMMAND", GitURL: "github.com/foo/bar/dir-with-+-in"}},
	{"github.com/foo/bar:tag-with-\\+-in+COMMAND", Command{Command: "COMMAND", GitURL: "github.com/foo/bar", Tag: "tag-with-+-in"}},
}

var commandNegativeTests = []string{
	"+target", "./something+target", "nope", "NOPE", "ABC+DEF+EFG", "+COMMAND/artifact",
}

func TestCommandParser(t *testing.T) {
	for _, tt := range commandTests {
		t.Run(tt.in, func(t *testing.T) {
			out, err := ParseCommand(tt.in)
			NoError(t, err, "parse target failed")
			Equal(t, tt.out, out)
		})
	}
}

func TestCommandParserNegative(t *testing.T) {
	for _, tt := range commandNegativeTests {
		t.Run(tt, func(t *testing.T) {
			_, err := ParseCommand(tt)
			Error(t, err, "parse command should have failed")
		})
	}
}

func TestCommandToString(t *testing.T) {
	for _, tt := range commandTests {
		t.Run(tt.in, func(t *testing.T) {
			str := tt.out.String()
			Equal(t, tt.in, str)
		})
	}
}
