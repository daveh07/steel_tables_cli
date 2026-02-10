// Package ui provides terminal UI components.
package ui

// Color constants for the Tokyo Midnight theme.
const (
	// Background colors
	Reset   = "\033[0m"
	Bg      = "\033[48;2;26;27;38m"
	BgLight = "\033[48;2;36;40;59m"

	// Text colors
	Text       = "\033[38;2;192;202;245m"
	TextDim    = "\033[38;2;115;131;168m"
	TextBright = "\033[38;2;255;255;255m"

	// Accent colors
	Accent       = "\033[38;2;111;236;206m"
	AccentBright = "\033[38;2;122;162;247m"
	Blue         = "\033[38;2;125;162;206m"

	// Status colors
	Success = "\033[38;2;102;232;236m"
	Warning = "\033[38;2;255;158;100m"
	Error   = "\033[38;2;247;118;142m"

	// Border colors
	Border       = "\033[38;2;60;63;83m"
	BorderBright = "\033[38;2;122;162;247m"

	// Clear and positioning
	Clear     = "\033[2J\033[H"
	ClearLine = "\033[2K"
)
