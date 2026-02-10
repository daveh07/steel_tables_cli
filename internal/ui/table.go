package ui

import (
	"fmt"
	"strings"

	"steel_tables/internal/columns"
	"steel_tables/internal/models"
)

// DrawColumnHeaders draws the column header row with units.
func DrawColumnHeaders(currentColumns []columns.ColumnInfo) {
	termWidth := GetTerminalWidth()
	fmt.Printf("%s%s", BgLight, Accent)
	fmt.Printf("%-25s", "Section")
	for _, col := range currentColumns {
		headerText := columns.GetHeaderWithUnit(col.Name)
		fmt.Printf("%-18s", truncateString(headerText, 17))
	}
	usedSpace := 25 + (len(currentColumns) * 18)
	remainingSpace := termWidth - usedSpace
	if remainingSpace > 0 {
		fmt.Printf("%s", strings.Repeat(" ", remainingSpace))
	}
	fmt.Printf("%s\n", Reset)

	// Separator line
	fmt.Printf("%s%s", BgLight, BorderBright)
	fmt.Print(strings.Repeat("─", 25))
	for range currentColumns {
		fmt.Print(strings.Repeat("─", 18))
	}
	if remainingSpace > 0 {
		fmt.Print(strings.Repeat("─", remainingSpace))
	}
	fmt.Printf("%s\n", Reset)
}

// DrawDataRows draws property rows starting at index 0.
func DrawDataRows(properties []models.SteelProperty, currentColumns []columns.ColumnInfo) {
	DrawDataRowsOffset(properties, currentColumns, 0)
}

// DrawDataRowsOffset draws property rows with a base offset for alternating colors.
func DrawDataRowsOffset(properties []models.SteelProperty, currentColumns []columns.ColumnInfo, baseIndex int) {
	termWidth := GetTerminalWidth()
	for i, prop := range properties {
		globalIndex := baseIndex + i
		if globalIndex%2 == 0 {
			fmt.Printf("%s", Bg)
		} else {
			fmt.Printf("%s", BgLight)
		}

		cleanedSection := cleanSectionName(prop.Section)
		fmt.Printf("%s%-25s%s", TextBright, truncateString(cleanedSection, 24), Text)

		for _, col := range currentColumns {
			value := col.Formatter(prop)
			if value == "-" || value == "" {
				fmt.Printf("%s%-18s", TextDim, truncateString(value, 17))
			} else {
				fmt.Printf("%s%-18s", Text, truncateString(value, 17))
			}
		}

		usedSpace := 25 + (len(currentColumns) * 18)
		remainingSpace := termWidth - usedSpace
		if remainingSpace > 0 {
			fmt.Printf("%s", strings.Repeat(" ", remainingSpace))
		}
		fmt.Printf("%s\n", Reset)
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
