package html

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/josephspurrier/polarbearblog/app/lib/datastorage"
	"github.com/josephspurrier/polarbearblog/app/lib/envdetect"
	"github.com/josephspurrier/polarbearblog/app/lib/websession"
	"github.com/josephspurrier/polarbearblog/app/model"
)

// FuncMap returns a map of template functions that can be used in templates.
func FuncMap(r *http.Request, storage *datastorage.Storage,
	sess *websession.Session) template.FuncMap {
	fm := make(template.FuncMap)
	fm["Stamp"] = func(t time.Time) string {
		return t.Format("2006-01-02")
	}
	fm["StampFriendly"] = func(t time.Time) string {
		return t.Format("02 Jan, 2006")
	}
	fm["PublishedPages"] = func() []model.Post {
		return storage.Site.PublishedPages()
	}
	fm["SiteURL"] = func() string {
		return storage.Site.SiteURL()
	}
	fm["SiteTitle"] = func() string {
		return storage.Site.SiteTitle()
	}
	fm["SiteSubtitle"] = func() string {
		return storage.Site.SiteSubtitle()
	}
	fm["SiteDescription"] = func() string {
		return storage.Site.Description
	}
	fm["SiteAuthor"] = func() string {
		return storage.Site.Author
	}
	fm["SiteFavicon"] = func() string {
		return storage.Site.Favicon
	}
	fm["Authenticated"] = func() bool {
		// If user is not authenticated, don't allow them to access the page.
		_, loggedIn := sess.User(r)
		return loggedIn
	}
	fm["GoogleAnalyticsID"] = func() string {
		if envdetect.RunningLocalDev() {
			return ""
		}
		return storage.Site.GoogleAnalyticsID
	}
	fm["DisqusID"] = func() string {
		if envdetect.RunningLocalDev() {
			return ""
		}
		return storage.Site.DisqusID
	}
	fm["SiteFooter"] = func() string {
		return storage.Site.Footer
	}
	fm["MFAEnabled"] = func() bool {
		return len(os.Getenv("PBB_MFA_KEY")) > 0
	}
	fm["AssetStamp"] = func(f string) string {
		return assetTimePath(f)
	}
	fm["SiteStyles"] = func() template.CSS {
		return template.CSS(storage.Site.Styles)
	}
	fm["StylesAppend"] = func() bool {
		if len(storage.Site.Styles) == 0 {
			// If there are no style, then always append.
			return true
		} else if storage.Site.StylesAppend {
			// Else if there are style and it's append, then append.
			return true
		}
		return false
	}
	fm["EnableStackEdit"] = func() bool {
		return storage.Site.StackEdit
	}

	return fm
}
