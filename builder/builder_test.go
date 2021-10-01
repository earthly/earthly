package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/earthly/earthly/cleanup"
)

// TestTempEarthlyOutDir tests that tempEarthlyOutDir always returns the same directory
func TestTempEarthlyOutDir(t *testing.T) {
	b, _ := NewBuilder(nil, Opt{
		CleanCollection: cleanup.NewCollection(),
	})

	outDir1, err := b.tempEarthlyOutDir()
	assert.NoError(t, err)

	outDir2, err := b.tempEarthlyOutDir()
	assert.NoError(t, err)

	b.opt.CleanCollection.Close()

	assert.Equal(t, outDir1, outDir2)
}
