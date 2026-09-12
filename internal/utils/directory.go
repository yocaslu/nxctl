package utils

import (
	"errors"
	"log"
	"os"
)

func DirExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil // Directory does not exist
		}

		// Handle other potential errors here (e.g., permission denied)
		log.Printf("failed to read %s: %s\n", path, err)
		return false, err
	}

	return true, nil // Returns true if it is a directory, false if it's a file
}
