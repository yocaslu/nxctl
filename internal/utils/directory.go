package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"nxctl/internal/nxlog"
	"os"
)

var MODULE_NAME string = "Utils.Directory"

func DirExist(path string) (bool, error) {
	logger := nxlog.ForModule(MODULE_NAME + ".DirExist")
	logger.Info(
		"Reading directory information",
		slog.String("path", path),
	)
	_, err := os.Stat(path)
	if err != nil {
		logger.Error(
			"Failed to read directory",
			slog.String("path", path),
			slog.String("error", err.Error()),
		)

		logger.Info(
			"Checking if path is a directory",
			slog.String("path", path),
		)

		if errors.Is(err, os.ErrNotExist) {
			return false, nil // Directory does not exist
		}

		// Handle other potential errors here (e.g., permission denied)
		return false, fmt.Errorf("%s", err)
	}

	return true, nil // Returns true if it is a directory, false if it's a file
}
