package proc

import (
	"fmt"
	"os/exec"
	"strings"
)

func Run(_env []string, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = _env
	// fullCmd := command + " " + strings.Join(args, " ")

	var stderr strings.Builder
	var stdout strings.Builder
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Start(); err != nil {
		return stdout.String(), fmt.Errorf("%s\n", stderr.String())
	}

	if err := cmd.Wait(); err != nil {
		return stdout.String(), fmt.Errorf("%s\n", stderr.String())
	}

	return stdout.String(), nil
}
