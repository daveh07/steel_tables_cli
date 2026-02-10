package ui

import (
	"fmt"
	"strings"
)

// DrawHeader draws the title box with table name and page info.
func DrawHeader(filename string, currentPage, totalPages, totalEntries int) {
	termWidth := GetTerminalWidth()
	titleText := fmt.Sprintf("STEEL PROPERTIES: %s", strings.ToUpper(strings.TrimSuffix(filename, ".json")))
	infoText := fmt.Sprintf("Page %d/%d | %d entries", currentPage, totalPages, totalEntries)
	infoText = strings.TrimSpace(infoText)

	boxWidth := len(titleText)
	if len(infoText) > boxWidth {
		boxWidth = len(infoText)
	}
	boxWidth += 6
	if boxWidth < 60 {
		boxWidth = 60
	}
	if boxWidth > termWidth-4 {
		boxWidth = termWidth - 4
	}

	centerOffset := (termWidth - boxWidth - 2) / 2
	if centerOffset < 0 {
		centerOffset = 0
	}
	remainingSpace := termWidth - centerOffset - boxWidth - 2
	if remainingSpace < 0 {
		remainingSpace = 0
	}

	// Top border
	fmt.Printf("%s%s", Bg, strings.Repeat(" ", centerOffset))
	fmt.Printf("%s╔%s╗", BorderBright, strings.Repeat("═", boxWidth))
	fmt.Printf("%s%s%s\n", Bg, strings.Repeat(" ", remainingSpace), Reset)

	// Title row
	titlePadding := (boxWidth - len(titleText)) / 2
	titleRightPadding := boxWidth - len(titleText) - titlePadding
	fmt.Printf("%s%s", Bg, strings.Repeat(" ", centerOffset))
	fmt.Printf("%s║%s%s%s%s%s%s%s║",
		BorderBright, Bg, strings.Repeat(" ", titlePadding), Accent, titleText, Bg, strings.Repeat(" ", titleRightPadding), BorderBright)
	fmt.Printf("%s%s%s\n", Bg, strings.Repeat(" ", remainingSpace), Reset)

	// Info row
	infoPadding := (boxWidth - len(infoText)) / 2
	infoRightPadding := boxWidth - len(infoText) - infoPadding
	fmt.Printf("%s%s", Bg, strings.Repeat(" ", centerOffset))
	fmt.Printf("%s║%s%s%s%s%s%s%s║",
		BorderBright, Bg, strings.Repeat(" ", infoPadding), TextDim, infoText, Bg, strings.Repeat(" ", infoRightPadding), BorderBright)
	fmt.Printf("%s%s%s\n", Bg, strings.Repeat(" ", remainingSpace), Reset)

	// Bottom border
	fmt.Printf("%s%s", Bg, strings.Repeat(" ", centerOffset))
	fmt.Printf("%s╚%s╝", BorderBright, strings.Repeat("═", boxWidth))
	fmt.Printf("%s%s%s\n\n", Bg, strings.Repeat(" ", remainingSpace), Reset)
}

// DrawNavigationFooter draws the row info and keyboard shortcuts.
func DrawNavigationFooter(currentPage, totalPages, startRow, endRow, totalRows int) {
	termWidth := GetTerminalWidth()

	// Row info line
	rowInfo := fmt.Sprintf("Rows %d–%d of %d", startRow+1, endRow, totalRows)
	rowInfoColored := fmt.Sprintf("%sRows %s%d–%d%s of %s%d%s",
		TextDim, Accent, startRow+1, endRow, TextDim, Accent, totalRows, TextDim)
	rowPadding := (termWidth - len(rowInfo)) / 2
	if rowPadding < 0 {
		rowPadding = 0
	}
	rowRightPad := termWidth - len(rowInfo) - rowPadding
	if rowRightPad < 0 {
		rowRightPad = 0
	}
	fmt.Printf("%s%s%s%s%s\n",
		Bg, strings.Repeat(" ", rowPadding), rowInfoColored, strings.Repeat(" ", rowRightPad), Reset)

	// Keyboard shortcuts
	footerText := fmt.Sprintf("  %s←%s %s→%s pages  |  %s↑%s %s↓%s scroll  |  %sPgUp/PgDn%s jump  |  %sm%s menu  |  %sq%s quit  ",
		Accent, Text, Accent, Text, Accent, Text, Accent, Text, Accent, Text, Accent, Text, Error, Text)
	plainText := "  ← → pages  |  ↑ ↓ scroll  |  PgUp/PgDn jump  |  m menu  |  q quit  "
	padding := (termWidth - len(plainText)) / 2
	if padding < 0 {
		padding = 0
	}
	rightPad := termWidth - len(plainText) - padding
	if rightPad < 0 {
		rightPad = 0
	}
	fmt.Printf("%s%s%s%s%s\n", Bg, strings.Repeat(" ", padding), footerText, strings.Repeat(" ", rightPad), Reset)
}
