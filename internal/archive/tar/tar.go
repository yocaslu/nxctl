package targz

import (
	"fmt"
	"nxctl/internal/proc"
	"os"
)

type Tar struct {
	Dest   string
	Source string
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
		return fmt.Errorf("%s\n", stderr)
	}

	return nil
}
