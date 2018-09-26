package main

import (
	"fmt"
	"time"
)

type user struct {
	name string
	age  int8
}

var u = user{name: "Jon Deo", age: 100}
var g = &u

func modifyUser(pu *user) {
	fmt.Println("modifyUser Received Value", pu)
	pu.name = "Anand"
}
func printUser(u <-chan *user) {
	time.Sleep(3 * time.Second)
	fmt.Println("printUser goRoutine called", <-u)
}

func main() {
	c := make(chan *user)
	go printUser(c)
	go func(g *user) {
		c <- g
	}(g)
	fmt.Println(g)
	// sleep for sometime to let the anonymous goroutine run
	time.Sleep(1 * time.Second)
	// modify g
	g = &user{name: "Ankur", age: 10}
	go modifyUser(g)
	time.Sleep(3 * time.Second)
	fmt.Println(g)

}
