package middleware

import (
	"net/http"
)

// Head will return a 200 for the uptimerobot.
func Head(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == http.MethodHead {
			return
		}

		next.ServeHTTP(w, r)
	})
}
