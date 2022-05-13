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

package kustomizev3

import (
	"github.com/tenable/terrascan/pkg/iac-providers/kustomize/commons"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
)

const (
	versionSuffix = "V3"
)

// LoadIacDir loads the kustomize directory and returns the ResourceConfig mapping which is evaluated by the policy engine
func (k *KustomizeV3) LoadIacDir(absRootDir string, options map[string]interface{}) (output.AllResourceConfigs, error) {
	return commons.NewKustomizeDirectoryLoader(absRootDir, options, true, versionSuffix).LoadIacDir()
}

// Name returns name of the provider
func (k *KustomizeV3) Name() string {
	return "kustomize"
}
