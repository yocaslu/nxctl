package config

import (
	"fmt"
	"log/slog"
	"nxctl/internal/nxlog"
	"nxctl/internal/utils"
	"os"
	"path"
)

type Config struct {
	Date       string `json:"date"`
	BackupPath string `json:"backup_path"`
}

func (c *Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("date", c.Date),
		slog.String("backup_path", c.BackupPath),
	)
}

var MODULE_NAME string = "config"

func Load() (*Config, error) {
	logger := nxlog.ForModule(MODULE_NAME + ".Load")
	date := utils.GetDate()
	backup_path := os.Getenv("BACKUP_PATH")

	if backup_path == "" {
		return nil, fmt.Errorf("BACKUP_PATH is missing!")
	}

	// Creating dedicated backup directory based on date
	backup_path = path.Join(backup_path, "nextcloud_backup-"+date)

	exist, err := utils.DirExist(backup_path)
	if err != nil {
		return nil, fmt.Errorf("Failed to check if backup directory exist: %s", err)
	}

	if !exist {
		logger.Debug(
			"Backup directory does not exist",
			slog.String("backup_path", backup_path),
		)

		err = createBackupDir(backup_path)
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}
	}

	cfg := &Config{
		Date:       date,
		BackupPath: backup_path,
	}

	logger.Debug(
		"Config object created",
		slog.Any("config", cfg),
	)

	return cfg, nil
}

func createBackupDir(backup_path string) error {
	logger := nxlog.ForModule(MODULE_NAME + ".createBackupDir")
	logger.Debug("Creating backup directory")
	err := os.Mkdir(backup_path, 0775)

	if err != nil {
		return fmt.Errorf("Failed to create backup directory: %s", err)
	}

	return nil
}
