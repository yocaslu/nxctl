package restic

import (
	"fmt"
	"log/slog"
	"nxctl/internal/proc"
	"nxctl/internal/utils"
	"os"
)

type Restic struct {
	Repo     string
	Password string
	_env     []string
}

func New(repo string, password string) *Restic {
	// TODO: check if repo exists

	return &Restic{
		Repo:     repo,
		Password: password,
		_env:     os.Environ(),
	}
}

func Load() (*Restic, error) {
	repo := os.Getenv("RESTIC_REPOSITORY")
	if repo == "" {
		return nil, fmt.Errorf("RESTIC_REPOSITORY is missing!")
	}

	password := os.Getenv("RESTIC_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("RESTIC_PASSWORD is missing!")
	}

	return &Restic{
		Repo:     repo,
		Password: password,
		_env:     os.Environ(),
	}, nil
}

func (r *Restic) CreateRepo() error {
	exist, _ := utils.DirExist(r.Repo)
	if exist {
		return nil
	}

	stdout, err := proc.Run(os.Environ(), "restic", "-r", r.Repo, "init")
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}

	slog.Debug(stdout)
	return nil
}

func (r *Restic) CreateSnapshot(target string) error {
	exist, _ := utils.DirExist(r.Repo)
	if !exist {
		err := r.CreateRepo()
		if err != nil {
			return fmt.Errorf("%s\n", err)
		}
	}

	_, err := proc.Run(r._env, "restic", "-r", r.Repo, "backup", target)
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}

	return nil
}
