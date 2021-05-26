package config

import (
	"encoding/json"

	"github.com/awslabs/goformation/v4/cloudformation/elasticsearch"
)

const (
	// ElasticsearchDomainAccessPolicy represents subresource aws_elasticsearch_domain_policy for Policy attribute
	ElasticsearchDomainAccessPolicy = "Policy"
)

// ElasticsearchDomainConfig holds config for aws_elastisearch_domain
type ElasticsearchDomainConfig struct {
	EncryptionAtRest            interface{} `json:"encrypt_at_rest,omitempty"`
	LogPublishingOptions        interface{} `json:"log_publishing_options,omitempty"`
	NodeToNodeEncryptionOptions interface{} `json:"node_to_node_encryption,omitempty"`
	Config
}

// ElasticsearchDomainAccessPolicyConfig holds config for aws_elasticsearch_domain_policy
type ElasticsearchDomainAccessPolicyConfig struct {
	Config
	DomainName     string `json:"domain_name"`
	AccessPolicies string `json:"access_policies"`
}

// EncryptionAtRestConfig holds config for encrypt_at_rest attribute of aws_elastisearch_domain
type EncryptionAtRestConfig struct {
	KmsKeyID string `json:"kms_key_id,omitempty"`
	Enabled  bool   `json:"enabled"`
}

// LogPublishingOptionsConfig holds config for log_publishing_options attribute of aws_elastisearch_domain
type LogPublishingOptionsConfig struct {
	LogType string `json:"log_type,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

// NodeToNodeEncryptionOptionsConfig holds config for node_to_node_encryption attribute of aws_elastisearch_domain
type NodeToNodeEncryptionOptionsConfig struct {
	Enabled bool `json:"enabled,omitempty"`
}

// GetElasticsearchDomainConfig returns config for aws_elastisearch_domain and aws_elasticsearch_domain_policy
func GetElasticsearchDomainConfig(d *elasticsearch.Domain) []AWSResourceConfig {
	resourceConfigs := make([]AWSResourceConfig, 0)

	// add domain config
	esDomainConfig := ElasticsearchDomainConfig{
		Config: Config{
			Name: d.DomainName,
			Tags: d.Tags,
		},
	}

	if d.LogPublishingOptions != nil {
		lpConfig := make([]LogPublishingOptionsConfig, 0)
		for ltype, options := range d.LogPublishingOptions {
			lpConfig = append(lpConfig, LogPublishingOptionsConfig{
				Enabled: options.Enabled,
				LogType: ltype,
			})
		}
		esDomainConfig.LogPublishingOptions = lpConfig
	}

	if d.NodeToNodeEncryptionOptions != nil {
		esDomainConfig.NodeToNodeEncryptionOptions = []NodeToNodeEncryptionOptionsConfig{{
			Enabled: d.NodeToNodeEncryptionOptions.Enabled,
		}}
	}

	if d.EncryptionAtRestOptions != nil {
		enc := EncryptionAtRestConfig{
			KmsKeyID: d.EncryptionAtRestOptions.KmsKeyId,
			Enabled:  d.EncryptionAtRestOptions.Enabled,
		}
		esDomainConfig.EncryptionAtRest = []EncryptionAtRestConfig{enc}
	}

	resourceConfigs = append(resourceConfigs, AWSResourceConfig{
		Resource: esDomainConfig,
	})

	// add domain access policy as aws_elasticsearch_domain_policy
	if d.AccessPolicies != nil {
		policyConfig := ElasticsearchDomainAccessPolicyConfig{
			Config: Config{
				Name: d.DomainName,
			},
		}
		policies, err := json.Marshal(d.AccessPolicies)
		if err == nil {
			policyConfig.AccessPolicies = string(policies)
		}
		resourceConfigs = append(resourceConfigs, AWSResourceConfig{
			Resource: policyConfig,
			Type:     ElasticsearchDomainAccessPolicy,
			Name:     d.DomainName,
		})
	}

	return resourceConfigs
}
