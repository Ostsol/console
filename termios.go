// Copyright 2013 Daniel Jo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

type termios syscall.Termios

var defaultTermios *termios

func init() {
	var err error
	defaultTermios, err = getTermios()
	if err != nil {
		panic(err)
	}
}

func getTermios() (*termios, error) {
	term := new(termios)
	if err := term.get(); err != nil {
		return nil, err
	}

	return term, nil
}

func (t *termios) get() error {
	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		0, syscall.TCGETS,
		uintptr(unsafe.Pointer(t)))

	if errno != 0 {
		return os.NewSyscallError("SYS_IOCTL", errno)
	}

	if r1 != 0 {
		return errors.New("Termios.Get(): unhandled error")
	}

	return nil
}

func (t *termios) set() error {
	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		0, syscall.TCSETS,
		uintptr(unsafe.Pointer(t)))
	if errno != 0 {
		return os.NewSyscallError("SYS_IOCTL", errno)
	}

	if r1 != 0 {
		return errors.New("Termios.Get(): unhandled error")
	}

	return nil
}

func (t *termios) rawMode() {
	t.Iflag &= ^uint32(syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON)
	t.Oflag &= ^uint32(syscall.OPOST)
	t.Lflag &= ^uint32(syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.IEXTEN)
	t.Cflag &= ^uint32(syscall.CSIZE | syscall.PARENB)
	t.Cflag |= syscall.CS8

	t.Cc[syscall.VMIN] = 1
	t.Cc[syscall.VTIME] = 0
}

