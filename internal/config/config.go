package config

import (
	"fmt"
	"log/slog"
	"nxctl/internal/utils"
	"nxctl/internal/utils/nxlog"
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
var BACKUP_DIRNAME string = "backup_"

func Load() (*Config, error) {
	logger := nxlog.ForModule(MODULE_NAME + ".Load")
	date := utils.GetDate()
	backupPath := os.Getenv("BACKUP_PATH")

	if backupPath == "" {
		return nil, fmt.Errorf("Failed to read BACKUP_PATH environment variable")
	}

	exist, err := utils.DirExist(backupPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to check if [%s] exist: [%w]", backupPath, err)
	}

	if !exist {
		err = utils.CreateDatedPath(backupPath, BACKUP_DIRNAME, date)
		if err != nil {
			return nil, fmt.Errorf(
				"Failed to create dated backup directory [%s/%s%s]: [%w]",
				backupPath,
				BACKUP_DIRNAME,
				date,
				err,
			)
		}
	}

	cfg := &Config{
		Date:       date,
		BackupPath: backupPath,
	}

	logger.Debug(
		"Config object created",
		slog.Any("config", cfg),
	)

	return cfg, nil
}

func CreateDatedBackupDir(backup_path string, prefix string, date string) error {
	datedpath := path.Join(backup_path, prefix+date)
	if err := utils.CreateDir(datedpath); err != nil {
		return fmt.Errorf("Failed to create directory [%s]: [%w]", datedpath, err)
	}

	return nil
}
