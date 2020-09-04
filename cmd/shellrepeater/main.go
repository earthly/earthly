package main

import (
	"github.com/earthly/earthly/debugger/server"
)

const addr = "0.0.0.0:5000"

func main() {
	x := server.NewServer(addr)
	x.Start()
}
