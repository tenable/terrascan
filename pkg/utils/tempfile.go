/*
    Copyright (C) 2022 Tenable, Inc.

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

	"go.uber.org/zap"
)

// CreateTempFile creates a file with provided contents in the temp directory
func CreateTempFile(content []byte, ext string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", fmt.Sprintf("terrascan-*.%s", ext))
	if err != nil {
		zap.S().Errorf("failed to create temp file: '%v'", err)
		return nil, err
	}

	zap.S().Debugf("created temp config file at '%s'", tempFile.Name())

	_, err = tempFile.Write(content)
	if err != nil {
		zap.S().Errorf("failed to write to temp file: '%v'", err)
		return nil, err
	}

	return tempFile, nil
}
