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
	arm_subnets              = "subnets"
	arm_properties           = "properties"
	arm_addressPrefix        = "addressPrefix"
	arm_networkSecurityGroup = "networkSecurityGroup"
)

const (
	tf_subnet        = "subnet"
	tf_addressPrefix = "address_prefix"
	tf_securityGroup = "security_group,omitempty"
)

// VirtualNetworkConfig returns config for azurerm_virtual_network
func VirtualNetworkConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tf_location: fn.LookUp(nil, params, r.Location).(string),
		tf_name:     fn.LookUp(nil, params, r.Name).(string),
		tf_tags:     r.Tags,
	}

	subs := convert.ToSlice(r.Properties, arm_subnets)
	subnets := make([]map[string]string, 0)
	for _, ss := range subs {
		s := ss.(map[string]interface{})
		prop := convert.ToMap(s, arm_properties)

		sub := map[string]string{
			tf_name:          fn.LookUp(vars, params, s[tf_name].(string)).(string),
			tf_addressPrefix: fn.LookUp(vars, params, prop[arm_addressPrefix].(string)).(string),
		}

		if nsg := convert.ToMap(prop, arm_networkSecurityGroup); nsg != nil {
			if sg, ok := fn.LookUp(vars, params, nsg["id"].(string)).(string); ok {
				sub[tf_securityGroup] = sg
			}
		}
		subnets = append(subnets, sub)
	}
	cf[tf_subnet] = subnets

	return cf
}
