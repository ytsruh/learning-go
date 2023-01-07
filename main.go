package main

import (
	"fmt"
	// Import local packages & assign a name to them
	"learning/lessons"
	types "learning/lessons/types"
)

func main()  {
	sum := types.Addition(7,8)
	fmt.Println(sum)
	lessons.Panic("efdcs")
}