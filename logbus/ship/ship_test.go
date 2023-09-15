package ship

import (
	"errors"
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_retryable(t *testing.T) {
	cases := []struct {
		note string
		err  error
		want bool
	}{
		{
			note: "not status error",
			err:  errors.New("fail"),
			want: false,
		},
		{
			note: "unavailable status error",
			err:  status.Error(codes.Unavailable, "unavailable"),
			want: true,
		},
		{
			note: "unknown error",
			err:  status.Error(codes.Unknown, "unknown"),
			want: true,
		},
		{
			note: "wrapped unknown error",
			err:  fmt.Errorf("error: %w", status.Error(codes.Unknown, "unknown")),
			want: true,
		},
		{
			note: "wrapped non-status error",
			err:  fmt.Errorf("error: %w", errors.New("failed")),
			want: false,
		},
		{
			note: "double-wrapped unknown error",
			err:  fmt.Errorf("error: %w", fmt.Errorf("error: %w", status.Error(codes.Unknown, "unknown"))),
			want: true,
		},
	}

	for _, c := range cases {
		t.Run(c.note, func(t *testing.T) {
			got := retryable(c.err)
			if got != c.want {
				t.Errorf("wanted %+v, got %+v", c.want, got)
			}
		})
	}
}
