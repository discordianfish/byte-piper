package pipeline

import (
	"errors"

	"github.com/rlmcpherson/s3gof3r"
)

func init() {
	inputMap["s3"] = newS3Input
}

func newS3Input(conf map[string]string) (input, error) {
	bucketName := conf["bucket"]
	if bucketName == "" {
		return nil, errors.New("No bucket specified")
	}
	fileName := conf["filename"]
	if fileName == "" {
		return nil, errors.New("No file name specified")
	}

	keys, err := s3gof3r.EnvKeys()
	if err != nil {
		return nil, err
	}
	s3 := s3gof3r.New(conf["endpoint"], keys)
	bucket := s3.Bucket(bucketName)

	r, _, err := bucket.GetReader(fileName, nil)
	return r, err
}
