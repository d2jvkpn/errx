package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func main() {
	f := exampleFn

	fmt.Println("Function name:", getFunctionName(f))
	// f()

	var e1, e2 error

	e1 = NewErr1()
	e2 = NewErr2()
	fmt.Printf("==> e1: %t, %v\n", e1 == nil, e1)
	fmt.Printf("==> e2: %t, %v\n", e2 == nil, e2)
}

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

type D1 struct {
	Item string
}

func (self D1) Error() string {
	return fmt.Sprintf("D1 Erorr: %s", self.Item)
}

type D2 struct {
	Item string
}

func (self *D2) Error() string {
	return fmt.Sprintf("D2 Erorr: %s", self.Item)
}

func NewErr1() (err error) {
	var d *D1

	return d
}

func NewErr2() (err error) {
	var d *D2

	return d
}
