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

const armNetworkInterfaces = "networkInterfaces"
const tfNetworkInterfaceIDs = "network_interface_ids"

// VirtualMachineConfig returns config for azurerm_virtual_machine
func VirtualMachineConfig(r types.Resource, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation: fn.LookUpString(nil, params, r.Location),
		tfName:     fn.LookUpString(nil, params, r.Name),
		tfTags:     functions.PatchAWSTags(r.Tags),
	}

	profile := convert.ToMap(r.Properties, armNetworkProfile)
	if interfaces, ok := profile[armNetworkInterfaces].([]interface{}); ok {
		iFaceIDs := []string{}
		for _, fs := range interfaces {
			iFace := fs.(map[string]interface{})
			iFaceIDs = append(iFaceIDs, iFace["id"].(string))
		}
		cf[tfNetworkInterfaceIDs] = iFaceIDs
	}
	return cf
}
