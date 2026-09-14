package postgres

import (
	"fmt"
	"nxctl/internal/proc"
	"nxctl/internal/utils"
	"os"
	"path"
)

type Postgres struct {
	UserName      string
	Password      string
	DatabaseName  string
	ContainerName string
}

type PgDump struct {
	Directory string
	Filename  string
	Extension string
	Date      string
}

func Load() (*Postgres, error) {
	username := os.Getenv("POSTGRES_USER")
	if username == "" {
		return nil, fmt.Errorf("POSTGRES_USER is missing!")
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD is missing!")
	}

	database := os.Getenv("POSTGRES_DB")
	if database == "" {
		return nil, fmt.Errorf("POSTGRES_DB is missing!")
	}

	container := os.Getenv("POSTGRES_CONTAINER_NAME")
	if container == "" {
		return nil, fmt.Errorf("POSTGRES_CONTAINER_NAME is missing!")
	}

	return &Postgres{
		UserName:      username,
		Password:      password,
		DatabaseName:  database,
		ContainerName: container,
	}, nil
}

func (d *PgDump) CreateDumpDir() error {
	exist, err := utils.DirExist(d.Directory)
	if err != nil {
		return err
	}

	if !exist {
		err = os.Mkdir(d.Directory, 0775)
		if err != nil {
			return fmt.Errorf("%s\n", err)
		}
	}

	return nil
}

func (p *Postgres) Dump(backup_path string) error {
	dump := PgDump{
		Directory: path.Join(backup_path, "postgresql"),
		Filename:  "nextcloud_pgdump-",
		Extension: ".sql",
		Date:      utils.GetDate(),
	}

	exist, err := utils.DirExist(dump.Directory)
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}

	if !exist {
		err = dump.CreateDumpDir()
		if err != nil {
			return fmt.Errorf("%s\n", err)
		}
	}

	args := []string{
		"exec", p.ContainerName,
		"pg_dump", p.DatabaseName,
		"-U", p.UserName,
		"-f", path.Join(dump.Directory, dump.Filename+dump.Date+dump.Extension),
	}

	_, err = proc.Run(os.Environ(), "docker", args...)
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}

	return nil
}
