package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// Terminal functions are implemented in terminal_unix.go and terminal_windows.go

// getTerminalState retrieves the current terminal attributes.
func getTerminalState() (*termios, error) {
	var state termios
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), uintptr(0x5401), uintptr(unsafe.Pointer(&state))); errno != 0 {
		return nil, errno
	}
	return &state, nil
}

// Color constants for TM theme
const (
	// Background colors
	ColorReset   = "\033[0m"
	ColorBg      = "\033[48;2;26;27;38m" // Background #1a1b26
	ColorBgLight = "\033[48;2;36;40;59m" // Lighter background #24283b

	// Text colors
	ColorText       = "\033[38;2;192;202;245m" // Main text #c0caf5
	ColorTextDim    = "\033[38;2;115;131;168m" // Dimmed text #7387a8
	ColorTextBright = "\033[38;2;255;255;255m" // Bright white

	// Accent colors
	ColorAccent       = "\033[38;2;111;236;206m" // Your accent color #6fecce
	ColorAccentBright = "\033[38;2;122;162;247m" // Bright accent #7aa2f7
	ColorBlue         = "\033[38;2;125;162;206m" // Blue #7da2ce

	// Status colors
	ColorSuccess = "\033[38;2;102;232;236m" // Green rgb(102, 232, 236)
	ColorWarning = "\033[38;2;255;158;100m" // Orange #ff9e64
	ColorError   = "\033[38;2;247;118;142m" // Error red #f7768e

	// Border colors
	ColorBorder       = "\033[38;2;60;63;83m"    // Border color #3c3f53
	ColorBorderBright = "\033[38;2;122;162;247m" // Bright border for focus

	// Clear and positioning
	ColorClear     = "\033[2J\033[H"
	ColorClearLine = "\033[2K"
)

// SteelProperty defines the structure for a single steel section property.
// Note the use of interface{} for fields that might be numeric or string (like "-").
type SteelProperty struct {
	Section string      `json:"Section"`
	Grade   int         `json:"Grade"`
	Weight  float64     `json:"Weight"`
	D       float64     `json:"d"`
	Bf      float64     `json:"bf"`
	Tf      float64     `json:"tf"`
	Tw      float64     `json:"tw"`
	R1      interface{} `json:"r1"`
	D1      float64     `json:"d1"`
	Tw1     interface{} `json:"tw__1"`
	Tf1     interface{} `json:"tf__1"`
	Ag      float64     `json:"Ag"`
	Ix      float64     `json:"Ix"`
	Zx      float64     `json:"Zx"`
	Sx      float64     `json:"Sx"`
	Rx      float64     `json:"rx"`
	Iy      float64     `json:"Iy"`
	Zy      float64     `json:"Zy"`
	Sy      float64     `json:"Sy"`
	Ry      float64     `json:"ry"`
	J       float64     `json:"J"`
	Iw      interface{} `json:"Iw"`
	Flange  interface{} `json:"flange"`
	Web     interface{} `json:"web"`
	Kf      interface{} `json:"kf"`
	CNS     interface{} `json:"-"` // Handled by custom unmarshaler
	Zex     float64     `json:"Zex"`
	CNS2    interface{} `json:"-"` // Handled by custom unmarshaler
	Zey     float64     `json:"Zey"`
	TwoTf   interface{} `json:"2tf"`

	// Additional fields for UA tables
	Zy5      float64     `json:"Zy5"`
	TanAlpha float64     `json:"Tan Alpha"`
	AlphaB   interface{} `json:"αb"`
	Fu       interface{} `json:"Fu"`
	R2       interface{} `json:"r2"`
	ZeyD     float64     `json:"ZeyD"`
	In       float64     `json:"In"`
	Ip       float64     `json:"Ip"`
	ZexC     float64     `json:"ZexC"`
	X5       interface{} `json:"x5"`
	Y5       float64     `json:"y5"`
	NL       float64     `json:"nL"`
	PB       float64     `json:"pB"`
	PT       interface{} `json:"pT"`
	Residual string      `json:"Residual"`
	Type     interface{} `json:"Type"`
}

// UnmarshalJSON is a custom unmarshaler for SteelProperty to handle JSON fields with commas.
func (sp *SteelProperty) UnmarshalJSON(data []byte) error {
	// Use an alias to avoid recursion
	type Alias SteelProperty
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(sp),
	}

	// First, unmarshal into the alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Now, unmarshal into a map to get the fields with commas
	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	// Manually assign the values from the map to the struct fields
	if val, ok := rawMap["C,N,S"]; ok {
		sp.CNS = val
	}
	if val, ok := rawMap["C,N,S__1"]; ok {
		sp.CNS2 = val
	}

	return nil
}

func main() {
	fmt.Print(ColorBg + ColorClear)
	// Defer the color reset as well.
	defer fmt.Print(ColorReset)

	if len(os.Args) < 2 {
		// Get the initial terminal state just once at the start for interactive mode.
		initialState, err := getTerminalState()
		if err != nil {
			// If we can't get the state, we can't safely run the interactive mode.
			log.Fatalf("Fatal: Could not get terminal state: %v", err)
		}
		// Defer the final restoration to ensure the terminal is clean on exit.
		defer restoreTerminal(initialState)

		// Run interactive mode with the known initial state.
		runInteractiveMode(initialState)
	} else {
		// Run command-line mode, which doesn't need state management.
		runCommandLineMode(os.Args[1])
	}

	// Final cleanup: clear the screen before the program fully exits.
	fmt.Print(ColorClear)
}

func runInteractiveMode(initialState *termios) {
	for {
		// CRITICAL: Always ensure we're in canonical mode for the menu
		// This allows typing and Enter to work correctly
		restoreTerminal(initialState)

		selectedFile := printWelcomeScreenInteractive()
		if selectedFile == "" { // User chose to quit from the menu
			return // Exit the function, defer in main() will handle cleanup.
		}

		filePath := filepath.Join("data", selectedFile)

		// Switch to raw mode for table navigation
		err := setRawMode()
		if err != nil {
			log.Printf("Error entering raw mode: %v. Returning to menu.", err)
			continue // Skip to the next loop iteration, which will restore the terminal.
		}

		// Display the table in raw mode
		returnToMenu := displayTable(filePath)

		// CRITICAL: Always restore to canonical mode after table view
		// This ensures the menu will work correctly on the next iteration
		restoreTerminal(initialState)

		if !returnToMenu { // User chose to quit from the table view
			break
		}
		// Loop continues - next iteration will ensure canonical mode for menu
	}
}

func runCommandLineMode(tableName string) {
	filename := strings.ToUpper(tableName)
	if !strings.HasSuffix(filename, "_PROPS") {
		filename += "_PROPS"
	}
	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}
	filePath := filepath.Join("data", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("Table '%s' not found.", tableName)
	}

	// For command-line mode, just print the table once without any interaction.
	printTableOnce(filePath)
}

// printTableOnce is a new function for non-interactive display.
func printTableOnce(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var properties []SteelProperty
	if err := json.Unmarshal(data, &properties); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	allColumns := getAllColumns()
	availableColumns := filterAvailableColumns(allColumns, properties)
	maxCols := getMaxCols()
	totalPages := (len(availableColumns) + maxCols - 1) / maxCols

	// Use the background color, but don't clear the whole screen like in interactive mode.
	fmt.Print(ColorBg)

	for i := 0; i < totalPages; i++ {
		startCol := i * maxCols
		endCol := startCol + maxCols
		if endCol > len(availableColumns) {
			endCol = len(availableColumns)
		}
		currentColumns := availableColumns[startCol:endCol]
		drawHeader(filepath.Base(filePath), i+1, totalPages, len(properties))
		drawColumnHeaders(currentColumns)
		drawDataRows(properties, currentColumns)
		fillRemainingSpace()
		if i < totalPages-1 {
			fmt.Println() // Add space between pages
		}
	}
	// Reset color at the very end.
	fmt.Print(ColorReset)
}

// displayTable now correctly assumes it is ONLY called when in RAW mode.
// It supports horizontal column paging and vertical row scrolling with fixed headers.
func displayTable(filePath string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var properties []SteelProperty
	if err := json.Unmarshal(data, &properties); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	allColumns := getAllColumns()
	availableColumns := filterAvailableColumns(allColumns, properties)

	currentPage := 0
	scrollRow := 0 // vertical scroll offset

	for {
		// Re-read terminal dimensions on every redraw (handles resize)
		termHeight := getTerminalHeight()
		maxCols := getMaxCols()

		// Header takes 5 lines (box + blank), col headers 2, footer 3, so data gets the rest
		visibleRows := termHeight - 10
		if visibleRows < 3 {
			visibleRows = 3
		}

		fmt.Print(ColorBg + ColorClear)

		startCol := currentPage * maxCols
		endCol := startCol + maxCols
		if endCol > len(availableColumns) {
			endCol = len(availableColumns)
		}

		currentColumns := availableColumns[startCol:endCol]
		totalPages := (len(availableColumns) + maxCols - 1) / maxCols

		// Clamp scroll offset
		maxScroll := len(properties) - visibleRows
		if maxScroll < 0 {
			maxScroll = 0
		}
		if scrollRow > maxScroll {
			scrollRow = maxScroll
		}

		// Determine visible slice of rows
		endRow := scrollRow + visibleRows
		if endRow > len(properties) {
			endRow = len(properties)
		}
		visibleProperties := properties[scrollRow:endRow]

		// Draw fixed header + column headers
		drawHeader(filepath.Base(filePath), currentPage+1, totalPages, len(properties))
		drawColumnHeaders(currentColumns)

		// Draw only the visible rows (with correct alternating colors based on original index)
		drawDataRowsOffset(visibleProperties, currentColumns, scrollRow)

		// Fill any remaining empty lines so the footer stays at the bottom
		drawnRows := len(visibleProperties)
		for i := drawnRows; i < visibleRows; i++ {
			termWidth := getTerminalWidth()
			fmt.Printf("%s%s%s\n", ColorBg, strings.Repeat(" ", termWidth), ColorReset)
		}

		drawNavigationFooter(currentPage, totalPages, scrollRow, endRow, len(properties))

		// Read a single byte in raw mode
		buffer := make([]byte, 128)
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			return false
		}
		input := buffer[:n]

		switch {
		case len(input) == 1 && (input[0] == 'q' || input[0] == 'Q' || input[0] == 3): // q, Q, or Ctrl+C
			return false // Quit
		case len(input) == 1 && (input[0] == 'm' || input[0] == 'M'):
			return true // Return to menu
		case len(input) == 1 && (input[0] == '>'):
			if endCol < len(availableColumns) {
				currentPage++
			}
		case len(input) == 1 && (input[0] == '<'):
			if currentPage > 0 {
				currentPage--
			}
		case len(input) == 3 && input[0] == 27 && input[1] == 91: // Arrow keys
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
		case len(input) == 4 && input[0] == 27 && input[1] == 91: // Page Up/Down
			switch input[3] {
			case 126:
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
}

func printWelcomeScreenInteractive() string {
	for {
		termWidth := getTerminalWidth()
		fmt.Print(ColorClear)

		// Title
		titleText := "STEEL TABLES VIEWER"
		titlePadding := 14
		titleBoxWidth := len(titleText) + (titlePadding * 2)
		if titleBoxWidth < 60 {
			titleBoxWidth = 60
		}
		if titleBoxWidth > termWidth-4 {
			titleBoxWidth = termWidth - 4
		}
		centerOffset := (termWidth - titleBoxWidth - 2) / 2
		if centerOffset < 0 {
			centerOffset = 0
		}
		remainingSpace := termWidth - centerOffset - titleBoxWidth - 2
		if remainingSpace < 0 {
			remainingSpace = 0
		}

		fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
		fmt.Printf("%s╔%s╗", ColorBorderBright, strings.Repeat("═", titleBoxWidth))
		fmt.Printf("%s%s%s\n", ColorBg, strings.Repeat(" ", remainingSpace), ColorReset)

		textPadding := (titleBoxWidth - len(titleText)) / 2
		leftPadding := textPadding
		rightPadding := titleBoxWidth - len(titleText) - leftPadding

		fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
		fmt.Printf("%s║%s%s%s%s%s%s%s║",
			ColorBorderBright, ColorBg, strings.Repeat(" ", leftPadding), ColorAccent, titleText, ColorBg, strings.Repeat(" ", rightPadding), ColorBorderBright)
		fmt.Printf("%s%s%s\n", ColorBg, strings.Repeat(" ", remainingSpace), ColorReset)

		fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
		fmt.Printf("%s╚%s╝", ColorBorderBright, strings.Repeat("═", titleBoxWidth))
		fmt.Printf("%s%s%s\n\n", ColorBg, strings.Repeat(" ", remainingSpace), ColorReset)

		// Available files
		fmt.Printf("%s%s▶ AVAILABLE STEEL TABLES:%s%s\n", ColorBg, ColorAccent,
			strings.Repeat(" ", termWidth-len("▶ AVAILABLE STEEL TABLES:")), ColorReset)
		listJSONFilesStyledFullWidth()

		// Instructions
		fmt.Printf("\n%s%s▶ INSTRUCTIONS:%s%s\n", ColorBg, ColorAccent,
			strings.Repeat(" ", termWidth-len("▶ INSTRUCTIONS:")), ColorReset)
		line1 := "  • Type the table name (e.g., PFC300, RHS450, UB350)"
		fmt.Printf("%s%s%s%s\n", ColorBg, line1,
			strings.Repeat(" ", termWidth-len(line1)), ColorReset)
		line2 := "  • Type q or quit to exit"
		fmt.Printf("%s%s%s%s\n", ColorBg, line2,
			strings.Repeat(" ", termWidth-len(line2)), ColorReset)
		line3 := "  • Press Enter to confirm your selection"
		fmt.Printf("%s%s%s%s\n\n", ColorBg, line3,
			strings.Repeat(" ", termWidth-len(line3)), ColorReset)

		// Input prompt
		fmt.Printf("%s▶ SELECT TABLE: %s", ColorAccent, ColorReset)

		// Use bufio.Reader to read a full line, which is more reliable than fmt.Scanln
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			// On EOF or other errors, treat as a quit command
			return ""
		}
		input = strings.TrimSpace(strings.ToUpper(input))

		if input == "Q" || input == "QUIT" || input == "EXIT" {
			return ""
		}

		if input == "" {
			continue // Empty input, show menu again
		}

		if isValidTable(input) {
			return input + "_PROPS.json"
		} else {
			fmt.Printf("\n%s✗ Table '%s' not found. Please try again...%s\n", ColorError, input, ColorReset)
			fmt.Printf("%sPress Enter to continue...%s", ColorTextDim, ColorReset)
			reader.ReadString('\n') // Wait for user to press Enter
		}
	}
}

func isValidTable(tableName string) bool {
	files, err := os.ReadDir("data")
	if err != nil {
		return false
	}
	expectedFilename := strings.ToUpper(tableName) + "_PROPS.json"
	for _, file := range files {
		if file.Name() == expectedFilename {
			return true
		}
	}
	return false
}

func listJSONFilesStyledFullWidth() {
	termWidth := getTerminalWidth()
	files, err := os.ReadDir("data")
	if err != nil {
		errorLine := fmt.Sprintf("✗ Error reading data directory: %v", err)
		fmt.Printf("%s%s%s%s%s\n", ColorBg, ColorError, errorLine,
			strings.Repeat(" ", termWidth-len(errorLine)), ColorReset)
		return
	}

	for i, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			displayName := strings.TrimSuffix(file.Name(), "_PROPS.json")
			line := fmt.Sprintf("  ● %s", displayName)
			padding := termWidth - len(line)
			if padding < 0 {
				padding = 0
			}

			if i%2 == 0 {
				fmt.Printf("%s%s  ● %s%s%s%s\n", ColorBg, ColorAccent, ColorTextBright, displayName,
					strings.Repeat(" ", padding), ColorReset)
			} else {
				fmt.Printf("%s%s  ● %s%s%s%s\n", ColorBg, ColorBlue, ColorText, displayName,
					strings.Repeat(" ", padding), ColorReset)
			}
		}
	}
}

type ColumnInfo struct {
	Name      string
	Formatter func(SteelProperty) string
}

func formatInterface(value interface{}) string {
	if value == nil {
		return "-"
	}
	switch v := value.(type) {
	case string:
		if v == "" || v == "-" {
			return "-"
		}
		return v
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%.0f", v)
		}
		return fmt.Sprintf("%.1f", v)
	case int:
		return fmt.Sprintf("%d", v)
	default:
		str := fmt.Sprintf("%v", v)
		if str == "" {
			return "-"
		}
		return str
	}
}

func cleanSectionName(section string) string {
	if idx := strings.Index(section, " (G"); idx != -1 {
		if endIdx := strings.Index(section[idx:], ")"); endIdx != -1 {
			return section[:idx] + section[idx+endIdx+1:]
		}
	}
	return section
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 2 {
		return s[:maxLen]
	}
	return s[:maxLen-1] + "."
}

func getTerminalWidth() int {
	cmd := exec.Command("tput", "cols")
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil {
		return 120 // Default fallback
	}
	width, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil || width < 80 {
		return 120 // Default fallback
	}
	if width > 200 {
		width = 200
	}
	return width
}

func getTerminalHeight() int {
	cmd := exec.Command("tput", "lines")
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil {
		return 40 // Default fallback
	}
	height, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil || height < 10 {
		return 40 // Default fallback
	}
	return height
}

// getMaxCols computes how many data columns fit in the terminal width.
func getMaxCols() int {
	termWidth := getTerminalWidth()
	sectionColWidth := 25
	colWidth := 18
	maxCols := (termWidth - sectionColWidth) / colWidth
	if maxCols < 1 {
		maxCols = 1
	}
	return maxCols
}

func drawHeader(filename string, currentPage, totalPages, totalEntries int) {
	termWidth := getTerminalWidth()
	titleText := fmt.Sprintf("STEEL PROPERTIES: %s", strings.ToUpper(strings.TrimSuffix(filename, ".json")))
	infoText := fmt.Sprintf("Page %d/%d | %d entries", currentPage, totalPages, totalEntries)

	// Ensure no trailing spaces in infoText
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
	fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
	fmt.Printf("%s╔%s╗", ColorBorderBright, strings.Repeat("═", boxWidth))
	fmt.Printf("%s%s%s\n", ColorBg, strings.Repeat(" ", remainingSpace), ColorReset)

	// Title row
	titlePadding := (boxWidth - len(titleText)) / 2
	titleRightPadding := boxWidth - len(titleText) - titlePadding
	fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
	fmt.Printf("%s║%s%s%s%s%s%s%s║",
		ColorBorderBright, ColorBg, strings.Repeat(" ", titlePadding), ColorAccent, titleText, ColorBg, strings.Repeat(" ", titleRightPadding), ColorBorderBright)
	fmt.Printf("%s%s%s\n", ColorBg, strings.Repeat(" ", remainingSpace), ColorReset)

	// Info row
	infoPadding := (boxWidth - len(infoText)) / 2
	infoRightPadding := boxWidth - len(infoText) - infoPadding
	fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
	fmt.Printf("%s║%s%s%s%s%s%s%s║",
		ColorBorderBright, ColorBg, strings.Repeat(" ", infoPadding), ColorTextDim, infoText, ColorBg, strings.Repeat(" ", infoRightPadding), ColorBorderBright)
	fmt.Printf("%s%s%s\n", ColorBg, strings.Repeat(" ", remainingSpace), ColorReset)

	// Bottom border
	fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
	fmt.Printf("%s╚%s╝", ColorBorderBright, strings.Repeat("═", boxWidth))
	fmt.Printf("%s%s%s\n\n", ColorBg, strings.Repeat(" ", remainingSpace), ColorReset)
}

func drawColumnHeaders(currentColumns []ColumnInfo) {
	termWidth := getTerminalWidth()
	fmt.Printf("%s%s", ColorBgLight, ColorAccent)
	fmt.Printf("%-25s", "Section")
	for _, col := range currentColumns {
		headerText := getColumnHeaderWithUnit(col.Name)
		fmt.Printf("%-18s", truncateString(headerText, 17))
	}
	usedSpace := 25 + (len(currentColumns) * 18)
	remainingSpace := termWidth - usedSpace
	if remainingSpace > 0 {
		fmt.Printf("%s", strings.Repeat(" ", remainingSpace))
	}
	fmt.Printf("%s\n", ColorReset)

	fmt.Printf("%s%s", ColorBgLight, ColorBorderBright)
	fmt.Print(strings.Repeat("─", 25))
	for range currentColumns {
		fmt.Print(strings.Repeat("─", 18))
	}
	if remainingSpace > 0 {
		fmt.Print(strings.Repeat("─", remainingSpace))
	}
	fmt.Printf("%s\n", ColorReset)
}

// getColumnHeaderWithUnit returns the column name with appropriate unit
func getColumnHeaderWithUnit(columnName string) string {
	unitMap := map[string]string{
		"Weight": "kg/m", "d": "mm", "bf": "mm", "tf": "mm", "tw": "mm", "r1": "mm", "d1": "mm",
		"tw__1": "mm", "tf__1": "mm", "Ag": "mm²", "Ix": "10⁶mm⁴", "Zx": "10³mm³", "Sx": "10³mm³",
		"rx": "mm", "Iy": "10⁶mm⁴", "Zy": "mm³", "Sy": "mm³", "ry": "mm", "J": "10³mm⁴", "Iw": "mm⁶",
		"flange": "mm", "web": "mm", "Zex": "mm³", "Zey": "mm³", "Zy5": "mm³", "Fu": "MPa", "r2": "mm",
		"ZeyD": "mm³", "In": "10³mm⁴", "Ip": "10³mm⁴", "ZexC": "mm³", "x5": "mm", "y5": "mm", "nL": "mm",
		"pB": "mm", "pT": "mm",
	}
	if unit, exists := unitMap[columnName]; exists {
		return fmt.Sprintf("%s (%s)", columnName, unit)
	}
	return columnName
}

func drawDataRows(properties []SteelProperty, currentColumns []ColumnInfo) {
	drawDataRowsOffset(properties, currentColumns, 0)
}

// drawDataRowsOffset draws data rows with a base offset for correct alternating row colors.
func drawDataRowsOffset(properties []SteelProperty, currentColumns []ColumnInfo, baseIndex int) {
	termWidth := getTerminalWidth()
	for i, prop := range properties {
		globalIndex := baseIndex + i
		if globalIndex%2 == 0 {
			fmt.Printf("%s", ColorBg)
		} else {
			fmt.Printf("%s", ColorBgLight)
		}

		cleanedSection := cleanSectionName(prop.Section)
		fmt.Printf("%s%-25s%s", ColorTextBright, truncateString(cleanedSection, 24), ColorText)

		for _, col := range currentColumns {
			value := col.Formatter(prop)
			if value == "-" || value == "" {
				fmt.Printf("%s%-18s", ColorTextDim, truncateString(value, 17))
			} else {
				fmt.Printf("%s%-18s", ColorText, truncateString(value, 17))
			}
		}

		usedSpace := 25 + (len(currentColumns) * 18)
		remainingSpace := termWidth - usedSpace
		if remainingSpace > 0 {
			fmt.Printf("%s", strings.Repeat(" ", remainingSpace))
		}
		fmt.Printf("%s\n", ColorReset)
	}
}

func fillRemainingSpace() {
	// This function can be used to fill the rest of the screen if needed,
	// but the current drawing logic handles it row by row.
}

func drawNavigationFooter(currentPage, totalPages, startRow, endRow, totalRows int) {
	termWidth := getTerminalWidth()

	// Row info line
	rowInfo := fmt.Sprintf("Rows %d–%d of %d", startRow+1, endRow, totalRows)
	rowInfoColored := fmt.Sprintf("%sRows %s%d–%d%s of %s%d%s",
		ColorTextDim, ColorAccent, startRow+1, endRow, ColorTextDim, ColorAccent, totalRows, ColorTextDim)
	rowPadding := (termWidth - len(rowInfo)) / 2
	if rowPadding < 0 {
		rowPadding = 0
	}
	rowRightPad := termWidth - len(rowInfo) - rowPadding
	if rowRightPad < 0 {
		rowRightPad = 0
	}
	fmt.Printf("%s%s%s%s%s\n",
		ColorBg, strings.Repeat(" ", rowPadding), rowInfoColored, strings.Repeat(" ", rowRightPad), ColorReset)

	// Build the footer with colored text
	footerText := fmt.Sprintf("  %s←%s %s→%s pages  |  %s↑%s %s↓%s scroll  |  %sPgUp/PgDn%s jump  |  %sm%s menu  |  %sq%s quit  ",
		ColorAccent, ColorText, // ←
		ColorAccent, ColorText, // →
		ColorAccent, ColorText, // ↑
		ColorAccent, ColorText, // ↓
		ColorAccent, ColorText, // PgUp/PgDn
		ColorAccent, ColorText, // m
		ColorError, ColorText) // q

	// Calculate plain text length for centering (without color codes)
	plainText := "  ← → pages  |  ↑ ↓ scroll  |  PgUp/PgDn jump  |  m menu  |  q quit  "
	padding := (termWidth - len(plainText)) / 2
	if padding < 0 {
		padding = 0
	}
	rightPad := termWidth - len(plainText) - padding
	if rightPad < 0 {
		rightPad = 0
	}

	fmt.Printf("%s%s%s%s%s\n",
		ColorBg, strings.Repeat(" ", padding), footerText, strings.Repeat(" ", rightPad), ColorReset)
}

func filterAvailableColumns(allColumns []ColumnInfo, properties []SteelProperty) []ColumnInfo {
	var availableColumns []ColumnInfo
	for _, col := range allColumns {
		hasNonDashData := false
		hasRealData := false

		for _, p := range properties {
			val := col.Formatter(p)

			// Check if any value is not a dash
			if val != "-" {
				hasNonDashData = true

				// Check if the value represents actual meaningful data (not empty or zero)
				if val != "" {
					// Check if it's a numeric zero in various formats
					if val == "0" || val == "0.0" || val == "0.00" || val == "0.000" {
						continue // This is still considered "no real data"
					}
					// Check if it's a float that formats to zero
					if f, err := strconv.ParseFloat(val, 64); err == nil && f == 0.0 {
						continue // This is still considered "no real data"
					}
					hasRealData = true
				}
			}
		}

		// Include column if it has non-dash values AND real data
		// This will exclude columns that are all dashes OR all zeros
		if hasNonDashData && hasRealData {
			availableColumns = append(availableColumns, col)
		}
	}
	return availableColumns
}

func getAllColumns() []ColumnInfo {
	return []ColumnInfo{
		{"Grade", func(p SteelProperty) string { return fmt.Sprintf("%d", p.Grade) }},
		{"Weight", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.Weight) }},
		{"d", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.D) }},
		{"bf", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.Bf) }},
		{"tf", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.Tf) }},
		{"tw", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.Tw) }},
		{"r1", func(p SteelProperty) string { return formatInterface(p.R1) }},
		{"d1", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.D1) }},
		{"tw__1", func(p SteelProperty) string { return formatInterface(p.Tw1) }},
		{"tf__1", func(p SteelProperty) string { return formatInterface(p.Tf1) }},
		{"Ag", func(p SteelProperty) string { return fmt.Sprintf("%.0f", p.Ag) }},
		{"Ix", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.Ix) }},
		{"Zx", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.Zx) }},
		{"Sx", func(p SteelProperty) string { return fmt.Sprintf("%.0f", p.Sx) }},
		{"rx", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.Rx) }},
		{"Iy", func(p SteelProperty) string { return fmt.Sprintf("%.2f", p.Iy) }},
		{"Zy", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.Zy) }},
		{"Sy", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.Sy) }},
		{"ry", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.Ry) }},
		{"J", func(p SteelProperty) string { return fmt.Sprintf("%.0f", p.J) }},
		{"Iw", func(p SteelProperty) string { return formatInterface(p.Iw) }},
		{"flange", func(p SteelProperty) string { return formatInterface(p.Flange) }},
		{"web", func(p SteelProperty) string { return formatInterface(p.Web) }},
		{"kf", func(p SteelProperty) string { return formatInterface(p.Kf) }},
		{"C,N,S", func(p SteelProperty) string { return formatInterface(p.CNS) }},
		{"Zex", func(p SteelProperty) string { return fmt.Sprintf("%.0f", p.Zex) }},
		{"C,N,S__1", func(p SteelProperty) string { return formatInterface(p.CNS2) }},
		{"Zey", func(p SteelProperty) string { return fmt.Sprintf("%.1f", p.Zey) }},
		{"2tf", func(p SteelProperty) string { return formatInterface(p.TwoTf) }},
		// Additional UA-specific columns
		{"Zy5", func(p SteelProperty) string {
			if p.Zy5 == 0 {
				return "-"
			}
			return fmt.Sprintf("%.1f", p.Zy5)
		}},
		{"TanAlpha", func(p SteelProperty) string {
			if p.TanAlpha == 0 {
				return "-"
			}
			return fmt.Sprintf("%.3f", p.TanAlpha)
		}},
		{"αb", func(p SteelProperty) string { return formatInterface(p.AlphaB) }},
		{"Fu", func(p SteelProperty) string { return formatInterface(p.Fu) }},
		{"r2", func(p SteelProperty) string { return formatInterface(p.R2) }},
		{"ZeyD", func(p SteelProperty) string {
			if p.ZeyD == 0 {
				return "-"
			}
			return fmt.Sprintf("%.1f", p.ZeyD)
		}},
		{"In", func(p SteelProperty) string {
			if p.In == 0 {
				return "-"
			}
			return fmt.Sprintf("%.2f", p.In)
		}},
		{"Ip", func(p SteelProperty) string {
			if p.Ip == 0 {
				return "-"
			}
			return fmt.Sprintf("%.2f", p.Ip)
		}},
		{"ZexC", func(p SteelProperty) string {
			if p.ZexC == 0 {
				return "-"
			}
			return fmt.Sprintf("%.0f", p.ZexC)
		}},
		{"x5", func(p SteelProperty) string { return formatInterface(p.X5) }},
		{"y5", func(p SteelProperty) string {
			if p.Y5 == 0 {
				return "-"
			}
			return fmt.Sprintf("%.1f", p.Y5)
		}},
		{"nL", func(p SteelProperty) string {
			if p.NL == 0 {
				return "-"
			}
			return fmt.Sprintf("%.1f", p.NL)
		}},
		{"pB", func(p SteelProperty) string {
			if p.PB == 0 {
				return "-"
			}
			return fmt.Sprintf("%.1f", p.PB)
		}},
		{"pT", func(p SteelProperty) string { return formatInterface(p.PT) }},
		{"Residual", func(p SteelProperty) string {
			if p.Residual == "" {
				return "-"
			}
			return p.Residual
		}},
		{"Type", func(p SteelProperty) string { return formatInterface(p.Type) }},
	}
}
