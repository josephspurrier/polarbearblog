package datastorage

import (
	"io/ioutil"
)

// LocalStorage -
type LocalStorage struct {
	path string
}

// NewLocalStorage -
func NewLocalStorage(storagePath string) *LocalStorage {
	return &LocalStorage{
		path: storagePath,
	}
}

// Load -
func (s *LocalStorage) Load() ([]byte, error) {
	b, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Save -
func (s *LocalStorage) Save(b []byte) error {
	err := ioutil.WriteFile(s.path, b, 0644)
	if err != nil {
		return err
	}

	return nil
}
