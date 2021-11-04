package earthfile2llb

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_validateCopySources(t *testing.T) {
	testcases := []struct {
		name string
		srcs []string
		err  error
	}{
		{
			name: "all from current build context",
			srcs: []string{"foo", "bar"},
			err:  nil,
		},
		{
			name: "source with parent dir",
			srcs: []string{"..", "bar"},
			err:  fmt.Errorf("COPY does not support whole parent directories for local source: %q", ".."),
		},
		{
			name: "source with parent dir trailing slash",
			srcs: []string{"../", "bar"},
			err:  fmt.Errorf("COPY does not support whole parent directories for local source: %q", "../"),
		},
		{
			name: "source with parent dir and glob",
			srcs: []string{"../foo", "../*"},
			err:  fmt.Errorf("COPY does not support glob patterns using parent directories for local source: %q", "../*"),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			err := validateCopySources(testcase.srcs)
			if !reflect.DeepEqual(err, testcase.err) {
				t.Logf("actual err: %v", err)
				t.Logf("expected err: %v", testcase.err)
				t.Error("unexpected error")
			}
		})
	}
}

func Test_buildContextForSources(t *testing.T) {
	testcases := []struct {
		name         string
		srcs         []string
		basePath     string
		buildContext string
		err          error
	}{
		{
			name:         "all from current build context",
			srcs:         []string{"foo", "bar"},
			basePath:     ".",
			buildContext: ".",
			err:          nil,
		},
		{
			name:         "all from local path build context",
			srcs:         []string{"foo", "bar"},
			basePath:     "./test",
			buildContext: "test",
			err:          nil,
		},
		{
			name:         "all from parent build context",
			srcs:         []string{"../foo", "../bar"},
			basePath:     ".",
			buildContext: "..",
			err:          nil,
		},
		{
			name:         "one from current build context and one from parent",
			srcs:         []string{"foo", "../bar"},
			basePath:     ".",
			buildContext: "",
			err:          fmt.Errorf("COPY command only supports a single build context, detected two: %q and %q", ".", ".."),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			bc, err := buildContextForSources(testcase.srcs, testcase.basePath)
			if bc != testcase.buildContext {
				t.Logf("actual build context: %q", bc)
				t.Logf("expected build context: %q", testcase.buildContext)
				t.Error("unexpected build context")
			}
			if !reflect.DeepEqual(err, testcase.err) {
				t.Logf("actual err: %v", err)
				t.Logf("expected err: %v", testcase.err)
				t.Error("unexpected error")
			}
		})
	}
}

func Test_buildContextFromPath(t *testing.T) {
	testcases := []struct {
		name         string
		path         string
		base         string
		buildContext string
	}{
		{
			name:         "from current directory",
			path:         "foo/bar",
			base:         ".",
			buildContext: ".",
		},
		{
			name:         "from current directory starting with ./",
			path:         "./foo/bar",
			base:         ".",
			buildContext: ".",
		},
		{
			name:         "from current directory starting with ./ with child base path",
			path:         "./foo/bar",
			base:         "./subdir",
			buildContext: "subdir",
		},
		{
			name:         "from parent directory",
			path:         "../foo/bar",
			base:         ".",
			buildContext: "..",
		},
		{
			name:         "from parent directory with child base path",
			path:         "../foo/bar",
			base:         "./subdir",
			buildContext: "subdir/..",
		},
		{
			name:         "from parent directory starting with ./",
			path:         "./../foo/bar",
			base:         ".",
			buildContext: "..",
		},
		{
			name:         "from two parent directory",
			path:         "../../foo/bar",
			base:         ".",
			buildContext: "../..",
		},
		{
			name:         "from two parent directory starting with ./",
			path:         "./../../foo/bar",
			base:         ".",
			buildContext: "../..",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			bc := buildContextFromPath(testcase.path, testcase.base)
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
