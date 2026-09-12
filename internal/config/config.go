package config

import (
	"fmt"
	"log"
	"nxctl/internal/utils"
	"os"
	"path"
)

type Config struct {
	Date       string
	BackupPath string
}

func Load() (*Config, error) {
	backup_path := os.Getenv("BACKUP_PATH")
	if backup_path == "" {
		return nil, fmt.Errorf("BACKUP_PATH is missing!")
	}

	date := utils.GetDate()
	return &Config{
		Date:       date,
		BackupPath: backup_path,
	}, nil
}

func (c *Config) CreateBackupDir() error {
	log.Printf("Creating backup directory: %s\n", c.BackupPath)
	err := os.Mkdir(path.Join(c.BackupPath, utils.GetDate()), 0775)

	if err != nil {
		log.Printf("failed to create backup directory: %s\n", err)
		return err
	}

	return nil
}
