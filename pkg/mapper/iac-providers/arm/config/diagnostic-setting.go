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
	armStorageAccountID = "storageAccountId"
	armCategory         = "category"
	armLogs             = "logs"
)

const (
	tfTargetResourceID = "target_resource_id"
	tfStorageAccountID = "storage_account_id"
	tfLog              = "log"
	tfCategory         = "category"
	tfRetentionPolicy  = "retention_policy"
	tfDays             = "days"
)

// DiagnosticSettingConfig returns config for azurerm_monitor_diagnostic_setting
func DiagnosticSettingConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation:         fn.LookUp(nil, params, r.Location).(string),
		tfName:             fn.LookUp(nil, params, r.Name).(string),
		tfTags:             r.Tags,
		tfTargetResourceID: fn.LookUp(vars, params, getTargetResourceID(r.DependsOn)).(string),
		tfStorageAccountID: fn.LookUp(vars, params, convert.ToString(r.Properties, armStorageAccountID)).(string),
	}

	logs := convert.ToSlice(r.Properties, armLogs)
	if len(logs) > 0 {
		tfLogMap := make([]map[string]interface{}, 0)
		for _, lg := range logs {
			mp := lg.(map[string]interface{})
			policy := convert.ToMap(mp, armRetentionPolicy)

			l := map[string]interface{}{
				tfEnabled:  convert.ToBool(mp, armEnabled),
				tfCategory: convert.ToString(mp, armCategory),
			}

			isEnabled := convert.ToBool(policy, armEnabled)
			if isEnabled {
				l[tfRetentionPolicy] = map[string]interface{}{
					tfEnabled: isEnabled,
					tfDays:    fn.LookUp(vars, params, convert.ToString(policy, armDays)).(float64),
				}
			} else {
				l[tfRetentionPolicy] = map[string]interface{}{
					tfEnabled: isEnabled,
				}
			}
		}
		cf[tfLog] = tfLogMap
	}
	return cf
}

func getTargetResourceID(deps []string) string {
	for _, d := range deps {
		if strings.Contains(d, "vault") {
			return d
		}
	}
	return ""
}
