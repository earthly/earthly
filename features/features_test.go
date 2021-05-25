package features

import (
	"testing"

	. "github.com/stretchr/testify/assert"
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
