# go-callprivate

Golang is a pretty useful language, but I detest the access policy of private (unexported) fields and functions in Golang, because that limits restricts flexibility of representing our ideas freely. It's very easy to access private field with memory address offset, then the go-callprivate library could give you freedom to call private methods.

Support: Tested on windows/linux/osx/freebsd

Examples:

```go
type obj int

func (o obj) private() {
	fmt.Println("private func LOL.", runtime.GOOS, runtime.GOARCH)
}

func main() {
	// example 1
	var o obj
	method := reflect.ValueOf(o).MethodByName("private")
	private.SetAccessible(method)
	method.Call(nil) // stdout...
	
	// example 2
	f, _ := os.Open(".")
	isdir := reflect.ValueOf(f).MethodByName("isdir")
	private.SetAccessible(isdir)
	fmt.Println(isdir.Call(nil)[0].Interface()) // stdout true
}
```