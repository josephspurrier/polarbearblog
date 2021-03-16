package middleware

import (
	"net/http"
	"strings"
)

// DisallowAuth does not allow authenticated users to access the page.
func (c *Handler) DisallowAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If user is authenticated, don't allow them to access the page.
		if _, loggedIn := c.Sess.User(r); !loggedIn {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// DisallowAnon does not allow anonymous users to access the page.
func (c *Handler) DisallowAnon(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't allow anon users to access the dashboard.
		if strings.HasPrefix(r.URL.Path, "/dashboard") {
			// If user is not authenticated, don't allow them to access the page.
			if _, loggedIn := c.Sess.User(r); !loggedIn {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
