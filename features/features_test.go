package features

import (
	"fmt"
	"testing"
)

func TestFeaturesStringEnabled(t *testing.T) {
	fts := &Features{
		Major:              0,
		Minor:              5,
		ReferencedSaveOnly: true,
	}
	s := fts.String()
	Equal(t, "VERSION --referenced-save-only 0.5", s)
}

func TestFeaturesStringDisabled(t *testing.T) {
	fts := &Features{
		Major:              1,
		Minor:              1,
		ReferencedSaveOnly: false,
	}
	s := fts.String()
	Equal(t, "VERSION 1.1", s)
}

func TestApplyFlagOverrides(t *testing.T) {
	fts := &Features{}
	err := ApplyFlagOverrides(fts, "referenced-save-only")
	Nil(t, err)
	Equal(t, true, fts.ReferencedSaveOnly)
	Equal(t, false, fts.UseCopyIncludePatterns)
	Equal(t, false, fts.ForIn)
	Equal(t, false, fts.RequireForceForUnsafeSaves)
	Equal(t, false, fts.NoImplicitIgnore)
}

func TestApplyFlagOverridesWithDashDashPrefix(t *testing.T) {
	fts := &Features{}
	err := ApplyFlagOverrides(fts, "--referenced-save-only")
	Nil(t, err)
	Equal(t, true, fts.ReferencedSaveOnly)
	Equal(t, false, fts.UseCopyIncludePatterns)
	Equal(t, false, fts.ForIn)
	Equal(t, false, fts.RequireForceForUnsafeSaves)
	Equal(t, false, fts.NoImplicitIgnore)
}

func TestApplyFlagOverridesMultipleFlags(t *testing.T) {
	fts := &Features{}
	err := ApplyFlagOverrides(fts, "referenced-save-only,use-copy-include-patterns,no-implicit-ignore")
	Nil(t, err)
	Equal(t, true, fts.ReferencedSaveOnly)
	Equal(t, true, fts.UseCopyIncludePatterns)
	Equal(t, false, fts.ForIn)
	Equal(t, false, fts.RequireForceForUnsafeSaves)
	Equal(t, true, fts.NoImplicitIgnore)
}

func TestApplyFlagOverridesEmptyString(t *testing.T) {
	fts := &Features{}
	err := ApplyFlagOverrides(fts, "")
	Nil(t, err)
	Equal(t, false, fts.ReferencedSaveOnly)
	Equal(t, false, fts.UseCopyIncludePatterns)
	Equal(t, false, fts.ForIn)
	Equal(t, false, fts.RequireForceForUnsafeSaves)
	Equal(t, false, fts.NoImplicitIgnore)
}

func TestVersionAtLeast(t *testing.T) {
	tests := []struct {
		earthlyVer Features
		major      int
		minor      int
		expected   bool
	}{
		{
			earthlyVer: Features{Major: 0, Minor: 6},
			major:      0,
			minor:      5,
			expected:   true,
		},
		{
			earthlyVer: Features{Major: 0, Minor: 6},
			major:      0,
			minor:      7,
			expected:   false,
		},
		{
			earthlyVer: Features{Major: 0, Minor: 6},
			major:      1,
			minor:      2,
			expected:   false,
		},
		{
			earthlyVer: Features{Major: 1, Minor: 2},
			major:      1,
			minor:      2,
			expected:   true,
		},
	}
	for _, test := range tests {
		title := fmt.Sprintf("earthly version %d.%d is at least %d.%d",
			test.earthlyVer.Major, test.earthlyVer.Minor, test.major, test.minor)
		t.Run(title, func(t *testing.T) {
			actual := versionAtLeast(test.earthlyVer, test.major, test.minor)
			Equal(t, test.expected, actual)
		})
	}
}
