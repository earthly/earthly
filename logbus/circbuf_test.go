package logbus

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_circBuf_grow(t *testing.T) {
	b := &circBuf{maxSize: 5}
	n, err := io.WriteString(b, "foo")
	require.NoError(t, err)
	require.Equal(t, 3, n)
	require.Len(t, b.data, 3)
}

func Test_circBuf_overflow(t *testing.T) {
	b := &circBuf{maxSize: 5}
	n, err := io.WriteString(b, "foobarbaz")
	require.NoError(t, err)
	require.Equal(t, 9, n)
	require.Equal(t, "arbaz", string(b.data))
}
