package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	tableName := "WC400"
	expectedFilename := strings.ToUpper(tableName) + "_PROPS.json"
	fmt.Println("Looking for:", expectedFilename)

	files, err := os.ReadDir("data")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Files in data directory:")
	for _, file := range files {
		fmt.Printf("  %s\n", file.Name())
		if file.Name() == expectedFilename {
			fmt.Println("MATCH FOUND:", file.Name())
		}
	}
}
