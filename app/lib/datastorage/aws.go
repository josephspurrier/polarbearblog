package datastorage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func EnsureBaseDir(fpath string) error {
	baseDir := path.Dir(fpath)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(baseDir, 0755)
}

// GCPStorage represents a S3 Storage object.
type S3Storage struct {
	bucket string
	object string
	region string
}

// NewGCPStorage returns a S3 storage item given a bucket and an object
// path.
func NewS3Storage(bucket string, object string, region string) *S3Storage {
	return &S3Storage{
		bucket: bucket,
		object: object,
		region: region,
	}
}

// Load downloads an object from a bucket and returns an error if it cannot
// be read.
func (s *S3Storage) Load() ([]byte, error) {
	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s.region)},
	)
	// Log any session errors
	if err != nil {
		log.Fatalf("Unable to create a new s3 session: %v", err)
		return nil, err
	}
	// create a new AWS S3 downloader
	downloader := s3manager.NewDownloader(sess)

	// Download the item from the bucket. If an error occurs, log it and exit. Otherwise, notify the user that the download succeeded.
	EnsureBaseDir("/tmp/" + s.object)
	file, err := os.Create("/tmp/" + s.object)
	if err != nil {
		log.Fatalf("Unable to create file. %v", err)
		return nil, err
	}

	defer file.Close()

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(s.object),
		})

	if err != nil {
		log.Fatalf("Unable to download item %q, %v", s.object, err)
		return nil, err
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Save uploads an object to a bucket and returns an error if it cannot be
// written.
func (s *S3Storage) Save(b []byte) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s.region)},
	)
	if err != nil {
		// Print the error and exit.
		fmt.Printf("Unable to create session %v", err)
		return err
	}
	uploader := s3manager.NewUploader(sess)

	file := bytes.NewReader(b)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.object),
		Body:   file,
	})

	if err != nil {
		// Print the error and exit.
		fmt.Printf("Unable to upload file to %q, %v", s.bucket, err)
		return err
	}

	fmt.Printf("Successfully uploaded file to %q\n", s.bucket)

	return nil
}
