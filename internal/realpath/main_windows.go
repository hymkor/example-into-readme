package realpath

import (
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

const _FILE_NAME_NORMALIZED = 0x0

func fromHandle(h windows.Handle) (string, error) {
	buf := make([]uint16, windows.MAX_PATH)
	n, err := windows.GetFinalPathNameByHandle(
		h,
		&buf[0],
		uint32(len(buf)),
		_FILE_NAME_NORMALIZED,
	)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(windows.UTF16ToString(buf[:n]), `\\?\`), nil
}

func fromFile(fd *os.File) (string, error) {
	return fromHandle(windows.Handle(fd.Fd()))
}
