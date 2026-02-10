//go:build !windows
// +build !windows

package ui

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// Termios holds terminal I/O settings for Unix systems.
type Termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]uint8
	Ispeed uint32
	Ospeed uint32
}

// GetTerminalState retrieves the current terminal attributes.
func GetTerminalState() (*Termios, error) {
	var state Termios
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), uintptr(0x5401), uintptr(unsafe.Pointer(&state))); errno != 0 {
		return nil, errno
	}
	return &state, nil
}

// SetRawMode puts the terminal into raw mode for single-key input.
func SetRawMode() error {
	var oldState Termios
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), uintptr(0x5401), uintptr(unsafe.Pointer(&oldState))); errno != 0 {
		return errno
	}
	newState := oldState
	newState.Lflag &^= 0x0000000A // Disable ECHO and ICANON
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), uintptr(0x5402), uintptr(unsafe.Pointer(&newState))); errno != 0 {
		return errno
	}
	return nil
}

// RestoreTerminal restores the terminal to its original state.
func RestoreTerminal(oldState *Termios) {
	if oldState == nil {
		return
	}
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), uintptr(0x5402), uintptr(unsafe.Pointer(oldState)))
}

// GetTerminalWidth returns the current terminal width in columns.
func GetTerminalWidth() int {
	cmd := exec.Command("tput", "cols")
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil {
		return 120
	}
	width, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil || width < 80 {
		return 120
	}
	if width > 200 {
		width = 200
	}
	return width
}

// GetTerminalHeight returns the current terminal height in rows.
func GetTerminalHeight() int {
	cmd := exec.Command("tput", "lines")
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil {
		return 40
	}
	height, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil || height < 10 {
		return 40
	}
	return height
}

// GetMaxCols computes how many data columns fit in the terminal width.
func GetMaxCols() int {
	termWidth := GetTerminalWidth()
	sectionColWidth := 25
	colWidth := 18
	maxCols := (termWidth - sectionColWidth) / colWidth
	if maxCols < 1 {
		maxCols = 1
	}
	return maxCols
}
