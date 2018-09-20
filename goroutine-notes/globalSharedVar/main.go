package main

import (
	"fmt"
	"time"
)

func main() {
	keepRunning := true

	goF := func() {
		for keepRunning {
			fmt.Println("Gang of goroutine running")
			time.Sleep(1 * time.Second)
		}
		fmt.Println("Gang of goroutine Exiting :(")
	}
	go goF()
	go goF()
	go goF()
	time.Sleep(2 * time.Second)
	keepRunning = false
	time.Sleep(3 * time.Second)
	fmt.Println("main goroutine exit")
}
