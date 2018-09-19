package main

import (
	"fmt"
	"time"
)

func main() {
	// a global variable that make a signal to the `goroutine`
	// to end itself.
	// Is stop thread safe ?
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("Received Signal to stop")
				return
			case <-time.After(10 * time.Second):
				// timeout
				fmt.Println("timing out")
			default:
				fmt.Println("Printing in 2 second")
				time.Sleep(2 * time.Second)
			}
		}
	}()
	time.Sleep(8 * time.Second)
	stop <- true
	time.Sleep(3 * time.Second)
	fmt.Println("program stopped")
}
