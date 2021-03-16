package route

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/josephspurrier/polarbearblog/app/lib/datastorage"
	"github.com/josephspurrier/polarbearblog/app/lib/htmltemplate"
	"github.com/josephspurrier/polarbearblog/app/lib/router"
	"github.com/josephspurrier/polarbearblog/app/lib/websession"
)

// Core -
type Core struct {
	Router  *router.Mux
	Storage *datastorage.Storage
	Render  *htmltemplate.Engine
	Sess    *websession.Session
}

// Register all routes.
func Register(storage *datastorage.Storage, sess *websession.Session, tmpl *htmltemplate.Engine) (*Core, error) {
	// Create core app.
	c := &Core{
		Router:  setupRouter(tmpl),
		Storage: storage,
		Render:  tmpl,
		Sess:    sess,
	}

	// Register routes.
	registerHomePost(&HomePost{c})
	registerAuthUtil(&AuthUtil{c})
	registerXMLUtil(&XMLUtil{c})
	registerAdminPost(&AdminPost{c})
	registerPost(&Post{c})

	return c, nil
}

func setupRouter(tmpl *htmltemplate.Engine) *router.Mux {
	// Set the handling of all responses.
	customServeHTTP := func(w http.ResponseWriter, r *http.Request, status int, err error) {
		// Handle only errors.
		if status >= 400 {
			vars := make(map[string]interface{})
			vars["title"] = fmt.Sprint(status)
			errTemplate := "400"
			if status == 404 {
				errTemplate = "404"
			}
			status, err = tmpl.ErrorTemplate(w, r, "base", errTemplate, vars)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
		}

		// Display server errors.
		if status >= 500 {
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

	// Send all 404 to the customer handler.
	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customServeHTTP(w, r, http.StatusNotFound, nil)
	})

	// Set up the router.
	rr := router.New(customServeHTTP, notFound)

	// Static assets.
	folder := filepath.FromSlash("assets")
	rr.Get("/assets...", func(w http.ResponseWriter, r *http.Request) (status int, err error) {
		// Don't allow directory browsing.
		if strings.HasSuffix(r.URL.Path, "/") {
			return http.StatusNotFound, nil
		}

		http.ServeFile(w, r, filepath.Join(folder, strings.TrimPrefix(r.URL.Path, "/assets/")))
		return
	})

	return rr
}
