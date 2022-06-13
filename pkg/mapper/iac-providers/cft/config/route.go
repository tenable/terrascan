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
	"github.com/awslabs/goformation/v5/cloudformation/ec2"
)

// RouteConfig holds config for aws_route
type RouteConfig struct {
	Config
	CarrierGatewayId            string `json:"carrier_gateway_id"`
	DestinationCidrBlock        string `json:"destination_cidr_block"`
	DestinationIpv6CidrBlock    string `json:"destination_ipv6_cidr_block"`
	EgressOnlyInternetGatewayId string `json:"egress_only_gateway_id"`
	GatewayId                   string `json:"gateway_id"`
	InstanceId                  string `json:"instance_id"`
	LocalGatewayId              string `json:"local_gateway_id"`
	NatGatewayId                string `json:"nat_gateway_id"`
	NetworkInterfaceId          string `json:"network_interface_id"`
	RouteTableId                string `json:"route_table_id"`
	TransitGatewayId            string `json:"transit_gateway_id"`
	VpcEndpointId               string `json:"vpc_endpoint_id"`
	VpcPeeringConnectionId      string `json:"vpc_peering_connection_id"`
}

// RouteTable returns config for aws_route
func GetRouteConfig(e *ec2.Route) []AWSResourceConfig {
	cf := RouteConfig{
		Config:                      Config{},
		CarrierGatewayId:            e.CarrierGatewayId,
		DestinationCidrBlock:        e.DestinationCidrBlock,
		DestinationIpv6CidrBlock:    e.DestinationIpv6CidrBlock,
		EgressOnlyInternetGatewayId: e.EgressOnlyInternetGatewayId,
		GatewayId:                   e.GatewayId,
		InstanceId:                  e.InstanceId,
		LocalGatewayId:              e.LocalGatewayId,
		NatGatewayId:                e.NatGatewayId,
		NetworkInterfaceId:          e.NetworkInterfaceId,
		RouteTableId:                e.RouteTableId,
		TransitGatewayId:            e.TransitGatewayId,
		VpcEndpointId:               e.VpcEndpointId,
		VpcPeeringConnectionId:      e.VpcPeeringConnectionId,
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: e.AWSCloudFormationMetadata,
	}}
}
