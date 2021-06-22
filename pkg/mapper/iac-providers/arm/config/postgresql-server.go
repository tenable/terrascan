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

package config

import (
	"strings"

	"github.com/accurics/terrascan/pkg/mapper/convert"
	fn "github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/types"
)

// PostgreSQLServerConfig returns config for azurerm_postgresql_server
func PostgreSQLServerConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation: fn.LookUp(nil, params, r.Location).(string),
		tfName:     fn.LookUp(nil, params, r.Name).(string),
		tfTags:     r.Tags,
		tfSkuName:  fn.LookUp(vars, params, r.SKU.Name).(string),
		tfVersion:  fn.LookUp(vars, params, convert.ToString(r.Properties, armVersion)).(string),
	}

	if profile := convert.ToMap(r.Properties, armStorageProfile); profile != nil {
		status := fn.LookUp(vars, params, convert.ToString(profile, armGeoRedundantBackup))
		cf[tfGeoRedundantBackupEnabled] = strings.EqualFold(strings.ToUpper(status.(string)), armStatusEnabled)

		value := fn.LookUp(vars, params, convert.ToString(profile, armBackupRetentionDays))
		cf[tfBackupRetentionDays] = value.(float64)

		value = fn.LookUp(vars, params, convert.ToString(profile, armStorageMB))
		cf[tfStorageMB] = value.(float64)

		status = fn.LookUp(vars, params, convert.ToString(profile, armSslEnforcement))
		cf[tfSslEnforcementEnabled] = strings.EqualFold(strings.ToUpper(status.(string)), armStatusEnabled)
	}
	return cf
}
