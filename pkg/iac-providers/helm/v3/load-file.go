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

package helmv3

import (
	"fmt"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"go.uber.org/zap"
)

var (
	errLoadIacFileNotSupported = fmt.Errorf("load iac file is not supported for helm")
)

// LoadIacFile is not supported for helm. Only loading chart directories are supported
func (h *HelmV3) LoadIacFile(absRootPath string) (allResourcesConfig output.AllResourceConfigs, err error) {
	zap.S().Errorf("load iac file is not supported for helm")
	return make(map[string][]output.ResourceConfig), errLoadIacFileNotSupported
}
