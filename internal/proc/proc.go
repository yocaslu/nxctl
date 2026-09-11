package proc

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func Run(_env []string, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = _env
	fullCmd := command + " " + strings.Join(args, " ")

	var stderr strings.Builder
	var stdout strings.Builder
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Start(); err != nil {
		errString := stderr.String()
		return stdout.String(), fmt.Errorf("%s\n", errString)
	}

	log.Printf("Executing: [%s]", fullCmd)
	if err := cmd.Wait(); err != nil {
		errString := stderr.String()
		return stdout.String(), fmt.Errorf("%s\n", errString)
	}

	return stdout.String(), nil
}
