package main

import (
	"fmt"
	"log"
	"net/http"
)

type pounds float32

func (p pounds) String() string {
	return fmt.Sprintf("₹%.2f", p)
}

type database map[string]pounds

func (d database) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for item, price := range d {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func main() {
	db := database{
		"foo": 1,
		"bar": 2,
	}

	log.Fatal(http.ListenAndServe("localhost:8000", db))
}
