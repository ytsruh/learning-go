package main

import (
	"fmt"
)

func flow()  {
	a := 5
	// if-else & else if
	if a == 3 {
		fmt.Println("Variable is equal to 3")
	} else if a >= 3{
		fmt.Println("Variable is greater than 3")
	} else {
		fmt.Println("Else statement")
	}

	switch a {
		case 2:
			fmt.Println("Variable is equal to 2")
		case 3:
			fmt.Println("Variable is equal to 3")
		case 4:
			fmt.Println("Variable is equal to 4")
		default:
			fmt.Println("This is the default case")
	}

	for i := 0; i < a; i++ {
		fmt.Println(i)
	}

	numbers := []int{1,2,3,4}
	for i, v := range numbers {
		fmt.Println("The index of the array is : ",i)
		fmt.Println("The value from the array is : ",v)
	}
	// Defer keyword can be used to run after function has completed. Multiple defer have odd behaviour
	defer fmt.Println("...or is it?")
	defer fmt.Println("...maybe...")
	fmt.Println("End of the exercise")
}