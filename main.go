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
	"syscall"
	"unsafe"
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

	currentPage := 0
	maxCols := 11 // 11 additional columns + Section = 12 total

	for {
		// Clear screen and set background
		fmt.Print(ColorBg + ColorClear)

		// Calculate column range for current page
		startCol := currentPage * maxCols
		endCol := startCol + maxCols
		if endCol > len(allColumns) {
			endCol = len(allColumns)
		}

		currentColumns := allColumns[startCol:endCol]
		totalPages := (len(allColumns) + maxCols - 1) / maxCols

		// Beautiful header with accent colors
		drawHeader(filepath.Base(filePath), currentPage+1, totalPages, len(properties))

		// Column headers with styling
		drawColumnHeaders(currentColumns)

		// Data rows with alternating colors
		drawDataRows(properties, currentColumns)

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

			// Update file path and reset to first page
			filePath = newFilePath
			currentPage = 0

			// Set raw mode again for the new table
			oldState = setRawMode()
			continue
		case '>':
			if endCol < len(allColumns) {
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
						if endCol < len(allColumns) {
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
	return s[:maxLen-3] + "..."
}

func readKey() byte {
	b := make([]byte, 1)
	syscall.Syscall(syscall.SYS_READ, uintptr(0), uintptr(unsafe.Pointer(&b[0])), 1)
	return b[0]
}

func readKeyNonBlocking() (byte, bool) {
	b := make([]byte, 1)
	n, _, _ := syscall.Syscall(syscall.SYS_READ, uintptr(0), uintptr(unsafe.Pointer(&b[0])), 1)
	if n > 0 {
		return b[0], true
	}
	return 0, false
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
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), uintptr(0x5402), uintptr(unsafe.Pointer(oldState)))
}

func drawHeader(filename string, currentPage, totalPages, totalEntries int) {
	// Get actual terminal width and subtract margin for borders
	termWidth := getTerminalWidth() - 4
	if termWidth < 80 {
		termWidth = 116 // Safe minimum
	}

	// Top border with background
	fmt.Printf("%s%s╔%s╗%s\n", ColorBg, ColorBorderBright,
		strings.Repeat("═", termWidth), ColorReset)

	// Title row with full background
	title := fmt.Sprintf("STEEL PROPERTIES: %s", strings.ToUpper(filename))
	if len(title) > termWidth-4 {
		title = title[:termWidth-7] + "..."
	}
	padding := (termWidth - len(title)) / 2
	remainingPadding := termWidth - len(title) - padding

	fmt.Printf("%s%s║%s", ColorBg, ColorBorderBright, ColorBg)
	fmt.Printf("%s", strings.Repeat(" ", padding))
	fmt.Printf("%s%s%s", ColorAccent, title, ColorBg)
	fmt.Printf("%s", strings.Repeat(" ", remainingPadding))
	fmt.Printf("%s║%s\n", ColorBorderBright, ColorReset)

	// Info row with full background
	info := fmt.Sprintf("Page %d/%d │ %d entries │ 12 columns per page", currentPage, totalPages, totalEntries)
	if len(info) > termWidth-4 {
		info = fmt.Sprintf("Page %d/%d │ %d entries", currentPage, totalPages, totalEntries)
	}
	infoPadding := (termWidth - len(info)) / 2
	remainingInfoPadding := termWidth - len(info) - infoPadding

	fmt.Printf("%s%s║%s", ColorBg, ColorBorderBright, ColorBg)
	fmt.Printf("%s", strings.Repeat(" ", infoPadding))
	fmt.Printf("%s%s%s", ColorTextDim, info, ColorBg)
	fmt.Printf("%s", strings.Repeat(" ", remainingInfoPadding))
	fmt.Printf("%s║%s\n", ColorBorderBright, ColorReset)

	// Bottom border with background
	fmt.Printf("%s%s╚%s╝%s\n\n", ColorBg, ColorBorderBright,
		strings.Repeat("═", termWidth), ColorReset)
}

func drawColumnHeaders(currentColumns []ColumnInfo) {
	// Header background
	fmt.Printf("%s%s", ColorBgLight, ColorAccent)

	// Section header (fixed column)
	fmt.Printf("%-25s", "Section")

	// Dynamic columns
	for _, col := range currentColumns {
		fmt.Printf("%-10s", truncateString(col.Name, 9))
	}
	fmt.Printf("%s\n", ColorReset)

	// Header separator with accent color
	fmt.Printf("%s", ColorBorderBright)
	fmt.Print(strings.Repeat("─", 25))
	for range currentColumns {
		fmt.Print(strings.Repeat("─", 10))
	}
	fmt.Printf("%s\n", ColorReset)
}

func drawDataRows(properties []SteelProperty, currentColumns []ColumnInfo) {
	for i, prop := range properties {
		// Alternate row colors for better readability
		if i%2 == 0 {
			fmt.Printf("%s", ColorBg) // Default background
		} else {
			fmt.Printf("%s", ColorBgLight) // Slightly lighter background
		}

		// Section column (always shown, highlighted)
		cleanedSection := cleanSectionName(prop.Section)
		fmt.Printf("%s%-25s%s", ColorTextBright, truncateString(cleanedSection, 24), ColorText)

		// Data columns
		for _, col := range currentColumns {
			value := col.Formatter(prop)

			// Color coding for different value types
			if value == "-" || value == "" {
				fmt.Printf("%s%-10s", ColorTextDim, truncateString(value, 9))
			} else {
				fmt.Printf("%s%-10s", ColorText, truncateString(value, 9))
			}
		}
		fmt.Printf("%s\n", ColorReset)
	}
}

func drawNavigationFooter(currentPage, totalPages int) {
	fmt.Printf("\n%s", ColorBg)

	// Navigation bar with beautiful styling
	fmt.Printf("%s╔", ColorBorder)
	fmt.Print(strings.Repeat("═", 118))
	fmt.Printf("╗%s\n", ColorReset)

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
	fmt.Printf("%sq%s %squit%s", ColorError, ColorText, ColorTextDim, ColorText)

	// Right-align page indicator
	pageInfo := fmt.Sprintf("Page %d/%d", currentPage+1, totalPages)
	remainingSpace := 118 - 95 - len(pageInfo) // Adjusted calculation for single line
	if remainingSpace > 0 {
		fmt.Printf("%s%s%s%s", strings.Repeat(" ", remainingSpace), ColorAccent, pageInfo, ColorText)
	}

	fmt.Printf(" %s║%s\n", ColorBorder, ColorReset)
	fmt.Printf("%s╚", ColorBorder)
	fmt.Print(strings.Repeat("═", 118))
	fmt.Printf("╝%s\n", ColorReset)
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
		fmt.Print(ColorClear)

		// Title with accent color
		fmt.Printf("%s╔══════════════════════════════════════════════════════════════════════════════╗%s\n", ColorBorderBright, ColorReset)
		fmt.Printf("%s║%s                              %sSTEEL TABLES VIEWER%s                             %s║%s\n", ColorBorderBright, ColorBg, ColorAccent, ColorBg, ColorBorderBright, ColorReset)
		fmt.Printf("%s╚══════════════════════════════════════════════════════════════════════════════╝%s\n\n", ColorBorderBright, ColorReset)

		// Available files section
		fmt.Printf("%s%s▶ AVAILABLE STEEL TABLES:%s\n", ColorBg, ColorAccent, ColorReset)
		listJSONFilesStyled()

		// Instructions
		fmt.Printf("\n%s%s▶ INSTRUCTIONS:%s\n", ColorBg, ColorAccent, ColorReset)
		fmt.Printf("%s  • Type the table name (e.g., %sPFC300%s, %sRHS450%s, %sUB350%s)\n", ColorBg, ColorSuccess, ColorText, ColorSuccess, ColorText, ColorSuccess, ColorText)
		fmt.Printf("%s  • Type %sq%s or %squit%s to exit\n", ColorBg, ColorError, ColorText, ColorError, ColorText)
		fmt.Printf("%s  • Press %sEnter%s to confirm your selection\n\n", ColorBg, ColorAccent, ColorText)

		// Input prompt
		fmt.Printf("%s%s▶ SELECT TABLE:%s ", ColorBg, ColorAccent, ColorReset)
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
			// Show error and continue loop
			fmt.Printf("\n%s%s✗ Table '%s' not found. Please try again...%s\n", ColorBg, ColorError, input, ColorReset)
			fmt.Printf("%s%sPress Enter to continue...%s", ColorBg, ColorTextDim, ColorReset)
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
