//go:build !windows
// +build !windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

func readKey() byte {
	var buf [1]byte
	_, err := os.Stdin.Read(buf[:])
	if err != nil {
		return 0
	}
	return buf[0]
}

func readKeyNonBlocking() (byte, bool) {
	var buf [1]byte
	n, err := os.Stdin.Read(buf[:])
	if err != nil || n == 0 {
		return 0, false
	}
	return buf[0], true
}

type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]uint8
	Ispeed uint32
	Ospeed uint32
}

func setRawMode() *termios {
	var oldState termios
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), uintptr(0x5401), uintptr(unsafe.Pointer(&oldState)))

	newState := oldState
	newState.Lflag &^= 0x0000000A // Disable ECHO and ICANON

	syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), uintptr(0x5402), uintptr(unsafe.Pointer(&newState)))
	return &oldState
}

func restoreTerminal(oldState *termios) {
	if oldState == nil {
		return
	}
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), uintptr(0x5402), uintptr(unsafe.Pointer(oldState)))
}
