//go:build windows
// +build windows

package ui

// Termios is a dummy struct for Windows compatibility.
type Termios struct{}

// GetTerminalState returns nil on Windows.
func GetTerminalState() (*Termios, error) {
	return &Termios{}, nil
}

// SetRawMode is a no-op on Windows.
func SetRawMode() error {
	return nil
}

// RestoreTerminal is a no-op on Windows.
func RestoreTerminal(oldState *Termios) {}

// GetTerminalWidth returns a default width on Windows.
func GetTerminalWidth() int {
	return 120
}

// GetTerminalHeight returns a default height on Windows.
func GetTerminalHeight() int {
	return 40
}

// GetMaxCols computes column count based on default width.
func GetMaxCols() int {
	return 5
}
