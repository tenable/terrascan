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
	"github.com/accurics/terrascan/pkg/mapper/convert"
	fn "github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/types"
)

const (
	arm_encryptionSettingsCollection = "encryptionSettingsCollection"
	arm_creationData                 = "creationData"
	arm_createOption                 = "createOption"
	arm_diskSizeGB                   = "diskSizeGB"
	arm_sourceResourceID             = "sourceResourceId"
)

const (
	tf_createOption       = "create_option"
	tf_diskSizeGB         = "disk_size_gb"
	tf_sourceResourceID   = "source_resource_id"
	tf_storageAccountType = "storage_account_type"
	tf_encryptionSettings = "encryption_settings"
)

// ManagedDiskConfig returns config for azurerm_managed_disk.
func ManagedDiskConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tf_location:           fn.LookUp(nil, params, r.Location).(string),
		tf_name:               fn.LookUp(nil, params, r.Name).(string),
		tf_tags:               r.Tags,
		tf_storageAccountType: r.SKU.Name,
		tf_encryptionSettings: convert.ToMap(r.Properties, arm_encryptionSettingsCollection),
		tf_diskSizeGB:         fn.LookUp(vars, params, convert.ToString(r.Properties, arm_diskSizeGB)).(float64),
	}

	data := convert.ToMap(r.Properties, arm_creationData)
	cf[tf_createOption] = convert.ToString(data, arm_createOption)
	cf[tf_sourceResourceID] = convert.ToString(data, arm_sourceResourceID)
	return cf
}
