package main

import (
	"fmt"
	"log"
	"net/http"
)

type message struct {
	text string
}

func (m *message) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, m.text)
}

func main() {
	mh1 := &message{"Hello world From Go- Default Mux Type"}
	http.Handle("/", mh1)
	log.Fatal(http.ListenAndServe(":3004", nil))
}
