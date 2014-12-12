package pipeline

import (
	"errors"
	"os"
)

type fileInput struct {
	filename string
}

func init() {
	inputMap["file"] = newFileInput
}

func newFileInput(conf map[string]string) (input, error) {
	if conf["path"] == "" {
		return nil, errors.New("path required")
	}
	return os.Open(conf["path"])
}
