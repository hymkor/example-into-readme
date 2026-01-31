//go:build !windows

package realpath

import (
	"os"
)

func fromFile(fd *os.File) (string, error) {
	return fd.Name(), nil
}
