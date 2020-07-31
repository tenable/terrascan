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
		} else {
			return homeDir, nil
		}
	}

	// get absolute file path
	path, _ = filepath.Abs(path)
	return path, nil
}
