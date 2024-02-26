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

// IngressEgress holds config for SecurityGroupEgress, SecurityGroupIngress attributes of SecurityGroupConfig
type IngressEgress struct {
	IPProtocol  string   `json:"protocol"`
	Description string   `json:"description"`
	CidrIP      []string `json:"cidr_blocks"`
	CidrIpv6    []string `json:"ipv6_cidr_blocks"`
	FromPort    int      `json:"from_port"`
	ToPort      int      `json:"to_port"`
}

// SecurityGroupConfig holds config for aws_security_group
type SecurityGroupConfig struct {
	Config
	GroupDescription     string          `json:"description,omitempty"`
	GroupName            string          `json:"name"`
	SecurityGroupEgress  []IngressEgress `json:"egress"`
	SecurityGroupIngress []IngressEgress `json:"ingress"`
}

// GetSecurityGroupConfig returns config for aws_security_group
// aws_security_group
func GetSecurityGroupConfig(s *ec2.SecurityGroup) []AWSResourceConfig {
	cf := SecurityGroupConfig{
		Config: Config{
			Name: functions.GetVal(s.GroupName),
			Tags: functions.PatchAWSTags(s.Tags),
		},
		GroupName:        functions.GetVal(s.GroupName),
		GroupDescription: s.GroupDescription,
	}

	ingresses := make([]IngressEgress, 0)
	for _, i := range s.SecurityGroupIngress {
		ingress := getIngressEgress(i)
		ingresses = append(ingresses, ingress)
	}
	cf.SecurityGroupIngress = ingresses

	egresses := make([]IngressEgress, 0)
	for _, e := range s.SecurityGroupEgress {
		egress := getIngressEgress(e)
		egresses = append(egresses, egress)
	}
	cf.SecurityGroupEgress = egresses

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: s.AWSCloudFormationMetadata,
	}}
}

func getIngressEgress(ie any) IngressEgress {
	if egress, ok := ie.(*ec2.SecurityGroup_Egress); ok {
		return getEgress(egress)
	}
	if egress, ok := ie.(ec2.SecurityGroup_Egress); ok {
		return getEgress(&egress)
	}

	if ingress, ok := ie.(*ec2.SecurityGroup_Ingress); ok {
		return getIngress(ingress)
	}

	if ingress, ok := ie.(ec2.SecurityGroup_Ingress); ok {
		return getIngress(&ingress)
	}
	return IngressEgress{}
}

func getEgress(egress *ec2.SecurityGroup_Egress) IngressEgress {
	return IngressEgress{
		IPProtocol:  egress.IpProtocol,
		Description: functions.GetVal(egress.Description),
		CidrIP:      []string{functions.GetVal(egress.CidrIp)},
		CidrIpv6:    []string{functions.GetVal(egress.CidrIpv6)},
		FromPort:    functions.GetVal(egress.FromPort),
		ToPort:      functions.GetVal(egress.ToPort),
	}
}

func getIngress(egress *ec2.SecurityGroup_Ingress) IngressEgress {
	return IngressEgress{
		IPProtocol:  egress.IpProtocol,
		Description: functions.GetVal(egress.Description),
		CidrIP:      []string{functions.GetVal(egress.CidrIp)},
		CidrIpv6:    []string{functions.GetVal(egress.CidrIpv6)},
		FromPort:    functions.GetVal(egress.FromPort),
		ToPort:      functions.GetVal(egress.ToPort),
	}
}
