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
	"strings"

	IacProvider "github.com/accurics/terrascan/pkg/iac-providers"
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
)

var (
	errEmptyIacPath      = fmt.Errorf("empty iac path, either use '-f' or '-d' option")
	errDirNotExists      = fmt.Errorf("directory does not exist")
	errFileNotExists     = fmt.Errorf("file does not exist")
	errIacNotSupported   = fmt.Errorf("iac type or version not supported")
	errCloudNotSupported = fmt.Errorf("cloud type not supported")
)

// ValidateInputs validates the inputs to the executor object
func (e *Executor) ValidateInputs() error {

	var err error

	// terrascan can accept either a file or a directory
	if e.filePath == "" && e.dirPath == "" {
		zap.S().Errorf("no IaC path specified; use '-f' for file or '-d' for directory")
		return errEmptyIacPath
	}

	if e.filePath != "" {
		// if file path, check if file exists
		e.filePath, err = utils.GetAbsPath(e.filePath)
		if err != nil {
			return err
		}

		if _, err := os.Stat(e.filePath); err != nil {
			zap.S().Errorf("file '%s' does not exist", e.filePath)
			return errFileNotExists
		}
		zap.S().Debugf("file '%s' exists", e.filePath)
	} else {
		// if directory, check if directory exists
		e.dirPath, err = utils.GetAbsPath(e.dirPath)
		if err != nil {
			return err
		}

		if _, err := os.Stat(e.dirPath); err != nil {
			zap.S().Errorf("directory '%s' does not exist", e.dirPath)
			return errDirNotExists
		}
		zap.S().Debugf("directory '%s' exists", e.dirPath)
	}

	// set default iac type/version if not already set
	if e.iacType == "" {
		// TODO: handle more than cloudType[0]
		e.iacType = policy.GetDefaultIacType(e.cloudType[0])
	}

	if e.iacVersion == "" {
		e.iacVersion = IacProvider.GetDefaultIacVersion(e.iacType)
	}

	// check if IaC type is supported
	if !IacProvider.IsIacSupported(e.iacType, e.iacVersion) {
		zap.S().Errorf("iac type '%s', version '%s' not supported", e.iacType, e.iacVersion)
		return errIacNotSupported
	}
	zap.S().Debugf("iac type '%s', version '%s' is supported", e.iacType, e.iacVersion)

	// check if cloud type is supported
	for _, ct := range e.cloudType {
		if !policy.IsCloudProviderSupported(ct) {
			zap.S().Errorf("cloud type '%s' not supported", ct)
			return errCloudNotSupported
		}
	}
	zap.S().Debugf("cloud type '%s' is supported", strings.Join(e.cloudType, ","))
	if len(e.policyPath) == 0 {
		e.policyPath = policy.GetDefaultPolicyPaths(e.cloudType)
	}
	zap.S().Debugf("using policy path %v", e.policyPath)

	// successful
	zap.S().Debug("input validation successful")
	return nil
}
