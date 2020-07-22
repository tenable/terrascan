package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
		errMsg := fmt.Sprintf("unable to resolve %s to absolute path. error: '%s'", path, err)
		log.Println(errMsg)
		return path, fmt.Errorf(errMsg)
	}
	return path, nil
}
