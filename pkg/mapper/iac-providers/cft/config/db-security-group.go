package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/rds"
)

// DBSecurityGroupConfig holds config for aws_db_security_group
type DBSecurityGroupConfig struct {
	Config
	Ingress []map[string]interface{} `json:"ingress"`
}

// GetDBSecurityGroupConfig returns config for aws_db_security_group
func GetDBSecurityGroupConfig(dbsg *rds.DBSecurityGroup) []AWSResourceConfig {
	cf := DBSecurityGroupConfig{
		Config: Config{
			Tags: dbsg.Tags,
		},
	}
	for _, dbsgi := range dbsg.DBSecurityGroupIngress {
		i := make(map[string]interface{})
		i["cidr"] = dbsgi.CIDRIP
		cf.Ingress = append(cf.Ingress, i)
	}
	return []AWSResourceConfig{{Resource: cf}}
}
