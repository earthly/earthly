package buildcontext

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_readExcludes(t *testing.T) {
	testcases := []struct {
		name                  string
		earthIgnoreContents   string
		earthlyIgnoreContents string
		expectedExcludes      []string
		expectedErr           error
	}{
		{
			name:                  "only .earthlyignore",
			earthlyIgnoreContents: `foobar/`,
			expectedExcludes:      []string{"foobar", ".tmp-earthly-out/", "build.earth", "Earthfile", ".earthignore", ".earthlyignore"},
		},
		{
			name:                "only .earthignore",
			earthIgnoreContents: `foobar/`,
			expectedExcludes:    []string{"foobar", ".tmp-earthly-out/", "build.earth", "Earthfile", ".earthignore", ".earthlyignore"},
		},
		{
			name:             "no ignore file, default to implicit rules",
			expectedExcludes: ImplicitExcludes,
		},
		{
			name:                  "both .earthignore and .earthlyignore results in error",
			earthlyIgnoreContents: `foobar/`,
			earthIgnoreContents:   `foobar/`,
			expectedExcludes:      ImplicitExcludes,
			expectedErr:           errDuplicateIgnoreFile,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "earthly-test-read-excludes")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			if testcase.earthIgnoreContents != "" {
				earthIgnoreFile, err := os.Create(filepath.Join(dir, earthIgnoreFile))
				if err != nil {
					t.Fatalf("failed to create .earthignore file")
				}

				_, err = earthIgnoreFile.WriteString(testcase.earthIgnoreContents)
				if err != nil {
					t.Fatalf("failed to write .earthignore file")
				}
			}

			if testcase.earthlyIgnoreContents != "" {
				earthlyIgnoreFile, err := os.Create(filepath.Join(dir, earthlyIgnoreFile))
				if err != nil {
					t.Fatalf("failed to create .earthlyignore file")
				}

				_, err = earthlyIgnoreFile.WriteString(testcase.earthlyIgnoreContents)
				if err != nil {
					t.Fatalf("failed to write .earthlyignore file")
				}
			}

			excludes, err := readExcludes(dir)
			if err != testcase.expectedErr {
				t.Logf("actual err: %v", err)
				t.Logf("expected err: %v", testcase.expectedErr)
				t.Error("unexpected error getting excludes")
			}

			if !reflect.DeepEqual(excludes, testcase.expectedExcludes) {
				t.Logf("actual excludes: %v", excludes)
				t.Logf("expected excludes: %v", testcase.expectedExcludes)
				t.Error("unexpected excludes list")
			}
		})
	}
}
