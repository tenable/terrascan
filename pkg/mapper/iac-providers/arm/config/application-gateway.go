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

const armWafConfiguration = "webApplicationFirewallConfiguration"

const (
	tfWafConfiguration = "waf_configuration"
	tfEnabled          = "enabled"
)

// ApplicationGatewayConfig returns config for azurerm_application_gateway
func ApplicationGatewayConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfName:     fn.LookUp(nil, params, r.Name).(string),
		tfLocation: fn.LookUp(nil, params, tfLocation).(string),
		tfTags:     r.Tags,
	}

	w := convert.ToMap(r.Properties, armWafConfiguration)
	cf[tfWafConfiguration] = map[string]interface{}{
		tfEnabled: convert.ToBool(w, armEnabled),
	}
	return cf
}
