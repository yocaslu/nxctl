package cmd

import (
	"log"
	"nxctl/internal/config"
	"nxctl/internal/docker/nextcloud"
	"nxctl/internal/docker/postgres"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup Nextcloud volume, storage and database",
	Long:  "Backup user files and Nextcloud Docker volume and dump database",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: create full backup of nextcloud using tar, compress using xz?
		nx, err := nextcloud.Load()
		if err != nil {
			log.Printf("failed to load Nextcloud informations: %s\n", err)
			return
		}

		cfg, err := config.Load()
		if err != nil {
			log.Printf("failed to load config: %s\n", err)
			return
		}

		db, err := postgres.Load()
		if err != nil {
			log.Printf("failed to load database informations: %s\n", err)
			return
		}

		nx.SetMaintenanceMode(true)

		err = db.Dump(cfg.BackupPath)
		if err != nil {
			log.Printf("failed to dump database: %s\n", err)
		}

		err = nx.Backup(cfg.BackupPath)
		if err != nil {
			log.Printf("failed to backup Nextcloud: %s\n", err)
			return
		}

		defer nx.SetMaintenanceMode(false)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
