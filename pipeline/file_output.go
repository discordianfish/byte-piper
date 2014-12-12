package pipeline

import (
	"errors"
	"os"
)

func init() {
	outputMap["file"] = newFileOutput
}

func newFileOutput(conf map[string]string) (output, error) {
	if conf["path"] == "" {
		return nil, errors.New("path required")
	}
	return os.Create(conf["path"])
}
