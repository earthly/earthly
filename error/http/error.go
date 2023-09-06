package http

import (
	"fmt"
	"net/http"
)

type Error struct {
	msg  string
	code int
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s, code: %d", e.msg, e.code)
}

func (e *Error) Code() int {
	return e.code
}

func New(code int, msg string) error {
	return &Error{
		msg:  msg,
		code: code,
	}
}

func Code(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if httErr, ok := err.(interface {
		Code() int
	}); ok {
		return httErr.Code()
	}
	return http.StatusInternalServerError
}
