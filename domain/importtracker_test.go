package domain

import (
	"testing"

	"github.com/earthly/earthly/conslogging"

	"github.com/stretchr/testify/assert"
)

func TestImports(t *testing.T) {
	var tests = []struct {
		importStr string
		as        string
		ref       string
		expected  string
		ok        bool
	}{
		{"github.com/foo/bar", "", "bar+abc", "github.com/foo/bar+abc", true},
		{"github.com/foo/bar", "buz", "buz+abc", "github.com/foo/bar+abc", true},
		{"github.com/foo/bar", "buz", "bar+abc", "", false},
		{"github.com/foo/bar:v1.2.3", "", "bar+abc", "github.com/foo/bar:v1.2.3+abc", true},
		{"github.com/foo/bar:v1.2.3", "buz", "buz+abc", "github.com/foo/bar:v1.2.3+abc", true},
		{"github.com/foo/bar:v1.2.3", "buz", "bar+abc", "", false},
		{"./foo/bar", "", "bar+abc", "./foo/bar+abc", true},
		{"./foo/bar", "buz", "buz+abc", "./foo/bar+abc", true},
		{"./foo/bar", "buz", "bar+abc", "", false},
		{"../foo/bar", "", "bar+abc", "../foo/bar+abc", true},
		{"../foo/bar", "buz", "buz+abc", "../foo/bar+abc", true},
		{"../foo/bar", "buz", "bar+abc", "", false},
		{"/foo/bar", "", "bar+abc", "/foo/bar+abc", true},
		{"/foo/bar", "buz", "buz+abc", "/foo/bar+abc", true},
		{"/foo/bar", "buz", "bar+abc", "", false},
	}

	var console conslogging.ConsoleLogger

	for _, tt := range tests {
		ir := NewImportTracker(console, nil)
		err := ir.Add(tt.importStr, tt.as, false, false, false)
		assert.NoError(t, err, "add import error")

		ref, err := ParseTarget(tt.ref)
		assert.NoError(t, err, "parse test case ref") // check that the test data is good
		assert.Equal(t, tt.ref, ref.String())         // sanity check

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
