package router

import (
	"net/http"
)

// handler is a internal handler.
type handler struct {
	HandlerFunc
	CustomServeHTTP func(w http.ResponseWriter, r *http.Request, status int, err error)
}

// ServeHTTP handles all the errors from the HTTP handlers.
func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := fn.HandlerFunc(w, r)
	fn.CustomServeHTTP(w, r, status, err)
}

// HandlerFunc is used to wrapper all endpoint functions so they work with generic
// routers.
type HandlerFunc func(http.ResponseWriter, *http.Request) (int, error)
