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

package iacprovider

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"go.uber.org/zap"
)

var (
	errIacNotSupported = fmt.Errorf("iac not supported")
)

// NewIacProvider returns a new IacProvider
func NewIacProvider(iacType, iacVersion string) (iacProvider IacProvider, err error) {
	// get IacProvider from supportedIacProviders
	iacProviderObject, supported := supportedIacProviders[supportedIacType(iacType)][supportedIacVersion(iacVersion)]
	if !supported {
		zap.S().Errorf("IaC type:'%s', version: '%s' not supported", iacType, iacVersion)
		return iacProvider, errIacNotSupported
	}

	return reflect.New(iacProviderObject).Interface().(IacProvider), nil
}

// IsIacSupported returns true/false depending on whether the IaC
// provider is supported in terrascan or not
func IsIacSupported(iacType, iacVersion string) bool {
	if _, supported := supportedIacProviders[supportedIacType(iacType)]; !supported {
		return false
	}
	if _, supported := supportedIacProviders[supportedIacType(iacType)][supportedIacVersion(iacVersion)]; !supported {
		return false
	}
	return true
}

// SupportedIacProviders returns list of Iac Providers supported in terrascan
func SupportedIacProviders() []string {
	var iacTypes []string
	for k := range supportedIacProviders {
		iacTypes = append(iacTypes, string(k))
	}
	sort.Strings(iacTypes)
	return iacTypes
}

// SupportedIacVersions returns a string of Iac providers and corresponding supported versions
func SupportedIacVersions() []string {
	var iacVersions []string
	for iac, versions := range supportedIacProviders {
		var versionSlice []string
		for k := range versions {
			versionSlice = append(versionSlice, string(k))
		}
		sort.Strings(versionSlice)
		versionString := strings.Join(versionSlice, ", ")
		iacVersions = append(iacVersions, fmt.Sprintf("%s: %s", string(iac), versionString))
	}
	sort.Strings(iacVersions)
	return iacVersions
}

// GetProviderIacVersions returns list of Iac Provider versions for the given IaC type
func GetProviderIacVersions(iacType string) []string {
	var versions []string

	for version := range supportedIacProviders[supportedIacType(iacType)] {
		versions = append(versions, string(version))
	}
	sort.Strings(versions)
	return versions
}

// GetDefaultIacVersion returns the default IaC version for the given IaC type
func GetDefaultIacVersion(iacType string) string {
	return string(defaultIacVersions[supportedIacType(iacType)])
}
