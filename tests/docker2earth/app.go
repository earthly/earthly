package main

import (
	"fmt"
	"os"
)

func main() {
	hello, err := os.ReadFile("/root/hello.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hello))
}
