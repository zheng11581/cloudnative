package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 10)

	// producer
	go func() {
		ticker := time.NewTicker(time.Second)
		for i := 0; i < 10; i++ {
			<-ticker.C
			ch <- i
		}
		defer ticker.Stop()
		defer close(ch)
	}()

	// consumer
	for item := range ch {
		fmt.Printf("recived %d\n", item)
	}

}
