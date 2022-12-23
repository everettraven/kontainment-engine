package router

import "net/http"

type Router interface {
	Routes() []Route
}

type Route interface {
	Method() string
	Pattern() string
	Handler() http.HandlerFunc
}

// route is a generic implementation of the Route interface
type route struct {
	method  string
	pattern string
	handler http.HandlerFunc
}

func (r *route) Method() string {
	return r.method
}

func (r *route) Pattern() string {
	return r.pattern
}

func (r *route) Handler() http.HandlerFunc {
	return r.handler
}

// NewRoute returns a new route for the given method, pattern, and handler
func NewRoute(method string, pattern string, handler http.HandlerFunc) Route {
	return &route{
		method:  method,
		pattern: pattern,
		handler: handler,
	}
}

// NewGetRoute returns a new GET route for the given pattern and handler
func NewGetRoute(pattern string, handler http.HandlerFunc) Route {
	return &route{
		method:  http.MethodGet,
		pattern: pattern,
		handler: handler,
	}
}

// NewPostRoute returns a new POST route for the given pattern and handler
func NewPostRoute(pattern string, handler http.HandlerFunc) Route {
	return &route{
		method:  http.MethodPost,
		pattern: pattern,
		handler: handler,
	}
}

// NewDeleteRoute returns a new DELETE route for the given pattern and handler
func NewDeleteRoute(pattern string, handler http.HandlerFunc) Route {
	return &route{
		method:  http.MethodDelete,
		pattern: pattern,
		handler: handler,
	}
}
