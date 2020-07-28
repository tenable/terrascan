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

	tfv12 "github.com/accurics/terrascan/pkg/iac-providers/terraform/v12"
)

// SupportedIacType data type for supported IaC provider
type supportedIacType string

// supported IaC providers
const (
	terraform supportedIacType = "terraform"
)

// supportedIacVersion data type for supported Iac provider
type supportedIacVersion string

// supported Iac versions
const (
	defaultVersion supportedIacVersion = "default"
	terraformV12   supportedIacVersion = "v12"
)

// map of supported IaC providers
var supportedIacProviders map[supportedIacType](map[supportedIacVersion]reflect.Type)

// initializes a map of supported IaC providers
func init() {
	supportedIacProviders = make(map[supportedIacType](map[supportedIacVersion]reflect.Type))

	// terraform support
	supportedTerraformVersions := make(map[supportedIacVersion]reflect.Type)
	supportedTerraformVersions[terraformV12] = reflect.TypeOf(tfv12.TfV12{})
	supportedIacProviders[terraform] = supportedTerraformVersions
}
