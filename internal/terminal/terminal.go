package terminal

import (
	"errors"
	"strconv"
	"syscall"
	"unsafe"
)

type size struct {
	Row  uint16
	Col  uint16
	Xpix uint16
	Ypix uint16
}

func GetSize(fd uintptr) (width, height int, err error) {
	s := &size{}
	retCode, _, errNo := syscall.Syscall(syscall.SYS_IOCTL,
		fd,
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(s)))

	if int(retCode) == -1 {
		return 0, 0, errors.New(strconv.Itoa(int(errNo)))
	}

	return int(s.Col), int(s.Row), nil
}
