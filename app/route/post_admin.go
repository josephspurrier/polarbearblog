package route

import (
	"net/http"
	"time"

	"github.com/josephspurrier/polarbearblog/app/lib/uuid"
	"github.com/josephspurrier/polarbearblog/app/model"
	"github.com/matryer/way"
)

// AdminPost -
type AdminPost struct {
	*Core
}

func registerAdminPost(c *AdminPost) {
	c.Router.Get("/dashboard/posts", c.index)
	c.Router.Get("/dashboard/posts/new", c.create)
	c.Router.Post("/dashboard/posts/new", c.store)
	c.Router.Get("/dashboard/posts/:id", c.edit)
	c.Router.Post("/dashboard/posts/:id", c.update)
	c.Router.Get("/dashboard/posts/:id/delete", c.destroy)
}

func (c *AdminPost) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	vars := make(map[string]interface{})
	vars["title"] = "Posts"
	vars["posts"] = c.Storage.Site.PostsAndPages(false)

	return c.Render.Template(w, r, "dashboard", "bloglist_edit", vars)
}

func (c *AdminPost) create(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	vars := make(map[string]interface{})
	vars["title"] = "New post"
	vars["token"] = c.Sess.SetCSRF(r)

	return c.Render.Template(w, r, "dashboard", "post_create", vars)
}

func (c *AdminPost) store(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	ID, err := uuid.Generate()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	r.ParseForm()

	// CSRF protection.
	success := c.Sess.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	now := time.Now()

	var p model.Post
	p.Title = r.FormValue("title")
	p.URL = r.FormValue("slug")
	p.Canonical = r.FormValue("canonical_url")
	p.Created = now
	p.Updated = now
	pubDate := r.FormValue("published_date")
	if pubDate == "" {
		pubDate = now.Format("2006-01-02")
	}
	ts, err := time.Parse("2006-01-02", pubDate)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	p.Timestamp = ts
	p.Content = r.FormValue("content")
	p.Tags = p.Tags.Split(r.FormValue("tags"))
	p.Page = r.FormValue("is_page") == "on"
	p.Published = r.FormValue("publish") == "on"

	// Save to storage.
	c.Storage.Site.Posts[ID] = p
	err = c.Storage.Save()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/dashboard/posts/"+ID, http.StatusFound)
	return
}

func (c *AdminPost) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	vars := make(map[string]interface{})
	vars["title"] = "Edit post"
	vars["token"] = c.Sess.SetCSRF(r)

	ID := way.Param(r.Context(), "id")

	var p model.Post
	var ok bool
	if p, ok = c.Storage.Site.Posts[ID]; !ok {
		return http.StatusNotFound, nil
	}

	vars["id"] = ID
	vars["ptitle"] = p.Title
	vars["url"] = p.URL
	vars["canonical"] = p.Canonical
	vars["timestamp"] = p.Timestamp
	vars["body"] = p.Content
	vars["tags"] = p.Tags.String()
	vars["page"] = p.Page
	vars["published"] = p.Published

	return c.Render.Template(w, r, "dashboard", "post_edit", vars)
}

func (c *AdminPost) update(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	ID := way.Param(r.Context(), "id")

	var p model.Post
	var ok bool
	if p, ok = c.Storage.Site.Posts[ID]; !ok {
		return http.StatusNotFound, nil
	}

	// Save the site.
	r.ParseForm()

	// CSRF protection.
	success := c.Sess.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	now := time.Now()

	p.Title = r.FormValue("title")
	p.URL = r.FormValue("slug")
	p.Canonical = r.FormValue("canonical_url")
	p.Updated = now
	pubDate := r.FormValue("published_date")
	ts, err := time.Parse("2006-01-02", pubDate)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	p.Timestamp = ts
	p.Content = r.FormValue("content")
	p.Tags = p.Tags.Split(r.FormValue("tags"))
	p.Page = r.FormValue("is_page") == "on"
	p.Published = r.FormValue("publish") == "on"

	c.Storage.Site.Posts[ID] = p

	err = c.Storage.Save()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/dashboard/posts/"+ID, http.StatusFound)
	return
}

func (c *AdminPost) destroy(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	ID := way.Param(r.Context(), "id")

	var ok bool
	if _, ok = c.Storage.Site.Posts[ID]; !ok {
		return http.StatusNotFound, nil
	}

	delete(c.Storage.Site.Posts, ID)

	err = c.Storage.Save()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/dashboard/posts", http.StatusFound)
	return
}
