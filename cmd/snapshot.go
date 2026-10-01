package cmd

import (
	"fmt"
	"log/slog"
	"nxctl/internal/config"
	"nxctl/internal/docker/nextcloud"
	"nxctl/internal/docker/postgres"
	"nxctl/internal/utils/nxlog"
	"os"

	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Long:  "Use restic to create a snapshot of Nextcloud, volume, storage and database",
	Short: "Snapshot Nextcloud user files, volume, storage and database.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		nxlog.InitLogger(debug)
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger := nxlog.ForModule(MODULE_NAME + ".snapshot")
		fmt.Println("Loading Nextcloud environment variables")
		nx, err := nextcloud.Load()
		if err != nil {
			logger.Error("Failed to load Nextcloud environment variables", slog.Any("error", err))
			os.Exit(1)
		}

		fmt.Println("Loading nxctl environment variables")
		cfg, err := config.Load()
		if err != nil {
			logger.Error("Failed to load config environment variables", slog.Any("error", err))
			os.Exit(1)
		}

		fmt.Println("Loading database environment variables")
		db, err := postgres.Load()
		if err != nil {
			logger.Error("Failed to load database environment variables", slog.Any("error", err))
			os.Exit(1)
		}

		fmt.Println("Enabling Nextcloud maintenance mode")
		nx.SetMaintenanceMode(true)

		fmt.Println("Dumping database")
		err = db.Dump(cfg.BackupPath)
		if err != nil {
			logger.Error("Failed to dump database", slog.Any("error", err))
		}

		fmt.Println("Creating Nextcloud snapshot")
		err = nx.CreateSnapshot()
		if err != nil {
			logger.Error("Failed to snapshot Nextcloud", slog.Any("error", err))
		}

		fmt.Println("Disabling Nextcloud maintenance mode")
		defer nx.SetMaintenanceMode(false)
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug information")
}
