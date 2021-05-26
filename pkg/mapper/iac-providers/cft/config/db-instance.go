package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/rds"
)

// DBInstanceConfig holds config for aws_db_instance
type DBInstanceConfig struct {
	Config
	EnabledCloudWatchLogsExports []string `json:"enabled_cloudwatch_logs_exports"`
	AutoMinorVersionUpgrade      bool     `json:"auto_minor_version_upgrade"`
	CaCertIdentifier             string   `json:"ca_cert_identifier"`
	StorageEncrypted             bool     `json:"storage_encrypted"`
	KmsKeyID                     string   `json:"kms_key_id,omitempty"`
	IamDBAuthEnabled             bool     `json:"iam_database_authentication_enabled"`
	PubliclyAccessible           bool     `json:"publicly_accessible"`
}

// GetDBInstanceConfig returns config for aws_db_instance
func GetDBInstanceConfig(d *rds.DBInstance) []AWSResourceConfig {
	cf := DBInstanceConfig{
		Config: Config{
			Tags: d.Tags,
			Name: d.DBName,
		},
		EnabledCloudWatchLogsExports: d.EnableCloudwatchLogsExports,
		AutoMinorVersionUpgrade:      d.AutoMinorVersionUpgrade,
		StorageEncrypted:             d.StorageEncrypted,
		KmsKeyID:                     d.KmsKeyId,
		CaCertIdentifier:             d.CACertificateIdentifier,
		IamDBAuthEnabled:             d.EnableIAMDatabaseAuthentication,
		PubliclyAccessible:           d.PubliclyAccessible,
	}
	return []AWSResourceConfig{{Resource: cf}}
}
