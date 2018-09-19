package main

import (
	"context"
	"fmt"
	"time"
)

// multiple go Routine context
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go watch(ctx, "[1]")
	go watch(ctx, "[2]")
	go watch(ctx, "[3]")

	time.Sleep(10 * time.Second)
	fmt.Println("Cacelling all goroutines")
	cancel()
	time.Sleep(5 * time.Second)
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "Cancelled")
			return
		default:
			fmt.Println(name, "goroutine running")
			time.Sleep(2 * time.Second)
		}
	}
}
