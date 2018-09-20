package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// stop is read only channel for the consumer.
func consumer(stop <-chan bool) {
	for {
		select {
		case <-stop:
			fmt.Println("exit the consumer go routine")
			return
		default:
			fmt.Println("Running the consumer go routine")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	stop := make(chan bool)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(stop <-chan bool) {
			defer wg.Done()
			consumer(stop)
		}(stop)
	}
	waitForKill()
	close(stop)
	fmt.Println("Stopping all Jobs!")
	wg.Wait()
}

func waitForKill() {
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt)
	signal.Notify(sign, syscall.SIGTERM)
	// block and wait as reading channel is synchronous
	// and blocking
	<-sign
}
