package main

import (
	"log"
	"net/http"
)

// net/http package provides serveral Function that implement
// the http.Handler interface

func main() {
	mux := http.NewServeMux()
	// FileServer function implements the http.Handler interface
	// https://golang.org/src/net/http/fs.go?s=20537:20577#L703
	// FileServer returns a type fileHandler,
	// which implements the http.Handler interface

	// func (f *fileHandler) ServeHTTP(w ResponseWriter, r *Request) {
	// 	upath := r.URL.Path
	// 	if !strings.HasPrefix(upath, "/") {
	// 		upath = "/" + upath
	// 		r.URL.Path = upath
	// 	}
	// 	serveFile(w, r, f.root, path.Clean(upath), true)
	// }
	file := http.FileServer(http.Dir("public"))
	mux.Handle("/", file)
	log.Fatal(http.ListenAndServe(":3001", mux))
}
