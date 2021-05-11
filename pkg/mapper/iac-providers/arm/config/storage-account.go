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
	arm_supportsHTTPSTrafficOnly = "supportsHttpsTrafficOnly"
	arm_networkAcls              = "networkAcls"
	arm_defaultAction            = "defaultAction"
	arm_bypass                   = "bypass"
)

const (
	tf_accountTier            = "account_tier"
	tf_accountReplicationType = "account_replication_type"
	tf_enableHTTPSTrafficOnly = "enable_https_traffic_only"
	tf_networkRules           = "network_rules"
	tf_defaultAction          = "default_action"
	tf_byPass                 = "bypass"
	tf_ipRules                = "ip_rules,omitempty"
)

// StorageAccountConfig returns config for azurerm_storage_account
func StorageAccountConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tf_location:               fn.LookUp(nil, params, r.Location).(string),
		tf_name:                   fn.LookUp(nil, params, r.Name).(string),
		tf_tags:                   r.Tags,
		tf_accountTier:            r.SKU.Tier,
		tf_accountReplicationType: r.SKU.Name,
		tf_enableHTTPSTrafficOnly: convert.ToBool(r.Properties, arm_supportsHTTPSTrafficOnly),
	}

	if acls := convert.ToMap(r.Properties, arm_networkAcls); acls != nil {
		rules := convert.ToSlice(acls, arm_ipRules)
		ipRules := []string{}
		for _, rule := range rules {
			r := rule.(map[string]string)
			ipRules = append(ipRules, r[arm_value])
		}
		cf[tf_networkRules] = map[string]interface{}{
			tf_defaultAction: fn.LookUp(vars, params, convert.ToString(acls, arm_defaultAction)).(string),
			tf_byPass:        fn.LookUp(vars, params, convert.ToString(acls, arm_bypass)).(string),
			tf_ipRules:       ipRules,
		}
	}
	return cf
}
