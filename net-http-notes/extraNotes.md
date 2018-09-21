Pipeline way

```Go
    type Pipeline struct {
        curr http.Handler
        next http.Handler
    }
    func (this Pipeline) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        curr(w,r)
        next(w,r)
    }

    func main(){
        ...
        p := Pipeline{LogHandler, IndexHandler}
        http.ListenAndServe(*addr, p)
    }
```

**not consider the ResponseWriter as a storage system**

`type RichHandler func(w http.ResponseWriter, req *http.Request) (int, error)`
**and make it implement the http.Handler interface if needed.**
