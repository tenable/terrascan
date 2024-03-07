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

const (
	armVersion             = "version"
	armStorageProfile      = "storageProfile"
	armStorageMB           = "storageMB"
	armBackupRetentionDays = "backupRetentionDays"
	armGeoRedundantBackup  = "geoRedundantBackup"
	armSslEnforcement      = "sslEnforcement"
	armStatusEnabled       = "ENABLED"
)

const (
	tfSkuName                   = "sku_name"
	tfStorageMB                 = "storage_mb"
	tfVersion                   = "version"
	tfBackupRetentionDays       = "backup_retention_days"
	tfGeoRedundantBackupEnabled = "geo_redundant_backup_enabled"
	tfSslEnforcementEnabled     = "ssl_enforcement_enabled"
)

// MySQLServerConfig returns config for azurerm_mysql_server
func MySQLServerConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation:                   fn.LookUpString(nil, params, r.Location),
		tfName:                       fn.LookUpString(nil, params, r.Name),
		tfSkuName:                    fn.LookUpString(vars, params, r.SKU.Name),
		tfTags:                       functions.PatchAWSTags(r.Tags),
		tfVersion:                    fn.LookUpString(vars, params, convert.ToString(r.Properties, armVersion)),
		tfAdministratorLogin:         fn.LookUpString(vars, params, convert.ToString(r.Properties, armAdministratorLogin)),
		tfAdministratorLoginPassword: fn.LookUpString(vars, params, convert.ToString(r.Properties, armAdministratorLoginPassword)),
	}

	profile := convert.ToMap(r.Properties, armStorageProfile)
	cf[tfStorageMB] = fn.LookUpFloat64(vars, params, convert.ToString(profile, armStorageMB))

	cf[tfBackupRetentionDays] = convert.ToFloat64(profile, armBackupRetentionDays)

	status := strings.ToUpper(convert.ToString(profile, armGeoRedundantBackup))
	cf[tfGeoRedundantBackupEnabled] = strings.EqualFold(status, armStatusEnabled)

	status = strings.ToUpper(convert.ToString(r.Properties, armSslEnforcement))
	cf[tfSslEnforcementEnabled] = strings.EqualFold(status, armStatusEnabled)
	return cf
}
