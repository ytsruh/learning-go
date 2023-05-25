package lessons

import (
	"fmt"
	"time"
)

func Now() {
	fmt.Println(time.Now())
}

func Sleep() {
	time.Sleep(3 * time.Second)
	fmt.Println("I'm awake!")
}

func Timeout() {
	fmt.Println("You have two seconds to calculate 19 * 4")
	for {
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("Time's up! The answer is 74.")
			return
		}
	}
}

func Ticker() {
	c := time.Tick(5 * time.Second)
	for t := range c {
		fmt.Printf("The time is now %v\n", t)
	}
}
