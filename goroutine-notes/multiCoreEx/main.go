package main

import (
	"fmt"
	"runtime"
)

var quit = make(chan int)

func loop(id int) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", id)
	}
	quit <- 0
}

func main() {
	runtime.GOMAXPROCS(2)
	for i := 0; i < 3; i++ {
		go loop(i) // go routine id
	}

	for i := 0; i < 3; i++ {
		<-quit
	}
}

/*
Preemptive output sometimes occurs
(indicating that Go has opened more than one native thread and achieved true parallelism)
0 1 1 0 0 0 0 1 1 1 1 1 0 0 0 0 0 1 1 1 2 2 2 2 2 2 2 2 2 2

// Sometimes it will output sequentially,
// print 0 and then print 1 and then print 2
// (indicating that Go opens a native thread,
// the goroutine on a single thread does not block and does not release the CPU)

2 2 2 2 2 2 2 2 2 2 0 0 0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1 1 1

2 2 2 2 2 2 2 2 2 2 1 1 1 1 1 1 1 1 1 1 0 0 0 0 0 0 0 0 0 0
Preemptive output
0 1 1 1 1 1 1 1 1 1 1 0 0 0 0 0 0 0 0 0 2 2 2 2 2 2 2 2 2 2
*/
