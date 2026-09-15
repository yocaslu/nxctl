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

	logger.Debug(
		"Created Tar object",
		slog.Any("Tar", t),
	)

	return t
}

func (t *Tar) Compress() error {
	logger := nxlog.ForModule(MODULE_NAME + ".Compress")
	logger.Debug(
		"Starting Tar compress",
		slog.Any("Tar", t),
	)

	_, stderr := proc.Run(os.Environ(), "tar", "-czpf", t.Dest, t.Source)
	if stderr != nil {
		return fmt.Errorf("Failed to compress: %s", stderr)
	}

	return nil
}
