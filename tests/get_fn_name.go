package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func exampleFn() {
	fmt.Println("This is an example function.")
}

func getFunctionName(i interface{}) string {
	pc := reflect.ValueOf(i).Pointer()
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	name := fn.Name()

	return name[strings.LastIndex(name, ".")+1:]
}

func main() {
	f := exampleFn

	fmt.Println("Function name:", getFunctionName(f))
	// f()
}
