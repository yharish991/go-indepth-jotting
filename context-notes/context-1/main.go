package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context has been cancelled")
				return
			default:
				fmt.Println("go-routine is executing with proper time")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx)
	time.Sleep(8 * time.Second)
	fmt.Println("Cancelling the Context")
	cancel()
	time.Sleep(2 * time.Second)
	fmt.Println("program exit")
}
