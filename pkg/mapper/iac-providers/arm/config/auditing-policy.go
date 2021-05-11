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
	"github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/types"
)

const (
	arm_storageEndpoint            = "storageEndpoint"
	arm_storageAccountAccessKey    = "storageAccountAccessKey"
	arm_isStorageSecondaryKeyInUse = "isStorageSecondaryKeyInUse"
	arm_retentionDays              = "retentionDays"
)

const (
	tf_storageEndpoint                    = "storage_endpoint,omitempty"
	tf_storageAccountAccessKey            = "storage_account_access_key,omitempty"
	tf_storageAccountAccessKeyIsSecondary = "storage_account_access_key_is_secondary,omitempty"
	tf_retentionInDays                    = "retention_in_days,omitempty"
)

// AuditingPolicyConfig returns config for azurerm_mssql_database_extended_auditing_policy
func AuditingPolicyConfig(r types.Resource) map[string]interface{} {
	return map[string]interface{}{
		tf_storageEndpoint:                    convert.ToString(r.Properties, arm_storageEndpoint),
		tf_storageAccountAccessKey:            convert.ToString(r.Properties, arm_storageAccountAccessKey),
		tf_storageAccountAccessKeyIsSecondary: convert.ToBool(r.Properties, arm_isStorageSecondaryKeyInUse),
		tf_retentionInDays:                    convert.ToFloat64(r.Properties, arm_retentionDays),
	}
}
