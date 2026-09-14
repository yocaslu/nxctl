package cmd

import (
	"fmt"
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
		fmt.Println("Loading Nextcloud environment variables")
		nx, err := nextcloud.Load()
		if err != nil {
			fmt.Printf("Failed to load Nextcloud environment variables:\n%s\n", err)
			return
		}

		fmt.Println("Loading nxctl environment variables")
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("failed to load config environment variables:\n%s\n", err)
			return
		}

		fmt.Println("Loading database environment variables")
		db, err := postgres.Load()
		if err != nil {
			fmt.Printf("Failed to load PostgreSQL environment variables:\n%s\n", err)
			return
		}

		fmt.Println("Enabling Nextcloud maintenance mode")
		nx.SetMaintenanceMode(true)

		fmt.Println("Dumping database...")
		err = db.Dump(cfg.BackupPath)
		if err != nil {
			fmt.Printf("Failed to dump database:\n%s\n", err)
		}

		fmt.Println("Creating Nextcloud backup")
		err = nx.Backup(cfg.BackupPath)
		if err != nil {
			fmt.Printf("Failed to backup Nextcloud:\n%s\n", err)
			return
		}

		fmt.Println("Disabling Nextcloud maintenance mode")
		defer nx.SetMaintenanceMode(false)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
