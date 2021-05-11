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
	arm_storageAccountID = "storageAccountId"
	arm_category         = "category"
	arm_logs             = "logs"
)

const (
	tf_targetResourceID = "target_resource_id"
	tf_storageAccountID = "storage_account_id"
	tf_log              = "log"
	tf_category         = "category"
	tf_retentionPolicy  = "retention_policy"
	tf_days             = "days"
)

// DiagnosticSettingConfig returns config for azurerm_monitor_diagnostic_setting
func DiagnosticSettingConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tf_location:         fn.LookUp(nil, params, r.Location).(string),
		tf_name:             fn.LookUp(nil, params, r.Name).(string),
		tf_tags:             r.Tags,
		tf_targetResourceID: fn.LookUp(vars, params, getTargetResourceID(r.DependsOn)).(string),
		tf_storageAccountID: fn.LookUp(vars, params, convert.ToString(r.Properties, arm_storageAccountID)).(string),
	}

	logs := convert.ToSlice(r.Properties, arm_logs)
	if len(logs) > 0 {
		tfLog := make([]map[string]interface{}, 0)
		for _, lg := range logs {
			mp := lg.(map[string]interface{})
			policy := convert.ToMap(mp, arm_retentionPolicy)

			l := map[string]interface{}{
				tf_enabled:  convert.ToBool(mp, arm_enabled),
				tf_category: convert.ToString(mp, arm_category),
			}

			isEnabled := convert.ToBool(policy, arm_enabled)
			if isEnabled {
				l[tf_retentionPolicy] = map[string]interface{}{
					tf_enabled: isEnabled,
					tf_days:    fn.LookUp(vars, params, convert.ToString(policy, arm_days)).(float64),
				}
			} else {
				l[tf_retentionPolicy] = map[string]interface{}{
					tf_enabled: isEnabled,
				}
			}
		}
		cf[tf_log] = tfLog
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
