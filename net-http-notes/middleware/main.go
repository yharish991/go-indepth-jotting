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
	fmt.Println(r.Header)
	fmt.Fprintf(w, m.text)
}

// middleware that modifes the request
// https://godoc.org/net/http#Handler
// Clearly Says
// Except for reading the body, handlers should not modify the provided Request.
// we want to set a USER ID header on every request, for func purpose
func userIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r2 := new(http.Request)
		*r2 = *r
		r2.Header.Set("USER-ID", "1234567")
		next.ServeHTTP(w, r2)
	})
}

// middleware that writes Header in the Response
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-id", "1234")
		next.ServeHTTP(w, r)
	})
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

// type logHandler func(w http.ResponseWriter, req *http.Request) (int, error)

// func log(h http.Hander) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// }
func main() {
	mux := http.NewServeMux()
	mh1 := &message{"Hello world From Middleware- hanlder Type"}
	mux.Handle("/", userIDMiddleware(mh1))
	mux.Handle("/count", &noOfTimes{next: mh1, count: 0})
	mux.Handle("/id", requestIDMiddleware(mh1))
	log.Fatal(http.ListenAndServe(":3003", loggingMiddleWare(mux)))
}
