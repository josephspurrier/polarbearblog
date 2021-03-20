package route_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josephspurrier/polarbearblog/app/lib/router"
	"github.com/josephspurrier/polarbearblog/app/route"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *router.Mux {
	// Set the handling of all responses.
	customServeHTTP := func(w http.ResponseWriter, r *http.Request, status int, err error) {}

	// Send all 404 to the customer handler.
	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customServeHTTP(w, r, http.StatusNotFound, nil)
	})

	// Set up the router.
	return router.New(customServeHTTP, notFound)
}

func TestXML(t *testing.T) {
	mux := setupRouter()

	// Create core app.
	c := &route.Core{}
	x := &route.XMLUtil{c}
	mux.Get("/robots.txt", x.Robots)
	r := httptest.NewRequest("GET", "/robots.txt", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	b, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, "User-agent: *\nAllow: /", string(b))
}
