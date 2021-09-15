package route

import (
	"net/http"
	"strings"

	"github.com/josephspurrier/polarbearblog/app/lib/htmltemplate"
	"github.com/josephspurrier/polarbearblog/app/model"
	"github.com/matryer/way"
)

// Post -
type Post struct {
	*Core
}

func registerPost(c *Post) {
	c.Router.Get("/blog", c.index)
	c.Router.Get("/:slug", c.show)
}

func (c *Post) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	vars := make(map[string]interface{})
	vars["tags"] = c.Storage.Site.Tags(true)

	// Determine if there is query.
	if q := r.URL.Query().Get("q"); len(q) > 0 {
		vars["query"] = q
		// Don't show tags when there is a filter.
		delete(vars, "tags")

		posts := make([]model.PostWithID, 0)
		for _, v := range c.Storage.Site.PostsAndPages(true) {
			match := false
			for _, tag := range v.Tags {
				if tag.Name == q {
					match = true
					break
				}
			}

			if match {
				posts = append(posts, v)
			}
		}

		vars["posts"] = posts
	} else {
		vars["posts"] = c.Storage.Site.PublishedPosts()
	}

	return c.Render.Template(w, r, "base", "bloglist_index", vars)
}

func (c *Post) show(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	slug := way.Param(r.Context(), "slug")
	p := c.Storage.Site.PostBySlug(slug)

	// Determine if in preview mode.
	preview := false
	if q := r.URL.Query().Get("preview"); len(q) > 0 && strings.ToLower(q) == "true" {
		preview = true
	}

	// Show 404 if not published and not in preview mode.
	if !p.Published && !preview {
		return http.StatusNotFound, nil
	}

	vars := make(map[string]interface{})
	// Don't show certain items on pages.
	if !p.Page {
		vars["title"] = p.Title
		vars["pubdate"] = p.Timestamp
	}

	vars["tags"] = p.Tags
	vars["canonical"] = p.Canonical
	vars["id"] = p.ID
	vars["posturl"] = p.URL
	vars["metadescription"] = htmltemplate.PlaintextBlurb(p.Content)

	return c.Render.Post(w, r, "base", p.Post, vars)
}
