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

const armEnableSoftDelete = "enableSoftDelete"

const (
	tfTenantID          = "tenant_id"
	tfSoftDeleteEnabled = "soft_delete_enabled"
)

// KeyVaultConfig returns config for azurerm_key_vault
func KeyVaultConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		tfLocation:          fn.LookUpString(nil, params, r.Location),
		tfName:              fn.LookUpString(nil, params, r.Name),
		tfTags:              functions.PatchAWSTags(r.Tags),
		tfSoftDeleteEnabled: convert.ToBool(r.Properties, armEnableSoftDelete),
		tfTenantID:          convert.ToString(params, armTenantID),
	}
}
