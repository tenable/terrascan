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

// SubnetMapping is internal struct as per terrform schema
type SubnetMapping struct {
	AllocationID       string `json:"allocation_id"`
	IPv6Address        string `json:"ipv6_address"`
	PrivateIPv4Address string `json:"private_ipv4_Address"`
	SubnetID           string `json:"subnet_id"`
}

// ElasticLoadBalancingV2LoadBalancerConfig holds config for aws_lb as per terraform schema
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
	awsconfig := make([]AWSResourceConfig, 0)
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
				subnetMapping.AllocationID = *mapping.AllocationId
			}
			subnetMapping.SubnetID = mapping.SubnetId
			cf.SubnetMappings = append(cf.SubnetMappings, subnetMapping)
		}
	}

	if len(e.LoadBalancerAttributes) != 0 {
		for _, attrib := range e.LoadBalancerAttributes {
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
	var awsconfigElb AWSResourceConfig
	awsconfigElb.Resource = cf
	awsconfigElb.Metadata = e.AWSCloudFormationMetadata
	awsconfig = append(awsconfig, awsconfigElb)
	return awsconfig
}
