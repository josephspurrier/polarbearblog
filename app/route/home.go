package route

import (
	"net/http"
	"time"

	"github.com/josephspurrier/polarbearblog/app/model"
)

// HomePost -
type HomePost struct {
	*Core
}

func registerHomePost(c *HomePost) {
	c.Router.Get("/", c.show)
	c.Router.Get("/dashboard", c.edit)
	c.Router.Post("/dashboard", c.update)
	c.Router.Get("/dashboard/reload", c.reload)
}

func (c *HomePost) show(w http.ResponseWriter, r *http.Request) (status int, err error) {
	p := model.Post{
		Content: c.Storage.Site.Content,
		URL:     "/",
	}

	if p.Content == "" {
		p.Content = "*No content yet.*"
	}

	vars := make(map[string]interface{})
	return c.Render.Post(w, r, "base", p, vars)
}

func (c *HomePost) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Edit site"
	vars["homeContent"] = c.Storage.Site.Content
	vars["ptitle"] = c.Storage.Site.Title
	vars["subtitle"] = c.Storage.Site.Subtitle
	vars["token"] = c.Sess.SetCSRF(r)

	// Help the user set the domain based off the current URL.
	if c.Storage.Site.URL == "" {
		vars["domain"] = r.Host
	} else {
		vars["domain"] = c.Storage.Site.URL
	}

	vars["scheme"] = c.Storage.Site.Scheme
	vars["pauthor"] = c.Storage.Site.Author
	vars["pdescription"] = c.Storage.Site.Description
	vars["favicon"] = c.Storage.Site.Favicon
	vars["loginurl"] = c.Storage.Site.LoginURL
	vars["googleanalytics"] = c.Storage.Site.GoogleAnalyticsID
	vars["disqus"] = c.Storage.Site.DisqusID
	vars["footer"] = c.Storage.Site.Footer

	return c.Render.Template(w, r, "dashboard", "home_edit", vars)
}

func (c *HomePost) update(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.ParseForm()

	// CSRF protection.
	success := c.Sess.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	c.Storage.Site.Title = r.FormValue("title")
	c.Storage.Site.Subtitle = r.FormValue("subtitle")
	c.Storage.Site.URL = r.FormValue("domain")
	c.Storage.Site.Content = r.FormValue("content")
	c.Storage.Site.Scheme = r.FormValue("scheme")
	c.Storage.Site.Author = r.FormValue("author")
	c.Storage.Site.Favicon = r.FormValue("favicon")
	c.Storage.Site.Description = r.FormValue("pdescription")
	c.Storage.Site.LoginURL = r.FormValue("loginurl")
	c.Storage.Site.GoogleAnalyticsID = r.FormValue("googleanalytics")
	c.Storage.Site.DisqusID = r.FormValue("disqus")
	c.Storage.Site.Footer = r.FormValue("footer")
	c.Storage.Site.Updated = time.Now()

	err = c.Storage.Save()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
	return
}

func (c *HomePost) reload(w http.ResponseWriter, r *http.Request) (status int, err error) {
	err = c.Storage.Load()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
	return
}
