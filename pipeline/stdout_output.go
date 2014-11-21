package pipeline

import "os"

func init() {
	outputMap["stdout"] = newStdoutOutput
}

func newStdoutOutput(conf map[string]string) (output, error) {
	return os.Stdout, nil
}
