package app

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/josephspurrier/polarbearblog/app/lib/datastorage"
	"github.com/josephspurrier/polarbearblog/app/lib/envdetect"
	"github.com/josephspurrier/polarbearblog/app/lib/htmltemplate"
	"github.com/josephspurrier/polarbearblog/app/lib/websession"
	"github.com/josephspurrier/polarbearblog/app/middleware"
	"github.com/josephspurrier/polarbearblog/app/model"
	"github.com/josephspurrier/polarbearblog/app/route"
)

const (
	storageSitePath    = "storage/site.json"
	storageSessionPath = "storage/session.json"
	sessionName        = "session"
)

// Boot -
func Boot() {
	// Get the environment variables.
	secretKey := os.Getenv("SS_SESSION_KEY")
	if len(secretKey) == 0 {
		log.Fatalln("Environment variable missing:", "SS_SESSION_KEY")
	}

	bucket := os.Getenv("GCP_BUCKET_NAME")
	if len(bucket) == 0 {
		log.Fatalln("Environment variable missing:", "GCP_BUCKET_NAME")
	}

	allowHTML, err := strconv.ParseBool(os.Getenv("SS_ALLOW_HTML"))
	if err != nil {
		log.Fatalln("Environment variable not able to parse as bool:", "SS_ALLOW_HTML")
	}

	// Create the local files if in development mode.
	if envdetect.RunningLocalDev() {
		if _, err := os.Stat(storageSitePath); os.IsNotExist(err) {
			ioutil.WriteFile(storageSitePath, []byte("{}"), 0644)
		}
		if _, err := os.Stat(storageSessionPath); os.IsNotExist(err) {
			ioutil.WriteFile(storageSessionPath, []byte("{}"), 0644)
		}
	}

	// Decode the key for encrypting sesstions.
	decodedSecretKey, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		log.Fatalln(err)
	}

	// Create new store object with the defaults.
	site := &model.Site{}

	var ds datastorage.Datastorer
	var ss websession.Sessionstorer

	if !envdetect.RunningLocalDev() {
		// Use Google when running in GCP.
		ds = datastorage.NewGCPStorage(bucket, storageSitePath)
		ss = datastorage.NewGCPStorage(bucket, storageSessionPath)
	} else {
		// Use local filesytem when developing.
		ds = datastorage.NewLocalStorage(storageSitePath)
		ss = datastorage.NewLocalStorage(storageSessionPath)
	}

	// Set up the data storage provider.
	storage, err := datastorage.New(ds, site)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Set up the session storage provider.
	store, err := websession.NewJSONSession(ss, decodedSecretKey)
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize a new session manager and configure the session lifetime.
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = false
	sessionManager.Store = store
	sess := websession.New(sessionName, sessionManager)

	// Set up the template engine.
	tmpl := htmltemplate.New(storage, sess, allowHTML)

	// Setup the routes.
	c, err := route.Register(storage, sess, tmpl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Set up the router and middleware.
	var mw http.Handler
	mw = c.Router
	h := middleware.NewHandler(c.Render, c.Sess, c.Router, c.Storage.Site.URL, c.Storage.Site.Scheme)
	mw = h.Redirect(mw)
	mw = h.Head(mw)
	mw = h.DisallowAnon(mw)
	mw = middleware.Gzip(mw)
	mw = h.LogRequest(mw)
	mw = sessionManager.LoadAndSave(mw)

	// Start the web server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Web server running on port:", port)
	log.Fatalln(http.ListenAndServe(":"+port, mw))
}
