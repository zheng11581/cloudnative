package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- int, load int) {
	//time.Sleep(time.Second)
	ch <- load
	time.Sleep(time.Second)
	fmt.Printf("Multi-Produce a load: %d\n", load)
}

func consumer(ch <-chan int) {
	if load, notClosed := <-ch; notClosed {
		fmt.Printf("Multi-Consume a load: %d\n", load)
	} else {
		fmt.Println("Channel Closed")
	}

}

func main() {
	ch := make(chan int, 10)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for i := 0; i < 1000; i++ {
		go producer(ch, i)
		//go consumer(ch)
	}
	for {
		select {
		case <-ticker.C:
			fmt.Println("Timeout waiting from channel ch")
		case v := <-ch:
			fmt.Printf("Receive %d from ch\n", v)

		}
	}

}
