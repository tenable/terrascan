package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
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
func GetSecurityGroupConfig(s *ec2.SecurityGroup) []AWSResourceConfig {
	cf := SecurityGroupConfig{
		Config: Config{
			Name: s.GroupName,
		},
		GroupName:        s.GroupName,
		GroupDescription: s.GroupDescription,
	}

	ingresses := make([]IngressEgress, 0)
	for _, i := range s.SecurityGroupIngress {
		ingress := IngressEgress{
			IPProtocol:  i.IpProtocol,
			Description: i.Description,
			CidrIP:      []string{i.CidrIp},
			CidrIpv6:    []string{i.CidrIpv6},
			FromPort:    i.FromPort,
			ToPort:      i.ToPort,
		}
		ingresses = append(ingresses, ingress)
	}
	cf.SecurityGroupIngress = ingresses

	egresses := make([]IngressEgress, 0)
	for _, e := range s.SecurityGroupEgress {
		egress := IngressEgress{
			IPProtocol:  e.IpProtocol,
			Description: e.Description,
			CidrIP:      []string{e.CidrIp},
			CidrIpv6:    []string{e.CidrIpv6},
			FromPort:    e.FromPort,
			ToPort:      e.ToPort,
		}
		egresses = append(egresses, egress)
	}
	cf.SecurityGroupEgress = egresses

	return []AWSResourceConfig{{Resource: cf}}
}
