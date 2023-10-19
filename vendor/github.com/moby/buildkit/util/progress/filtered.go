package progress

import (
	"context"
)

// this entire file is earthly-specific, it implements a filtered Reader which can be used to skip statsStream data

type FilteredReaderSkipFn func(ctx context.Context, p *Progress) (bool, error)

type filteredReader struct {
	r      Reader
	skipFn FilteredReaderSkipFn
}

func NewFilteredReader(r Reader, skipFn FilteredReaderSkipFn) Reader {
	return &filteredReader{
		r:      r,
		skipFn: skipFn,
	}
}

func (fr *filteredReader) Read(ctx context.Context) ([]*Progress, error) {
	progress, err := fr.r.Read(ctx)
	if err != nil {
		return nil, err
	}
	filteredProgress := []*Progress{}
	for _, p := range progress {
		shouldSkip, err := fr.skipFn(ctx, p)
		if err != nil {
			return nil, err
		}
		if shouldSkip {
			continue
		}
		filteredProgress = append(filteredProgress, p)
	}
	return filteredProgress, nil
}
