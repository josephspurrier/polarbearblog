package htmltemplate

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/josephspurrier/polarbearblog/app/lib/datastorage"
	"github.com/josephspurrier/polarbearblog/app/lib/websession"
	"github.com/josephspurrier/polarbearblog/app/model"
)

// New -
func New(ds *datastorage.Storage, sess *websession.Session) *Engine {
	return &Engine{
		storage: ds,
		sess:    sess,
	}
}

// Engine -
type Engine struct {
	storage *datastorage.Storage
	sess    *websession.Session
}

// Template -
func (te *Engine) Template(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.internalTemplate(w, r, mainTemplate, partialTemplate, http.StatusOK, vars)
}

// ErrorTemplate -
func (te *Engine) ErrorTemplate(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, vars map[string]interface{}) (status int, err error) {
	return te.internalTemplate(w, r, mainTemplate, partialTemplate, http.StatusNotFound, vars)
}

// internalTemplate -
func (te *Engine) internalTemplate(w http.ResponseWriter, r *http.Request, mainTemplate string,
	partialTemplate string, statusCode int, vars map[string]interface{}) (status int, err error) {
	// Functions available in the templates.
	fm := te.funcMap(r)

	// Load main template.
	base, err := ioutil.ReadFile(fmt.Sprintf("html/%v.tmpl", mainTemplate))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the main template with the functions.
	t, err := template.New("root").Funcs(fm).Parse(string(base))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Read the partial.
	content, err := ioutil.ReadFile(fmt.Sprintf("html/partial/%v.tmpl", partialTemplate))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the partial with the main template.
	t, err = t.Parse(string(content))
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

// Post -
func (te *Engine) Post(w http.ResponseWriter, r *http.Request, mainTemplate string,
	post model.Post, vars map[string]interface{}) (status int, err error) {
	// Functions available in the templates.
	fm := te.funcMap(r)

	// Display 404 if not found.
	if post.URL == "" {
		return http.StatusNotFound, nil
	}

	// Load main template.
	base, err := ioutil.ReadFile(fmt.Sprintf("html/%v.tmpl", mainTemplate))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the main template with the functions.
	t, err := template.New("root").Funcs(fm).Parse(string(base))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Parse the content.
	t, err = sanitizedContent(t, post.Content)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Execute the template.
	if err := t.Execute(w, vars); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
