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
		} else {
			return homeDir, nil
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

// FindAllDirectories Walks the file path and returns a list of all directories within
func FindAllDirectories(basePath string) ([]string, error) {
	dirList := make([]string, 0)
	err := filepath.Walk(basePath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if fileInfo != nil && fileInfo.IsDir() {
			dirList = append(dirList, filePath)
		}
		return err
	})
	return dirList, err
}

// FilterFileInfoBySuffix Given a list of files, returns a subset of files containing a suffix which matches the input filter
func FilterFileInfoBySuffix(allFileList *[]os.FileInfo, filter string) *[]string {
	fileList := make([]string, 0)

	for i := range *allFileList {
		if strings.HasSuffix((*allFileList)[i].Name(), filter) {
			fileList = append(fileList, (*allFileList)[i].Name())
		}
	}
	return &fileList
}
