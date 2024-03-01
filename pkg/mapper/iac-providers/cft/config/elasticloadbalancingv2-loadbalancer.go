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
	"github.com/awslabs/goformation/v7/cloudformation/elasticloadbalancingv2"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// AccessLogs holds inned terraform structure for access_log attribute
type AccessLogs struct {
	Bucket  string
	Prefix  string
	Enabled bool
}

type SubnetMapping struct {
	AllocationId       string `json:"allocation_id"`
	IPv6Address        string `json:"ipv6_address"`
	PrivateIPv4Address string `json:"private_ipv4_Address"`
	SubnetId           string `json:"subnet_id"`
}

// ElasticLoadBalancingV2LoadBalancerConfig holds config for aws_lb
type ElasticLoadBalancingV2LoadBalancerConfig struct {
	SecurityGroups           []string        `json:"security_groups"`
	Subnets                  []string        `json:"subnets"`
	LoadBalancerType         string          `json:"load_balancer_type"`
	EnableDeletionProtection bool            `json:"enable_deletion_protection"`
	Internal                 bool            `json:"internal"`
	AccessLogs               AccessLogs      `json:"access_logs"`
	EnableCrossZone          bool            `json:"enable_cross_zone_load_balancing"`
	EnableHTTP               bool            `json:"enable_http2"`
	SubnetMappings           []SubnetMapping `json:"subnet_mapping"`
	Config
}

// GetElasticLoadBalancingV2LoadBalancerConfig returns config for aws_lb
// aws_lb
func GetElasticLoadBalancingV2LoadBalancerConfig(e *elasticloadbalancingv2.LoadBalancer, elbname string) []AWSResourceConfig {
	// elbpolicies := make([]ElasticLoadBalancingLoadBalancerPoliciesConfig, len(e.Policies))
	awsconfig := make([]AWSResourceConfig, 0)

	// for i, policy := range e.Policies {
	// 	indexedElbName := fmt.Sprintf("%s%d", elbname, i)

	// 	elbpolicies[i].LoadBalancerName = indexedElbName
	// 	elbpolicies[i].PolicyName = policy.PolicyName
	// 	elbpolicies[i].PolicyTypeName = policy.PolicyType

	// 	elbpolicies[i].PolicyAttribute = make([]PolicyAttributeBlock, len(policy.Attributes))
	// 	for ai := range policy.Attributes {
	// 		attribVals, ok := policy.Attributes[ai].(map[string]interface{})
	// 		if !ok {
	// 			continue
	// 		}

	// 		elbpolicies[i].PolicyAttribute[ai].Name, ok = attribVals["Name"].(string)
	// 		if !ok {
	// 			continue
	// 		}

	// 		elbpolicies[i].PolicyAttribute[ai].Value, ok = attribVals["Value"].(string)
	// 		if !ok {
	// 			continue
	// 		}

	// 		// variable "ok" is only used for safe type conversion
	// 		_ = ok
	// 	}

	// 	awsconfig[i].Type = GetPolicies
	// 	awsconfig[i].Name = indexedElbName
	// 	awsconfig[i].Resource = elbpolicies[i]
	// 	awsconfig[i].Metadata = e.AWSCloudFormationMetadata
	// }

	cf := ElasticLoadBalancingV2LoadBalancerConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(e.Tags),
			Name: *e.Name,
		},
		Internal: false,
	}

	cf.LoadBalancerType = "application"

	if len(e.SecurityGroups) != 0 {
		cf.SecurityGroups = e.SecurityGroups
	}
	if len(e.Subnets) != 0 {
		cf.Subnets = e.Subnets
	}
	cf.SubnetMappings = make([]SubnetMapping, 0)
	if len(e.SubnetMappings) != 0 {
		for _, mapping := range e.SubnetMappings {
			var subnetMapping SubnetMapping
			if mapping.IPv6Address != nil {
				subnetMapping.IPv6Address = *mapping.IPv6Address
			}
			if mapping.PrivateIPv4Address != nil {
				subnetMapping.PrivateIPv4Address = *mapping.PrivateIPv4Address
			}

			if mapping.AllocationId != nil {
				subnetMapping.AllocationId = *mapping.AllocationId
			}
			subnetMapping.SubnetId = mapping.SubnetId
			cf.SubnetMappings = append(cf.SubnetMappings, subnetMapping)
		}
	} else {

	}

	if len(e.LoadBalancerAttributes) != 0 {
		for _, attrib := range e.LoadBalancerAttributes {

			/*
						- Key: load_balancing.cross_zone.enabled
				          Value: true
				        - Key: deletion_protection.enabled
				          Value: false

				        - Key: ipv6.deny_all_igw_traffic
				          Value: false
				        - Key: routing.http2.enabled
				          Value: true

				        - Key: access_logs.s3.bucket
				          Value: bucketVal
				        - Key: access_logs.s3.prefix
				          Value: somePrefix
						- Key: access_logs.s3.enabled
				          Value: true
			*/
			if attrib.Key != nil && attrib.Value != nil {
				switch *attrib.Key {
				case "load_balancing.cross_zone.enabled":
					cf.EnableCrossZone = functions.GetBoolValueFromString(*attrib.Value)
				case "deletion_protection.enabled":
					cf.EnableDeletionProtection = functions.GetBoolValueFromString(*attrib.Value)
				case "routing.http2.enabled":
					cf.EnableHTTP = functions.GetBoolValueFromString(*attrib.Value)
				case "access_logs.s3.enabled":
					cf.AccessLogs.Enabled = functions.GetBoolValueFromString(*attrib.Value)
				case "access_logs.s3.bucket":
					cf.AccessLogs.Bucket = *attrib.Value
				case "access_logs.s3.prefix":
					cf.AccessLogs.Prefix = *attrib.Value
				default:
				}
			}

		}
	}
	// if e.CrossZone != nil {
	// 	cf.CrossZoneLoadBalancing = *e.CrossZone
	// } else {
	// 	cf.CrossZoneLoadBalancing = false
	// }
	// if e.AccessLoggingPolicy != nil {
	// 	cf.AccessLoggingPolicy = ELBAccessLoggingPolicyConfig{
	// 		Enabled: e.AccessLoggingPolicy.Enabled,
	// 	}
	// }

	// if e.Listeners != nil {
	// 	lc := make([]ELBListenerConfig, 0)
	// 	for _, listener := range e.Listeners {
	// 		lc = append(lc, ELBListenerConfig{
	// 			InstanceProtocol: functions.GetVal(listener.InstanceProtocol),
	// 			LBProtocol:       listener.Protocol,
	// 		})
	// 	}
	// 	cf.Listeners = lc
	// }

	var awsconfigElb AWSResourceConfig
	awsconfigElb.Resource = cf
	awsconfigElb.Metadata = e.AWSCloudFormationMetadata
	awsconfig = append(awsconfig, awsconfigElb)

	return awsconfig
}
