package main

import (
	"fmt"
)

func shorthandVariables() {
	// Type 1 - typically used outside of functions
	var s,t = "hello" ,"world"
	// Type 2 - typically used outside of functions
	var ( 
		x = "testing new"
		z = "variables"
	)
	// Type 3 - typically used inside of functions
	u := "new type of variables"
	//Run main
	fmt.Println(s + " " + t)
	fmt.Println(x + " " + z)
	fmt.Println(u)
}

var glob string = "Globally scoped variable"

func lexicalScope(){
	fmt.Println("Print global variable : " + glob)
	a := true
	if a {
		fmt.Println("Printing 'a' variable from outer block : ", a)
		i:= 673
		if a != false {
			fmt.Println("Printing 'i' variable from outer block : ", i)
		}
	}
}

func pointer(){
	x := "test variable"
	// Use & before a variable to get the pointer or the location variable is held in memory
	fmt.Println(&x)
	// Variables passed into a function will have a different memory location than the original as a new variable is created
}



func main()  {
	shorthandVariables()
	lexicalScope()
	pointer()
}