// Steel Tables Viewer - A terminal-based viewer for structural steel section properties.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"steel_tables/internal/config"
	"steel_tables/internal/ui"
	"steel_tables/internal/viewer"
)

func main() {
	fmt.Print(ui.Bg + ui.Clear)
	defer fmt.Print(ui.Reset)

	if len(os.Args) < 2 {
		runInteractiveMode()
	} else {
		runCLIMode(os.Args[1])
	}

	fmt.Print(ui.Clear)
}

func runInteractiveMode() {
	initialState, err := ui.GetTerminalState()
	if err != nil {
		log.Fatalf("Fatal: Could not get terminal state: %v", err)
	}
	defer ui.RestoreTerminal(initialState)

	for {
		ui.RestoreTerminal(initialState)

		selectedFile := ui.PrintWelcomeScreen()
		if selectedFile == "" {
			return
		}

		filePath := config.DataFile(selectedFile)

		if err := ui.SetRawMode(); err != nil {
			log.Printf("Error entering raw mode: %v. Returning to menu.", err)
			continue
		}

		returnToMenu := viewer.DisplayTable(filePath)
		ui.RestoreTerminal(initialState)

		if !returnToMenu {
			break
		}
	}
}

func runCLIMode(tableName string) {
	filename := strings.ToUpper(tableName)
	if !strings.HasSuffix(filename, "_PROPS") {
		filename += "_PROPS"
	}
	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}
	filePath := config.DataFile(filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("Table '%s' not found.", tableName)
	}

	viewer.PrintTableOnce(filePath)
}
