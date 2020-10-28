package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	hello, err := ioutil.ReadFile("/root/hello.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hello))
}
