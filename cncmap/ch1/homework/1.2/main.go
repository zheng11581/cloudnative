package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 10)

	go func() {
		ticker := time.NewTicker(time.Second)
		for i := 0; i < 10; i++ {
			<-ticker.C
			select {
			case ch <- i:
			}
		}
		defer close(ch)
		defer ticker.Stop()
	}()

	for item := range ch {
		fmt.Printf("recived %d\n", item)
	}
}
