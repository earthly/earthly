package flagutil

import (
	"github.com/dustin/go-humanize"
)

type ByteSizeValue uint64

func (b *ByteSizeValue) Set(s string) error {
	v, err := humanize.ParseBytes(s)
	if err != nil {
		return err
	}
	*b = ByteSizeValue(v)
	return nil
}

func (b *ByteSizeValue) String() string { return humanize.Bytes(uint64(*b)) }
