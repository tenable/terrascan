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

package tfv15

import (
	commons "github.com/accurics/terrascan/pkg/iac-providers/terraform/commons"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// LoadIacDir starts traversing from the given rootDir and traverses through
// all the descendant modules present to create an output list of all the
// resources present in rootDir and descendant modules
func (*TfV15) LoadIacDir(absRootDir string, nonRecursive, useTerraformCache bool) (allResourcesConfig output.AllResourceConfigs, err error) {
	return commons.NewTerraformDirectoryLoader(absRootDir, nonRecursive, useTerraformCache).LoadIacDir()
}
