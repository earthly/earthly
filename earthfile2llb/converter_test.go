package earthfile2llb

import (
	"testing"
)

func Test_buildContextFromPath(t *testing.T) {
	testcases := []struct {
		name         string
		path         string
		buildContext string
	}{
		{
			name:         "from current directory",
			path:         "foo/bar",
			buildContext: ".",
		},
		{
			name:         "from current directory starting with ./",
			path:         "./foo/bar",
			buildContext: ".",
		},
		{
			name:         "from parent directory",
			path:         "../foo/bar",
			buildContext: "..",
		},
		{
			name:         "from parent directory starting with ./",
			path:         "./../foo/bar",
			buildContext: "..",
		},
		{
			name:         "from two parent directory",
			path:         "../../foo/bar",
			buildContext: "../..",
		},
		{
			name:         "from two parent directory starting with ./",
			path:         "./../../foo/bar",
			buildContext: "../..",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			bc := buildContextFromPath(testcase.path)
			if bc != testcase.buildContext {
				t.Logf("actual build context: %q", bc)
				t.Logf("expected build context: %q", testcase.buildContext)
				t.Error("unexpected build context")
			}
		})
	}
}

func Test_pathFromBuildContext(t *testing.T) {
	testcases := []struct {
		name             string
		path             string
		fromBuildContext string
	}{
		{
			name:             "from current directory",
			path:             "foo/bar",
			fromBuildContext: "foo/bar",
		},
		{
			name:             "from current directory starting with ./",
			path:             "./foo/bar",
			fromBuildContext: "foo/bar",
		},
		{
			name:             "from parent directory",
			path:             "../foo/bar",
			fromBuildContext: "foo/bar",
		},
		{
			name:             "from parent directory starting with ./",
			path:             "./../foo/bar",
			fromBuildContext: "foo/bar",
		},
		{
			name:             "from two parent directory",
			path:             "../../foo/bar",
			fromBuildContext: "foo/bar",
		},
		{
			name:             "from two parent directory starting with ./",
			path:             "./../../foo/bar",
			fromBuildContext: "foo/bar",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			path := pathFromBuildContext(testcase.path)
			if path != testcase.fromBuildContext {
				t.Logf("actual path from build context: %q", path)
				t.Logf("expected path from build context: %q", testcase.fromBuildContext)
				t.Error("unexpected path from build context")
			}
		})
	}
}
