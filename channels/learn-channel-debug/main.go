package main

import "fmt"

func goroutineA(c1 chan int, c2 chan int) {
	_ = <-c2
	c1 <- 1
	return
}

func goroutineB(c1 chan int, c2 chan int) {
	fmt.Println("c1", <-c1)
	c2 <- 2
	return
}

func goroutineC(c1 chan int, c2 chan int) {
	_ = <-c2
	c1 <- 1
	return
}

func goroutineD(c1 chan int, c2 chan int) {
	_ = <-c2
	c1 <- 1
	return
}

func goroutineE(c1 chan int, c2 chan int) {
	_ = <-c2
	c1 <- 1
	return
}

func main() {
	c1 := make(chan int, 3)
	c2 := make(chan int)
	//c1 <- 2
	// c2 <- 2
	go goroutineA(c1, c2)
	go goroutineC(c1, c2)
	go goroutineD(c1, c2)
	go goroutineE(c1, c2)
	c2 <- 2
	go goroutineB(c1, c2)

	for {
	}
}
