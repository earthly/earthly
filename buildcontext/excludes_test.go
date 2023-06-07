package buildcontext

import (
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
		dockerIgnoreContents  string
		useDockerIgnore       bool
		noImplicitIgnore      bool
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
			name:                 "only .dockerignore",
			dockerIgnoreContents: `foobar/`,
			useDockerIgnore:      true,
			expectedExcludes:     []string{"foobar", ".tmp-earthly-out/", "build.earth", "Earthfile", ".earthignore", ".earthlyignore"},
		},
		{
			name:                  "only .earthlyignore with no implicit ignore",
			earthlyIgnoreContents: `foobar/`,
			noImplicitIgnore:      true,
			expectedExcludes:      []string{"foobar"},
		},
		{
			name:                "only .earthignore with no implicit ignore",
			earthIgnoreContents: `foobar/`,
			noImplicitIgnore:    true,
			expectedExcludes:    []string{"foobar"},
		},
		{
			name:                 "only .dockerignore with no implicit ignore",
			dockerIgnoreContents: `foobar/`,
			noImplicitIgnore:     true,
			useDockerIgnore:      true,
			expectedExcludes:     []string{"foobar"},
		},
		{
			name:             "no ignore file, default to implicit rules",
			expectedExcludes: ImplicitExcludes,
		},
		{
			name:             "no ignore file and no implicit ignore",
			noImplicitIgnore: true,
			expectedExcludes: []string{},
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
			dir := t.TempDir()

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

			if testcase.dockerIgnoreContents != "" {
				dockerIgnoreFile, err := os.Create(filepath.Join(dir, dockerIgnoreFile))
				if err != nil {
					t.Fatalf("failed to create .dockerignore file")
				}

				_, err = dockerIgnoreFile.WriteString(testcase.dockerIgnoreContents)
				if err != nil {
					t.Fatalf("failed to write .dockerignore file")
				}
			}

			excludes, err := readExcludes(dir, testcase.noImplicitIgnore, testcase.useDockerIgnore)
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
