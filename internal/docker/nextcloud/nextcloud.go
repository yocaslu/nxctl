package nextcloud

import (
	"fmt"
	"log/slog"
	"nxctl/internal/archive/restic"
	targz "nxctl/internal/archive/tar"
	"nxctl/internal/docker"
	"nxctl/internal/nxlog"
	"nxctl/internal/proc"
	"nxctl/internal/utils"
	"os"
	"path"
)

var MODULE_NAME string = "nextcloud"

type Nextcloud struct {
	ContainerName   string `json:"container_name"` // nome do container docker em execucao
	VolumePath      string `json:"volume_path"`    // armazenamento da instalacao do nextcloud
	DataDir         string `json:"data_dir"`       // armazenamento dos usuarios
	MaintenanceMode bool   `json:"maintenance_mode"`
}

func (nx *Nextcloud) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("container_name", nx.ContainerName),
		slog.String("volume_path", nx.VolumePath),
		slog.String("data_dir", nx.DataDir),
		slog.Bool("maintenance_mode", nx.MaintenanceMode),
	)
}

func Load() (*Nextcloud, error) {
	logger := nxlog.ForModule(MODULE_NAME + ".Load")

	logger.Info("Reading Nextcloud environment variables")
	container := os.Getenv("NEXTCLOUD_CONTAINER_NAME")
	if container == "" {
		logger.Error(
			"Failed to read NEXTCLOUD_CONTAINER_NAME environment variable",
			slog.Any("env", os.Environ()),
		)
		return nil, fmt.Errorf("NEXTCLOUD_CONTAINER_NAME is missing!")
	}

	volume_name := os.Getenv("NEXTCLOUD_VOLUME_NAME")
	if volume_name == "" {
		logger.Error(
			"Failed to read NEXTCLOUD_VOLUME_NAME environment variable",
			slog.Any("env", os.Environ()),
		)
		return nil, fmt.Errorf("NEXTCLOUD_VOLUME_NAME is missing!")
	}

	data := os.Getenv("NEXTCLOUD_DATA_DIR")
	if data == "" {
		logger.Error(
			"Failed to read NEXTCLOUD_DATA_DIR environment variable",
			slog.Any("env", os.Environ()),
		)
		return nil, fmt.Errorf("NEXTCLOUD_DATA_DIR is missing!")
	}

	var nx *Nextcloud = &Nextcloud{
		ContainerName:   container,
		VolumePath:      path.Join(docker.VOLUMES_PATH, volume_name),
		DataDir:         data,
		MaintenanceMode: false,
	}

	logger.Info(
		"Nextcloud environment variables were successfully read",
		slog.Any("nexcloud_data", nx),
	)

	return nx, nil
}

func (nx *Nextcloud) SetMaintenanceMode(flag bool) error {
	logger := nxlog.ForModule(MODULE_NAME + ".SetMaintenanceMode")

	var mode string = "--off"
	if flag {
		mode = "--on"
	}

	logger.Info(
		"Switching Nextcloud maintenance mode",
		slog.Bool("maintenance_mode", flag),
	)

	_, err := proc.Run(os.Environ(), "docker", "exec", "-u", "www-data", nx.ContainerName, "php", "occ", "maintenance:mode", mode)
	if err != nil {
		logger.Error(
			"Failed to switch maintenance mode",
			slog.Bool("maintenance_mode", flag),
			slog.Any("nextcloud_data", nx),
		)

		return fmt.Errorf("%s", err)
	}

	return nil
}

func (nx *Nextcloud) CreateSnapshot() error {
	logger := nxlog.ForModule(MODULE_NAME + ".CreateSnapshot")

	logger.Info(
		"Reading Restic environment variables",
		slog.Any("nextcloud_data", nx),
	)

	restic, err := restic.Load()
	if err != nil {
		logger.Error(
			"Failed to read Restic environment variables",
			slog.Any("nextcloud_data", nx),
			slog.Any("env", os.Environ()),
		)

		return fmt.Errorf("%s", err)
	}

	targetPaths := []string{
		nx.DataDir,
		nx.VolumePath,
	}

	// TODO: Use Goroutines!
	for _, target := range targetPaths {
		logger.Info(
			"Creating snapshot",
			slog.String("target", target),
			slog.Any("nextcloud_data", nx),
		)

		err := restic.CreateSnapshot(target)
		if err != nil {
			logger.Error(
				"Failed to create snapshot",
				slog.String("target", target),
				slog.Any("nextcloud_data", nx),
				slog.Any("env", os.Environ()),
			)
			return fmt.Errorf("%s", err)
		}
	}

	logger.Info("Disabling maintenance mode")
	defer nx.SetMaintenanceMode(false)
	return nil
}

// Create an full backup of everything (volume, database and storage)
func (nx *Nextcloud) Backup(backup_path string) error {
	logger := nxlog.ForModule(MODULE_NAME + ".Backup")

	type Target struct {
		Source string
		Dest   string
	}
	date := utils.GetDate()
	targets := []targz.Tar{
		{
			Source: nx.VolumePath,
			Dest:   path.Join(backup_path, "nextcloud_volume"+"-"+date),
		},
		{
			Source: nx.DataDir,
			Dest:   path.Join(backup_path, "nextcloud_storage"+"-"+date),
		},
	}

	// TODO: Use Goroutines!
	for _, target := range targets {
		logger.Info(
			"Creating snapshot",
			slog.Any("target", target),
			slog.Any("nextcloud_data", nx),
		)

		err := target.Compress()
		if err != nil {
			logger.Error(
				"Failed to create snapshot",
				slog.Any("target", target),
				slog.Any("nextcloud_data", nx),
				slog.Any("env", os.Environ()),
			)
			return fmt.Errorf("%s", err)
		}
	}

	return nil
}
