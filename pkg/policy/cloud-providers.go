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
	"sort"

	"github.com/accurics/terrascan/pkg/config"
)

// cloudProviderType data type for supported cloud types in terrascan
type cloudProviderType struct {
	isIndirect bool
	// policyPaths only populated if isIndirect == false
	policyPaths func() []string
	// policyNames only populated if isIndirect == true
	policyNames func() []string
}

// supportedCloudType data type for supported cloud types in terrascan
type supportedCloudType string

// supportedIacType data type for supported iac types
type supportedIacType string

// supportedIacVersion data type for supported iac versions
type supportedIacVersion string

// supportedCloudProvider map of supported cloud provider and its default policy path
var supportedCloudProvider = make(map[supportedCloudType]cloudProviderType)

// defaultIacType map of default IaC type for a given policy/cloud provider
var defaultIacType = make(map[supportedCloudType]supportedIacType)

// defaultIacVersion map of default IaC version for a given policy/cloud provider
var defaultIacVersion = make(map[supportedCloudType]supportedIacVersion)

var (
	basePolicyPath = config.GetPolicyBasePath()
)

func registerActualCloudProvider(cloudType supportedCloudType, iacTypeDefault supportedIacType, iacVersionDefault supportedIacVersion, isIndirect bool, getPolicyPaths func() []string) {
	if isIndirect {
		supportedCloudProvider[cloudType] = cloudProviderType{
			isIndirect:  true,
			policyNames: getPolicyPaths,
		}
	} else {
		supportedCloudProvider[cloudType] = cloudProviderType{
			isIndirect:  false,
			policyPaths: getPolicyPaths,
		}
	}

	defaultIacType[cloudType] = iacTypeDefault
	defaultIacVersion[cloudType] = iacVersionDefault
}

// RegisterIndirectCloudProvider registers a cloud provider with terrascan
func RegisterIndirectCloudProvider(cloudType supportedCloudType, iacTypeDefault supportedIacType, iacVersionDefault supportedIacVersion, getPolicyNames func() []string) {
	registerActualCloudProvider(cloudType, iacTypeDefault, iacVersionDefault, true, getPolicyNames)
}

// RegisterCloudProvider registers a cloud provider with terrascan
func RegisterCloudProvider(cloudType supportedCloudType, iacTypeDefault supportedIacType, iacVersionDefault supportedIacVersion) {
	registerActualCloudProvider(cloudType, iacTypeDefault, iacVersionDefault, false, func() []string { return []string{basePolicyPath + "/" + string(cloudType)} })
}

// IsCloudProviderSupported returns whether a cloud provider is supported in terrascan
func IsCloudProviderSupported(cloudType string) bool {
	_, supported := supportedCloudProvider[supportedCloudType(cloudType)]
	return supported
}

// GetDefaultPolicyPaths returns the paths to default policies for the given cloud providers
func GetDefaultPolicyPaths(cloudTypes []string) []string {
	var providers []string

	// Expand any indirect names
	var names []string
	for _, x := range cloudTypes {
		def := supportedCloudProvider[supportedCloudType(x)]
		if def.isIndirect {
			names = append(names, def.policyNames()...)
		} else {
			names = append(names, x)
		}
	}

	for _, x := range names {
		paths := (supportedCloudProvider[supportedCloudType(x)]).policyPaths()
		providers = append(providers, paths...)
	}
	return providers
}

// GetDefaultIacType returns the default IaC type for the given cloudType
// on the command line, the minimum arg required is the policy type (called cloudType here, so it's misleading)
// thus, for a given policy type, we need to specify a default IaC type
func GetDefaultIacType(cloudType string) string {
	return string(defaultIacType[supportedCloudType(cloudType)])
}

// SupportedPolicyTypes returns the list of policies supported in terrascan
func SupportedPolicyTypes(includeIndirect bool) []string {
	var policyTypes []string
	for k, v := range supportedCloudProvider {
		if !includeIndirect && v.isIndirect {
			continue
		}
		policyTypes = append(policyTypes, string(k))
	}
	sort.Strings(policyTypes)
	return policyTypes
}
