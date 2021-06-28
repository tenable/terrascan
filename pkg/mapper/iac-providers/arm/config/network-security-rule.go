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
	armAccess               = "access"
	armDirection            = "direction"
	armProtocol             = "protocol"
	armSourceAddressPrefix  = "sourceAddressPrefix"
	armSourcePortRange      = "sourcePortRange"
	armDestinationPortRange = "destinationPortRange"
)

const (
	tfAccess               = "access"
	tfDirection            = "direction"
	tfProtocol             = "protocol"
	tfSourceAddressPrefix  = "source_address_prefix"
	tfSourcePortRange      = "source_port_range"
	tfDestinationPortRange = "destination_port_range"
)

// NetworkSecurityRuleConfig returns config for azurerm_network_security_rule
func NetworkSecurityRuleConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		tfLocation:             fn.LookUpString(nil, params, r.Location),
		tfName:                 fn.LookUpString(nil, params, r.Name),
		tfTags:                 r.Tags,
		tfAccess:               convert.ToString(r.Properties, armAccess),
		tfDirection:            convert.ToBool(r.Properties, armDirection),
		tfProtocol:             convert.ToString(r.Properties, armProtocol),
		tfSourceAddressPrefix:  convert.ToString(r.Properties, armSourceAddressPrefix),
		tfSourcePortRange:      convert.ToString(r.Properties, armSourcePortRange),
		tfDestinationPortRange: convert.ToString(r.Properties, armDestinationPortRange),
	}
}
