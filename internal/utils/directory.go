package utils

import (
	"errors"
	"fmt"
	"os"
)

func DirExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil // Directory does not exist
		}

		// Handle other potential errors here (e.g., permission denied)
		return false, fmt.Errorf("%s", err)
	}

	return true, nil // Returns true if it is a directory, false if it's a file
}
