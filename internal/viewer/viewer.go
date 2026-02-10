// Package viewer handles the interactive table display.
package viewer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"steel_tables/internal/columns"
	"steel_tables/internal/models"
	"steel_tables/internal/ui"
)

// DisplayTable shows an interactive table view with scrolling and paging.
// Returns true if user wants to return to menu, false to quit.
func DisplayTable(filePath string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var properties []models.SteelProperty
	if err := json.Unmarshal(data, &properties); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	allColumns := columns.GetAll()
	availableColumns := columns.FilterAvailable(allColumns, properties)

	currentPage := 0
	scrollRow := 0

	for {
		termHeight := ui.GetTerminalHeight()
		maxCols := ui.GetMaxCols()

		visibleRows := termHeight - 10
		if visibleRows < 3 {
			visibleRows = 3
		}

		fmt.Print(ui.Bg + ui.Clear)

		startCol := currentPage * maxCols
		endCol := startCol + maxCols
		if endCol > len(availableColumns) {
			endCol = len(availableColumns)
		}

		currentColumns := availableColumns[startCol:endCol]
		totalPages := (len(availableColumns) + maxCols - 1) / maxCols

		maxScroll := len(properties) - visibleRows
		if maxScroll < 0 {
			maxScroll = 0
		}
		if scrollRow > maxScroll {
			scrollRow = maxScroll
		}

		endRow := scrollRow + visibleRows
		if endRow > len(properties) {
			endRow = len(properties)
		}
		visibleProperties := properties[scrollRow:endRow]

		ui.DrawHeader(filepath.Base(filePath), currentPage+1, totalPages, len(properties))
		ui.DrawColumnHeaders(currentColumns)
		ui.DrawDataRowsOffset(visibleProperties, currentColumns, scrollRow)

		// Fill empty lines
		drawnRows := len(visibleProperties)
		termWidth := ui.GetTerminalWidth()
		for i := drawnRows; i < visibleRows; i++ {
			fmt.Printf("%s%s%s\n", ui.Bg, strings.Repeat(" ", termWidth), ui.Reset)
		}

		ui.DrawNavigationFooter(currentPage, totalPages, scrollRow, endRow, len(properties))

		// Handle input
		buffer := make([]byte, 128)
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			return false
		}
		input := buffer[:n]

		switch {
		case len(input) == 1 && (input[0] == 'q' || input[0] == 'Q' || input[0] == 3):
			return false
		case len(input) == 1 && (input[0] == 'm' || input[0] == 'M'):
			return true
		case len(input) == 1 && input[0] == '>':
			if endCol < len(availableColumns) {
				currentPage++
			}
		case len(input) == 1 && input[0] == '<':
			if currentPage > 0 {
				currentPage--
			}
		case len(input) == 3 && input[0] == 27 && input[1] == 91:
			switch input[2] {
			case 65: // Up
				if scrollRow > 0 {
					scrollRow--
				}
			case 66: // Down
				if scrollRow < maxScroll {
					scrollRow++
				}
			case 67: // Right
				if endCol < len(availableColumns) {
					currentPage++
				}
			case 68: // Left
				if currentPage > 0 {
					currentPage--
				}
			}
		case len(input) == 4 && input[0] == 27 && input[1] == 91 && input[3] == 126:
			if input[2] == 53 { // Page Up
				scrollRow -= visibleRows
				if scrollRow < 0 {
					scrollRow = 0
				}
			} else if input[2] == 54 { // Page Down
				scrollRow += visibleRows
				if scrollRow > maxScroll {
					scrollRow = maxScroll
				}
			}
		}
	}
}

// PrintTableOnce prints the table non-interactively (for CLI mode).
func PrintTableOnce(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var properties []models.SteelProperty
	if err := json.Unmarshal(data, &properties); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	allColumns := columns.GetAll()
	availableColumns := columns.FilterAvailable(allColumns, properties)
	maxCols := ui.GetMaxCols()
	totalPages := (len(availableColumns) + maxCols - 1) / maxCols

	fmt.Print(ui.Bg)

	for i := 0; i < totalPages; i++ {
		startCol := i * maxCols
		endCol := startCol + maxCols
		if endCol > len(availableColumns) {
			endCol = len(availableColumns)
		}
		currentColumns := availableColumns[startCol:endCol]
		ui.DrawHeader(filepath.Base(filePath), i+1, totalPages, len(properties))
		ui.DrawColumnHeaders(currentColumns)
		ui.DrawDataRows(properties, currentColumns)
		if i < totalPages-1 {
			fmt.Println()
		}
	}
	fmt.Print(ui.Reset)
}
