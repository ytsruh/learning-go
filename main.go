package main

import (
	"fmt"
	// Import local packages & assign a name to them
	types "learning/lessons/types"
)

func main()  {
	sum := types.Addition(7,8)
	fmt.Println(sum)
}