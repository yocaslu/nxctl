package restic

import (
	"fmt"
	"log/slog"
	"nxctl/internal/proc"
	"nxctl/internal/utils"
	"nxctl/internal/utils/nxlog"
	"os"
)

type Restic struct {
	Repo     string `json:"repo"`
	Password string
	_env     []string
}

var MODULE_NAME string = "Restic"

func (r *Restic) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("repo", r.Repo),
	)
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
	logger := nxlog.ForModule(MODULE_NAME + ".Load")
	logger.Debug("Creating Restic object")

	repo := os.Getenv("RESTIC_REPOSITORY")
	if repo == "" {
		return nil, fmt.Errorf("RESTIC_REPOSITORY is missing!")
	}

	password := os.Getenv("RESTIC_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("RESTIC_PASSWORD is missing!")
	}

	r := &Restic{
		Repo:     repo,
		Password: password,
		_env:     os.Environ(),
	}

	logger.Debug("Restic object created", slog.Any("restic", r))
	return r, nil
}

func (r *Restic) CreateRepo() error {
	logger := nxlog.ForModule(MODULE_NAME + ".Load")

	exist, _ := utils.DirExist(r.Repo)
	if exist {
		return nil
	}

	logger.Debug("Creating repository", slog.Any("restic", r))
	stdout, err := proc.Run(os.Environ(), "restic", "-r", r.Repo, "init")

	if err != nil {
		return fmt.Errorf("Failed to create repository: [%w]", err)
	}

	slog.Debug(stdout)
	return nil
}

func (r *Restic) CreateSnapshot(target string) error {
	logger := nxlog.ForModule(MODULE_NAME + ".CreateSnapshot")
	exist, _ := utils.DirExist(r.Repo)

	if !exist {
		logger.Debug(
			"Creating Restic repository",
			slog.Any("restic", r),
		)

		err := r.CreateRepo()
		if err != nil {
			return fmt.Errorf("Failed to create repository: [%w]", err)
		}
	}

	logger.Debug("Creating snapshot")
	_, err := proc.Run(r._env, "restic", "-r", r.Repo, "backup", target)
	if err != nil {
		return fmt.Errorf("Failed to create snapshot: [%w]", err)
	}

	return nil
}
