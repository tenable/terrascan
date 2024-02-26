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

const armWafConfiguration = "webApplicationFirewallConfiguration"

const (
	tfWafConfiguration = "waf_configuration"
	tfEnabled          = "enabled"
)

// ApplicationGatewayConfig returns config for azurerm_application_gateway
func ApplicationGatewayConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfName:     fn.LookUpString(nil, params, r.Name),
		tfLocation: fn.LookUpString(nil, params, tfLocation),
		tfTags:     functions.PatchAWSTags(r.Tags),
	}

	w := convert.ToMap(r.Properties, armWafConfiguration)
	cf[tfWafConfiguration] = map[string]interface{}{
		tfEnabled: convert.ToBool(w, armEnabled),
	}
	return cf
}
