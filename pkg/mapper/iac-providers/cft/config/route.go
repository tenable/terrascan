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
	"github.com/awslabs/goformation/v7/cloudformation/ec2"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// RouteConfig holds config for aws_route
type RouteConfig struct {
	Config
	CarrierGatewayID            string `json:"carrier_gateway_id"`
	DestinationCidrBlock        string `json:"destination_cidr_block"`
	DestinationIpv6CidrBlock    string `json:"destination_ipv6_cidr_block"`
	EgressOnlyInternetGatewayID string `json:"egress_only_gateway_id"`
	GatewayID                   string `json:"gateway_id"`
	InstanceID                  string `json:"instance_id"`
	LocalGatewayID              string `json:"local_gateway_id"`
	NatGatewayID                string `json:"nat_gateway_id"`
	NetworkInterfaceID          string `json:"network_interface_id"`
	RouteTableID                string `json:"route_table_id"`
	TransitGatewayID            string `json:"transit_gateway_id"`
	VpcEndpointID               string `json:"vpc_endpoint_id"`
	VpcPeeringConnectionID      string `json:"vpc_peering_connection_id"`
}

// GetRouteConfig returns config for aws_route
// aws_route no pilicy
func GetRouteConfig(e *ec2.Route) []AWSResourceConfig {
	cf := RouteConfig{
		CarrierGatewayID:            functions.GetVal(e.CarrierGatewayId),
		DestinationCidrBlock:        functions.GetVal(e.DestinationCidrBlock),
		DestinationIpv6CidrBlock:    functions.GetVal(e.DestinationIpv6CidrBlock),
		EgressOnlyInternetGatewayID: functions.GetVal(e.EgressOnlyInternetGatewayId),
		GatewayID:                   functions.GetVal(e.GatewayId),
		InstanceID:                  functions.GetVal(e.InstanceId),
		LocalGatewayID:              functions.GetVal(e.LocalGatewayId),
		NatGatewayID:                functions.GetVal(e.NatGatewayId),
		NetworkInterfaceID:          functions.GetVal(e.NetworkInterfaceId),
		RouteTableID:                e.RouteTableId,
		TransitGatewayID:            functions.GetVal(e.TransitGatewayId),
		VpcEndpointID:               functions.GetVal(e.VpcEndpointId),
		VpcPeeringConnectionID:      functions.GetVal(e.VpcPeeringConnectionId),
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: e.AWSCloudFormationMetadata,
	}}
}
