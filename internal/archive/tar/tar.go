package targz

import (
	"fmt"
	"log/slog"
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

func New(dest string, source string) *Tar {
	// TODO: verify if this paths exists
	return &Tar{
		Dest:   dest,
		Source: source,
	}
}

func (t *Tar) Compress() error {
	_, stderr := proc.Run(os.Environ(), "tar", "-czpf", t.Dest, t.Source)
	if stderr != nil {
		return fmt.Errorf("%s", stderr)
	}

	return nil
}
