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
	mux := http.NewServeMux()
	mh1 := &message{"Hello world From Go- hanlder Type"}
	mux.Handle("/", mh1)
	log.Fatal(http.ListenAndServe(":3003", mux))
}
