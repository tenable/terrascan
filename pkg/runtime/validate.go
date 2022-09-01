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

package runtime

import (
	"fmt"
	"os"
	"strings"

	IacProvider "github.com/tenable/terrascan/pkg/iac-providers"
	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"
)

var (
	errEmptyIacPath           = fmt.Errorf("empty iac path, either use '-f' or '-d' option")
	errDirNotExists           = fmt.Errorf("directory does not exist")
	errFileNotExists          = fmt.Errorf("file does not exist")
	errNotValidFile           = fmt.Errorf("not a valid file")
	errNotValidDir            = fmt.Errorf("not a valid directory")
	errIacNotSupported        = fmt.Errorf("iac type or version not supported")
	errCloudNotSupported      = fmt.Errorf("cloud type not supported")
	errSeverityNotSupported   = fmt.Errorf("severity level not supported")
	errCategoryNotSupported   = "category not supported : %v"
	errCategoriesNotSupported = "categories not supported : %v"
)

// ValidateInputs validates the inputs to the executor object
func (e *Executor) ValidateInputs() error {

	var err error
	var fi os.FileInfo
	var mode os.FileMode

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

		if fi, err = os.Stat(e.filePath); err != nil {
			zap.S().Errorf("file '%s' does not exist", e.filePath)
			return errFileNotExists
		}

		mode = fi.Mode()
		if !mode.IsRegular() {
			zap.S().Errorf("input path '%s' is not a valid file", e.filePath)
			return errNotValidFile
		}

		// the default value of dirPath is '.', make it empty
		e.dirPath = ""

		zap.S().Debugf("file '%s' exists", e.filePath)
	} else {
		// if directory, check if directory exists
		e.dirPath, err = utils.GetAbsPath(e.dirPath)
		if err != nil {
			return err
		}

		if fi, err = os.Stat(e.dirPath); err != nil {
			zap.S().Errorf("directory '%s' does not exist", e.dirPath)
			return errDirNotExists
		}

		mode = fi.Mode()
		if !mode.IsDir() {
			zap.S().Errorf("input path '%s' is not a valid directory", e.dirPath)
			return errNotValidDir
		}

		zap.S().Debugf("directory '%s' exists", e.dirPath)
	}

	// set default iac type/version if not already set
	// when iac type is not provided and -d option is used, all iac providers would be used for scanning
	if e.iacType == "" {
		if e.dirPath != "" {
			e.iacType = "all"
		} else {
			// TODO: handle more than policyTypes[0]
			e.iacType = policy.GetDefaultIacType(e.policyTypes[0])
		}
	}

	if e.iacVersion == "" {
		// ignore a version provided, when iacType resolves to 'all'
		if e.iacType != "all" {
			e.iacVersion = IacProvider.GetDefaultIacVersion(e.iacType)
		}
	}

	// check if cloud type(policy type) is supported
	for _, ct := range e.policyTypes {
		if !policy.IsCloudProviderSupported(ct) {
			zap.S().Errorf("cloud type '%s' not supported", ct)
			return errCloudNotSupported
		}
	}
	zap.S().Debugf("cloud type '%s' is supported", strings.Join(e.policyTypes, ","))

	if len(e.policyPath) == 0 {
		e.policyPath = policy.GetDefaultPolicyPaths(e.policyTypes)
	}

	zap.S().Debugf("using policy path %v", e.policyPath)

	// check if IaC type is supported
	// check is not required when it is 'all', iac providers with default version would be created later
	if e.iacType != "all" {
		if !IacProvider.IsIacSupported(e.iacType, e.iacVersion) {
			zap.S().Errorf("iac type '%s', version '%s' not supported", e.iacType, e.iacVersion)
			return errIacNotSupported
		}
		zap.S().Debugf("iac type '%s', version '%s' is supported", e.iacType, e.iacVersion)
	}

	if len(e.categories) > 0 {
		if isValid, invalidInputs := utils.ValidateCategoryInput(e.categories); !isValid {

			if len(invalidInputs) == 1 {
				return fmt.Errorf(errCategoryNotSupported, invalidInputs)
			}

			return fmt.Errorf(errCategoriesNotSupported, invalidInputs)
		}
	}
	zap.S().Debugf("using categories %v", e.categories)

	if len(e.severity) > 0 && !utils.ValidateSeverityInput(e.severity) {
		return errSeverityNotSupported
	}
	zap.S().Debugf("using severity level %v", e.severity)

	// successful
	zap.S().Debug("input validation successful")
	return nil
}
