package pipeline

import (
	"errors"

	"github.com/rlmcpherson/s3gof3r"
)

func init() {
	outputMap["s3"] = newS3Output
}

func newS3Output(conf map[string]string) (output, error) {
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

	return bucket.PutWriter(fileName, nil, nil)
}
