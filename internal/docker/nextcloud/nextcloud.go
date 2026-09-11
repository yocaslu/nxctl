package nextcloud

import (
	"fmt"
	"log"
	"nxctl/internal/archive/restic"
	targz "nxctl/internal/archive/tar"
	"nxctl/internal/docker"
	"nxctl/internal/proc"
	"nxctl/internal/utils"
	"os"
	"path"
)

type Nextcloud struct {
	ContainerName   string // nome do container docker em execucao
	VolumePath      string // armazenamento da instalacao do nextcloud
	DataDir         string // armazenamento dos usuarios
	MaintenanceMode bool
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

	var mode string
	if flag {
		mode = "--on"
	} else {
		mode = "--off"
	}

	_, err := proc.Run(os.Environ(), "docker", "exec", "-u", "www-data", nx.ContainerName, "php", "occ", "maintenance:mode", mode)
	if err != nil {
		return err
	}

	return nil
}

func (nx *Nextcloud) CreateSnapshot() error {

	restic, err := restic.Load()
	if err != nil {
		log.Printf("failed to create Restic instance:\n%s\n", err)
		return err
	}

	targetPaths := []string{
		nx.DataDir,
		nx.VolumePath,
	}

	// TODO: Use Goroutines!
	for index, target := range targetPaths {
		log.Printf("Restic: [%d of %d] Backing up %s.\n", index+1, len(targetPaths), target)
		err := restic.CreateSnapshot(target)
		if err != nil {
			log.Printf("Failed to backup %s:\n%s\n", target, err)
			return err
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
	for index, target := range targets {
		log.Printf("Tar: [%d of %d] Backing up %s.\n", index+1, len(targets), target)
		err := target.Compress()
		if err != nil {
			log.Printf("Failed to backup %s:\n%s\n", target, err)
			return err
		}
	}

	return nil
}
