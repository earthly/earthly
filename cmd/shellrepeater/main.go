package main

import (
	"github.com/earthly/earthly/debugger/server"
)

const addr = "0.0.0.0:8373"

func main() {
	x := server.NewServer(addr)
	x.Start()
}
