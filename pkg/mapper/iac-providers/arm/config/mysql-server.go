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

const (
	arm_version             = "version"
	arm_storageProfile      = "storageProfile"
	arm_storageMB           = "storageMB"
	arm_backupRetentionDays = "backupRetentionDays"
	arm_geoRedundantBackup  = "geoRedundantBackup"
	arm_sslEnforcement      = "sslEnforcement"
	arm_statusEnabled       = "ENABLED"
)

const (
	tf_skuName                   = "sku_name"
	tf_storageMB                 = "storage_mb"
	tf_version                   = "version"
	tf_backupRetentionDays       = "backup_retention_days"
	tf_geoRedundantBackupEnabled = "geo_redundant_backup_enabled"
	tf_sslEnforcementEnabled     = "ssl_enforcement_enabled"
)

// MySQLServerConfig returns config for azurerm_mysql_server
func MySQLServerConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tf_location:                   fn.LookUp(nil, params, r.Location).(string),
		tf_name:                       fn.LookUp(nil, params, r.Name).(string),
		tf_skuName:                    fn.LookUp(vars, params, r.SKU.Name).(string),
		tf_tags:                       r.Tags,
		tf_version:                    fn.LookUp(vars, params, convert.ToString(r.Properties, arm_version)).(string),
		tf_administratorLogin:         fn.LookUp(vars, params, convert.ToString(r.Properties, arm_administratorLogin)).(string),
		tf_administratorLoginPassword: fn.LookUp(vars, params, convert.ToString(r.Properties, arm_administratorLoginPassword)).(string),
	}

	profile := convert.ToMap(r.Properties, arm_storageProfile)
	cf[tf_storageMB] = fn.LookUp(vars, params, convert.ToString(profile, arm_storageMB)).(float64)

	cf[tf_backupRetentionDays] = convert.ToFloat64(profile, arm_backupRetentionDays)

	status := strings.ToUpper(convert.ToString(profile, arm_geoRedundantBackup))
	cf[tf_geoRedundantBackupEnabled] = strings.EqualFold(status, arm_statusEnabled)

	status = strings.ToUpper(convert.ToString(r.Properties, arm_sslEnforcement))
	cf[tf_sslEnforcementEnabled] = strings.EqualFold(status, arm_statusEnabled)
	return cf
}
