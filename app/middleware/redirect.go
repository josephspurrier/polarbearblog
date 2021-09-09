package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/josephspurrier/polarbearblog/app/lib/envdetect"
)

// Redirect will handle all redirects required.
func (c *Handler) Redirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Redirect to the correct website.
		if !envdetect.RunningLocalDev() && len(c.SiteURL) > 0 && !strings.Contains(r.Host, c.SiteURL) {
			http.Redirect(w, r, fmt.Sprintf("%v://%v%v", c.SiteScheme, c.SiteURL, r.URL.Path), http.StatusPermanentRedirect)
			return
		}

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
