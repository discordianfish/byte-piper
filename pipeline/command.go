package pipeline

import "github.com/flynn/go-shlex"

func parseCommand(line string) (string, []string, error) {
	parts, err := shlex.Split(line)
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}
	return parts[0], args, err
}
