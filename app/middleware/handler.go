package middleware

import (
	"github.com/josephspurrier/polarbearblog/app/lib/htmltemplate"
	"github.com/josephspurrier/polarbearblog/app/lib/router"
	"github.com/josephspurrier/polarbearblog/app/lib/websession"
)

// Handler -
type Handler struct {
	Router     *router.Mux
	Render     *htmltemplate.Engine
	Sess       *websession.Session
	SiteURL    string
	SiteScheme string
}

// NewHandler -
func NewHandler(te *htmltemplate.Engine, sess *websession.Session, mux *router.Mux, siteURL string, siteScheme string) *Handler {
	return &Handler{
		Render:     te,
		Router:     mux,
		Sess:       sess,
		SiteURL:    siteURL,
		SiteScheme: siteScheme,
	}

}
