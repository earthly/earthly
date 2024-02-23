module github.com/earthly/earthly/examples/go-monorepo/services/one

go 1.21

require (
	github.com/earthly/earthly/examples/go-monorepo/libs/hello v0.0.0
	github.com/labstack/echo/v4 v4.9.0
)

replace github.com/earthly/earthly/examples/go-monorepo/libs/hello v0.0.0 => ../../libs/hello

require (
	github.com/labstack/gommon v0.3.1 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
