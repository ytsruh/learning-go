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

// use reflect package to check the type
func checkType (v any) {
	fmt.Println(reflect.TypeOf(v))
}

func main()  {
	// fmt.Println(sayHello("world"))
	// fmt.Println(addition(1,4))
	// fmt.Printf("d for Integer: %d\n", a)
	// fmt.Printf("6d for Integer: %6d\n", a)
	// fmt.Printf("t for Boolean: %t\n", b)
	// fmt.Printf("g for Float: %g\n", c)
	// fmt.Printf("e for Scientific Notation: %e\n", d)
	// fmt.Printf("E for Scientific Notation: %E\n", d)
	// fmt.Printf("s for String: %s\n", e)
	// fmt.Printf("G for Complex: %G\n", c)
	// fmt.Printf("15s String: %15s\n", e)
	// fmt.Printf("-10s String: %-10s\n",e)
	t:= fmt.Sprintf("Print from right: %[3]d %[2]d %[1]d", 11, 22, 33)
	fmt.Println(t)	
	checkType(a)
}