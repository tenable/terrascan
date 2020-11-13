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

package iacprovider

import (
	"reflect"
)

// map of supported IaC providers
var supportedIacProviders = make(map[supportedIacType]map[supportedIacVersion]reflect.Type)

// map of default IaC versions for each IaC provider type
var defaultIacVersions = make(map[supportedIacType]supportedIacVersion)

// RegisterIacProvider registers an IaC provider for terrascan
// if the Iac provider does not have a version, it can be kept empty
func RegisterIacProvider(iacType supportedIacType, iacVersion, defaultIacVersion supportedIacVersion, iacProvider reflect.Type) {

	if iacVersion == "" {
		iacVersion = defaultIacVersion
	}

	// version support
	supportedTerraformVersions := make(map[supportedIacVersion]reflect.Type)
	supportedTerraformVersions[iacVersion] = iacProvider

	// default version
	defaultIacVersions[iacType] = defaultIacVersion

	// type support
	supportedIacProviders[iacType] = supportedTerraformVersions
}
