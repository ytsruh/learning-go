package lessons

import (
	"fmt"
	"time"
)

func slowFunc(c chan string){
	time.Sleep(time.Second * 2)
	c <- "slowFunc() has finished"
}

func write(ch chan int) {
    for i := 1; i < 6; i++ {
        ch <- i
        fmt.Println("successfully wrote", i, "to ch")
    }
    close(ch)
}

func pinger(ch chan string) {
	t := time.NewTicker(1 * time.Second)
	for {
		ch <- "ping"
		<- t.C
	}
}

func sender(c chan string) {
	t := time.NewTicker(1 * time.Second)
	for {
		c <- "I'm sending a message"
		<- t.C
	}
}

func Channels()  {
	// Regular channel
	channel := make(chan string)
	go slowFunc(channel)
	msg := <- channel
	fmt.Println(msg)
	// Buffered channel
    ch := make(chan int, 2) // Creates capacity of 2
    go write(ch)
    time.Sleep(2 * time.Second)
    for v := range ch {
        fmt.Println("read value", v, "from ch")
        time.Sleep(2 * time.Second)
 
    }
	// Blocking channel
	// messages := make(chan string)
	// go pinger(messages)
	// for {
	// 	m := <-messages
	// 	fmt.Println(m)
	// }
	// Quitting a channel
	list := make(chan string)
	stop := make(chan bool)
	go sender(list)
	go func(){
		time.Sleep(time.Second * 2)
		fmt.Printf("Time is up")
		stop <- true
	}()
	for {
		select {
		case <- stop:
			return
		case newmsg := <- list:
			fmt.Println(newmsg)
		}
	}
}