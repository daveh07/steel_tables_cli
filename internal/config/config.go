// Package config provides application configuration and paths.
package config

import (
	"os"
	"path/filepath"
)

var dataDir string

func init() {
	// Try to find data directory relative to the real executable path
	// (resolving symlinks so it works when installed via symlink too)
	execPath, err := os.Executable()
	if err == nil {
		realPath, err := filepath.EvalSymlinks(execPath)
		if err == nil {
			execPath = realPath
		}
		execDir := filepath.Dir(execPath)
		candidate := filepath.Join(execDir, "data")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			dataDir = candidate
			return
		}
		// Also check one level up (for cmd/steel_tables/ layout)
		candidate = filepath.Join(execDir, "..", "..", "data")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			dataDir, _ = filepath.Abs(candidate)
			return
		}
	}

	// Fallback: try current working directory
	if info, err := os.Stat("data"); err == nil && info.IsDir() {
		abs, err := filepath.Abs("data")
		if err == nil {
			dataDir = abs
		} else {
			dataDir = "data"
		}
		return
	}

	// Last resort
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
