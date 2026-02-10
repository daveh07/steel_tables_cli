// Package config provides application configuration and paths.
package config

import (
	"os"
	"path/filepath"
)

var dataDir string

func init() {
	// Try to find data directory relative to executable
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		candidate := filepath.Join(execDir, "data")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			dataDir = candidate
			return
		}
	}

	// Fallback: try current working directory
	if info, err := os.Stat("data"); err == nil && info.IsDir() {
		dataDir = "data"
		return
	}

	// Last resort: assume relative path
	dataDir = "data"
}

// DataDir returns the path to the data directory.
func DataDir() string {
	return dataDir
}

// DataFile returns the full path to a file in the data directory.
func DataFile(filename string) string {
	return filepath.Join(dataDir, filename)
}
