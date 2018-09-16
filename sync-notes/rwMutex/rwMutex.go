package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.RWMutex
	a := 0

	for i := 0; i < 10; i++ {
		go func(id int) {

			mu.RLock()
			fmt.Println("Acqured Read Lock in", id)
			defer mu.RUnlock()
			// try to acquire the write Lock at i = 50000
			if id == 5 {
				go func(id int) {
					mu.Lock()
					fmt.Println("Acqured Write Lock in", id)
					a++
					fmt.Println(a)
					time.Sleep(time.Second * 5)
					mu.Unlock()
				}(id)
			}
			fmt.Println(a)
		}(i)
	}

	time.Sleep(time.Second * 1)
	fmt.Println(a)
}
