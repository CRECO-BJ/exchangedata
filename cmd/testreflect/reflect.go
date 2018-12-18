package main

import (
	"fmt"
	"reflect"
)

type user struct {
	name string
	age  int
}

func main() {
	var a interface{}
	var b = &user{name: "testb", age: 2}

	a = b
	fmt.Println(reflect.TypeOf(a))
	fmt.Println(reflect.TypeOf(b))
	fmt.Println(reflect.ValueOf(a).Elem())
	fmt.Println(reflect.ValueOf(b).Elem())
}
