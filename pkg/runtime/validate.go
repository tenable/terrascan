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

package runtime

import (
	"fmt"
	"os"

	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"

	IacProvider "github.com/accurics/terrascan/pkg/iac-providers"
)

var (
	errEmptyIacPath     = fmt.Errorf("empty iac path, either use '-f' or '-d' option")
	errIncorrectIacPath = fmt.Errorf("cannot accept both '-f' and '-d' options together")
	errDirNotExists     = fmt.Errorf("directory does not exist")
	errFileNotExists    = fmt.Errorf("file does not exist")
	errIacNotSupported  = fmt.Errorf("iac type or version not supported")
)

// ValidateInputs validates the inputs to the executor object
func (e *Executor) ValidateInputs() error {

	// terrascan can accept either a file or a directory
	if e.filePath == "" && e.dirPath == "" {
		zap.S().Errorf("no IaC path specified; use '-f' for file or '-d' for directory")
		return errEmptyIacPath
	}
	if e.filePath != "" && e.dirPath != "" {
		zap.S().Errorf("cannot accept both '-f %s' and '-d %s' options together", e.filePath, e.dirPath)
		return errIncorrectIacPath
	}

	if e.dirPath != "" {
		// if directory, check if directory exists
		absDirPath, err := utils.GetAbsPath(e.dirPath)
		if err != nil {
			return err
		}

		if _, err := os.Stat(absDirPath); err != nil {
			zap.S().Errorf("directory '%s' does not exist", absDirPath)
			return errDirNotExists
		}
		zap.S().Debugf("directory '%s' exists", absDirPath)
	} else {

		// if file path, check if file exists
		absFilePath, err := utils.GetAbsPath(e.filePath)
		if err != nil {
			return err
		}

		if _, err := os.Stat(absFilePath); err != nil {
			zap.S().Errorf("file '%s' does not exist", absFilePath)
			return errFileNotExists
		}
		zap.S().Debugf("file '%s' exists", absFilePath)
	}

	// check if Iac type is supported
	if !IacProvider.IsIacSupported(e.iacType, e.iacVersion) {
		zap.S().Errorf("iac type '%s', version '%s' not supported", e.iacType, e.iacVersion)
		return errIacNotSupported
	}
	zap.S().Debugf("iac type '%s', version '%s' is supported", e.iacType, e.iacVersion)

	// check if policy type is supported

	// successful
	zap.S().Debug("input validation successful")
	return nil
}
