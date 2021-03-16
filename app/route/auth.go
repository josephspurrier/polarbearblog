package route

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/josephspurrier/polarbearblog/app/lib/envdetect"
	"github.com/josephspurrier/polarbearblog/app/lib/passhash"
	"github.com/josephspurrier/polarbearblog/app/lib/totp"
	"github.com/matryer/way"
)

// AuthUtil -
type AuthUtil struct {
	*Core
}

func registerAuthUtil(c *AuthUtil) {
	c.Router.Get("/login/:slug", c.login)
	c.Router.Post("/login/:slug", c.loginPost)
	c.Router.Get("/dashboard/logout", c.logout)
}

func (c *AuthUtil) login(w http.ResponseWriter, r *http.Request) (status int, err error) {
	slug := way.Param(r.Context(), "slug")
	if slug != c.Storage.Site.LoginURL {
		return http.StatusNotFound, nil
	}

	vars := make(map[string]interface{})
	vars["title"] = "Login"
	vars["token"] = c.Sess.SetCSRF(r)

	return c.Render.Template(w, r, "base", "login", vars)
}

func (c *AuthUtil) loginPost(w http.ResponseWriter, r *http.Request) (status int, err error) {
	slug := way.Param(r.Context(), "slug")
	if slug != c.Storage.Site.LoginURL {
		return http.StatusNotFound, nil
	}

	r.ParseForm()

	// CSRF protection.
	success := c.Sess.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	mfa := r.FormValue("mfa")
	remember := r.FormValue("remember")

	username := os.Getenv("SS_USERNAME")
	if len(username) == 0 {
		log.Println("Environment variable missing:", "SS_USERNAME")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	hash := os.Getenv("SS_PASSWORD_HASH")
	if len(hash) == 0 {
		log.Println("Environment variable missing:", "SS_PASSWORD_HASH")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	mfakey := os.Getenv("SS_MFA_KEY")
	if len(mfakey) == 0 {
		log.Println("Environment variable missing:", "SS_MFA_KEY")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	imfa, err := strconv.Atoi(mfa)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	mfaSuccess, err := totp.Authenticate(imfa, mfakey)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Decode the hash - this is to allow it to be stored easily since dollar
	// signs are difficult to work with.
	hashDecoded, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	passMatch := passhash.MatchString(string(hashDecoded), password)

	// When running locally, let any MFA pass.
	if envdetect.RunningLocalDev() {
		mfaSuccess = true
	}

	// If the username and password don't match, then just redirect.
	if email != username || !passMatch || !mfaSuccess {
		fmt.Printf("Login attempt failed. Username: %v | Password match: %v | MFA success: %v\n", username, passMatch, mfaSuccess)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fmt.Printf("Login attempt successful.\n")

	c.Sess.SetUser(r, username)
	if remember == "on" {
		c.Sess.RememberMe(r, true)
	} else {
		c.Sess.RememberMe(r, false)
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
	return
}

func (c *AuthUtil) logout(w http.ResponseWriter, r *http.Request) (status int, err error) {
	c.Sess.Logout(r)

	http.Redirect(w, r, "/", http.StatusFound)
	return
}
