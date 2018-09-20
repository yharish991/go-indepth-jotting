package main

import (
	"fmt"
	"sync"
	"time"
)

// What are some negative points ?
// 1 . use of global function variable, wg of main inside various go routine.
func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("1st goroutine done")
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("2nd goroutine done")
		wg.Done()
	}()
	// wait for all the unfinished goroutine
	wg.Wait()
	fmt.Println("All goroutine executed")
}
