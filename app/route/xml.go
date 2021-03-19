package route

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/josephspurrier/polarbearblog/app/lib/htmltemplate"
)

// XMLUtil -
type XMLUtil struct {
	*Core
}

func registerXMLUtil(c *XMLUtil) {
	c.Router.Get("/robots.txt", c.robots)
	c.Router.Get("/sitemap.xml", c.sitemap)
	c.Router.Get("/rss.xml", c.rss)
}

func (c *XMLUtil) robots(w http.ResponseWriter, r *http.Request) (status int, err error) {
	w.Header().Set("Content-Type", "text/plain")
	text :=
		`User-agent: *
Allow: /`
	fmt.Fprintf(w, text)
	return
}

func (c *XMLUtil) sitemap(w http.ResponseWriter, r *http.Request) (status int, err error) {
	// Resource: https://www.sitemaps.org/protocol.html
	// Resource: https://golang.org/src/encoding/xml/example_test.go

	type URL struct {
		Location     string `xml:"loc"`
		LastModified string `xml:"lastmod"`
	}

	type Sitemap struct {
		XMLName xml.Name `xml:"urlset"`
		XMLNS   string   `xml:"xmlns,attr"`
		XHTML   string   `xml:"xmlns:xhtml,attr"`
		URL     []URL    `xml:"url"`
	}

	m := &Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		XHTML: "http://www.w3.org/1999/xhtml",
	}

	// Home page
	m.URL = append(m.URL, URL{
		Location:     c.Storage.Site.SiteURL(),
		LastModified: c.Storage.Site.Updated.Format("2006-01-02"),
	})

	// Posts and pages
	for _, v := range c.Storage.Site.PostsAndPages(true) {
		m.URL = append(m.URL, URL{
			Location:     c.Storage.Site.SiteURL() + "/" + v.FullURL(),
			LastModified: v.Timestamp.Format("2006-01-02"),
		})
	}

	// Tags
	for _, v := range c.Storage.Site.Tags(true) {
		m.URL = append(m.URL, URL{
			Location:     c.Storage.Site.SiteURL() + "/blog?q=" + v.Name,
			LastModified: v.Timestamp.Format("2006-01-02"),
		})
	}

	output, err := xml.MarshalIndent(m, "  ", "    ")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	header := []byte(xml.Header)
	output = append(header[:], output[:]...)

	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprintf(w, string(output))
	return
}

func (c *XMLUtil) rss(w http.ResponseWriter, r *http.Request) (status int, err error) {
	// Resource: https://www.rssboard.org/rss-specification
	// Rsource: https://validator.w3.org/feed/check.cgi

	type Item struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		PubDate     string `xml:"pubDate"`
		GUID        string `xml:"guid"`
		Description string `xml:"description"`
	}

	type AtomLink struct {
		Href string `xml:"href,attr"`
		Rel  string `xml:"rel,attr"`
		Type string `xml:"type,attr"`
	}

	type Sitemap struct {
		XMLName       xml.Name `xml:"rss"`
		Version       string   `xml:"version,attr"`
		Atom          string   `xml:"xmlns:atom,attr"`
		Title         string   `xml:"channel>title"`
		Link          string   `xml:"channel>link"`
		Description   string   `xml:"channel>description"`
		Generator     string   `xml:"channel>generator"`
		Language      string   `xml:"channel>language"`
		LastBuildDate string   `xml:"channel>lastBuildDate"`
		AtomLink      AtomLink `xml:"channel>atom:link"`
		Items         []Item   `xml:"channel>item"`
	}

	m := &Sitemap{
		Version:       "2.0",
		Atom:          "http://www.w3.org/2005/Atom",
		Title:         c.Storage.Site.SiteTitle(),
		Link:          c.Storage.Site.SiteURL(),
		Description:   c.Storage.Site.Description,
		Generator:     "Polar Bear Blog",
		Language:      "en-us",
		LastBuildDate: time.Now().Format(time.RFC1123Z),
		AtomLink: AtomLink{
			Href: c.Storage.Site.SiteURL() + "/rss.xml",
			Rel:  "self",
			Type: "application/rss+xml",
		},
	}

	for _, v := range c.Storage.Site.PostsAndPages(true) {
		plaintext := htmltemplate.PlaintextBlurb(v.Post.Content)
		m.Items = append(m.Items, Item{
			Title:       v.Title,
			Link:        c.Storage.Site.SiteURL() + "/" + v.FullURL(),
			PubDate:     v.Timestamp.Format(time.RFC1123Z),
			GUID:        c.Storage.Site.SiteURL() + "/" + v.FullURL(),
			Description: plaintext,
		})
	}

	output, err := xml.MarshalIndent(m, "  ", "    ")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	header := []byte(xml.Header)
	output = append(header[:], output[:]...)

	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprintf(w, string(output))
	return
}
