package htmltemplate

import (
	"net/http"

	"github.com/josephspurrier/polarbearblog/app/model"
	"github.com/josephspurrier/polarbearblog/html"
)

// New returns a HTML template engine.
func New(manager *html.TemplateManager, allowUnsafeHTML bool) *Engine {
	return &Engine{
		allowUnsafeHTML: allowUnsafeHTML,
		manager:         manager,
	}
}

// Engine represents a HTML template engine.
type Engine struct {
	allowUnsafeHTML bool
	manager         *html.TemplateManager
}

// Template renders HTML to a response writer and returns a 200 status code and
// an error if one occurs.
func (te *Engine) Template(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.partialTemplate(w, r, mainTemplate, partialTemplate, http.StatusOK, vars)
}

// ErrorTemplate renders HTML to a response writer and returns a 404 status code
// and an error if one occurs.
func (te *Engine) ErrorTemplate(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.partialTemplate(w, r, mainTemplate, partialTemplate, http.StatusNotFound, vars)
}

// partialTemplate converts content from markdown to HTML and then outputs to
// a response writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) partialTemplate(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, statusCode int, vars map[string]interface{}) (status int, err error) {
	// Parse the template.
	t, err := te.manager.PartialTemplate(r, mainTemplate, partialTemplate)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Output the status code.
	w.WriteHeader(statusCode)

	// Execute the template.
	if err := t.Execute(w, vars); err != nil {
		return http.StatusInternalServerError, err
	}

	return statusCode, nil
}

// Post converts a site post from markdown to HTML and then outputs to response
// writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) Post(w http.ResponseWriter, r *http.Request, mainTemplate string,
	post model.Post, vars map[string]interface{}) (status int, err error) {
	// Display 404 if not found.
	if post.URL == "" {
		return http.StatusNotFound, nil
	}

	// Parse the template.
	t, err := te.manager.PostTemplate(r, mainTemplate)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the content.
	t, err = te.sanitizedContent(t, post.Content)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Execute the template.
	if err := t.Execute(w, vars); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
