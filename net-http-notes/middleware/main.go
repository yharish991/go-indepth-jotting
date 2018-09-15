package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type message struct {
	text string
}

func (m *message) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, m.text)
}

// accept an http.Handler as an argument, and return an http.Handler.
// This makes it easy to chain middlewares
func loggingMiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer fmt.Printf("[%s] - Response Time: %v\n", r.URL.String(), time.Since(start))
		h.ServeHTTP(w, r)
	})
}

// Any Type can also be used to make an http Handler
// track the no of times this URL has been accessed.
type noOfTimes struct {
	next  http.Handler
	count int
}

// type needs to implement the ServeHTTP Method to be the Handler Type
func (t *noOfTimes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.count++
	fmt.Printf("counter = %d\n", t.count)
	fmt.Println("get", r.URL.Path, " from ", r.RemoteAddr)
	t.next.ServeHTTP(w, r)
}

func main() {
	mux := http.NewServeMux()
	mh1 := &message{"Hello world From Middleware- hanlder Type"}
	mux.Handle("/", mh1)
	mux.Handle("/count", &noOfTimes{next: mh1, count: 0})
	log.Fatal(http.ListenAndServe(":3003", loggingMiddleWare(mux)))
}
