package conslogging

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		f := NewPrefixFormatter(truncateSha)
		require.Zero(t, getCacheSize(&f.cache))
		formatter.Format(random, DefaultPadding)
		size := getCacheSize(&formatter.cache)
		assert.Equal(t, 1, size, "cache size should have incremented by 1")
		formatter.Format(random, DefaultPadding)
		size = getCacheSize(&formatter.cache)
		assert.Equal(t, 1, size, "cache size should have stayed the same")
		formatter.Format(random, 3)
		size = getCacheSize(&formatter.cache)
		assert.Equal(t, 2, size, "cache size should have incremented since padding is different")
		formatter.Format(otherRandom, DefaultPadding)
		size = getCacheSize(&formatter.cache)
		assert.Equal(t, 3, size, "cache size should have incremented since prefix is different")
	})
	t.Run("executes all options", func(t *testing.T) {
		optsCallNum := 0
		prefix := "123"
		expectedLen := len(prefix)

		optFunc := func(add string) func(str string, padding int, curLen int) string {
			return func(str string, padding int, curLen int) string {
				optsCallNum++
				assert.Equal(t, expectedLen, curLen)
				expectedLen++
				return str + add
			}
		}
		f := NewPrefixFormatter(optFunc("4"), optFunc("5"))
		newPrefix := f.Format(prefix, 1)
		assert.Equal(t, "12345", newPrefix)
		assert.Equal(t, 2, optsCallNum)
	})
	t.Run("does not execute options when no padding", func(t *testing.T) {
		optsCallNum := 0
		prefix := "123"
		expectedLen := len(prefix)

		optFunc := func(add string) func(str string, padding int, curLen int) string {
			return func(str string, padding int, curLen int) string {
				optsCallNum++
				assert.Equal(t, expectedLen, curLen)
				expectedLen++
				return str + add
			}
		}
		f := NewPrefixFormatter(optFunc("4"), optFunc("5"))
		newPrefix := f.Format(prefix, NoPadding)
		assert.Equal(t, prefix, newPrefix)
		assert.Equal(t, 0, optsCallNum)
	})
}
