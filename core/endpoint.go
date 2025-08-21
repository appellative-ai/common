package core

import "net/http"

// HttpHandler - extend the http.HandlerFunc to include the http.Response
type HttpHandler func(w http.ResponseWriter, req *http.Request, resp *http.Response)
type HttpInit func(r *http.Request) *http.Request

type Endpoint interface {
	Pattern() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type endpoint struct {
	pattern string
	handler HttpHandler
	init    HttpInit
	chain   Exchange
}

func NewEndpoint(pattern string, handler HttpHandler, init HttpInit, operatives []any) Endpoint {
	e := new(endpoint)
	e.pattern = pattern
	e.handler = handler
	e.init = init
	e.chain = BuildNetwork(operatives)
	return e
}

func (e *endpoint) Pattern() string {
	return e.pattern
}

func (e *endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e.handler == nil || e.chain == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if e.init != nil {
		r = e.init(r)
	}
	resp, _ := e.chain(r)
	e.handler(w, r, resp)
}
