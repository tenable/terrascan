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
	armDNSPrefix         = "dnsPrefix"
	armAgentPoolProfiles = "agentPoolProfiles"
	armPoolName          = "name"
	armNodeCount         = "count"
	armVMSize            = "vmSize"
	armAddonProfiles     = "addonProfiles"
	armNetworkProfile    = "networkProfile"
	armNetworkPlugin     = "networkPlugin"
	armNetworkPolicy     = "networkPolicy"
)

const (
	tfDNSPrefix       = "dns_prefix"
	tfDefaultNodePool = "default_node_pool"
	tfNodeCount       = "node_count"
	tfVMSize          = "vm_size"
	tfAddonProfile    = "addon_profile"
	tfConfig          = "config"
	tfNetworkProfile  = "network_profile"
	tfNetworkPlugin   = "network_plugin"
	tfNetworkPolicy   = "network_policy"
)

// KubernetesClusterConfig returns config for azurerm_kubernetes_cluster.
func KubernetesClusterConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation:  fn.LookUpString(nil, params, r.Location),
		tfName:      fn.LookUpString(nil, params, r.Name),
		tfTags:      functions.PatchAWSTags(r.Tags),
		tfDNSPrefix: fn.LookUpString(vars, params, convert.ToString(r.Properties, armDNSPrefix)),
	}

	poolProfiles := convert.ToSlice(r.Properties, armAgentPoolProfiles)
	dnp := make([]map[string]interface{}, 0)
	for _, p := range poolProfiles {
		profile := p.(map[string]interface{})
		newPool := map[string]interface{}{
			tfName:      fn.LookUpString(vars, params, convert.ToString(profile, armPoolName)),
			tfVMSize:    fn.LookUpString(vars, params, convert.ToString(profile, armVMSize)),
			tfNodeCount: fn.LookUpFloat64(vars, params, convert.ToString(profile, armNodeCount)),
		}
		dnp = append(dnp, newPool)
	}
	cf[tfDefaultNodePool] = dnp

	addonProfiles := convert.ToMap(r.Properties, armAddonProfiles)
	aps := make(map[string]interface{})
	for key, value := range addonProfiles {
		addon := value.(map[string]interface{})
		profile := map[string]interface{}{
			tfEnabled: convert.ToBool(addon, "enabled"),
		}

		if cfg, ok := addon["config"]; ok {
			profile[tfConfig] = cfg.(map[string]interface{})
		}

		if key == "kubeDashboard" {
			aps["kube_dashboard"] = profile
		}
		aps[key] = profile
	}
	cf[tfAddonProfile] = aps

	netProfile := convert.ToMap(r.Properties, armNetworkProfile)
	cf[tfNetworkProfile] = map[string]string{
		tfNetworkPlugin: fn.LookUpString(vars, params, convert.ToString(netProfile, armNetworkPlugin)),
		tfNetworkPolicy: fn.LookUpString(vars, params, convert.ToString(netProfile, armNetworkPolicy)),
	}
	return cf
}
