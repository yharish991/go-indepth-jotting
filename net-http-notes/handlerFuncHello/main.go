package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// if we don't use NewServeMux, Default ServeMux is Used.
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world From Go- HanldeFunc Type")
	})
	log.Fatal(http.ListenAndServe(":3002", mux))
}
