package lessons

import (
	"fmt"
)

func GetPrizes() (string, int){
	i := "goldfish"
	x := 3

	return i, x
}

// This function will accept any number of integers using the ... syntax
func SumNumbers(numbers ...int)int  {
	total := 0
	for _ , number :=range numbers{
		total += number
	}
	return total
}

// Example of a recursive function that calls itself until a condition is met
func FeedMe(portion int, eaten int)int{
	eaten = portion + eaten
	if eaten >= 5 {
		fmt.Println("I'm full! I've eaten", eaten)
		return eaten
	}
	fmt.Println("I'm still hungry! I've eaten", eaten)
	return FeedMe(portion, eaten)
}

func AnotherFunction(f func() string) string{
	return f()
}