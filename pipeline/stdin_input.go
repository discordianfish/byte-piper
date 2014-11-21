package pipeline

import "os"

func init() {
	inputMap["stdin"] = newStdinInput
}

func newStdinInput(conf map[string]string) (input, error) {
	return os.Stdin, nil
}
