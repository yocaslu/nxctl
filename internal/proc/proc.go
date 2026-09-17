package proc

import (
	"fmt"
	"os/exec"
	"strings"
)

var MODULE_NAME string = "proc"

func Run(_env []string, command string, args ...string) (string, error) {
	// logger := nxlog.ForModule(MODULE_NAME + ".Run")

	cmd := exec.Command(command, args...)
	cmd.Env = _env

	bOutput, err := cmd.CombinedOutput() // Run cmd and collect stdout and stderr
	strOutput := strings.TrimSpace(string(bOutput))

	if err != nil {
		// logger.Debug(
		// 	"Failed while running command",
		// 	slog.String("cmd", cmd.String()),
		// 	slog.String("output", strOutput),
		// 	slog.Any("error", err),
		// )

		if strOutput != "" {
			return strOutput, fmt.Errorf(
				"Failed to run command [%v] (%w): %s",
				command,
				err,
				strOutput,
			)
		}

		return strOutput, fmt.Errorf(
			"Failed to run command [%v]: %w",
			command,
			err,
		)
	}

	return strOutput, nil
}
