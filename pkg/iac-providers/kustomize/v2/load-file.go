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

package kustomizev2

import (
	"github.com/tenable/terrascan/pkg/iac-providers/kustomize/commons"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
)

// LoadIacFile is not supported for kustomize. Only loading directories that have kustomization.y(a)ml file are supported
func (k *KustomizeV2) LoadIacFile(absRootPath string, options map[string]interface{}) (allResourcesConfig output.AllResourceConfigs, err error) {
	return commons.LoadIacFile()
}
