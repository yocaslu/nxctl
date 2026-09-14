package cmd

import (
	"fmt"
	"log"
	"nxctl/internal/config"
	"nxctl/internal/docker/nextcloud"
	"nxctl/internal/docker/postgres"

	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Long:  "Use restic to create a snapshot of Nextcloud, volume, storage and database",
	Short: "Snapshot Nextcloud user files, volume, storage and database.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Loading Nextcloud environment variables")
		nx, err := nextcloud.Load()
		if err != nil {
			log.Printf("Failed to load Nextcloud environment variables:\n%s\n", err)
			return
		}

		fmt.Println("Loading nxctl environment variables")
		cfg, err := config.Load()
		if err != nil {
			log.Printf("Failed to load config environment variables:\n%s\n", err)
			return
		}

		fmt.Println("Loading database environment variables")
		db, err := postgres.Load()
		if err != nil {
			log.Printf("Failed to load database environment variables:\n%s\n", err)
			return
		}

		fmt.Println("Enabling Nextcloud maintenance mode")
		nx.SetMaintenanceMode(true)

		fmt.Println("Dumping database")
		err = db.Dump(cfg.BackupPath)
		if err != nil {
			log.Printf("Failed to dump database:\n%s\n", err)
		}

		fmt.Println("Creating Nextcloud snapshot")
		err = nx.CreateSnapshot()
		if err != nil {
			log.Printf("Failed to snapshot Nextcloud:\n%s\n", err)
			return
		}

		fmt.Println("Disabling Nextcloud maintenance mode")
		defer nx.SetMaintenanceMode(false)
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
}
