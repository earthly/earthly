package buildkitskipper

import (
	"context"
)

func NewNoSkip() *NoSkip {
	return &NoSkip{}
}

type NoSkip struct{}

func (no *NoSkip) Add(ctx context.Context, data []byte) error {
	return nil
}

func (no *NoSkip) Exists(ctx context.Context, data []byte) (bool, error) {
	return false, nil
}
