package model

import (
	"strings"
	"time"
)

// PostWithIDList -
type PostWithIDList []PostWithID

func (t PostWithIDList) Len() int {
	return len(t)
}
func (t PostWithIDList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t PostWithIDList) Less(i, j int) bool {
	if t[i].Timestamp.Equal(t[j].Timestamp) {
		return t[i].Title > t[j].Title // Sort by title ASC
	} else if t[i].Timestamp.Before(t[j].Timestamp) {
		return true // Sort by timestamp, DESC
	}

	return false
}

// PostList -
type PostList []Post

func (t PostList) Len() int {
	return len(t)
}
func (t PostList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t PostList) Less(i, j int) bool {
	if t[i].Timestamp.Equal(t[j].Timestamp) {
		return t[i].Title > t[j].Title // Sort by title ASC
	} else if t[i].Timestamp.Before(t[j].Timestamp) {
		return true // Sort by timestamp, DESC
	}

	return false
}

// Post -
type Post struct {
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Canonical string    `json:"canonical"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	Page      bool      `json:"page"`
	Tags      TagList   `json:"tags"`
}

// PostWithID -
type PostWithID struct {
	Post
	ID string `json:"id"`
}

// FullURL -
func (p *Post) FullURL() string {
	return p.URL
}

// TagList -
type TagList []Tag

func (t TagList) Len() int {
	return len(t)
}
func (t TagList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t TagList) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// String -
func (t TagList) String() string {
	arr := make([]string, 0)
	for _, v := range t {
		arr = append(arr, v.Name)
	}

	return strings.Join(arr, ",")
}

// Split -
func (t TagList) Split(s string) TagList {
	trimmed := strings.TrimSpace(s)

	// Return an empty object since split returns 1 element when empty.
	if len(trimmed) == 0 {
		return TagList{}
	}

	ts := time.Now()

	arrTags := make([]Tag, 0)
	tags := strings.Split(trimmed, ",")
	for _, v := range tags {
		arrTags = append(arrTags, Tag{
			Name:      strings.TrimSpace(v),
			Timestamp: ts,
		})
	}

	return arrTags
}
