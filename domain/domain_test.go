package domain

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

var targetTests = []struct {
	in  string
	out Target
}{
	{"+target", Target{Target: "target", LocalPath: "."}},
	{"+another-target", Target{Target: "another-target", LocalPath: "."}},
	{"./a/local/dir+target", Target{Target: "target", LocalPath: "./a/local/dir"}},
	{"/abs/local/dir+target", Target{Target: "target", LocalPath: "/abs/local/dir"}},
	{"../rel/local/dir+target", Target{Target: "target", LocalPath: "../rel/local/dir"}},
	{"github.com/foo/bar+target", Target{Target: "target", Registry: "github.com", ProjectPath: "foo/bar"}},
	{"github.com/foo/bar:tag+target", Target{Target: "target", Registry: "github.com", ProjectPath: "foo/bar", Tag: "tag"}},
	{"github.com/foo/bar:tag/with/slash+target", Target{Target: "target", Registry: "github.com", ProjectPath: "foo/bar", Tag: "tag/with/slash"}},
	// \+
	{"./a/local/dir-with-\\+-in-it+target", Target{Target: "target", LocalPath: "./a/local/dir-with-+-in-it"}},
	{"/abs/local/dir-with-\\+-in+target", Target{Target: "target", LocalPath: "/abs/local/dir-with-+-in"}},
	{"../rel/local/dir-with-\\+-in+target", Target{Target: "target", LocalPath: "../rel/local/dir-with-+-in"}},
	{"github.com/foo/bar/dir-with-\\+-in+target", Target{Target: "target", Registry: "github.com", ProjectPath: "foo/bar/dir-with-+-in"}},
	{"github.com/foo/bar:tag-with-\\+-in+target", Target{Target: "target", Registry: "github.com", ProjectPath: "foo/bar", Tag: "tag-with-+-in"}},
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
	{"github.com/foo/bar+target/artifact", Artifact{Target: Target{Target: "target", Registry: "github.com", ProjectPath: "foo/bar"}, Artifact: "/artifact"}},
	{"github.com/foo/bar:tag+target/artifact", Artifact{Target: Target{Target: "target", Registry: "github.com", ProjectPath: "foo/bar", Tag: "tag"}, Artifact: "/artifact"}},
	{"github.com/foo/bar:tag/with/slash+target/artifact", Artifact{Target: Target{Target: "target", Registry: "github.com", ProjectPath: "foo/bar", Tag: "tag/with/slash"}, Artifact: "/artifact"}},
	// \+ in target
	{"./a/local/dir-with-\\+-in-it+target/artifact", Artifact{Target: Target{Target: "target", LocalPath: "./a/local/dir-with-+-in-it"}, Artifact: "/artifact"}},
	{"/abs/local/dir-with-\\+-in+target/artifact", Artifact{Target: Target{Target: "target", LocalPath: "/abs/local/dir-with-+-in"}, Artifact: "/artifact"}},
	{"../rel/local/dir-with-\\+-in+target/artifact", Artifact{Target: Target{Target: "target", LocalPath: "../rel/local/dir-with-+-in"}, Artifact: "/artifact"}},
	{"github.com/foo/bar/dir-with-\\+-in+target/artifact", Artifact{Target: Target{Target: "target", Registry: "github.com", ProjectPath: "foo/bar/dir-with-+-in"}, Artifact: "/artifact"}},
	{"github.com/foo/bar:tag-with-\\+-in+target/artifact", Artifact{Target: Target{Target: "target", Registry: "github.com", ProjectPath: "foo/bar", Tag: "tag-with-+-in"}, Artifact: "/artifact"}},
	// \+ in artifact
	{"+target/artifact-with-\\+", Artifact{Target: Target{Target: "target", LocalPath: "."}, Artifact: "/artifact-with-+"}},
	{"+another-target/deep/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "another-target", LocalPath: "."}, Artifact: "/deep/artifact-with-+/in/it"}},
	{"+another-target/deep/artifact/with-\\+/and/*", Artifact{Target: Target{Target: "another-target", LocalPath: "."}, Artifact: "/deep/artifact/with-+/and/*"}},
	// \+ in target and artifact
	{"./a/local/dir-with-\\+-in-it+target/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "target", LocalPath: "./a/local/dir-with-+-in-it"}, Artifact: "/artifact-with-+/in/it"}},
	{"/abs/local/dir-with-\\+-in+target/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "target", LocalPath: "/abs/local/dir-with-+-in"}, Artifact: "/artifact-with-+/in/it"}},
	{"../rel/local/dir-with-\\+-in+target/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "target", LocalPath: "../rel/local/dir-with-+-in"}, Artifact: "/artifact-with-+/in/it"}},
	{"github.com/foo/bar/dir-with-\\+-in+target/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "target", Registry: "github.com", ProjectPath: "foo/bar/dir-with-+-in"}, Artifact: "/artifact-with-+/in/it"}},
	{"github.com/foo/bar:tag-with-\\+-in+target/artifact-with-\\+/in/it", Artifact{Target: Target{Target: "target", Registry: "github.com", ProjectPath: "foo/bar", Tag: "tag-with-+-in"}, Artifact: "/artifact-with-+/in/it"}},
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

func TestArtifactToString(t *testing.T) {
	for _, tt := range artifactTests {
		t.Run(tt.in, func(t *testing.T) {
			str := tt.out.String()
			Equal(t, tt.in, str)
		})
	}
}
