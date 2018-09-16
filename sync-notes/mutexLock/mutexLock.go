package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	a := 0
	var mu sync.Mutex
	for i := 0; i < 10000; i++ {
		go func(id int) {
			mu.Lock()
			defer mu.Unlock()
			a++
			fmt.Println("Loop", id, "a", a)
		}(i)
	}

	time.Sleep(time.Second)
	fmt.Println("Last Value of a", a)
}
