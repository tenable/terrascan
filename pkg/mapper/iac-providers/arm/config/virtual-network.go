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
	armSubnets              = "subnets"
	armProperties           = "properties"
	armAddressPrefix        = "addressPrefix"
	armNetworkSecurityGroup = "networkSecurityGroup"
)

const (
	tfSubnet        = "subnet"
	tfAddressPrefix = "address_prefix"
	tfSecurityGroup = "security_group,omitempty"
)

// VirtualNetworkConfig returns config for azurerm_virtual_network
func VirtualNetworkConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation: fn.LookUpString(nil, params, r.Location),
		tfName:     fn.LookUpString(nil, params, r.Name),
		tfTags:     functions.PatchAWSTags(r.Tags),
	}

	subs := convert.ToSlice(r.Properties, armSubnets)
	subnets := make([]map[string]string, 0)
	for _, ss := range subs {
		s := ss.(map[string]interface{})
		prop := convert.ToMap(s, armProperties)

		sub := map[string]string{
			tfName:          fn.LookUpString(vars, params, convert.ToString(s, tfName)),
			tfAddressPrefix: fn.LookUpString(vars, params, convert.ToString(prop, armAddressPrefix)),
		}

		if nsg := convert.ToMap(prop, armNetworkSecurityGroup); nsg != nil {
			if sg, ok := fn.LookUp(vars, params, nsg["id"].(string)).(string); ok {
				sub[tfSecurityGroup] = sg
			}
		}
		subnets = append(subnets, sub)
	}
	cf[tfSubnet] = subnets

	return cf
}
