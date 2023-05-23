package hasher_test

import (
	"context"
	"os"
	"testing"

	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/util/buildkitskipper/hasher"
)

var emptyHash = []byte{0xda, 0x39, 0xa3, 0xee, 0x5e, 0x6b, 0x4b, 0xd, 0x32, 0x55, 0xbf, 0xef, 0x95, 0x60, 0x18, 0x90, 0xaf, 0xd8, 0x7, 0x9}

func TestEmptyHasherIsNil(t *testing.T) {
	h := hasher.New()
	// empty string hash. e.g. running `true | sha1sum` in bash will output: "da39a3ee5e6b4b0d3255bfef95601890afd80709  -"
	Equal(t, h.GetHash(), emptyHash)
}

func TestNilHasherIsNil(t *testing.T) {
	var h *hasher.Hasher
	Nil(t, h.GetHash())
}

func TestHashCommand(t *testing.T) {
	h1 := hasher.New()
	h1.HashCommand(spec.Command{
		Name: "RUN",
		Args: []string{"ls", "/foo"},
	})
	hash1 := h1.GetHash()
	NotNil(t, hash1)
	NotEqual(t, hash1, emptyHash)

	h2 := hasher.New()
	h2.HashCommand(spec.Command{
		Name: "RUN",
		Args: []string{"ls", "/bar"},
	})
	hash2 := h2.GetHash()
	NotNil(t, hash2)
	NotEqual(t, hash2, emptyHash)

	NotEqual(t, hash1, hash2)
}

func TestHashEmptyFile(t *testing.T) {
	file, err := os.CreateTemp("", "file-to-hash")
	if err != nil {
		NoError(t, err)
	}
	defer os.Remove(file.Name())

	h := hasher.New()
	err = h.HashFile(context.Background(), file.Name())
	NoError(t, err)
	hash := h.GetHash()
	NotNil(t, hash)
	NotEqual(t, hash, emptyHash)
}

func TestHashFile(t *testing.T) {
	file, err := os.CreateTemp("", "file-to-hash")
	if err != nil {
		NoError(t, err)
	}
	defer os.Remove(file.Name())

	f, err := os.OpenFile(file.Name(), os.O_RDWR|os.O_TRUNC, 0)
	if err != nil {
		NoError(t, err)
	}
	_, err = f.Write([]byte("hello"))
	if err != nil {
		NoError(t, err)
	}
	err = f.Close()
	if err != nil {
		NoError(t, err)
	}

	h := hasher.New()
	err = h.HashFile(context.Background(), file.Name())
	NoError(t, err)
	hash := h.GetHash()
	NotNil(t, hash)
	NotEqual(t, hash, emptyHash)
}
