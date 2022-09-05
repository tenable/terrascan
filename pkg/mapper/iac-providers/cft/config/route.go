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
	"github.com/awslabs/goformation/v6/cloudformation/ec2"
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
func GetRouteConfig(e *ec2.Route) []AWSResourceConfig {
	cf := RouteConfig{
		CarrierGatewayID:            *e.CarrierGatewayId,
		DestinationCidrBlock:        *e.DestinationCidrBlock,
		DestinationIpv6CidrBlock:    *e.DestinationIpv6CidrBlock,
		EgressOnlyInternetGatewayID: *e.EgressOnlyInternetGatewayId,
		GatewayID:                   *e.GatewayId,
		InstanceID:                  *e.InstanceId,
		LocalGatewayID:              *e.LocalGatewayId,
		NatGatewayID:                *e.NatGatewayId,
		NetworkInterfaceID:          *e.NetworkInterfaceId,
		RouteTableID:                e.RouteTableId,
		TransitGatewayID:            *e.TransitGatewayId,
		VpcEndpointID:               *e.VpcEndpointId,
		VpcPeeringConnectionID:      *e.VpcPeeringConnectionId,
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: e.AWSCloudFormationMetadata,
	}}
}
