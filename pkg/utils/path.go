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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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
	path, _ = filepath.Abs(path)
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
func FilterFileInfoBySuffix(allFileList *[]os.FileInfo, filter []string) []*string {
	fileList := make([]*string, 0)

	for i := range *allFileList {
		for j := range filter {
			if strings.HasSuffix((*allFileList)[i].Name(), filter[j]) {
				filename := (*allFileList)[i].Name()
				fileList = append(fileList, &filename)
			}
		}
	}
	return fileList
}

// FindFilesBySuffix finds all files within a given directory that have the specified suffixes
// Returns a map with keys as directories and values as a list of files
func FindFilesBySuffix(basePath string, suffixes []string) (map[string][]*string, error) {
	retMap := make(map[string][]*string)

	// Walk the file path and find all directories
	dirList, err := FindAllDirectories(basePath)
	if err != nil {
		zap.S().Error("error encountered traversing directories", zap.String("base path", basePath), zap.Error(err))
		return retMap, err
	}

	if len(dirList) == 0 {
		return retMap, fmt.Errorf("no directories found for path %s", basePath)
	}

	sort.Strings(dirList)
	for i := range dirList {
		// Find all files in the current dir
		var fileInfo []os.FileInfo
		fileInfo, err = ioutil.ReadDir(dirList[i])
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				zap.S().Debug("error while searching for files", zap.String("dir", dirList[i]), zap.Error(err))
			}
			continue
		}

		fileList := FilterFileInfoBySuffix(&fileInfo, suffixes)
		if len(fileList) > 0 {
			retMap[dirList[i]] = fileList
		}
	}

	return retMap, nil
}

// FindFilesBySuffixInDir finds all the immediate files within a given directory that have the specified suffixes
// IT DOES NOT LOOK INTO ANY SUBDIRECTORY. JUST A SINGLE LEVEL FILE SEARCH.
// Returns an array for string pointers as a list of files
func FindFilesBySuffixInDir(basePath string, suffixes []string) ([]*string, error) {
	fileInfos, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	return FilterFileInfoBySuffix(&fileInfos, suffixes), nil
}

// AddFileExtension returns full file name string after adding the extension to the filename
func AddFileExtension(file, ext string) string {
	return fmt.Sprintf("%v.%v", file, ext)
}
