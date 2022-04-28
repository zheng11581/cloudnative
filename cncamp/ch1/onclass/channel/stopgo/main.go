package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				fmt.Print("goroutine is down...")
				return

			}
		}
	}()
	close(done)
	time.Sleep(10 * time.Second)

}
