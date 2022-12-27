package main

import (
	"fmt"
	"reflect"
)

var a int = 654
// Go booleans are false unless decalred true
var b bool
// Floats can be either 32 or 64 bit
var c float64 = 2.6541
var d complex128 = 4 + 1i
var e string = "Australia" 

func addition(x int, y int) int {
	return x + y
}

func sayHello(str string) string {
	return "Hello " + str
}

// use reflect package to check the type. interface type can be when unknown
func printType(v interface{}) {
	fmt.Println(reflect.TypeOf(v))
}