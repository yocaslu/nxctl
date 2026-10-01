package postgres

import (
	"fmt"
	"log/slog"
	"nxctl/internal/proc"
	"nxctl/internal/utils"
	"nxctl/internal/utils/nxlog"
	"os"
	"path"
)

var MODULE_NAME string = "postgres"

type Postgres struct {
	UserName      string `json:"username"`
	Password      string `json:"password"`
	DatabaseName  string `json:"database_name"`
	ContainerName string `json:"container_name"`
}

func (p *Postgres) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("username", p.UserName),
		slog.String("password", "CENSURED"),
		slog.String("database_name", p.UserName),
		slog.String("container_name", p.ContainerName),
	)
}

type PgDump struct {
	Directory string `json:"directory"`
	Filename  string `json:"filename"`
	Extension string `json:"extension"`
	Date      string `json:"date"`
}

func (p *PgDump) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("username", p.Directory),
		slog.String("filename", p.Filename),
		slog.String("extension", p.Extension),
		slog.String("date", p.Date),
	)
}

func Load() (*Postgres, error) {
	logger := nxlog.ForModule(MODULE_NAME + ".Load")

	logger.Debug("Reading PostgreSQL environment variables")
	username := os.Getenv("POSTGRES_USER")
	if username == "" {
		return nil, fmt.Errorf("Failed to read POSTGRES_USER environment variable")
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("Failed to read POSTGRES_PASSWORD environment variable")
	}

	database := os.Getenv("POSTGRES_DB")
	if database == "" {
		return nil, fmt.Errorf("Failed to read POSTGRES_DB environment variable")
	}

	container := os.Getenv("POSTGRES_CONTAINER_NAME")
	if container == "" {
		return nil, fmt.Errorf("Failed to read POSTGRES_CONTAINER_NAME environment variable")
	}

	var pg *Postgres = &Postgres{
		UserName:      username,
		Password:      password,
		DatabaseName:  database,
		ContainerName: container,
	}

	logger.Debug(
		"PostgreSQL object created",
		slog.Any("postgres_data", pg),
	)

	return pg, nil
}

func (d *PgDump) CreateDumpDir() error {
	logger := nxlog.ForModule(MODULE_NAME + ".CreateDumpDir")

	logger.Debug(
		"Checking if database dump directory exists",
		slog.Any("pgdump", d),
	)

	exist, err := utils.DirExist(d.Directory)
	if err != nil {
		logger.Debug(
			"Failed to check if database dump directory exist",
			slog.Any("pgdump", d),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("Failed to check if database dump directory exist: [%w]", err)
	}

	if !exist {
		logger.Debug("Creating database dump directory")
		err = os.Mkdir(d.Directory, 0775)
		if err != nil {
			return fmt.Errorf("Failed to create database dump directory: [%w]", err)
		}

		logger.Debug("Database dump directory successfully created")
	}

	return nil
}

func (p *Postgres) Dump(backup_path string) error {
	logger := nxlog.ForModule(MODULE_NAME + ".Dump")

	dump := PgDump{
		Directory: path.Join(backup_path, "postgresql"),
		Filename:  "nextcloud_pgdump-",
		Extension: ".sql",
		Date:      utils.GetDate(),
	}

	logger.Debug(
		"Created PgDump struct",
		slog.Any("pgdump", dump),
	)

	logger.Debug(
		"Checking if database dump directory exists",
		slog.Any("pgdump", dump),
	)

	exist, err := utils.DirExist(dump.Directory)
	if err != nil {
		logger.Debug(
			"Failed to check if database dump directory exist",
			slog.Any("pgdump", dump),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("Failed to check if database dump directory exist: [%w]", err)
	}

	if !exist {
		err = dump.CreateDumpDir()
		if err != nil {
			logger.Debug(
				"Failed to create database dump directory",
				slog.Any("pgdump", dump),
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("Failed to create dump directory: [%w]", err)
		}
	}

	args := []string{
		"exec", p.ContainerName,
		"pg_dump", p.DatabaseName,
		"-U", p.UserName,
		"-f", path.Join(dump.Directory, dump.Filename+dump.Date+dump.Extension),
	}

	logger.Debug("Dumping database")
	_, err = proc.Run(os.Environ(), "docker", args...)
	if err != nil {
		return fmt.Errorf("Failed to dump database: [%w]", err)
	}

	return nil
}
