package pipeline

import (
	"errors"
	"os/exec"
)

func init() {
	outputMap["command"] = newCommandOutput
}

func newCommandOutput(conf map[string]string) (output, error) {
	c := conf["command"]
	if c == "" {
		return nil, errors.New("Require command")
	}
	cmd, args, err := parseCommand(c)
	if err != nil {
		return nil, err
	}
	command := exec.Command(cmd, args...)
	out, err := command.StdinPipe()
	if err != nil {
		return nil, err
	}
	return out, command.Start()
}
