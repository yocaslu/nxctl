package restic

import (
	"fmt"
	"log"
	"nxctl/internal/proc"
	"os"
)

type Restic struct {
	Repo     string
	Password string
	_env     []string
}

func New(repo string, password string) (*Restic, error) {
	// TODO: check if repo exists

	return &Restic{
		Repo:     repo,
		Password: password,
		_env:     os.Environ(),
	}, nil
}

func Load() (*Restic, error) {
	repo := os.Getenv("RESTIC_REPOSITORY")
	if repo != "" {
		return nil, fmt.Errorf("RESTIC_REPOSITORY is missing!")
	}

	password := os.Getenv("RESTIC_PASSWORD")
	if password != "" {
		return nil, fmt.Errorf("RESTIC_PASSWORD is missing!")
	}

	return &Restic{
		Repo:     repo,
		Password: password,
		_env:     os.Environ(),
	}, nil
}

func (r *Restic) CreateSnapshot(target string) error {
	stdout, err := proc.Run(r._env, "restic", "-r", r.Repo, "backup", target)
	if err != nil {
		log.Printf("failed to backup %s", target)
		return err
	}

	log.Println(stdout)
	return nil
}
