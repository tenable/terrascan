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
	arm_targetResourceID                         = "targetResourceId"
	arm_storageID                                = "storageId"
	arm_enabled                                  = "enabled"
	arm_retentionPolicy                          = "retentionPolicy"
	arm_days                                     = "days"
	arm_flowAnalyticsConfiguration               = "flowAnalyticsConfiguration"
	arm_networkWatcherFlowAnalyticsConfiguration = "networkWatcherFlowAnalyticsConfiguration"
	arm_workspaceID                              = "workspaceId"
	arm_workspaceRegion                          = "workspaceRegion"
	arm_workspaceResourceID                      = "workspaceResourceId"
	arm_trafficAnalyticsInterval                 = "trafficAnalyticsInterval"
)

const (
	tf_networkSecurityGroupID = "network_security_group_id"
	tf_trafficAnalytics       = "traffic_analytics"
	tf_workspaceID            = "workspace_id,omitempty"
	tf_workspaceRegion        = "workspace_region,omitempty"
	tf_workspaceResourceID    = "workspace_resource_id,omitempty"
	tf_intervalInMinutes      = "interval_in_minutes,omitempty"
)

// NetworkWatcherFlowLogConfig returns config for azurerm_network_watcher_flow_log
func NetworkWatcherFlowLogConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tf_location:               fn.LookUp(nil, params, r.Location).(string),
		tf_name:                   fn.LookUp(nil, params, r.Name).(string),
		tf_tags:                   r.Tags,
		tf_networkSecurityGroupID: fn.LookUp(vars, params, convert.ToString(r.Properties, arm_targetResourceID)).(string),
		tf_storageAccountID:       fn.LookUp(vars, params, convert.ToString(r.Properties, arm_storageID)).(string),
		tf_enabled:                convert.ToBool(r.Properties, arm_enabled),
	}

	policy := convert.ToMap(r.Properties, arm_retentionPolicy)
	cf[tf_retentionPolicy] = map[string]interface{}{
		tf_enabled: convert.ToBool(policy, arm_enabled),
		tf_days:    fn.LookUp(vars, params, convert.ToString(policy, arm_days)).(float64),
	}

	flowConfig := convert.ToMap(r.Properties, arm_flowAnalyticsConfiguration)
	if flowConfig != nil {
		networkConfig := convert.ToMap(flowConfig, arm_networkWatcherFlowAnalyticsConfiguration)
		cf[tf_trafficAnalytics] = map[string]interface{}{
			tf_enabled:             convert.ToBool(networkConfig, arm_enabled),
			tf_workspaceID:         fn.LookUp(vars, params, arm_workspaceID).(string),
			tf_workspaceRegion:     fn.LookUp(vars, params, arm_workspaceRegion).(string),
			tf_workspaceResourceID: fn.LookUp(vars, params, arm_workspaceResourceID).(string),
			tf_intervalInMinutes:   convert.ToFloat64(networkConfig, arm_trafficAnalyticsInterval),
		}
	}
	return cf
}
