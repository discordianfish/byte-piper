package pipeline

import "os"

type fileInput struct {
	filename string
}

func init() {
	inputMap["file"] = newFileInput
}

func newFileInput(conf map[string]string) (input, error) {
	return os.Open(conf["path"])
}
