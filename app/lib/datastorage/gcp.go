package datastorage

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
)

// GCPStorage represents a Google Cloud Storage object.
type GCPStorage struct {
	bucket string
	object string
}

// NewGCPStorage returns a Google Cloud storage item given a bucket and an object
// path.
func NewGCPStorage(bucket string, object string) *GCPStorage {
	return &GCPStorage{
		bucket: bucket,
		object: object,
	}
}

// Load downloads an object from a bucket and returns an error if it cannot
// be read.
func (s *GCPStorage) Load() ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(s.bucket).Object(s.object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Save uploads an object to a bucket and returns an error if it cannot be
// written.
func (s *GCPStorage) Save(b []byte) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	f := bytes.NewReader(b)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(s.bucket).Object(s.object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}
