package main

import "fmt"

func main() {
	ch := make(chan int, 5)

	select {
	case msg := <-ch:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}

}
