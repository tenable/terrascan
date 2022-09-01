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

package tfplan

import (
	"fmt"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
)

var (
	errIacDirNotSupport = fmt.Errorf("tfplan should always be a file, not a directory. Please specify path to tfplan file  with '-f' option")
)

// LoadIacDir is not supported for tfplan IacType. Terraform plan should always
// be a file and not a directory
func (k *TFPlan) LoadIacDir(absRootDir string, options map[string]interface{}) (output.AllResourceConfigs, error) {
	return output.AllResourceConfigs{}, errIacDirNotSupport
}

// Name returns name of the provider
func (*TFPlan) Name() string {
	return "tfplan"
}
