package cmd

import (
	"errors"
	"fmt"
	"nxctl/internal/component/nextcloud"
	"nxctl/internal/component/postgres"
	"nxctl/internal/config"
	"nxctl/internal/utils/nxlog"

	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Long:  "Use restic to create a snapshot of Nextcloud, volume, storage and database",
	Short: "Snapshot Nextcloud user files, configuration and dump the database.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		nxlog.InitLogger(debug)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return Snapshot()	
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug information")
}

func Snapshot() error {
	fmt.Println("Loading Nextcloud environment variables")
	nx, err := nextcloud.Load()
	if err != nil {
		fmt.Errorf("Failed to load Nextcloud environment variables: [%w]",err)
	}

	fmt.Println("Loading nxctl environment variables")
	cfg, err := config.Load()
	if err != nil {
		fmt.Errorf("Failed to load config environment variables: [%w]", err)
	}

	fmt.Println("Loading database environment variables")
	db, err := postgres.Load()
	if err != nil {
		fmt.Errorf("Failed to load database environment variables: [%w]", err)
	}

	fmt.Println("Enabling Nextcloud maintenance mode")
	if err = nx.SetMaintenanceMode(true); err != nil {
		return fmt.Errorf("failed to enable Nextcloud maintenance mode: [%w]", err)
	}

	defer func() {
		fmt.Println("Disabling Nextcloud maintenance mode")
		if merr := nx.SetMaintenanceMode(false); merr != nil {
			err = errors.Join(err, fmt.Errorf("failed to disable Nextcloud maintenance mode: [%w]", merr)) 	
		}
	}()

	fmt.Println("Dumping database")
	if err = db.Dump(cfg.BackupPath); err != nil {
		fmt.Errorf("Failed to dump database: [%w]", err)
	}

	fmt.Println("Creating Nextcloud snapshot")
	if err = nx.CreateSnapshot(); err != nil {
		fmt.Errorf("Failed to snapshot Nextcloud: [%w]", err)
	}

	return nil
}