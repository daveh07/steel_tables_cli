package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SteelProperty struct {
	Section string `json:"Section"`
}

type SectionResult struct {
	Section   string
	TableName string
	FilePath  string
}

func main() {
	query := "PFC"
	fmt.Printf("Testing filter with query: %s\n", query)

	results := filterTables(query)
	fmt.Printf("Found %d sections:\n", len(results))

	for i, result := range results {
		if i >= 10 { // Show only first 10 for brevity
			fmt.Printf("... and %d more\n", len(results)-10)
			break
		}
		fmt.Printf("  %s (from %s)\n", result.Section, result.TableName)
	}
}

func filterTables(query string) []SectionResult {
	files, err := os.ReadDir("data")
	if err != nil {
		return nil
	}

	filters := strings.Split(strings.ToUpper(query), "+")
	for i, filter := range filters {
		filters[i] = strings.TrimSpace(filter)
	}

	var matchingSections []SectionResult

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join("data", file.Name())
			tableName := strings.TrimSuffix(file.Name(), "_PROPS.json")

			data, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			var properties []SteelProperty
			if err := json.Unmarshal(data, &properties); err != nil {
				continue
			}

			for _, prop := range properties {
				sectionName := strings.ToUpper(prop.Section)

				for _, filter := range filters {
					if filter != "" && strings.Contains(sectionName, filter) {
						matchingSections = append(matchingSections, SectionResult{
							Section:   prop.Section,
							TableName: tableName,
							FilePath:  filePath,
						})
						break
					}
				}
			}
		}
	}

	return matchingSections
}
