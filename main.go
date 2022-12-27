package main

import (
	"fmt"
)

func getPrizes() (string, int){
	i := "goldfish"
	x := 3

	return i, x
}

// This function will accept any number of integers using the ... syntax
func sumNumbers(numbers ...int)int  {
	total := 0
	for _ , number :=range numbers{
		total += number
	}
	return total
}

// Example of a recursive function that calls itself until a condition is met
func feedMe(portion int, eaten int)int{
	eaten = portion + eaten
	if eaten >= 5 {
		fmt.Println("I'm full! I've eaten", eaten)
		return eaten
	}
	fmt.Println("I'm still hungry! I've eaten", eaten)
	return feedMe(portion, eaten)
}

func anotherFunction(f func() string) string{
	return f()
}

func main()  {
	fmt.Println(sumNumbers(10,20,33))
	fmt.Println(feedMe(1,2))
	fn := func(){
		fmt.Println("Function as value")
	}
	fn()
	x := func() string {
		return "This is returned from function called x"
	}

	fmt.Println(anotherFunction(x))
}