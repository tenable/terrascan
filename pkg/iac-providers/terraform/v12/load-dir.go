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

package tfv12

import (
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	commons "github.com/tenable/terrascan/pkg/iac-providers/terraform/commons"
	"go.uber.org/zap"
)

// LoadIacDir starts traversing from the given rootDir and traverses through
// all the descendant modules present to create an output list of all the
// resources present in rootDir and descendant modules
func (*TfV12) LoadIacDir(absRootDir string, options map[string]interface{}) (allResourcesConfig output.AllResourceConfigs, err error) {
	zap.S().Warn("There may be a few breaking changes while working with terraform v0.12 files. For further information, refer to https://github.com/tenable/terrascan/releases/v1.3.0")
	return commons.NewTerraformDirectoryLoader(absRootDir, version, options).LoadIacDir()
}

// Name returns name of the provider
func (*TfV12) Name() string {
	return "terraform"
}
