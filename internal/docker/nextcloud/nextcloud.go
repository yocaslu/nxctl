package nextcloud

import (
	"fmt"
	"log/slog"
	"nxctl/internal/archive/restic"
	targz "nxctl/internal/archive/tar"
	"nxctl/internal/docker"
	"nxctl/internal/proc"
	"nxctl/internal/utils"
	"os"
	"path"
)

type Nextcloud struct {
	ContainerName   string `json:"container_name"` // nome do container docker em execucao
	VolumePath      string `json:"volume_path"` // armazenamento da instalacao do nextcloud
	DataDir         string `json:"data_dir"` // armazenamento dos usuarios
	MaintenanceMode bool `json:"maintenance_mode"`
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
	container := os.Getenv("NEXTCLOUD_CONTAINER_NAME")
	if container == "" {
		return nil, fmt.Errorf("NEXTCLOUD_CONTAINER_NAME is missing!")
	}

	volume_name := os.Getenv("NEXTCLOUD_VOLUME_NAME")
	if volume_name == "" {
		return nil, fmt.Errorf("NEXTCLOUD_VOLUME_NAME is missing!")
	}

	data := os.Getenv("NEXTCLOUD_DATA_DIR")
	if data == "" {
		return nil, fmt.Errorf("NEXTCLOUD_DATA_DIR is missing!")
	}

	return &Nextcloud{
		ContainerName:   container,
		VolumePath:      path.Join(docker.VOLUMES_PATH, volume_name),
		DataDir:         data,
		MaintenanceMode: false,
	}, nil
}

func (nx *Nextcloud) SetMaintenanceMode(flag bool) error {

	var mode string = "--off"
	if flag {
		mode = "--on"
	}

	_, err := proc.Run(os.Environ(), "docker", "exec", "-u", "www-data", nx.ContainerName, "php", "occ", "maintenance:mode", mode)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}

func (nx *Nextcloud) CreateSnapshot() error {

	restic, err := restic.Load()
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	targetPaths := []string{
		nx.DataDir,
		nx.VolumePath,
	}

	// TODO: Use Goroutines!
	for _, target := range targetPaths {
		err := restic.CreateSnapshot(target)
		if err != nil {
			return fmt.Errorf("%s", err)
		}
	}

	defer nx.SetMaintenanceMode(false)
	return nil
}

// Create an full backup of everything (volume, database and storage)
func (nx *Nextcloud) Backup(backup_path string) error {

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
		err := target.Compress()
		if err != nil {
			return fmt.Errorf("%s", err)
		}
	}

	return nil
}
