//go:build windows
// +build windows

package main

import (
	"bufio"
	"os"
)

func readKey() byte {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadByte()
	return input
}

func readKeyNonBlocking() (byte, bool) {
	reader := bufio.NewReader(os.Stdin)
	if reader.Buffered() > 0 {
		input, _ := reader.ReadByte()
		return input, true
	}
	return 0, false
}

type termios struct {
	// Dummy struct for Windows compatibility
}

func setRawMode() *termios {
	// Windows doesn't support raw mode in the same way
	return nil
}

func restoreTerminal(oldState *termios) {
	// No action needed on Windows
}
