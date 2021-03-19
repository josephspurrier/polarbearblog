package htmltemplate

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/josephspurrier/polarbearblog/app/lib/datastorage"
	"github.com/josephspurrier/polarbearblog/app/lib/websession"
	"github.com/josephspurrier/polarbearblog/app/model"
)

// New returns a HTML template engine.
func New(ds *datastorage.Storage, sess *websession.Session, sanitize bool) *Engine {
	return &Engine{
		storage:  ds,
		sess:     sess,
		sanitize: sanitize,
	}
}

// Engine represents a HTML template engine.
type Engine struct {
	storage  *datastorage.Storage
	sess     *websession.Session
	sanitize bool
}

// Template renders HTML to a response writer and returns a 200 status code and
// an error if one occurs.
func (te *Engine) Template(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.internalTemplate(w, r, mainTemplate, partialTemplate, http.StatusOK, vars)
}

// ErrorTemplate renders HTML to a response writer and returns a 404 status code
// and an error if one occurs.
func (te *Engine) ErrorTemplate(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.internalTemplate(w, r, mainTemplate, partialTemplate, http.StatusNotFound, vars)
}

// internalTemplate converts content from markdown to HTML and then outputs to
// a response writer. Returns an HTTP status code and an error if one occurs.
func (te *Engine) internalTemplate(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, statusCode int, vars map[string]interface{}) (status int, err error) {
	// Functions available in the templates.
	fm := te.funcMap(r)

	baseTemplate := fmt.Sprintf("html/%v.tmpl", mainTemplate)
	headerTemplate := "html/partial/head.tmpl"
	contentTemplate := fmt.Sprintf("html/partial/%v.tmpl", partialTemplate)

	// Parse the main template with the functions.
	t, err := template.New(path.Base(baseTemplate)).Funcs(fm).ParseFiles(baseTemplate,
		headerTemplate, contentTemplate)
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
	// Functions available in the templates.
	fm := te.funcMap(r)

	// Display 404 if not found.
	if post.URL == "" {
		return http.StatusNotFound, nil
	}

	baseTemplate := fmt.Sprintf("html/%v.tmpl", mainTemplate)
	headerTemplate := "html/partial/head.tmpl"

	// Parse the main template with the functions.
	t, err := template.New(path.Base(baseTemplate)).Funcs(fm).ParseFiles(baseTemplate, headerTemplate)
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
