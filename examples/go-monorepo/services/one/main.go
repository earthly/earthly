package main

import (
	"net/http"

	"github.com/earthly/earthly/examples/go-monorepo/libs/hello"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/one/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, hello.Greet("World"))
	})
	_ = e.Start(":8080")
}
