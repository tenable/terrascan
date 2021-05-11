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
	arm_dnsPrefix         = "dnsPrefix"
	arm_agentPoolProfiles = "agentPoolProfiles"
	arm_poolName          = "name"
	arm_nodeCount         = "count"
	arm_vmSize            = "vmSize"
	arm_addonProfiles     = "addonProfiles"
	arm_networkProfile    = "networkProfile"
	arm_networkPlugin     = "networkPlugin"
	arm_networkPolicy     = "networkPolicy"
)

const (
	tf_dnsPrefix       = "dns_prefix"
	tf_defaultNodePool = "default_node_pool"
	tf_nodeCount       = "node_count"
	tf_vmSize          = "vm_size"
	tf_addonProfile    = "addon_profile"
	tf_config          = "config"
	tf_networkProfile  = "network_profile"
	tf_networkPlugin   = "network_plugin"
	tf_networkPolicy   = "network_policy"
)

// KubernetesClusterConfig returns config for azurerm_kubernetes_cluster.
func KubernetesClusterConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tf_location:  fn.LookUp(nil, params, r.Location).(string),
		tf_name:      fn.LookUp(nil, params, r.Name).(string),
		tf_tags:      r.Tags,
		tf_dnsPrefix: fn.LookUp(vars, params, convert.ToString(r.Properties, arm_dnsPrefix)).(string),
	}

	poolProfiles := convert.ToSlice(r.Properties, arm_agentPoolProfiles)
	dnp := make([]map[string]interface{}, 0)
	for _, p := range poolProfiles {
		profile := p.(map[string]interface{})
		newPool := map[string]interface{}{
			tf_name:      fn.LookUp(vars, params, convert.ToString(profile, arm_poolName)).(string),
			tf_vmSize:    fn.LookUp(vars, params, convert.ToString(profile, arm_vmSize)).(string),
			tf_nodeCount: fn.LookUp(vars, params, convert.ToString(profile, arm_nodeCount)).(float64),
		}
		dnp = append(dnp, newPool)
	}
	cf[tf_defaultNodePool] = dnp

	addonProfiles := convert.ToMap(r.Properties, arm_addonProfiles)
	aps := make(map[string]interface{})
	for key, value := range addonProfiles {
		addon := value.(map[string]interface{})
		profile := map[string]interface{}{
			tf_enabled: addon["enabled"].(bool),
		}

		if cfg, ok := addon["config"]; ok {
			profile[tf_config] = cfg.(map[string]interface{})
		}

		if key == "kubeDashboard" {
			aps["kube_dashboard"] = profile
		}
		aps[key] = profile
	}
	cf[tf_addonProfile] = aps

	netProfile := convert.ToMap(r.Properties, arm_networkProfile)
	cf[tf_networkProfile] = map[string]string{
		tf_networkPlugin: fn.LookUp(vars, params, convert.ToString(netProfile, arm_networkPlugin)).(string),
		tf_networkPolicy: fn.LookUp(vars, params, convert.ToString(netProfile, arm_networkPolicy)).(string),
	}
	return cf
}
