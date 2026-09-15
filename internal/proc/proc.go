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

	logger.Info("Starting command", slog.String("command", fullCmd))
	if err := cmd.Start(); err != nil {
		logger.Error("Failed to start command",
			slog.Any("error", err),
			slog.String("stderr", stderr.String()),
			slog.String("stdout", stdout.String()),
			slog.String("command", fullCmd),
			slog.Any("env", _env),
		)
		return stdout.String(), fmt.Errorf("%s", stderr.String())
	}

	logger.Info("Waiting command to finish", slog.String("command", fullCmd))
	if err := cmd.Wait(); err != nil {
		logger.Error("Command failed during execution",
			slog.Any("error", err),
			slog.String("stderr", stderr.String()),
			slog.String("stdout", stdout.String()),
			slog.String("command", fullCmd),
			slog.Any("env", _env),
		)
		return stdout.String(), fmt.Errorf("%s", stderr.String())
	}

	return stdout.String(), nil
}
