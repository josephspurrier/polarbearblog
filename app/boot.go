package app

import (
	"fmt"
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

var (
	storageSitePath    = "storage/site.json"
	storageSessionPath = "storage/session.bin"
	sessionName        = "session"
)

// Boot -
func Boot() (http.Handler, error) {
	// Set the storage and session environment variables.
	sitePath := os.Getenv("SS_SITE_PATH")
	if len(sitePath) > 0 {
		storageSitePath = sitePath
	}

	sessionPath := os.Getenv("SS_SESSION_PATH")
	if len(sessionPath) > 0 {
		storageSessionPath = sessionPath
	}

	sname := os.Getenv("SS_SESSION_NAME")
	if len(sname) > 0 {
		sessionName = sname
	}

	// Get the environment variables.
	secretKey := os.Getenv("SS_SESSION_KEY")
	if len(secretKey) == 0 {
		return nil, fmt.Errorf("environment variable missing: %v", "SS_SESSION_KEY")
	}

	bucket := os.Getenv("GCP_BUCKET_NAME")
	if len(bucket) == 0 {
		return nil, fmt.Errorf("environment variable missing: %v", "GCP_BUCKET_NAME")
	}

	allowHTML, err := strconv.ParseBool(os.Getenv("SS_ALLOW_HTML"))
	if err != nil {
		return nil, fmt.Errorf("environment variable not able to parse as bool: %v", "SS_ALLOW_HTML")
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
		return nil, err
	}

	// Set up the session storage provider.
	en := websession.NewEncryptedStorage(secretKey)
	store, err := websession.NewJSONSession(ss, en)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// Set up the router and middleware.
	var mw http.Handler
	mw = c.Router
	h := middleware.NewHandler(c.Render, c.Sess, c.Router, c.Storage.Site.URL, c.Storage.Site.Scheme)
	mw = h.Redirect(mw)
	mw = h.Head(mw)
	mw = h.DisallowAnon(mw)
	mw = sessionManager.LoadAndSave(mw)
	mw = middleware.Gzip(mw)
	mw = h.LogRequest(mw)

	return mw, nil
}
