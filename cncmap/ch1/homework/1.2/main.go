package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 10)
	go producer(ch)
	consumer(ch)

}

func producer(ch chan<- int) {
	ticker := time.NewTicker(time.Second)
	for i := 0; i < 10; i++ {
		<-ticker.C
		ch <- i
	}
	defer ticker.Stop()
	defer close(ch)
}

func consumer(ch <-chan int) {
	for item := range ch {
		fmt.Printf("recived %d\n", item)
	}
}
