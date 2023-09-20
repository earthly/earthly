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

func Test_prefixFormatter_Format(t *testing.T) {
	t.Run("uses cache correctly", func(t *testing.T) {
		random := uuid.NewString()
		otherRandom := uuid.NewString()
		origSize := getCacheSize(&formatter.cache)

		formatter.Format(random, DefaultPadding)
		size := getCacheSize(&formatter.cache)
		assert.Equal(t, origSize+1, size, "cache size should have incremented by 1")
		formatter.Format(random, DefaultPadding)
		assert.Equal(t, origSize+1, size, "cache size should have stayed the same")
		formatter.Format(random, 3)
		size = getCacheSize(&formatter.cache)
		assert.Equal(t, origSize+2, size, "cache size should have incremented since padding is different")
		formatter.Format(otherRandom, DefaultPadding)
		size = getCacheSize(&formatter.cache)
		assert.Equal(t, origSize+3, size, "cache size should have incremented since prefix is different")
	})
}
