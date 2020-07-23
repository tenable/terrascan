package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

// GetAbsPath returns absolute path from passed file path resolving even ~ to user home dir and any other such symbols that are only
// shell expanded can also be handled here
func GetAbsPath(path string) (string, error) {

	// Only shell resolves `~` to home so handle it specially
	if strings.HasPrefix(path, "~") {
		homeDir := os.Getenv("HOME")
		if len(path) > 1 {
			path = filepath.Join(homeDir, path[1:])
		}
	}

	// get absolute file path
	path, err := filepath.Abs(path)
	if err != nil {
		zap.S().Errorf("unable to resolve %s to absolute path. error: '%s'", path, err)
		return path, fmt.Errorf("failed to resolve absolute path")
	}
	return path, nil
}
