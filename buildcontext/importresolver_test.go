package buildcontext

import (
	"testing"

	"github.com/earthly/earthly/domain"
	. "github.com/stretchr/testify/assert"
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

	for _, tt := range tests {
		ir := NewImportResolver(nil, nil)
		err := ir.AddImport(tt.importStr, tt.as, false)
		NoError(t, err, "add import error")

		ref, err := domain.ParseTarget(tt.ref)
		NoError(t, err, "parse test case ref") // check that the test data is good
		Equal(t, tt.ref, ref.String())         // sanity check

		ref2, err := ir.DerefImport(ref)
		if tt.ok {
			NoError(t, err, "deref import")
			Equal(t, tt.expected, ref2.StringCanonical()) // StringCanonical shows its resolved form
			Equal(t, tt.ref, ref2.String())               // String shows its import form
		} else {
			Error(t, err, "deref should have error'd")
		}
	}
}
