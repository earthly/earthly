package hasher

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"os"
)

type Hasher struct {
	h hash.Hash
}

func New() *Hasher {
	return &Hasher{
		h: sha1.New(),
	}
}

func (h *Hasher) GetHash() []byte {
	if h == nil {
		return nil
	}
	return h.h.Sum(nil)
}

func (h *Hasher) HashInt(i int) {
	h.HashString(fmt.Sprintf("int:%d", i))
}

func (h *Hasher) HashJSONMarshalled(v any) {
	dt, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to hash command: %s", err)) // shouldn't happen
	}
	h.HashBytes(dt)
}

func (h *Hasher) HashString(s string) {
	h.HashBytes([]byte(fmt.Sprintf("str:%s", s)))
}

func (h *Hasher) HashBytes(b []byte) {
	h.h.Write([]byte(fmt.Sprintf("%d", len(b))))
	h.h.Write(b)
}

func (h *Hasher) HashFile(ctx context.Context, src string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}
	h.HashString(fmt.Sprintf("mode: %d;", stat.Mode()))
	h.HashString(fmt.Sprintf("size: %d;", stat.Size()))

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	readCh := make(chan error)
	buf := make([]byte, 1024*1024)
	r := bufio.NewReader(f)
	for {
		var n int
		go func() {
			var err error
			n, err = r.Read(buf)
			readCh <- err
		}()
		select {
		case err := <-readCh:
			if err == io.EOF {
				return nil
			} else if err != nil {
				return err
			}
			h.h.Write(buf[:n])
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
