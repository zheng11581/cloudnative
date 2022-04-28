package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	baseCtx := context.Background()

	type favContextKey string
	f := func(ctx context.Context, k favContextKey) (v interface{}) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return v
		}
		fmt.Println("key not found:", k)
		return ""
	}
	k1 := favContextKey("language")
	k2 := favContextKey("author")

	// WithValue
	valueCtx := context.WithValue(baseCtx, k1, "Golang")
	fmt.Println(f(valueCtx, k1))
	fmt.Println(f(valueCtx, k2))

	// WithTimeOut
	timeoutCtx, cancel := context.WithTimeout(baseCtx, time.Second)
	defer cancel()
	go func(ctx context.Context) {
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <-ctx.Done():
				fmt.Println("child process interrupt...")
				return
			default:
				fmt.Println("enter default")
			}
		}
	}(timeoutCtx)
	select {
	case <-timeoutCtx.Done():
		time.Sleep(1 * time.Second)
		fmt.Println("main process exit!")
	}

}
