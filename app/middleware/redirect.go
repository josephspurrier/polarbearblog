package middleware

import (
	"net/http"
	"strings"
)

// Redirect will handle all redirects required.
func (c *Handler) Redirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "text/html; charset=utf-8")
		// Don't allow access to files with a slash at the end.
		if strings.Contains(r.URL.Path, ".") && strings.HasSuffix(r.URL.Path, "/") {
			c.Router.NotFound(w, r)
			return
		}

		// Strip trailing slash.
		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, strings.TrimRight(r.URL.Path, "/"), http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
