/*
    Copyright (C) 2020 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package utils

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
)

// GetHomeDir returns the home directory path
func GetHomeDir() (terrascanDir string) {
	zap.S().Debug("looking up for the home directory path")

	terrascanDir, err := homedir.Dir()

	if err != nil {
		zap.S().Warnf("unable to determine the home directory: %v\n", err)
	}

	return
}

// GenerateTempDir generates a temporary directory
func GenerateTempDir() string {
	return filepath.Join(os.TempDir(), GenRandomString(6))
}

// IsDirExists checks wether the provided directory exists or not
func IsDirExists(dir string) bool {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		zap.S().Errorf("directory %s does not exist.", dir)
		return false
	}
	return true
}
