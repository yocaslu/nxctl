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

var BackupDirName string = "backup_"
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup Nextcloud volume, storage and database",
	Long:  "Backup user files and Nextcloud Docker volume and dump database",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		nxlog.InitLogger(debug)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return Backup()	
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug information")
}

func Backup() error {
	// TODO: create full backup of nextcloud using tar, compress using xz?
	fmt.Println("Loading Nextcloud environment variables")
	nx, err := nextcloud.Load()
	if err != nil {
		return fmt.Errorf("failed to load Nextcloud environment variables: [%w]", err)
	}

	fmt.Println("Loading nxctl environment variables")
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config environment variables: [%w]",err)
	}

	fmt.Println("Loading database environment variables")
	db, err := postgres.Load()
	if err != nil {
		return fmt.Errorf("failed to load PostgreSQL environment variables: [%w]",err)
	}

	fmt.Println("Enabling Nextcloud maintenance mode")
	err = nx.SetMaintenanceMode(true)
	if err != nil {
		return fmt.Errorf("failed to enable nextcloud maintenance mode: [%w]", err)
	}

	// Ensure maintenance mode is disabled
	defer func() {
		fmt.Println("Disabling Nextcloud maintenance mode")
		if merr := nx.SetMaintenanceMode(false); merr != nil {
			err = errors.Join(err, fmt.Errorf("failed to disable Nextcloud maintenance mode: [%w]", merr)) 	
		}
	}()

	fmt.Println("Dumping database...")
	if err = db.Dump(cfg.BackupPath); err != nil {
		return fmt.Errorf("failed to dump database: [%w]", err)
	}

	fmt.Println("Creating Nextcloud backup")
	if err = nx.CreateBackup(cfg.BackupPath); err != nil {
		return fmt.Errorf("failed to backup Nextcloud: [%w]", err)
	}

	return nil
}
