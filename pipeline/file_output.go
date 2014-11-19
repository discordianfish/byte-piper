package pipeline

import "os"

func init() {
	outputMap["file"] = newFileOutput
}

func newFileOutput(conf map[string]string) (output, error) {
	return os.Create(conf["path"])
}
