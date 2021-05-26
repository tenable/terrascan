package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/kinesis"
)

// KinesisStreamConfig holds config for aws_kinesis_stream
type KinesisStreamConfig struct {
	Config
	Name           string `json:"name"`
	KmsKeyID       string `json:"kms_key_id,omitempty"`
	EncryptionType string `json:"encryption_type,omitempty"`
}

// GetKinesisStreamConfig returns config for aws_kinesis_stream
func GetKinesisStreamConfig(k *kinesis.Stream) []AWSResourceConfig {
	cf := KinesisStreamConfig{
		Config: Config{
			Name: k.Name,
			Tags: k.Tags,
		},
		Name: k.Name,
	}

	if k.StreamEncryption != nil {
		cf.EncryptionType = k.StreamEncryption.EncryptionType
		cf.KmsKeyID = k.StreamEncryption.KeyId
	}

	return []AWSResourceConfig{{Resource: cf}}
}
