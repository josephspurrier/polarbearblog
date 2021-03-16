package datastorage

import (
	"encoding/json"

	"github.com/josephspurrier/polarbearblog/app/lib/envdetect"
	"github.com/josephspurrier/polarbearblog/app/model"
)

// Datastorer -
type Datastorer interface {
	Save([]byte) error
	Load() ([]byte, error)
}

// Storage -
type Storage struct {
	Site       *model.Site
	datastorer Datastorer
}

// New is a way to interact with the site storage.
func New(ds Datastorer, site *model.Site) (*Storage, error) {
	s := &Storage{
		Site:       site,
		datastorer: ds,
	}

	err := s.Load()
	if err != nil {
		return nil, err
	}

	// Save to storage. Ensure the posts exists first so it doesn't error.
	if s.Site.Posts == nil {
		s.Site.Posts = make(map[string]model.Post)
	}
	// Ensure redirects don't try to happen if the scheme is empty.
	if s.Site.Scheme == "" {
		s.Site.Scheme = "http"
	}
	// Ensure it's set to the login page works.
	if s.Site.LoginURL == "" {
		s.Site.LoginURL = "admin"
	}

	return s, nil
}

func (s *Storage) Save() error {
	var b []byte
	var err error

	if envdetect.RunningLocalDev() {
		// Indent so the data is easy to read.
		b, err = json.MarshalIndent(s.Site, "", "    ")
	} else {
		b, err = json.Marshal(s.Site)
	}

	if err != nil {
		return err
	}

	err = s.datastorer.Save(b)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Load() error {
	b, err := s.datastorer.Load()
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, s.Site)
	if err != nil {
		return err
	}

	return nil
}
