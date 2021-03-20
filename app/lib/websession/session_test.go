package websession_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/josephspurrier/polarbearblog/app/lib/datastorage"
	"github.com/josephspurrier/polarbearblog/app/lib/websession"
	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	// Set up the session storage provider.
	f := "data.bin"
	err := ioutil.WriteFile(f, []byte(""), 0644)
	assert.NoError(t, err)
	ss := datastorage.NewLocalStorage(f)
	secretkey := "82a18fbbfed2694bb15d512a70c53b1a088e669966918d3d474564b2ac44349b"
	en := websession.NewEncryptedStorage(secretkey)
	store, err := websession.NewJSONSession(ss, en)
	assert.NoError(t, err)

	// Initialize a new session manager and configure the session lifetime.
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = false
	sessionManager.Store = store
	sess := websession.New("session", sessionManager)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Test user
		u := "foo"
		sess.SetUser(r, u)
		user, found := sess.User(r)
		assert.True(t, found)
		assert.Equal(t, u, user)

		// Test Logout
		sess.Logout(r)
		_, found = sess.User(r)
		assert.False(t, found)

		// Test persistence
		assert.Equal(t, sessionManager.Cookie.Persist, false)
		sess.RememberMe(r, true)
		assert.Equal(t, sessionManager.Cookie.Persist, true)

		// Test CSRF
		assert.False(t, sess.CSRF(r))
		token := sess.SetCSRF(r)
		r.Form = url.Values{}
		r.Form.Set("token", token)
		assert.True(t, sess.CSRF(r))
	})

	mw := sessionManager.LoadAndSave(mux)
	mw.ServeHTTP(w, r)

	os.Remove(f)
}
