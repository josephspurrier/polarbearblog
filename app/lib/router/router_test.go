package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// defaultServeHTTP is the default ServeHTTP function that receives the status and error from
// the function call.
var defaultServeHTTP = func(w http.ResponseWriter, r *http.Request, status int,
	err error) {
	if status >= 400 {
		if err != nil {
			http.Error(w, err.Error(), status)
		} else {
			http.Error(w, "", status)
		}
	}
}

func TestParams(t *testing.T) {
	mux := New(defaultServeHTTP, nil)
	mux.Get("/user/:name", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			assert.Equal(t, "john", mux.Param(r, "name"))
			return http.StatusOK, nil
		}))

	r := httptest.NewRequest("GET", "/user/john", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
}

func TestInstance(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	mux.Get("/user/:name", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			assert.Equal(t, "john", mux.Param(r, "name"))
			return http.StatusOK, nil
		}))

	r := httptest.NewRequest("GET", "/user/john", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
}

func TestPostForm(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	form := url.Values{}
	form.Add("username", "jsmith")

	mux.Post("/user", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			r.ParseForm()
			assert.Equal(t, "jsmith", r.FormValue("username"))
			return http.StatusOK, nil
		}))

	r := httptest.NewRequest("POST", "/user", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
}

func TestPostJSON(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	j, err := json.Marshal(map[string]interface{}{
		"username": "jsmith",
	})
	assert.Nil(t, err)

	mux.Post("/user", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			b, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			r.Body.Close()
			assert.Equal(t, `{"username":"jsmith"}`, string(b))
			return http.StatusOK, nil
		}))

	r := httptest.NewRequest("POST", "/user", bytes.NewBuffer(j))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
}

func TestGet(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	called := false

	mux.Get("/user", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			called = true
			return http.StatusOK, nil
		}))

	r := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestDelete(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	called := false

	mux.Delete("/user", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			called = true
			return http.StatusOK, nil
		}))

	r := httptest.NewRequest("DELETE", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestHead(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	called := false

	mux.Head("/user", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			called = true
			return http.StatusOK, nil
		}))

	r := httptest.NewRequest("HEAD", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestOptions(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	called := false

	mux.Options("/user", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			called = true
			return http.StatusOK, nil
		}))

	r := httptest.NewRequest("OPTIONS", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestPatch(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	called := false

	mux.Patch("/user", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			called = true
			return http.StatusOK, nil
		}))

	r := httptest.NewRequest("PATCH", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestPut(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	called := false

	mux.Put("/user", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			called = true
			return http.StatusOK, nil
		}))

	r := httptest.NewRequest("PUT", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func Test404(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	called := false

	mux.Get("/user", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			called = true
			return http.StatusOK, nil
		}))

	r := httptest.NewRequest("GET", "/badroute", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, false, called)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test500NoError(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	called := true

	mux.Get("/user", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			called = true
			return http.StatusInternalServerError, nil
		}))

	r := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func Test500WithError(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	called := true
	specificError := errors.New("specific error")

	mux.Get("/user", HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) (status int, err error) {
			called = true
			return http.StatusInternalServerError, specificError
		}))

	r := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, w.Body.String(), specificError.Error()+"\n")
}

func Test400(t *testing.T) {
	notFound := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		},
	)

	mux := New(defaultServeHTTP, notFound)

	r := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestNotFound(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	r := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	mux.NotFound(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBadRequest(t *testing.T) {
	mux := New(defaultServeHTTP, nil)

	r := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	mux.BadRequest(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
