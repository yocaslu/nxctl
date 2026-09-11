package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	Date       string
	BackupPath string
}

func Load() (*Config, error) {

	backup_volume := os.Getenv("BACKUP_PATH")
	if backup_volume == "" {
		return nil, fmt.Errorf("BACKUP_PATH is missing!")
	}

	date := time.Now().Format("02-01-2006")
	backup_dir := backup_volume

	return &Config{
		Date:       date,
		BackupPath: backup_dir,
	}, nil
}

func (c *Config) CreateBackupDir() error {
	log.Printf("Creating backup directory: %s\n", c.BackupPath)
	err := os.Mkdir(c.BackupPath, 0775)
	if err != nil {
		return err
	}

	return nil
}
