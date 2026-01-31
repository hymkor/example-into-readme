package realpath

import (
	"os"
)

func FromFile(fd *os.File) (string, error) {
	return fromFile(fd)
}
