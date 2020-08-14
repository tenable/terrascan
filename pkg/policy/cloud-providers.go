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

package policy

import (
	"github.com/accurics/terrascan/pkg/config"
)

// supportedCloudType data type for supported cloud types in terrascan
type supportedCloudType string

// supportedCloudProvider map of supported cloud provider and its default policy path
var supportedCloudProvider = make(map[supportedCloudType]string)

var (
	basePolicyPath = config.GetPolicyBasePath()
)

// RegisterCloudProvider registers a cloud provider with terrascan
func RegisterCloudProvider(cloudType supportedCloudType) {
	policyPath := basePolicyPath + "/" + string(cloudType)
	supportedCloudProvider[cloudType] = policyPath
}

// IsCloudProviderSupported returns whether a cloud provider is supported in terrascan
func IsCloudProviderSupported(cloudType string) bool {
	_, supported := supportedCloudProvider[supportedCloudType(cloudType)]
	return supported
}

// GetDefaultPolicyPath returns the path to default policies for a given cloud provider
func GetDefaultPolicyPath(cloudType string) string {
	return supportedCloudProvider[supportedCloudType(cloudType)]
}
