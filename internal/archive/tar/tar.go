package targz

import (
	"fmt"
	"log/slog"
	"nxctl/internal/nxlog"
	"nxctl/internal/proc"
	"os"
)

type Tar struct {
	Dest   string `json:"dest"`
	Source string `json:"source"`
}

func (t *Tar) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("dest", t.Dest),
		slog.String("source", t.Source),
	)
}

var MODULE_NAME string = "targz"

func New(dest string, source string) *Tar {
	// TODO: verify if this paths exists
	logger := nxlog.ForModule(MODULE_NAME + ".New")
	t := &Tar{
		Dest:   dest,
		Source: source,
	}

	logger.Info(
		"Created Tar object",
		slog.Any("Tar", t),
	)

	return t
}

func (t *Tar) Compress() error {
	logger := nxlog.ForModule(MODULE_NAME + ".Compress")
	logger.Info(
		"Starting Tar compress",
		slog.Any("Tar", t),
	)

	_, stderr := proc.Run(os.Environ(), "tar", "-czpf", t.Dest, t.Source)
	if stderr != nil {
		logger.Error(
			"Failed to compress",
			slog.Any("Tar", t),
		)

		return fmt.Errorf("%s", stderr)
	}

	return nil
}
