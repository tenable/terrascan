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

package config

import (
	"github.com/tenable/terrascan/pkg/mapper/convert"
	fn "github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/types"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

const (
	armSource = "source"
	armValue  = "value"
)

const tfValue = "value"

// PostgreSQLConfigurationConfig returns config for azurerm_postgresql_configuration
func PostgreSQLConfigurationConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		tfLocation: fn.LookUpString(nil, params, r.Location),
		tfTags:     functions.PatchAWSTags(r.Tags),
		tfName:     convert.ToString(r.Properties, armSource),
		tfValue:    convert.ToString(r.Properties, armValue),
	}
}
