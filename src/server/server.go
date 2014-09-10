package server

import . "modinputs"
import "net/http"
import "fmt"

var servers map[int]Mux = make(map[int]Mux)

func Get(port int) Mux {
	server, found := servers[port]
	if !found {
		server = *new(Mux)
		server.port = port
		server.m = make(map[string]http.Handler)
		servers[port] = server
	}
	return server
}

type Mux struct {
	port int
	m    map[string]http.Handler
}

func (mux Mux) Register(path string, h http.Handler) {
	if _, exists := mux.m[path]; exists {
		Warn("Handler for path %s already exists (duplicate")
	} else {
		mux.m[path] = h
	}
}

func (mux Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Info("Incoming request on port=%d -> method=%s uri=%s length=%d from clientip=%s", mux.port, r.Method, r.URL.Path, r.ContentLength, r.RemoteAddr)

	if handler, found := mux.m[r.URL.Path]; found {
		handler.ServeHTTP(w, r)
	} else {
		Warn("No handler found for path=%s on port=%d", r.URL.Path, mux.port)
		http.Error(w, "404", 404)
	}
}

func StartListening() {
	for port, server := range servers {
		go http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	}
}
