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
	"strings"

	"github.com/tenable/terrascan/pkg/mapper/convert"
	fn "github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/types"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// PostgreSQLServerConfig returns config for azurerm_postgresql_server
func PostgreSQLServerConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation: fn.LookUpString(nil, params, r.Location),
		tfName:     fn.LookUpString(nil, params, r.Name),
		tfTags:     functions.PatchAWSTags(r.Tags),
		tfSkuName:  fn.LookUpString(vars, params, r.SKU.Name),
		tfVersion:  fn.LookUpString(vars, params, convert.ToString(r.Properties, armVersion)),
	}

	if profile := convert.ToMap(r.Properties, armStorageProfile); profile != nil {
		status := fn.LookUpString(vars, params, convert.ToString(profile, armGeoRedundantBackup))
		cf[tfGeoRedundantBackupEnabled] = strings.EqualFold(strings.ToUpper(status), armStatusEnabled)

		cf[tfBackupRetentionDays] = fn.LookUpFloat64(vars, params, convert.ToString(profile, armBackupRetentionDays))

		cf[tfStorageMB] = fn.LookUpFloat64(vars, params, convert.ToString(profile, armStorageMB))

		status = fn.LookUpString(vars, params, convert.ToString(profile, armSslEnforcement))
		cf[tfSslEnforcementEnabled] = strings.EqualFold(strings.ToUpper(status), armStatusEnabled)
	}
	return cf
}
