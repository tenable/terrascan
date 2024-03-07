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
	armTargetResourceID                         = "targetResourceId"
	armStorageID                                = "storageId"
	armEnabled                                  = "enabled"
	armRetentionPolicy                          = "retentionPolicy"
	armDays                                     = "days"
	armFlowAnalyticsConfiguration               = "flowAnalyticsConfiguration"
	armNetworkWatcherFlowAnalyticsConfiguration = "networkWatcherFlowAnalyticsConfiguration"
	armWorkspaceID                              = "workspaceId"
	armWorkspaceRegion                          = "workspaceRegion"
	armWorkspaceResourceID                      = "workspaceResourceId"
	armTrafficAnalyticsInterval                 = "trafficAnalyticsInterval"
)

const (
	tfNetworkSecurityGroupID = "network_security_group_id"
	tfTrafficAnalytics       = "traffic_analytics"
	tfWorkspaceID            = "workspace_id,omitempty"
	tfWorkspaceRegion        = "workspace_region,omitempty"
	tfWorkspaceResourceID    = "workspace_resource_id,omitempty"
	tfIntervalInMinutes      = "interval_in_minutes,omitempty"
)

// NetworkWatcherFlowLogConfig returns config for azurerm_network_watcher_flow_log
func NetworkWatcherFlowLogConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation:               fn.LookUpString(nil, params, r.Location),
		tfName:                   fn.LookUpString(nil, params, r.Name),
		tfTags:                   functions.PatchAWSTags(r.Tags),
		tfNetworkSecurityGroupID: fn.LookUpString(vars, params, convert.ToString(r.Properties, armTargetResourceID)),
		tfStorageAccountID:       fn.LookUpString(vars, params, convert.ToString(r.Properties, armStorageID)),
		tfEnabled:                convert.ToBool(r.Properties, armEnabled),
	}

	policy := convert.ToMap(r.Properties, armRetentionPolicy)
	cf[tfRetentionPolicy] = map[string]interface{}{
		tfEnabled: convert.ToBool(policy, armEnabled),
		tfDays:    fn.LookUpFloat64(vars, params, convert.ToString(policy, armDays)),
	}

	flowConfig := convert.ToMap(r.Properties, armFlowAnalyticsConfiguration)
	if flowConfig != nil {
		networkConfig := convert.ToMap(flowConfig, armNetworkWatcherFlowAnalyticsConfiguration)
		cf[tfTrafficAnalytics] = map[string]interface{}{
			tfEnabled:             convert.ToBool(networkConfig, armEnabled),
			tfWorkspaceID:         fn.LookUpString(vars, params, armWorkspaceID),
			tfWorkspaceRegion:     fn.LookUpString(vars, params, armWorkspaceRegion),
			tfWorkspaceResourceID: fn.LookUpString(vars, params, armWorkspaceResourceID),
			tfIntervalInMinutes:   convert.ToFloat64(networkConfig, armTrafficAnalyticsInterval),
		}
	}
	return cf
}
