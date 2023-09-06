package http

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "nil returns ok", args: args{err: nil}, want: http.StatusOK},
		{name: "http error returns the code", args: args{err: &Error{code: http.StatusNotExtended}}, want: http.StatusNotExtended},
		{name: "non http error returns 500", args: args{err: errors.New("not http")}, want: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Code(tt.args.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestError_Code(t *testing.T) {
	err := &Error{code: 123}
	want := 123
	got := err.Code()
	assert.Equal(t, want, got)
}

func TestError_Error(t *testing.T) {
	err := &Error{msg: "err message", code: 123}
	want := "err message, code: 123"
	got := err.Error()
	assert.Equal(t, want, got)
}

func TestNew(t *testing.T) {
	want := &Error{msg: "err message", code: 123}
	got := New(123, "err message")
	assert.Equal(t, want, got)
}
