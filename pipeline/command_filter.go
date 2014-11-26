package pipeline

import (
	"errors"
	"io"
	"log"
	"os/exec"

	shlex "github.com/flynn/go-shlex"
)

func init() {
	filterMap["command"] = newCommandFilter
}

type commandFilter struct {
	stdin  io.Writer
	stdout io.Reader
}

func newCommandFilter(conf map[string]string) (filter, error) {
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
	in, err := command.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	return &commandFilter{
		stdin:  in,
		stdout: out,
	}, command.Run()
}

func (f *commandFilter) Link(r io.Reader) error {
	go func() {
		if _, err := copy(f.stdin, r); err != nil {
			log.Print(err)
			return
		}
	}()
	return nil
}

func (f *commandFilter) Read(p []byte) (int, error) {
	return f.stdout.Read(p)
}
