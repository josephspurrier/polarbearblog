package datastorage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
)

// GCPStorage -
type GCPStorage struct {
	path      string
	projectID string
	bucket    string
}

// NewGCPStorage -
func NewGCPStorage(storagePath string, projectID string, bucket string) *GCPStorage {
	return &GCPStorage{
		path:      storagePath,
		projectID: projectID,
		bucket:    bucket,
	}
}

// Load -
func (s *GCPStorage) Load() ([]byte, error) {
	b, err := downloadFile(s.bucket, s.path)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Save -
func (s *GCPStorage) Save(b []byte) error {
	err := uploadFile(s.bucket, s.path, b)
	if err != nil {
		return err
	}

	return nil
}

// downloadFile downloads an object.
func downloadFile(bucket, object string) ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	return data, nil
}

// uploadFile uploads an object.
func uploadFile(bucket, object string, data []byte) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	f := bytes.NewReader(data)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
