package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josephspurrier/polarbearblog/app/middleware"
	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	r := httptest.NewRequest("HEAD", "/", nil)
	w := httptest.NewRecorder()
	mux := http.NewServeMux()
	mw := middleware.Head(mux)
	mw.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
