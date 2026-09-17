package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"nxctl/internal/utils/nxlog"
	"os"
	"path"
)

var MODULE_NAME string = "utils.directory"

func DirExist(path string) (bool, error) {
	logger := nxlog.ForModule(MODULE_NAME + ".DirExist")
	logger.Debug(
		"Reading directory information",
		slog.String("path", path),
	)
	_, err := os.Stat(path)
	if err != nil {
		logger.Debug(
			"Checking if path is a directory",
			slog.String("path", path),
		)

		if errors.Is(err, os.ErrNotExist) {
			return false, nil // Directory does not exist
		}

		// Handle other potential errors here (e.g., permission denied)
		return false, err
	}

	return true, nil // Returns true if it is a directory, false if it's a file
}

func CreateDir(path string) error {
	logger := nxlog.ForModule(MODULE_NAME + ".CreateDir")

	err := os.MkdirAll(path, 0775)
	if err != nil {
		return fmt.Errorf("Failed to create directory [%s]: %w", path, err)
	}

	logger.Debug("Created directory",
		slog.String("path", path),
	)

	return nil
}

func CreateDatedPath(pathdir string, dirname string, date string) error {
	logger := nxlog.ForModule(MODULE_NAME + ".CreatedDatedPath")

	datedPath := path.Join(pathdir, dirname+date)
	exist, err := DirExist(datedPath)
	if err != nil {
		return fmt.Errorf("Failed to check if [%s] exist: %w", datedPath, err)
	} else if exist {
		return nil
	}

	if err := CreateDir(datedPath); err != nil {
		return fmt.Errorf("Failed to create dated directory [%s]: %w", pathdir, err)
	}

	logger.Debug(
		"Created dated path",
		slog.String("pathdir", pathdir),
		slog.String("filename", dirname),
		slog.String("date", date),
	)

	return nil
}
