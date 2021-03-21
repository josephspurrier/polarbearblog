package route

import (
	"net/http"
)

// Styles -
type Styles struct {
	*Core
}

func registerStyles(c *Styles) {
	c.Router.Get("/dashboard/styles", c.edit)
	c.Router.Post("/dashboard/styles", c.update)
}

func (c *Styles) edit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Site styles"
	vars["token"] = c.Sess.SetCSRF(r)
	vars["favicon"] = c.Storage.Site.Favicon
	vars["styles"] = c.Storage.Site.Styles
	vars["stylesappend"] = c.Storage.Site.StylesAppend
	vars["stackedit"] = c.Storage.Site.StackEdit
	vars["prism"] = c.Storage.Site.Prism

	return c.Render.Template(w, r, "dashboard", "styles_edit", vars)
}

func (c *Styles) update(w http.ResponseWriter, r *http.Request) (status int, err error) {
	r.ParseForm()

	// CSRF protection.
	success := c.Sess.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	c.Storage.Site.Favicon = r.FormValue("favicon")
	c.Storage.Site.Styles = r.FormValue("styles")
	c.Storage.Site.StylesAppend = (r.FormValue("stylesappend") == "on")
	c.Storage.Site.StackEdit = (r.FormValue("stackedit") == "on")
	c.Storage.Site.Prism = (r.FormValue("prism") == "on")

	err = c.Storage.Save()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/dashboard/styles", http.StatusFound)
	return
}
