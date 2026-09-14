package config

import (
	"fmt"
	"nxctl/internal/utils"
	"os"
	"path"
)

type Config struct {
	Date       string
	BackupPath string
}

func Load() (*Config, error) {
	date := utils.GetDate()
	backup_path := os.Getenv("BACKUP_PATH")
	if backup_path == "" {
		return nil, fmt.Errorf("BACKUP_PATH is missing!")
	}

	// Creating dedicated backup directory based on date
	backup_path = path.Join(backup_path, "nextcloud_backup-"+date)

	exist, err := utils.DirExist(backup_path)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if !exist {
		err = createBackupDir(backup_path)
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}
	}

	return &Config{
		Date:       date,
		BackupPath: backup_path,
	}, nil
}

func createBackupDir(backup_path string) error {
	err := os.Mkdir(backup_path, 0775)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}
