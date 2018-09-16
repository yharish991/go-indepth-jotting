package main

import (
	"fmt"
	"time"
)

func main() {
	a := 0
	for i := 0; i < 10000; i++ {
		go func(id int) {
			a++
			fmt.Println("Loop", id, "a", a)
		}(i)
	}

	time.Sleep(time.Second)
	fmt.Println("Last Value of a", a)
}
