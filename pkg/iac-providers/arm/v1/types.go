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

package armv1

import "github.com/hashicorp/go-multierror"

// ARMV1 struct implements the IacProvider interface
type ARMV1 struct {
	templateParameters map[string]interface{}

	errIacLoadDirs *multierror.Error
	// absRootDir is the root directory being scanned.
	// if a file scan was initiated, absRootDir should be empty.
	absRootDir string
}

const (
	// JSONExtension json
	JSONExtension = ".json"

	// ParametersFileExtension .parameters.json
	ParametersFileExtension = ".parameters.json"

	// MetadataFileExtension .metadata.json
	MetadataFileExtension = ".metadata.json"

	// UnknownExtension unknown
	UnknownExtension = "unknown"
)

// ARMFileExtensions returns the valid extensions for Azure ARM (json)
func ARMFileExtensions() []string {
	return []string{JSONExtension}
}
