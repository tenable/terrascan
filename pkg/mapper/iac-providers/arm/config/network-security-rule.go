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
	arm_access               = "access"
	arm_direction            = "direction"
	arm_protocol             = "protocol"
	arm_sourceAddressPrefix  = "sourceAddressPrefix"
	arm_sourcePortRange      = "sourcePortRange"
	arm_destinationPortRange = "destinationPortRange"
)

const (
	tf_access               = "access"
	tf_direction            = "direction"
	tf_protocol             = "protocol"
	tf_sourceAddressPrefix  = "source_address_prefix"
	tf_sourcePortRange      = "source_port_range"
	tf_destinationPortRange = "destination_port_range"
)

// NetworkSecurityRuleConfig returns config for azurerm_network_security_rule
func NetworkSecurityRuleConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		tf_location:             fn.LookUp(nil, params, r.Location).(string),
		tf_name:                 fn.LookUp(nil, params, r.Name).(string),
		tf_tags:                 r.Tags,
		tf_access:               convert.ToString(r.Properties, arm_access),
		tf_direction:            convert.ToBool(r.Properties, arm_direction),
		tf_protocol:             convert.ToString(r.Properties, arm_protocol),
		tf_sourceAddressPrefix:  convert.ToString(r.Properties, arm_sourceAddressPrefix),
		tf_sourcePortRange:      convert.ToString(r.Properties, arm_sourcePortRange),
		tf_destinationPortRange: convert.ToString(r.Properties, arm_destinationPortRange),
	}
}
