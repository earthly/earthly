package conslogging

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func getCacheSize(m *sync.Map) int {
	size := 0
	m.Range(func(key, value any) bool {
		size++
		return true
	})
	return size
}

func Test_prefixBeautifier_Beautify(t *testing.T) {
	t.Run("uses cache correctly", func(t *testing.T) {
		random := uuid.NewString()
		otherRandom := uuid.NewString()
		origSize := getCacheSize(&beautifier.cache)

		beautifier.Beautify(random, DefaultPadding)
		size := getCacheSize(&beautifier.cache)
		assert.Equal(t, origSize+1, size, "cache size should have incremented by 1")
		beautifier.Beautify(random, DefaultPadding)
		assert.Equal(t, origSize+1, size, "cache size should have stayed the same")
		beautifier.Beautify(random, 3)
		size = getCacheSize(&beautifier.cache)
		assert.Equal(t, origSize+2, size, "cache size should have incremented since padding is different")
		beautifier.Beautify(otherRandom, DefaultPadding)
		size = getCacheSize(&beautifier.cache)
		assert.Equal(t, origSize+3, size, "cache size should have incremented since prefix is different")
	})
}
