package main

import (
	"fmt"
	"github.com/spance/go-callprivate/private"
	"os"
	"reflect"
	"runtime"
)

type obj int

func (o obj) private() {
	fmt.Println("private func LOL.", runtime.GOOS, runtime.GOARCH)
}

func main() {
	// example 1
	var o obj
	method := reflect.ValueOf(o).MethodByName("private")
	private.SetAccessible(method)
	method.Call(nil)

	// example 2
	f, _ := os.Open(".")
	isdir := reflect.ValueOf(f).MethodByName("isdir")
	private.SetAccessible(isdir)
	fmt.Println(isdir.Call(nil)[0].Interface())
}
