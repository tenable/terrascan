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

package azurev1

import (
	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// LoadIacFile loads the specified ARM template file.
// Note that a single ARM template json file may contain multiple resource definitions.
func (a *ARM) LoadIacFile(absFilePath string) (allResourcesConfig output.AllResourceConfigs, err error) {
	return nil, nil
}
