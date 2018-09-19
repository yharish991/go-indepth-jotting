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

	time.Sleep(1 * time.Second)
	fmt.Println("Cacelling all goroutines")
	cancel()
	time.Sleep(10 * time.Second)
}

func watch(ctx context.Context, name string) {
	go watchInner(ctx, name)
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "Cancelled")
			return
		default:
			fmt.Println(name, "goroutine running")

			time.Sleep(1 * time.Second)
		}
	}
}

func watchInner(ctx context.Context, name string) {
	dctx, _ := context.WithCancel(ctx)
	go watchInnerDerivedSubContext(dctx, name)
	// cancel()
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "Cancelled Inner")
			return
		default:
			fmt.Println(name, "Inner Goroutine running")
			time.Sleep(1 * time.Second)
		}
	}
}

func watchInnerDerivedSubContext(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "Cancelled DerivedSubContext")
			return
		default:
			fmt.Println(name, "DerivedSubContext Goroutine running")
			time.Sleep(10 * time.Second)
		}
	}
}
