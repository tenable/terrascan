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
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/types"
)

const (
	armStorageEndpoint            = "storageEndpoint"
	armStorageAccountAccessKey    = "storageAccountAccessKey"
	armIsStorageSecondaryKeyInUse = "isStorageSecondaryKeyInUse"
	armRetentionDays              = "retentionDays"
)

const (
	tfStorageEndpoint                    = "storage_endpoint,omitempty"
	tfStorageAccountAccessKey            = "storage_account_access_key,omitempty"
	tfStorageAccountAccessKeyIsSecondary = "storage_account_access_key_is_secondary,omitempty"
	tfRetentionInDays                    = "retention_in_days,omitempty"
)

// AuditingPolicyConfig returns config for azurerm_mssql_database_extended_auditing_policy
func AuditingPolicyConfig(r types.Resource) map[string]interface{} {
	return map[string]interface{}{
		tfStorageEndpoint:                    convert.ToString(r.Properties, armStorageEndpoint),
		tfStorageAccountAccessKey:            convert.ToString(r.Properties, armStorageAccountAccessKey),
		tfStorageAccountAccessKeyIsSecondary: convert.ToBool(r.Properties, armIsStorageSecondaryKeyInUse),
		tfRetentionInDays:                    convert.ToFloat64(r.Properties, armRetentionDays),
	}
}
