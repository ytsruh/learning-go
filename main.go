package main

import (
	// Import local packages & assign a name to them
	"fmt"
	"learning/lessons"
	types "learning/lessons/types"
	"time"
)

func slowFunc(c chan string){
	time.Sleep(time.Second * 5)
	c <- "slowFunc() has finished"
}


func main()  {
	// types is imported as with its own declaration
	sum := types.Addition(7,8)
	fmt.Println(sum)
	lessons.Channels()
}