package router

import (
	"net/http"

	"github.com/matryer/way"
)

// Mux contains the router.
type Mux struct {
	router *way.Router

	// customServeHTTP is the serve function.
	customServeHTTP func(w http.ResponseWriter, r *http.Request, status int, err error)
}

// New returns an instance of the router.
func New(csh func(w http.ResponseWriter, r *http.Request, status int, err error), notFound http.Handler) *Mux {
	r := way.NewRouter()
	if notFound != nil {
		r.NotFound = notFound
	}

	return &Mux{
		router:          r,
		customServeHTTP: csh,
	}
}

// ServeHTTP routes the incoming http.Request based on method and path
// extracting path parameters as it goes.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

// NotFound shows the 404 page.
func (m *Mux) NotFound(w http.ResponseWriter, r *http.Request) {
	m.customServeHTTP(w, r, http.StatusNotFound, nil)
}

// BadRequest shows the 400 page.
func (m *Mux) BadRequest(w http.ResponseWriter, r *http.Request) {
	m.customServeHTTP(w, r, http.StatusBadRequest, nil)
}

// Param returns a URL parameter.
func (m *Mux) Param(r *http.Request, param string) string {
	return way.Param(r.Context(), param)
}
