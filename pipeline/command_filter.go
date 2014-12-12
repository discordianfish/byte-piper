package pipeline

import (
	"errors"
	"io"
	"log"
	"os/exec"
)

func init() {
	filterMap["command"] = newCommandFilter
}

type commandFilter struct {
	stdout  io.Reader
	command *exec.Cmd
}

func newCommandFilter(conf map[string]string) (filter, error) {
	c := conf["command"]
	if c == "" {
		return nil, errors.New("Require command")
	}
	cmd, args, err := parseCommand(c)
	if err != nil {
		return nil, err
	}
	log.Printf("cmd: %s, args: %#v (from %s)", cmd, args, c)
	command := exec.Command(cmd, args...)
	out, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	return &commandFilter{
		stdout:  out,
		command: command,
	}, nil
}

func (f *commandFilter) Link(r io.Reader) error {
	f.command.Stdin = r
	go func() {
		if err := f.command.Run(); err != nil {
			log.Printf("Couldn't execute %s: %s", f.command.Path, err)
		}
	}()
	return nil
}

func (f *commandFilter) Read(p []byte) (int, error) {
	return f.stdout.Read(p)
}
