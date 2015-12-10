package main

import (
	"fmt"
	"github.com/spance/go-callprivate/private"
	"net/http"
	"reflect"
	"runtime"
)

type obj int

func (o *obj) private() {
	fmt.Println("private func LOL.", runtime.GOOS, runtime.GOARCH)
}

func main() {
	// example 1: Call *obj.private()
	var o obj
	method := reflect.ValueOf(&o).MethodByName("private")
	private.SetAccessible(method)
	method.Call(nil) // stdout ...

	// example 2: Call http.Header.clone()
	var h = http.Header{"k": {"v"}}
	clone := reflect.ValueOf(h).MethodByName("clone")
	private.SetAccessible(clone)
	fmt.Println(clone.Call(nil)[0]) // stdout map[k:[v]]
}
