package pipeline

import (
	"errors"
	"os/exec"

	shlex "github.com/flynn/go-shlex"
)

func init() {
	inputMap["command"] = newCommandInput
}

func newCommandInput(conf map[string]string) (input, error) {
	c := conf["command"]
	if c == "" {
		return nil, errors.New("Require command")
	}
	cmd, err := shlex.Split(c)
	if err != nil {
		return nil, err
	}
	args := []string{}
	if len(args) > 1 {
		args = cmd[1:]
	}
	command := exec.Command(cmd[0], args...)
	out, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	return out, command.Start()
}
