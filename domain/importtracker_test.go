package domain

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/hint"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImports(t *testing.T) {
	var tests = []struct {
		importStr                   string
		as                          string
		ref                         string
		expected                    string
		expectedPathResultFuncCalls int
		ok                          bool
	}{
		{"github.com/foo/bar", "", "bar+abc", "github.com/foo/bar+abc", 0, true},
		{"github.com/foo/bar", "buz", "buz+abc", "github.com/foo/bar+abc", 0, true},
		{"github.com/foo/bar", "buz", "bar+abc", "", 0, false},
		{"github.com/foo/bar:v1.2.3", "", "bar+abc", "github.com/foo/bar:v1.2.3+abc", 0, true},
		{"github.com/foo/bar:v1.2.3", "buz", "buz+abc", "github.com/foo/bar:v1.2.3+abc", 0, true},
		{"github.com/foo/bar:v1.2.3", "buz", "bar+abc", "", 0, false},
		{"./foo/bar", "", "bar+abc", "./foo/bar+abc", 2, true},
		{"./foo/bar", "buz", "buz+abc", "./foo/bar+abc", 2, true},
		{"./foo/bar", "buz", "bar+abc", "", 2, false},
		{"../foo/bar", "", "bar+abc", "../foo/bar+abc", 2, true},
		{"../foo/bar", "buz", "buz+abc", "../foo/bar+abc", 2, true},
		{"../foo/bar", "buz", "bar+abc", "", 2, false},
		{"/foo/bar", "", "bar+abc", "/foo/bar+abc", 2, true},
		{"/foo/bar", "buz", "buz+abc", "/foo/bar+abc", 2, true},
		{"/foo/bar", "buz", "bar+abc", "", 2, false},
	}

	var console conslogging.ConsoleLogger

	for _, tt := range tests {
		ir := NewImportTracker(console, nil)
		pathResultFuncCounter := 0
		ir.pathResultFunc = func(path string) pathResult {
			pathResultFuncCounter++
			if filepath.Base(path) == "Earthfile" {
				return file
			}
			return dir
		}
		err := ir.Add(tt.importStr, tt.as, false, false, false)
		require.NoError(t, err, "add import error")
		require.Equalf(t, tt.expectedPathResultFuncCalls, pathResultFuncCounter, "pathResultFunc was called an unexpected number of times")

		ref, err := ParseTarget(tt.ref)
		require.NoError(t, err, "parse test case ref") // check that the test data is good
		require.Equal(t, tt.ref, ref.String())         // sanity check

		ref2, _, _, err := ir.Deref(ref)
		if tt.ok {
			assert.NoError(t, err, "deref import")
			assert.Equal(t, tt.expected, ref2.StringCanonical()) // StringCanonical shows its resolved form
			assert.Equal(t, tt.ref, ref2.String())               // String shows its import form
		} else {
			assert.Error(t, err, "deref should have error'd")
		}
	}
}

func TestImportAdd(t *testing.T) {
	var tests = map[string]struct {
		importStr string
		f         pathResultFunc
		expected  error
	}{
		"path does not exist": {
			importStr: "./foo/bar",
			f: func(path string) pathResult {
				return notExist
			},
			expected: hint.Wrapf(errors.New(`path "./foo/bar" does not exist`), `Verify the path "./foo/bar" exists`),
		},
		"not a directory (ends with Earthfile)": {
			importStr: "./foo/bar/Earthfile",
			f: func(path string) pathResult {
				return file
			},
			expected: hint.Wrap(errors.New(`path "./foo/bar/Earthfile" is not a directory`), `Did you mean to import "./foo/bar"?`),
		},
		"not a directory": {
			importStr: "./foo/bar",
			f: func(path string) pathResult {
				return file
			},
			expected: hint.Wrap(errors.New(`path "./foo/bar/Earthfile" is not a directory`), "Please use a directory when using a local IMPORT path"),
		},
		"path ends with an Earthfile directory": {
			importStr: "./foo/bar/Earthfile",
			f: func(path string) pathResult {
				if filepath.Base(filepath.Dir(path)) == "Earthfile" {
					return dir
				}
				return notExist
			},
			expected: hint.Wrap(errors.New(`path "./foo/bar" does not contain an Earthfile`), `The path "./foo/bar" ends with an "Earthfile" which is a directory.\nDid you mean to create an "Earthfile" as a file instead?`),
		},
		"Earthfile does not exist": {
			importStr: "./foo/bar",
			f: func(path string) pathResult {
				if filepath.Base(path) == "Earthfile" {
					return notExist
				}
				return dir
			},
			expected: hint.Wrap(errors.New(`path "./foo/bar" does not contain an "Earthfile"`), `Verify the path "./foo/bar" contains an Earthfile`),
		},
		"Earthfile exists but is a directory": {
			importStr: "./foo/bar",
			f: func(path string) pathResult {
				return dir
			},
			expected: hint.Wrap(errors.New(`path "./foo/bar" does contains an "Earthfile" which is not a file`), `The local IMPORT path "./foo/bar" contains an "Earthfile" directory and not a file`),
		},
	}

	var console conslogging.ConsoleLogger

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ir := NewImportTracker(console, nil)
			ir.pathResultFunc = tt.f
			err := ir.Add(tt.importStr, "alias", false, false, false)
			assert.Error(t, err, "add import did not error")
		})
	}
}
