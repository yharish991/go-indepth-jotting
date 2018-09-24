package main

import (
	"fmt"
	"time"
)

func goRoutineA(a <-chan int) {
	val := <-a
	fmt.Println("goRoutineA received the data", val)
}

func goRoutineB(b <-chan int) {
	val := <-b
	fmt.Println("goRoutineB received the data", val)
}

func main() {
	ch := make(chan int)
	go goRoutineA(ch)
	go goRoutineB(ch)
	ch <- 3
	//close(ch)
	time.Sleep(time.Second * 1)
}
