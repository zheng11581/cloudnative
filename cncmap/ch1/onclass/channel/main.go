package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 10)
	defer close(ch)

	// consumer
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			default:
				fmt.Printf("send: %d\n", <-ch)
			}
		}
		defer ticker.Stop()
	}()

	// producer
	for i := 0; i < 10; i++ {
		ch <- i
	}
	time.Sleep(11 * time.Second)

}
