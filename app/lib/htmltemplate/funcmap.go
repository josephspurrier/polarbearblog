package htmltemplate

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/josephspurrier/polarbearblog/app/lib/envdetect"
	"github.com/josephspurrier/polarbearblog/app/model"
)

// funcMap returns a map of template functions that can be used in templates.
func (te *Engine) funcMap(r *http.Request) template.FuncMap {
	fm := make(template.FuncMap)
	fm["Stamp"] = func(t time.Time) string {
		return t.Format("2006-01-02")
	}
	fm["StampFriendly"] = func(t time.Time) string {
		return t.Format("02 Jan, 2006")
	}
	fm["PublishedPages"] = func() []model.Post {
		return te.storage.Site.PublishedPages()
	}
	fm["SiteURL"] = func() string {
		return te.storage.Site.SiteURL()
	}
	fm["SiteTitle"] = func() string {
		return te.storage.Site.SiteTitle()
	}
	fm["SiteSubtitle"] = func() string {
		return te.storage.Site.SiteSubtitle()
	}
	fm["SiteDescription"] = func() string {
		return te.storage.Site.Description
	}
	fm["SiteAuthor"] = func() string {
		return te.storage.Site.Author
	}
	fm["SiteFavicon"] = func() string {
		return te.storage.Site.Favicon
	}
	fm["Authenticated"] = func() bool {
		// If user is not authenticated, don't allow them to access the page.
		_, loggedIn := te.sess.User(r)
		return loggedIn
	}
	fm["GoogleAnalyticsID"] = func() string {
		if envdetect.RunningLocalDev() {
			return ""
		}
		return te.storage.Site.GoogleAnalyticsID
	}
	fm["DisqusID"] = func() string {
		if envdetect.RunningLocalDev() {
			return ""
		}
		return te.storage.Site.DisqusID
	}
	fm["SiteFooter"] = func() string {
		return te.storage.Site.Footer
	}
	fm["MFAEnabled"] = func() bool {
		return len(os.Getenv("PBB_MFA_KEY")) > 0
	}
	fm["AssetStamp"] = func(f string) string {
		return assetTimePath(f)
	}
	fm["SiteStyles"] = func() template.CSS {
		return template.CSS(te.storage.Site.Styles)
	}
	fm["StylesAppend"] = func() bool {
		if len(te.storage.Site.Styles) == 0 {
			// If there are no style, then always append.
			return true
		} else if te.storage.Site.StylesAppend {
			// Else if there are style and it's append, then append.
			return true
		}
		return false
	}

	return fm
}
