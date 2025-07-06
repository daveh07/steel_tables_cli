package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Color constants for Tokyo Midnight theme
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
	ColorAccent   = "\033[38;2;111;236;206m"                    // Your accent color #6fecce
	ColorAccentBg = "\033[48;2;111;236;206m\033[38;2;26;27;38m" // Accent background
	ColorBlue     = "\033[38;2;125;162;206m"                    // Blue #7da2ce

	// Status colors
	ColorSuccess = "\033[38;2;102;232;236m" // Green rgb(102, 232, 236)
	ColorWarning = "\033[38;2;255;158;100m" // Orange #ff9e64
	ColorError   = "\033[38;2;247;118;142m" // Red #f7768e

	// Border colors
	ColorBorder       = "\033[38;2;52;59;88m"    // Border #343b58
	ColorBorderBright = "\033[38;2;111;236;206m" // Bright border (accent)

	// Clear and positioning
	ColorClear     = "\033[2J\033[H"
	ColorClearLine = "\033[2K"
)

// SteelProperty represents a single steel property entry with all available fields
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
	CNS     interface{} `json:"C,N,S"`
	Zex     float64     `json:"Zex"`
	CNS2    interface{} `json:"C,N,S__1"`
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

func main() {
	// Set background color for the entire terminal
	fmt.Print(ColorBg + ColorClear)
	defer fmt.Print(ColorReset) // Reset colors when exiting

	if len(os.Args) < 2 {
		// Interactive mode - show menu and get user input
		selectedFile := printWelcomeScreenInteractive()
		if selectedFile == "" {
			return // User chose to quit
		}

		filename := selectedFile
		if !strings.HasSuffix(filename, ".json") {
			filename += ".json"
		}

		filePath := filepath.Join("data", filename)

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			printErrorScreen(fmt.Sprintf("File %s not found in data folder", filename))
			return
		}

		displayTable(filePath)
		return
	}

	// Command line mode - original behavior
	filename := os.Args[1]
	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}

	filePath := filepath.Join("data", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		printErrorScreen(fmt.Sprintf("File %s not found in data folder", filename))
		return
	}

	displayTable(filePath)
}

func printWelcomeScreen() {
	fmt.Print(ColorClear)

	// Title with accent color
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════════════════════╗%s\n", ColorBorderBright, ColorReset)
	fmt.Printf("%s║%s                              %sSTEEL TABLES VIEWER%s                       %s║%s\n", ColorBorderBright, ColorBg, ColorAccent, ColorBg, ColorBorderBright, ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════════════════╝%s\n\n", ColorBorderBright, ColorReset)

	// Usage section
	fmt.Printf("%s%s▶ USAGE:%s\n", ColorBg, ColorAccent, ColorReset)
	fmt.Printf("%s  %sgo run main.go <filename>%s\n\n", ColorBg, ColorTextBright, ColorReset)

	// Available files section
	fmt.Printf("%s%s▶ STEEL PROPERTIES:%s\n", ColorBg, ColorAccent, ColorReset)
	listJSONFilesStyled()

	// Instructions
	fmt.Printf("\n%s%s▶ NAVIGATION:%s\n", ColorBg, ColorAccent, ColorReset)
	fmt.Printf("%s  %s>%s or %s→%s  Next page of columns\n", ColorBg, ColorSuccess, ColorText, ColorSuccess, ColorText)
	fmt.Printf("%s  %s<%s or %s←%s  Previous page of columns\n", ColorBg, ColorSuccess, ColorText, ColorSuccess, ColorText)
	fmt.Printf("%s  %sm%s        Main Menu\n\n", ColorBg, ColorAccent, ColorText)
	fmt.Printf("%s  %sq%s        Quit application\n\n", ColorBg, ColorError, ColorText)

	fmt.Print(ColorReset)
}

func printErrorScreen(message string) {
	fmt.Print(ColorClear)

	// Error header
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════════════════════╗%s\n", ColorError, ColorReset)
	fmt.Printf("%s║%s                                   %sERROR%s                                      %s║%s\n", ColorError, ColorBg, ColorError, ColorBg, ColorError, ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════════════════╝%s\n\n", ColorError, ColorReset)

	// Error message
	fmt.Printf("%s%s✗ %s%s\n\n", ColorBg, ColorError, message, ColorReset)

	// Available files section
	fmt.Printf("%s%s▶ AVAILABLE FILES:%s\n", ColorBg, ColorAccent, ColorReset)
	listJSONFilesStyled()

	fmt.Print(ColorReset)
}

func listJSONFilesStyled() {
	files, err := os.ReadDir("data")
	if err != nil {
		fmt.Printf("%s%s✗ Error reading data directory: %v%s\n", ColorBg, ColorError, err, ColorReset)
		return
	}

	for i, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			// Extract the identifying part (remove _PROPS.json)
			displayName := file.Name()
			if strings.HasSuffix(displayName, "_PROPS.json") {
				displayName = strings.TrimSuffix(displayName, "_PROPS.json")
			} else if strings.HasSuffix(displayName, ".json") {
				displayName = strings.TrimSuffix(displayName, ".json")
			}

			// Alternate colors for better readability
			if i%2 == 0 {
				fmt.Printf("%s  %s● %s%s%s\n", ColorBg, ColorAccent, ColorTextBright, displayName, ColorReset)
			} else {
				fmt.Printf("%s  %s● %s%s%s\n", ColorBg, ColorBlue, ColorText, displayName, ColorReset)
			}
		}
	}
}

// listJSONFilesStyledFullWidth displays available JSON files with full terminal width background
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
			// Extract the identifying part (remove _PROPS.json)
			displayName := file.Name()
			if strings.HasSuffix(displayName, "_PROPS.json") {
				displayName = strings.TrimSuffix(displayName, "_PROPS.json")
			} else if strings.HasSuffix(displayName, ".json") {
				displayName = strings.TrimSuffix(displayName, ".json")
			}

			// Create the line content
			var line string
			if i%2 == 0 {
				line = fmt.Sprintf("  ● %s", displayName)
			} else {
				line = fmt.Sprintf("  ● %s", displayName)
			}

			// Alternate colors for better readability with full width background
			if i%2 == 0 {
				fmt.Printf("%s%s%s%s%s%s\n", ColorBg, ColorAccent, "  ● ", ColorTextBright, displayName,
					strings.Repeat(" ", termWidth-len(line))+ColorReset)
			} else {
				fmt.Printf("%s%s%s%s%s%s\n", ColorBg, ColorBlue, "  ● ", ColorText, displayName,
					strings.Repeat(" ", termWidth-len(line))+ColorReset)
			}
		}
	}
}

func displayTable(filePath string) {
	// Set terminal to raw mode for immediate key reading
	oldState := setRawMode()
	defer restoreTerminal(oldState)

	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	// Parse JSON
	var properties []SteelProperty
	if err := json.Unmarshal(data, &properties); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	// Define all available columns (excluding Section which is always shown)
	allColumns := []ColumnInfo{
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

	// Filter out empty columns dynamically
	availableColumns := filterAvailableColumns(allColumns, properties)

	currentPage := 0
	maxCols := 7 // 7 additional columns + Section = 8 total (adjusted for wider columns to fit all text)

	for {
		// Clear screen and set background
		fmt.Print(ColorBg + ColorClear)

		// Calculate column range for current page
		startCol := currentPage * maxCols
		endCol := startCol + maxCols
		if endCol > len(availableColumns) {
			endCol = len(availableColumns)
		}

		currentColumns := availableColumns[startCol:endCol]
		totalPages := (len(availableColumns) + maxCols - 1) / maxCols

		// Beautiful header with accent colors
		drawHeader(filepath.Base(filePath), currentPage+1, totalPages, len(properties))

		// Column headers with styling
		drawColumnHeaders(currentColumns)

		// Data rows with alternating colors
		drawDataRows(properties, currentColumns)

		// Fill any remaining vertical space with background color
		fillRemainingSpace()

		// Navigation footer
		drawNavigationFooter(currentPage, totalPages)

		// Read single character input
		key := readKey()

		switch key {
		case 'q', 'Q':
			fmt.Print(ColorReset)
			return
		case 'm', 'M':
			// Return to main menu
			fmt.Print(ColorReset)
			restoreTerminal(oldState) // Restore terminal before returning to menu

			// Call main menu again
			selectedFile := printWelcomeScreenInteractive()
			if selectedFile == "" {
				return // User chose to quit from menu
			}

			// Load the new selected file
			newFilePath := filepath.Join("data", selectedFile)
			if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
				printErrorScreen(fmt.Sprintf("File %s not found in data folder", selectedFile))
				return
			}

			// Read and parse the new file
			data, err = os.ReadFile(newFilePath)
			if err != nil {
				log.Fatal("Error reading file:", err)
			}

			if err := json.Unmarshal(data, &properties); err != nil {
				log.Fatal("Error parsing JSON:", err)
			}

			// Re-filter columns for the new dataset
			availableColumns = filterAvailableColumns(allColumns, properties)

			// Update file path and reset to first page
			filePath = newFilePath
			currentPage = 0

			// Set raw mode again for the new table
			oldState = setRawMode()
			continue
		case '>':
			if endCol < len(availableColumns) {
				currentPage++
			}
		case '<':
			if currentPage > 0 {
				currentPage--
			}
		case 27: // ESC sequence for arrow keys
			// Read the next characters to handle arrow keys
			if b1, ok := readKeyNonBlocking(); ok && b1 == 91 { // '['
				if b2, ok := readKeyNonBlocking(); ok {
					switch b2 {
					case 67: // Right arrow
						if endCol < len(availableColumns) {
							currentPage++
						}
					case 68: // Left arrow
						if currentPage > 0 {
							currentPage--
						}
					}
				}
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

// cleanSectionName removes grade information like (G300) or (G350) from section names
func cleanSectionName(section string) string {
	// Remove patterns like " (G300)", " (G350)", etc.
	// This regex matches " (G" followed by digits and closing parenthesis
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
	// Very minimal truncation - allow almost all text to show
	if maxLen <= 2 {
		return s[:maxLen]
	}
	return s[:maxLen-1] + "."
}

func drawHeader(filename string, currentPage, totalPages, totalEntries int) {
	// Get actual terminal width for proper centering
	termWidth := getTerminalWidth()
	if termWidth < 100 {
		termWidth = 150 // Wider minimum for better display
	}

	// Calculate box width based on content or minimum width
	titleText := fmt.Sprintf("STEEL PROPERTIES: %s", strings.ToUpper(filename))
	infoText := fmt.Sprintf("Page %d/%d │ %d entries │ 8 columns per page", currentPage, totalPages, totalEntries)

	// Use the longer of the two texts to determine box width, with padding
	boxWidth := len(titleText)
	if len(infoText) > boxWidth {
		boxWidth = len(infoText)
	}
	boxWidth += 8 // Add padding

	// Ensure minimum width
	if boxWidth < 80 {
		boxWidth = 80
	}

	// Calculate centering offset for the entire box
	centerOffset := (termWidth - boxWidth - 4) / 2 // -4 for borders
	if centerOffset < 0 {
		centerOffset = 0
	}

	// Top border with centering
	fmt.Printf("%s%s%s╔%s╗%s\n", ColorBg, strings.Repeat(" ", centerOffset), ColorBorderBright,
		strings.Repeat("═", boxWidth), ColorReset)

	// Title row with centering
	titlePadding := (boxWidth - len(titleText)) / 2
	remainingTitlePadding := boxWidth - len(titleText) - titlePadding

	fmt.Printf("%s%s%s║%s", ColorBg, strings.Repeat(" ", centerOffset), ColorBorderBright, ColorBg)
	fmt.Printf("%s", strings.Repeat(" ", titlePadding))
	fmt.Printf("%s%s%s", ColorAccent, titleText, ColorBg)
	fmt.Printf("%s", strings.Repeat(" ", remainingTitlePadding))
	fmt.Printf("%s║%s\n", ColorBorderBright, ColorReset)

	// Info row with centering
	infoPadding := (boxWidth - len(infoText)) / 2
	remainingInfoPadding := boxWidth - len(infoText) - infoPadding

	fmt.Printf("%s%s%s║%s", ColorBg, strings.Repeat(" ", centerOffset), ColorBorderBright, ColorBg)
	fmt.Printf("%s", strings.Repeat(" ", infoPadding))
	fmt.Printf("%s%s%s", ColorTextDim, infoText, ColorBg)
	fmt.Printf("%s", strings.Repeat(" ", remainingInfoPadding))
	fmt.Printf("%s║%s\n", ColorBorderBright, ColorReset)

	// Bottom border with centering
	fmt.Printf("%s%s%s╚%s╝%s\n\n", ColorBg, strings.Repeat(" ", centerOffset), ColorBorderBright,
		strings.Repeat("═", boxWidth), ColorReset)
}

func drawColumnHeaders(currentColumns []ColumnInfo) {
	termWidth := getTerminalWidth()

	// Header background - fill entire terminal width
	fmt.Printf("%s%s", ColorBgLight, ColorAccent)

	// Section header (fixed column) - wider for better readability
	fmt.Printf("%-35s", "Section")

	// Dynamic columns with units - wider to fit full text
	for _, col := range currentColumns {
		headerText := getColumnHeaderWithUnit(col.Name)
		fmt.Printf("%-18s", truncateString(headerText, 17))
	}

	// Fill remaining space to terminal width
	usedSpace := 35 + (len(currentColumns) * 18)
	remainingSpace := termWidth - usedSpace
	if remainingSpace > 0 {
		fmt.Printf("%s", strings.Repeat(" ", remainingSpace))
	}
	fmt.Printf("%s\n", ColorReset)

	// Header separator with accent color - fill entire terminal width
	fmt.Printf("%s%s", ColorBgLight, ColorBorderBright)
	fmt.Print(strings.Repeat("─", 35))
	for range currentColumns {
		fmt.Print(strings.Repeat("─", 18))
	}
	// Fill remaining separator space
	if remainingSpace > 0 {
		fmt.Print(strings.Repeat("─", remainingSpace))
	}
	fmt.Printf("%s\n", ColorReset)
}

// getColumnHeaderWithUnit returns the column name with appropriate unit
func getColumnHeaderWithUnit(columnName string) string {
	unitMap := map[string]string{
		"Grade":    "",
		"Weight":   "(kg/m)",
		"d":        "(mm)",
		"bf":       "(mm)",
		"tf":       "(mm)",
		"tw":       "(mm)",
		"r1":       "(mm)",
		"d1":       "(mm)",
		"tw__1":    "(mm)",
		"tf__1":    "(mm)",
		"Ag":       "(mm²)",
		"Ix":       "(10³mm⁴)",
		"Zx":       "(10³mm³)",
		"Sx":       "(10³mm³)",
		"rx":       "(mm)",
		"Iy":       "(10³mm⁴)",
		"Zy":       "(mm³)",
		"Sy":       "(mm³)",
		"ry":       "(mm)",
		"J":        "(mm⁴)",
		"Iw":       "(mm⁶)",
		"flange":   "(mm)",
		"web":      "(mm)",
		"kf":       "",
		"C,N,S":    "",
		"Zex":      "(mm³)",
		"C,N,S__1": "",
		"Zey":      "(mm³)",
		"2tf":      "",
		"Zy5":      "(mm³)",
		"TanAlpha": "",
		"αb":       "",
		"Fu":       "(MPa)",
		"r2":       "(mm)",
		"ZeyD":     "(mm³)",
		"In":       "(10³mm⁴)",
		"Ip":       "(10³mm⁴)",
		"ZexC":     "(mm³)",
		"x5":       "(mm)",
		"y5":       "(mm)",
		"nL":       "(mm)",
		"pB":       "(mm)",
		"pT":       "(mm)",
		"Residual": "",
		"Type":     "",
		"ZeyL":     "(mm³)",
		"ZyL":      "(mm³)",
		"ZeyR":     "(mm³)",
		"ZyR":      "(mm³)",
		"xL":       "(mm)",
		"Xo":       "(mm)",
	}

	if unit, exists := unitMap[columnName]; exists && unit != "" {
		return columnName + unit
	}
	return columnName
}

func drawDataRows(properties []SteelProperty, currentColumns []ColumnInfo) {
	termWidth := getTerminalWidth()

	for i, prop := range properties {
		// Alternate row colors for better readability
		if i%2 == 0 {
			fmt.Printf("%s", ColorBg) // Default background
		} else {
			fmt.Printf("%s", ColorBgLight) // Slightly lighter background
		}

		// Section column (always shown, highlighted) - wider for better readability
		cleanedSection := cleanSectionName(prop.Section)
		fmt.Printf("%s%-35s%s", ColorTextBright, truncateString(cleanedSection, 34), ColorText)

		// Data columns - wider to match headers
		for _, col := range currentColumns {
			value := col.Formatter(prop)

			// Color coding for different value types
			if value == "-" || value == "" {
				fmt.Printf("%s%-18s", ColorTextDim, truncateString(value, 17))
			} else {
				fmt.Printf("%s%-18s", ColorText, truncateString(value, 17))
			}
		}

		// Fill remaining space to terminal width with background color
		usedSpace := 35 + (len(currentColumns) * 18)
		remainingSpace := termWidth - usedSpace
		if remainingSpace > 0 {
			fmt.Printf("%s", strings.Repeat(" ", remainingSpace))
		}
		fmt.Printf("%s\n", ColorReset)
	}
}

func drawNavigationFooter(currentPage, totalPages int) {
	termWidth := getTerminalWidth()

	fmt.Printf("\n%s", ColorBg)

	// Calculate content for the navigation box
	pageInfo := fmt.Sprintf("Page %d/%d", currentPage+1, totalPages)

	// Calculate box width with padding - make it a reasonable size
	contentLength := len("NAVIGATION: < or ← prev │ > or → next │ m Main Menu │ q quit ") + len(pageInfo)
	boxWidth := contentLength + 6 // Add padding
	if boxWidth < 84 {
		boxWidth = 84 // Minimum width similar to title box
	}

	// Calculate centering offset
	centerOffset := (termWidth - boxWidth - 2) / 2 // -2 for borders
	if centerOffset < 0 {
		centerOffset = 0
	}

	// Fill background before the navigation box
	fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
	fmt.Printf("%s╔%s╗", ColorBorder, strings.Repeat("═", boxWidth))
	// Fill background after the navigation box
	remainingSpace := termWidth - centerOffset - boxWidth - 2
	if remainingSpace > 0 {
		fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", remainingSpace))
	}
	fmt.Printf("%s\n", ColorReset)

	// Content row with centering
	fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
	fmt.Printf("%s║%s ", ColorBorder, ColorBg)

	// Navigation instructions with color coding
	fmt.Printf("%sNAVIGATION:%s ", ColorAccent, ColorText)

	if currentPage > 0 {
		fmt.Printf("%s<%s or %s←%s %sprev%s │ ", ColorSuccess, ColorText, ColorSuccess, ColorText, ColorTextDim, ColorText)
	} else {
		fmt.Printf("%s<%s or %s←%s %sprev%s │ ", ColorTextDim, ColorText, ColorTextDim, ColorText, ColorTextDim, ColorText)
	}

	if currentPage < totalPages-1 {
		fmt.Printf("%s>%s or %s→%s %snext%s │ ", ColorSuccess, ColorText, ColorSuccess, ColorText, ColorTextDim, ColorText)
	} else {
		fmt.Printf("%s>%s or %s→%s %snext%s │ ", ColorTextDim, ColorText, ColorTextDim, ColorText, ColorTextDim, ColorText)
	}

	fmt.Printf("%sm%s %sMain Menu%s │ ", ColorAccent, ColorText, ColorTextDim, ColorText)
	fmt.Printf("%sq%s %squit%s %s%s%s", ColorError, ColorText, ColorTextDim, ColorText, ColorAccent, pageInfo, ColorText)

	// Calculate padding to fill the box
	usedContentSpace := contentLength + 1 // +1 for space after "NAVIGATION:"
	contentPadding := boxWidth - usedContentSpace
	if contentPadding > 0 {
		fmt.Printf("%s", strings.Repeat(" ", contentPadding))
	}

	fmt.Printf("%s║", ColorBorder)
	if remainingSpace > 0 {
		fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", remainingSpace))
	}
	fmt.Printf("%s\n", ColorReset)

	// Bottom border with centering
	fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
	fmt.Printf("%s╚%s╝", ColorBorder, strings.Repeat("═", boxWidth))
	if remainingSpace > 0 {
		fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", remainingSpace))
	}
	fmt.Printf("%s\n", ColorReset)

	// Fill final line to terminal width
	fmt.Printf("%s%s%s\n", ColorBg, strings.Repeat(" ", termWidth), ColorReset)
}

// getTerminalWidth returns the terminal width, defaulting to 120 if unable to detect
func getTerminalWidth() int {
	cmd := exec.Command("tput", "cols")
	output, err := cmd.Output()
	if err != nil {
		return 120 // Default fallback
	}

	width, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil || width < 80 {
		return 120 // Default fallback
	}

	// Ensure we don't exceed reasonable bounds
	if width > 200 {
		width = 200
	}

	return width
}

func printWelcomeScreenInteractive() string {
	for {
		termWidth := getTerminalWidth()
		fmt.Print(ColorClear)

		// Fill entire screen with background color first
		for i := 0; i < 30; i++ {
			fmt.Printf("%s%s%s\n", ColorBg, strings.Repeat(" ", termWidth), ColorReset)
		}

		// Move cursor back to top
		fmt.Print("\033[H")

		// Title with accent color - properly centered
		titleBoxWidth := 80
		centerOffset := (termWidth - titleBoxWidth) / 2
		if centerOffset < 0 {
			centerOffset = 0
		}

		// Fill background before the title box
		fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
		fmt.Printf("%s╔══════════════════════════════════════════════════════════════════════════════╗", ColorBorderBright)
		// Fill background after the title box
		remainingSpace := termWidth - centerOffset - titleBoxWidth
		if remainingSpace > 0 {
			fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", remainingSpace))
		}
		fmt.Printf("%s\n", ColorReset)

		fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
		fmt.Printf("%s║%s                              %sSTEEL TABLES VIEWER%s                                %s║",
			ColorBorderBright, ColorBg, ColorAccent, ColorBg, ColorBorderBright)
		if remainingSpace > 0 {
			fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", remainingSpace))
		}
		fmt.Printf("%s\n", ColorReset)

		fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", centerOffset))
		fmt.Printf("%s╚══════════════════════════════════════════════════════════════════════════════╝", ColorBorderBright)
		if remainingSpace > 0 {
			fmt.Printf("%s%s", ColorBg, strings.Repeat(" ", remainingSpace))
		}
		fmt.Printf("%s\n\n", ColorReset)

		// Available files section with full width background
		fmt.Printf("%s%s▶ AVAILABLE STEEL TABLES:%s%s\n", ColorBg, ColorAccent,
			strings.Repeat(" ", termWidth-len("▶ AVAILABLE STEEL TABLES:")), ColorReset)
		listJSONFilesStyledFullWidth()

		// Instructions with full width background
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

		// Input prompt with full width background
		promptLine := "▶ SELECT TABLE: "
		fmt.Printf("%s%s%s", ColorBg, ColorAccent, promptLine)
		fmt.Printf("%s%s", ColorBg, ColorTextBright)

		// Read user input
		var input string
		fmt.Scanln(&input)

		// Reset colors and clear the input line
		fmt.Print(ColorReset)

		// Handle input
		input = strings.TrimSpace(strings.ToUpper(input))

		if input == "Q" || input == "QUIT" || input == "EXIT" {
			return ""
		}

		if input == "" {
			continue // Empty input, show menu again
		}

		// Validate input against available files
		if isValidTable(input) {
			// Convert back to the proper filename format
			filename := input + "_PROPS.json"
			return filename
		} else {
			// Show error and continue loop with full width background
			errorLine := fmt.Sprintf("✗ Table '%s' not found. Please try again...", input)
			fmt.Printf("\n%s%s%s%s%s\n", ColorBg, ColorError, errorLine,
				strings.Repeat(" ", termWidth-len(errorLine)), ColorReset)

			continuePrompt := "Press Enter to continue..."
			fmt.Printf("%s%s%s%s%s", ColorBg, ColorTextDim, continuePrompt,
				strings.Repeat(" ", termWidth-len(continuePrompt)), ColorReset)
			fmt.Scanln() // Wait for user to press Enter
		}
	}
}

func isValidTable(tableName string) bool {
	files, err := os.ReadDir("data")
	if err != nil {
		return false
	}

	// Convert input to expected filename format (files are already uppercase)
	expectedFilename := strings.ToUpper(tableName) + "_PROPS.json"

	for _, file := range files {
		// Files are already in the correct case, just compare directly
		if file.Name() == expectedFilename {
			return true
		}
	}
	return false
}

// fillRemainingSpace fills any remaining vertical space with background color
func fillRemainingSpace() {
	termWidth := getTerminalWidth()
	// Add a few empty lines with full background color to fill remaining space
	for i := 0; i < 3; i++ {
		fmt.Printf("%s%s%s\n", ColorBg, strings.Repeat(" ", termWidth), ColorReset)
	}
}

// filterAvailableColumns removes columns that have no meaningful data
func filterAvailableColumns(allColumns []ColumnInfo, properties []SteelProperty) []ColumnInfo {
	var availableColumns []ColumnInfo

	for _, col := range allColumns {
		hasData := false

		// Check if any row has meaningful data for this column
		for _, prop := range properties {
			value := col.Formatter(prop)
			// Consider various "empty" representations
			if value != "-" && value != "" && value != "0" && value != "0.0" && value != "0.00" && value != "0.000" {
				hasData = true
				break
			}
		}

		// Only include columns that have at least some meaningful data
		if hasData {
			availableColumns = append(availableColumns, col)
		}
	}

	return availableColumns
}
