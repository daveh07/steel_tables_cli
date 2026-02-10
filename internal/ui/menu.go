package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PrintWelcomeScreen displays the main menu and returns the selected table filename.
// Returns empty string if user wants to quit.
func PrintWelcomeScreen() string {
	reader := bufio.NewReader(os.Stdin)

	for {
		termWidth := GetTerminalWidth()
		fmt.Print(Clear)

		// Title box
		titleText := "STEEL TABLES VIEWER"
		titleBoxWidth := len(titleText) + 28
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

		fmt.Printf("%s%s", Bg, strings.Repeat(" ", centerOffset))
		fmt.Printf("%s╔%s╗", BorderBright, strings.Repeat("═", titleBoxWidth))
		fmt.Printf("%s%s%s\n", Bg, strings.Repeat(" ", remainingSpace), Reset)

		textPadding := (titleBoxWidth - len(titleText)) / 2
		rightPadding := titleBoxWidth - len(titleText) - textPadding

		fmt.Printf("%s%s", Bg, strings.Repeat(" ", centerOffset))
		fmt.Printf("%s║%s%s%s%s%s%s%s║",
			BorderBright, Bg, strings.Repeat(" ", textPadding), Accent, titleText, Bg, strings.Repeat(" ", rightPadding), BorderBright)
		fmt.Printf("%s%s%s\n", Bg, strings.Repeat(" ", remainingSpace), Reset)

		fmt.Printf("%s%s", Bg, strings.Repeat(" ", centerOffset))
		fmt.Printf("%s╚%s╝", BorderBright, strings.Repeat("═", titleBoxWidth))
		fmt.Printf("%s%s%s\n\n", Bg, strings.Repeat(" ", remainingSpace), Reset)

		// Available tables
		printFullWidthLine("▶ AVAILABLE STEEL TABLES:", Accent, termWidth)
		listJSONFiles(termWidth)

		// Instructions
		fmt.Println()
		printFullWidthLine("▶ INSTRUCTIONS:", Accent, termWidth)
		printFullWidthLine("  • Type the table name (e.g., PFC300, RHS450, UB350)", Text, termWidth)
		printFullWidthLine("  • Type q or quit to exit", Text, termWidth)
		printFullWidthLine("  • Press Enter to confirm your selection", Text, termWidth)
		fmt.Println()

		// Prompt
		fmt.Printf("%s▶ SELECT TABLE: %s", Accent, Reset)

		input, err := reader.ReadString('\n')
		if err != nil {
			return ""
		}
		input = strings.TrimSpace(strings.ToUpper(input))

		if input == "Q" || input == "QUIT" || input == "EXIT" {
			return ""
		}
		if input == "" {
			continue
		}
		if isValidTable(input) {
			return input + "_PROPS.json"
		}

		fmt.Printf("\n%s✗ Table '%s' not found. Please try again...%s\n", Error, input, Reset)
		fmt.Printf("%sPress Enter to continue...%s", TextDim, Reset)
		reader.ReadString('\n')
	}
}

func printFullWidthLine(text, color string, termWidth int) {
	padding := termWidth - len(text)
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("%s%s%s%s%s\n", Bg, color, text, strings.Repeat(" ", padding), Reset)
}

func listJSONFiles(termWidth int) {
	files, err := os.ReadDir("data")
	if err != nil {
		fmt.Printf("%s%s✗ Error reading data directory: %v%s\n", Bg, Error, err, Reset)
		return
	}

	i := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			displayName := strings.TrimSuffix(file.Name(), "_PROPS.json")
			line := fmt.Sprintf("  ● %s", displayName)
			padding := termWidth - len(line)
			if padding < 0 {
				padding = 0
			}

			bulletColor := Accent
			textColor := TextBright
			if i%2 == 1 {
				bulletColor = Blue
				textColor = Text
			}
			fmt.Printf("%s%s  ● %s%s%s%s\n", Bg, bulletColor, textColor, displayName, strings.Repeat(" ", padding), Reset)
			i++
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
