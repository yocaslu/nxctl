package cmd

import (
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
		nx, err := nextcloud.Load()
		if err != nil {
			log.Printf("failed to load Nextcloud informations:\n%s\n", err)
			return
		}

		cfg, err := config.Load()
		if err != nil {
			log.Printf("failed to load config:\n%s\n", err)
			return
		}

		db, err := postgres.Load()
		if err != nil {
			log.Printf("failed to load database informations:\n%s\n", err)
			return
		}

		nx.SetMaintenanceMode(true)

		err = db.Dump(cfg.BackupPath)
		if err != nil {
			log.Printf("failed to dump database:\n%s\n", err)
		}

		err = nx.CreateSnapshot()
		if err != nil {
			log.Printf("failed to snapshot Nextcloud: \n%s\n", err)
			return
		}

		defer nx.SetMaintenanceMode(false)
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
}
