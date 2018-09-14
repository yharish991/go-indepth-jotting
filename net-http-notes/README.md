## net/http Package In-Depth Learning Notes

### Basic

> Major Components for processing HTTP requests.

- ServeMux
- Handler

### Handler

`http.Handler` Interface - all requested processors

```Go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

**`ServeHTTP` should write reply headers and data to the ResponseWriter and then return**

The [`http.ResponseWriter`](https://golang.org/pkg/net/http/#ResponseWriter) is also an interface.

`http.ResponseWriter` Interface

```Go
type ResponseWriter interface {
    Header() Header
    Write([]byte) (int, error)
    WriteHeader(statusCode int)
}
```

### ServeMux

[Doc](https://golang.org/pkg/net/http/#ServeMux)

- ServeMux is an HTTP request multiplexer.
- It matches the URL of each incoming request against a list of registered patterns and **calls the handler** for the pattern that most closely matches the URL.
- The http package has a package level variable DefaultServeMux, indicating the default route: var DefaultServeMux = NewServeMux(), which is registered to the route when registering the processor using the package-level `http.Handle() and http.HandleFunc()` methods.**It poses a security risk** as is stored in a global variable, any package is able to access it and register a route – including any third-party packages that your application imports.

`http.ServeMux`

```Go
type ServeMux struct {
}
```

> Methods of `ServeMux` receiver type.

```Go
func (mux *ServeMux) Handle(pattern string, handler Handler)
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)
```

ServeMux Also Implements `http.Handler` interface

```Go
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)
```

ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL.

`HandleFunc` function takes second agruments a function with following signature,
`func(ResponseWriter, *Request)`.

How does this function becomes `handler` that implements `http.Handler` Interface ?

Example handleFunc1:

```Go
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world From Go- HanldeFunc Type")
	})
	log.Fatal(http.ListenAndServe(":3002", mux))
}
```

As `ServeMux`

matches the URL of each incoming request against a list of registered patterns and **calls the handler** for the pattern that most closely matches the URL.

[src](https://golang.org/src/net/http/server.go?s=72834:72921#L2381)

```Go
// HandleFunc registers the handler function for the given pattern.
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerFunc(handler))
}
```

> so `func(ResponseWriter, *Request)` is converted to a HandlerFunc type

### `http.HandlerFunc`

- The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers.

[Doc](https://golang.org/pkg/net/http/#HandlerFunc)

`type HandlerFunc` do implements `http.Handler` Interface

```Go
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
```

So Example handleFunc1, can also be written as.

```Go
func main() {
    mux := http.NewServeMux()

    // Explicitly Converting function to a HandlerFunc type
    mh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world From Go- HanldeFunc Type")
    })

	mux.Handle("/", mh)
	log.Fatal(http.ListenAndServe(":3002", mux))
}
```

### Advanced

Go net/http request-response flow

> Client -> Requests -> [Multiplexer(router)] -> handler -> Response -> Client

**Multiplexer in Go is based on ServeMux structure**

### ServeMux

[src](https://golang.org/src/net/http/server.go?s=72834:72921#L2149)

```Go
type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	hosts bool // whether any patterns contain hostnames
}

type muxEntry struct {
	h       Handler
	pattern string
}
```

Focus:

- `m` in ServeMux, is of type `map[string]muxEntry` : which is the key to URL matching, It takes the URL Path as the key and the corresponding Handler `muxEntry` as the Value.

### Process to registers the handler for the given pattern

[src](https://golang.org/src/net/http/server.go?s=72834:72921#L2366)
If the pattern is not registered, the handler will be registered to this pattern

Overview:

```Go
func (mux *ServeMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
    }
    // pattern ie /url/ becomes the key in the ServeMux struct
	mux.m[pattern] = muxEntry{h: handler, pattern: pattern}

	if pattern[0] != '/' {
		mux.hosts = true
	}
}
```

### Process to match the handler for the given pattern

- After registering the route, starting the web service also requires server monitoring.
- `http.ListenAndServer` method can be seen to create a Server object, and call the same name method of the `Server` type:

[src](https://golang.org/src/net/http/server.go?s=59784:59844#L3002)

```Go
// ListenAndServe always returns a non-nil error.
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```

##### `Server` type

[src](https://golang.org/src/net/http/server.go?s=59784:59844#L2445)

```Go
type Server struct {
	Addr    string  // TCP address to listen on, ":http" if empty
	Handler Handler // handler to invoke, http.DefaultServeMux if nil
	......
}
```

##### `ListenAndServe` method of receiver type `*Server`

The `server.ListenAndServe()` method internally calls `net.Listen("tcp", addr)`, which internally calls `net.ListenTCP()` to create and return a listener `net.Listener`, such as ln;

[src](https://golang.org/src/net/http/server.go?s=59784:59844#L2750)

```Go
func (srv *Server) ListenAndServe() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}
```

Finally, the monitored TCP object is passed to the Serve method:

#### `Serve` method of receiver type `*Server`

[src](https://golang.org/src/net/http/server.go?s=59784:59844#L2795)

Serve accepts incoming connections on the Listener l, creating a new service goroutine for each. The service goroutines read requests and then call `srv.Handler` to reply to them.

Remember `Server` type was struct with a field Handler.

```Go
func (srv *Server) Serve(l net.Listener) error {
	...
	defer l.Close()

	...
	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
		rw, e := l.Accept()
		...
		c := srv.newConn(rw)
		c.setState(c.rwc, StateNew) // before Serve can return
		go c.serve(ctx)
	}
}
```

It uses the newConn method to create the connection object.

```Go
// Create new connection from rwc.
func (srv *Server) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: srv,
		rwc:    rwc,
	}
	if debugServerConnections {
		c.rwc = newLoggingConn("server", c.rwc)
	}
	return c
}
```

Finally, the connection request is processed using the goroutine.

[src](https://golang.org/src/net/http/server.go?s=59784:59844#L1738)

```Go
// Serve a new connection.
func (c *conn) serve(ctx context.Context) {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	...

	c.r = &connReader{conn: c}
	c.bufr = newBufioReader(c.r)
	c.bufw = newBufioWriterSize(checkConnErrorWriter{c}, 4<<10)

	for {
		w, err := c.readRequest(ctx)
		...
		serverHandler{c.server}.ServeHTTP(w, w.req)
		...
	}
}
```

**The next step is to call the `serverHandler{c.server}.ServeHTTP(w, w.req)` method to process the request**

#### `serverHandler` type

```Go
// serverHandler delegates to either the server's Handler or
// DefaultServeMux and also handles "OPTIONS *" requests.
type serverHandler struct {
	srv *Server
}
```

The serverHandler is an important structure. It has a field nearby, that is, the Server structure. It also implements the Handler interface method ServeHTTP, and does an important thing in the interface method to initialize the multiplexer route multiplexer. If the server object does not specify a Handler, the default DefaultServeMux is used as the route multiplexer. And call the ServeHTTP method that initializes the Handler.

[src](https://golang.org/src/net/http/server.go?s=59784:59844#L2733)

```Go
func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	if req.RequestURI == "*" && req.Method == "OPTIONS" {
		handler = globalOptionsHandler{}
	}
	handler.ServeHTTP(rw, req)
}
```

[src](https://golang.org/src/net/http/server.go?s=59784:59844#L2350)

```Go
// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}
```

[src](https://golang.org/src/net/http/server.go?s=59784:59844#L2281)

```Go
// Handler returns the handler to use for the given request,
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {
	...
	return mux.handler(host, r.URL.Path)
}
```

[src](https://golang.org/src/net/http/server.go?s=59784:59844#L2331)

```Go
// handler is the main implementation of Handler.
func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	// Host-specific pattern takes precedence over generic ones
	if mux.hosts {
		h, pattern = mux.match(host + path)
	}
	if h == nil {
		h, pattern = mux.match(path)
	}
	if h == nil {
		h, pattern = NotFoundHandler(), ""
	}
	return
}
```

[src](https://golang.org/src/net/http/server.go?s=72834:72921#L2218)

Overview:

```Go
// Find a handler on a handler map given a path string.
// Most-specific (longest) pattern wins.
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
	// Check for exact match first.
	v, ok := mux.m[path]
	if ok {
		return v.h, v.pattern
	}

	// Check for longest valid match.
	var n = 0
	for k, v := range mux.m {
		if !pathMatch(k, path) {
			continue
		}
		if h == nil || len(k) > n {
			n = len(k)
			h = v.h
			pattern = v.pattern
		}
	}
	return
}
```

So **exact path set is the best match** because it is judged by the length of the path.
Of course, it explains why `/` can match all but this is the last to be successfully matched.
