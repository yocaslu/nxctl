package postgres

import (
	"fmt"
	"log"
	"nxctl/internal/proc"
	"nxctl/utils"
	"os"
	"path"
)

type Postgres struct {
	UserName      string
	Password      string
	DatabaseName  string
	ContainerName string
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

func (p *Postgres) Dump(backup_path string) error {
	backup_path = path.Join(backup_path, "postgresql")                             // adiciona o diretorio especifico de dump do pgsql
	filepath := path.Join(backup_path, "nextcloud_pgdump_"+utils.GetDate()+".sql") // constroi o nome do arquivo baseado na data
	_, stderr := proc.Run(os.Environ(), "docker", "exec", p.ContainerName, "pg_dump", p.DatabaseName, "-U", p.UserName, "-f", filepath)
	if stderr != nil {
		log.Printf("Failed to dump PostgreSQL %s database due to: %s", p.DatabaseName, stderr)
		return stderr
	}

	return nil
}
