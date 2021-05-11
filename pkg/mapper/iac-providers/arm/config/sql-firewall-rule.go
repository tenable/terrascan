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
	arm_startIPAddress = "startIpAddress"
	arm_endIPAddress   = "endIpAddress"
)

const (
	tf_startIPAddress = "start_ip_address"
	tf_endIPAddress   = "end_ip_address"
)

// SQLFirewallRuleConfig returns config for azurerm_sql_firewall_rule
func SQLFirewallRuleConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		tf_location:       fn.LookUp(nil, params, r.Location).(string),
		tf_name:           fn.LookUp(nil, params, r.Name).(string),
		tf_tags:           r.Tags,
		tf_startIPAddress: fn.LookUp(nil, params, convert.ToString(r.Properties, arm_startIPAddress)).(string),
		tf_endIPAddress:   fn.LookUp(nil, params, convert.ToString(r.Properties, arm_endIPAddress)).(string),
	}
}
