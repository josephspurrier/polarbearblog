package model

import (
	"fmt"
	"sort"
	"time"
)

// Site -
type Site struct {
	Title             string          `json:"title"`
	Subtitle          string          `json:"subtitle"`
	Author            string          `json:"author"`
	Favicon           string          `json:"favicon"`
	Description       string          `json:"description"`
	Footer            string          `json:"footer"`
	Scheme            string          `json:"scheme"`
	URL               string          `json:"url"`
	LoginURL          string          `json:"loginurl"`
	GoogleAnalyticsID string          `json:"googleanalytics"`
	DisqusID          string          `json:"disqus"`
	Created           time.Time       `json:"created"`
	Updated           time.Time       `json:"updated"`
	Content           string          `json:"content"` // Home content.
	Styles            string          `json:"styles"`
	StylesAppend      bool            `json:"stylesappend"`
	StackEdit         bool            `json:"stackedit"`
	Prism             bool            `json:"prism"`
	Posts             map[string]Post `json:"posts"`
}

// SiteURL -
func (s Site) SiteURL() string {
	return fmt.Sprintf("%v://%v", s.Scheme, s.URL)
}

// SiteTitle -
func (s Site) SiteTitle() string {
	return fmt.Sprintf("%v", s.Title)
}

// SiteSubtitle -
func (s Site) SiteSubtitle() string {
	return fmt.Sprintf("%v", s.Subtitle)
}

// PublishedPosts -
func (s Site) PublishedPosts() []Post {
	arr := make(PostList, 0)
	for _, v := range s.Posts {
		if v.Published && !v.Page {
			arr = append(arr, v)
		}
	}

	sort.Sort(sort.Reverse(arr))

	return arr
}

// PublishedPages -
func (s Site) PublishedPages() []Post {
	arr := make(PostList, 0)
	for _, v := range s.Posts {
		if v.Published && v.Page {
			arr = append(arr, v)
		}
	}

	sort.Sort(sort.Reverse(arr))

	return arr
}

// PostsAndPages -
func (s Site) PostsAndPages(onlyPublished bool) PostWithIDList {
	arr := make(PostWithIDList, 0)
	for k, v := range s.Posts {
		if onlyPublished && !v.Published {
			continue
		}

		p := PostWithID{Post: v, ID: k}
		arr = append(arr, p)
	}

	sort.Sort(sort.Reverse(arr))

	return arr
}

// Tags -
func (s Site) Tags(onlyPublished bool) TagList {
	// Get unique values.
	m := make(map[string]Tag)
	for _, v := range s.Posts {
		if onlyPublished && !v.Published {
			continue
		}

		for _, t := range v.Tags {
			m[t.Name] = t
		}
	}

	// Create unsorted tag list.
	arr := make(TagList, 0)
	for _, v := range m {
		arr = append(arr, v)
	}

	// Sort by name.
	sort.Sort(arr)

	return arr
}

// PostBySlug -
func (s Site) PostBySlug(slug string) PostWithID {
	// FIXME: This needs to be optimized.
	var p PostWithID
	for k, v := range s.Posts {
		if v.URL == slug {
			p = PostWithID{
				Post: v,
				ID:   k,
			}
			break
		}
	}

	return p
}
