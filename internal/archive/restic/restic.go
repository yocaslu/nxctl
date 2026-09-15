package restic

import (
	"fmt"
	"log/slog"
	"nxctl/internal/nxlog"
	"nxctl/internal/proc"
	"nxctl/internal/utils"
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
	logger.Info("Creating Restic object")

	repo := os.Getenv("RESTIC_REPOSITORY")
	if repo == "" {
		logger.Error(
			"Failed to find RESTIC_REPOSITORY environment variable",
			slog.Any("env", os.Environ()),
		)

		return nil, fmt.Errorf("RESTIC_REPOSITORY is missing!")
	}

	password := os.Getenv("RESTIC_PASSWORD")
	if password == "" {
		logger.Error(
			"Failed to find RESTIC_PASSWORD environment variable",
			slog.Any("env", os.Environ()),
		)

		return nil, fmt.Errorf("RESTIC_PASSWORD is missing!")
	}

	r := &Restic{
		Repo:     repo,
		Password: password,
		_env:     os.Environ(),
	}

	logger.Info("Restic object created", slog.Any("restic", r))
	return r, nil
}

func (r *Restic) CreateRepo() error {
	logger := nxlog.ForModule(MODULE_NAME + ".Load")

	exist, _ := utils.DirExist(r.Repo)
	if exist {
		return nil
	}

	logger.Info("Creating repository", slog.Any("restic", r))
	stdout, err := proc.Run(os.Environ(), "restic", "-r", r.Repo, "init")

	if err != nil {
		logger.Error(
			"Failed to create repository",
			slog.Any("restic", r),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s", err)
	}

	slog.Debug(stdout)
	return nil
}

func (r *Restic) CreateSnapshot(target string) error {
	logger := nxlog.ForModule(MODULE_NAME + ".CreateSnapshot")
	exist, _ := utils.DirExist(r.Repo)

	if !exist {
		logger.Warn("Repository directory not found, creating it.")
		err := r.CreateRepo()

		if err != nil {
			logger.Error(
				"Failed to create repository",
				slog.Any("restic", r),
			)

			return fmt.Errorf("%s", err)
		}
	}

	logger.Info("Creating snapshot")
	_, err := proc.Run(r._env, "restic", "-r", r.Repo, "backup", target)
	if err != nil {
		logger.Error(
			"Failed to create snapshot",
			slog.Any("restic", r),
		)

		return fmt.Errorf("%s", err)
	}

	return nil
}
