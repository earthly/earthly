package builder

import (
	"context"
	"testing"

	"github.com/earthly/earthly/cleanup"

	"github.com/stretchr/testify/assert"
)

// TestTempEarthlyOutDir tests that tempEarthlyOutDir always returns the same directory
func TestTempEarthlyOutDir(t *testing.T) {
	b, _ := NewBuilder(context.Background(), Opt{
		CleanCollection: cleanup.NewCollection(),
	})

	outDir1, err := b.tempEarthlyOutDir()
	assert.NoError(t, err)

	outDir2, err := b.tempEarthlyOutDir()
	assert.NoError(t, err)

	b.opt.CleanCollection.Close()

	assert.Equal(t, outDir1, outDir2)
}
