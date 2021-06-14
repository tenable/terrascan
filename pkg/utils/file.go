package utils

import (
	"go.uber.org/zap"
	"os"
)

// GetFileMode fetches the filemode from a file path
func GetFileMode(path string) os.FileMode {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			zap.S().Errorf("file %s does not exist.", path)
		} else {
			zap.S().Errorf("unable to fetch file info for path %s.", path)
		}
	}
	return fi.Mode()
}
