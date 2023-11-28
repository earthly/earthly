package circbuf

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	_, err := NewBuffer(-1)
	require.Error(t, err)
}

func TestWriteGrow(t *testing.T) {
	b := &Buffer{maxSize: 5}
	n, err := io.WriteString(b, "foo")
	r := require.New(t)
	r.NoError(err)
	r.Equal(3, n)
	r.Len(b.data, 3)
	r.Equal("foo", string(b.Bytes()))
}

func TestWriteOverflow(t *testing.T) {
	b := &Buffer{maxSize: 5}
	n, err := io.WriteString(b, "foobarbaz")

	r := require.New(t)
	r.NoError(err)
	r.Equal(9, n)
	r.Equal("arbaz", string(b.data))
	r.Equal("arbaz", string(b.Bytes()))
}

func TestWriteMulti(t *testing.T) {
	b := &Buffer{maxSize: 5}
	r := require.New(t)

	n, err := io.WriteString(b, "mr")
	r.NoError(err)
	r.Equal(2, n)

	n, err = io.WriteString(b, "world")
	r.NoError(err)
	r.Equal(5, n)

	n, err = io.WriteString(b, "wide")
	r.NoError(err)
	r.Equal(4, n)

	r.Equal("edwid", string(b.data))
	r.Equal(1, b.offset)

	r.Equal("dwide", string(b.Bytes()))
}
