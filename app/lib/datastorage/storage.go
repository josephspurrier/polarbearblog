package datastorage

import (
	"encoding/json"

	"github.com/josephspurrier/polarbearblog/app/lib/envdetect"
	"github.com/josephspurrier/polarbearblog/app/model"
)

// Datastorer reads and writes data to an object.
type Datastorer interface {
	Save([]byte) error
	Load() ([]byte, error)
}

// Storage represents a writable and readable object.
type Storage struct {
	Site       *model.Site
	datastorer Datastorer
}

// New returns a writable and readable site object. Returns an error if the
// object cannot be initially read.
func New(ds Datastorer, site *model.Site) (*Storage, error) {
	s := &Storage{
		Site:       site,
		datastorer: ds,
	}

	err := s.Load()
	if err != nil {
		return nil, err
	}

	// Set the defaults for the site object.
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

// Save writes the site object to the data storage and returns an error if it
// cannot be written.
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

// Load reads the site object from the data storage and returns an error if
// it cannot be read.
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
