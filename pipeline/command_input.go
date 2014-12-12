package pipeline

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

func init() {
	inputMap["command"] = newCommandInput
}

func newCommandInput(conf map[string]string) (input, error) {
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
	command.Stderr = os.Stderr
	out, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := command.Start(); err != nil {
		return nil, err
	}

	go func() {
		if err := command.Wait(); err != nil {
			log.Print(err)
		}
		out.Close()
	}()
	return out, nil

}
