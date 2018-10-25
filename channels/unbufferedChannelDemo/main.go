package main

func goroutineA(c2 chan int) {
	c2 <- 2
}

func main() {
	c2 := make(chan int)
	go goroutineA(c2)

	for {
	}
}
