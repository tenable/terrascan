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
	armEncryptionSettingsCollection = "encryptionSettingsCollection"
	armCreationData                 = "creationData"
	armCreateOption                 = "createOption"
	armDiskSizeGB                   = "diskSizeGB"
	armSourceResourceID             = "sourceResourceId"
)

const (
	tfCreateOption       = "create_option"
	tfDiskSizeGB         = "disk_size_gb"
	tfSourceResourceID   = "source_resource_id"
	tfStorageAccountType = "storage_account_type"
	tfEncryptionSettings = "encryption_settings"
)

// ManagedDiskConfig returns config for azurerm_managed_disk.
func ManagedDiskConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation:           fn.LookUp(nil, params, r.Location).(string),
		tfName:               fn.LookUp(nil, params, r.Name).(string),
		tfTags:               r.Tags,
		tfStorageAccountType: r.SKU.Name,
		tfEncryptionSettings: convert.ToMap(r.Properties, armEncryptionSettingsCollection),
		tfDiskSizeGB:         fn.LookUp(vars, params, convert.ToString(r.Properties, armDiskSizeGB)).(float64),
	}

	data := convert.ToMap(r.Properties, armCreationData)
	cf[tfCreateOption] = convert.ToString(data, armCreateOption)
	cf[tfSourceResourceID] = convert.ToString(data, armSourceResourceID)
	return cf
}
