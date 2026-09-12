package targz

import (
	"fmt"
	"log"
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
	stdout, stderr := proc.Run(os.Environ(), "tar", "-czpf", t.Dest, t.Source)
	if stderr != nil {
		return fmt.Errorf("Failed to compact due to: %s\n", stderr)
	} else {
		log.Println(stdout)
	}

	return nil
}
