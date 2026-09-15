package proc

import (
	"fmt"
	"log/slog"
	"nxctl/internal/nxlog"
	"os/exec"
	"strings"
)

var MODULE_NAME string = "proc"

func Run(_env []string, command string, args ...string) (string, error) {
	logger := nxlog.ForModule(MODULE_NAME + ".Run")

	cmd := exec.Command(command, args...)
	cmd.Env = _env
	fullCmd := command + " " + strings.Join(args, " ")

	var stderr strings.Builder
	var stdout strings.Builder
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	logger.Debug("Starting command", slog.String("command", fullCmd))
	if err := cmd.Start(); err != nil {
		return stdout.String(), fmt.Errorf("%s", stderr.String())
	}

	logger.Debug("Waiting command to finish", slog.String("command", fullCmd))
	if err := cmd.Wait(); err != nil {
		return stdout.String(), fmt.Errorf("%s", stderr.String())
	}

	return stdout.String(), nil
}
