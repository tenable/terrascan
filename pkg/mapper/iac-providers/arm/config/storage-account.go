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
	fn "github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/types"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

const (
	armSupportsHTTPSTrafficOnly = "supportsHttpsTrafficOnly"
	armNetworkAcls              = "networkAcls"
	armDefaultAction            = "defaultAction"
	armBypass                   = "bypass"
)

const (
	tfAccountTier            = "account_tier"
	tfAccountReplicationType = "account_replication_type"
	tfEnableHTTPSTrafficOnly = "enable_https_traffic_only"
	tfNetworkRules           = "network_rules"
	tfDefaultAction          = "default_action"
	tfByPass                 = "bypass"
	tfIPRules                = "ip_rules,omitempty"
)

// StorageAccountConfig returns config for azurerm_storage_account
func StorageAccountConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation:               fn.LookUpString(nil, params, r.Location),
		tfName:                   fn.LookUpString(nil, params, r.Name),
		tfTags:                   functions.PatchAWSTags(r.Tags),
		tfAccountTier:            r.SKU.Tier,
		tfAccountReplicationType: r.SKU.Name,
		tfEnableHTTPSTrafficOnly: convert.ToBool(r.Properties, armSupportsHTTPSTrafficOnly),
	}

	if acls := convert.ToMap(r.Properties, armNetworkAcls); acls != nil {
		rules := convert.ToSlice(acls, armIPRules)
		ipRules := []string{}
		for _, rule := range rules {
			r := rule.(map[string]string)
			ipRules = append(ipRules, r[armValue])
		}
		cf[tfNetworkRules] = map[string]interface{}{
			tfDefaultAction: fn.LookUpString(vars, params, convert.ToString(acls, armDefaultAction)),
			tfByPass:        fn.LookUpString(vars, params, convert.ToString(acls, armBypass)),
			tfIPRules:       ipRules,
		}
	}
	return cf
}
