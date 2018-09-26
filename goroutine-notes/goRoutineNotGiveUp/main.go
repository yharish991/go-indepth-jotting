package main

import (
	"fmt"
	"runtime"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		fmt.Println(s)
	}
}
func main() {
	// Go using a single core.
	// two goroutines are in one thread
	runtime.GOMAXPROCS(1)
	// all goroutines will run in a native thread, that is, only one CPU core is used.
	// but When a goroutine blocks, such as by calling a blocking system call, the
	// run-time automatically moves other goroutines on the same operating system
	// thread to a different, runnable thread so they won't be blocked
	go say("hello") // never gets a chance to run
	//  infinite loop, it takes up all the resources of the single-core CPU
	for {

	}
}
