/*
    Copyright (C) 2021 Accurics, Inc.

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
	"fmt"
	"strings"

	"github.com/awslabs/goformation/v5/cloudformation/s3"
)

const (
	// PublicAccessBlock represents subresource aws_s3_bucket_public_access_block for attribute PublicAccessBlockConfiguration
	PublicAccessBlock = "PublicAccessBlock"
)

// S3BucketConfig holds config for aws_s3_bucket
type S3BucketConfig struct {
	Config
	Bucket               string                       `json:"bucket"`
	AccessControl        string                       `json:"acl"`
	BucketEncryption     []ServerSideEncryptionConfig `json:"server_side_encryption_configuration,omitempty"`
	Logging              []LoggingConfig              `json:"logging"`
	WebsiteConfiguration []WebsiteConfig              `json:"website"`
	Versioning           []VersioningConfig           `json:"versioning,omitempty"`
}

// ServerSideEncryptionConfig holds config for server_side_encryption_configuration
type ServerSideEncryptionConfig struct {
	ServerSideEncryptionConfiguration []ServerSideEncryptionRule `json:"rule"`
}

// ServerSideEncryptionRule holds config for rule
type ServerSideEncryptionRule struct {
	ServerSideEncryptionByDefault []DefaultSSEConfig `json:"apply_server_side_encryption_by_default,omitempty"`
	BucketKeyEnabled              bool               `json:"bucket_key_enabled"`
}

// DefaultSSEConfig holds config for apply_server_side_encryption_by_default
type DefaultSSEConfig struct {
	KMSMasterKeyID string `json:"kms_master_key_id"`
	SSEAlgorithm   string `json:"sse_algorithm"`
}

// LoggingConfig holds config for logging
type LoggingConfig struct {
	DestinationBucketName string `json:"target_bucket"`
	LogFilePrefix         string `json:"target_prefix"`
}

// WebsiteConfig holds config for website
type WebsiteConfig struct {
	RedirectAllRequestsTo interface{} `json:"redirect_all_requests_to"`
	RoutingRules          interface{} `json:"routing_rules"`
	ErrorDocument         string      `json:"error_document"`
	IndexDocument         string      `json:"index_document"`
}

// VersioningConfig holds config for versioning
type VersioningConfig struct {
	Status bool `json:"enabled"`
}

// S3BucketPublicAccessBlockConfig holds config for aws_s3_bucket_public_access_block
type S3BucketPublicAccessBlockConfig struct {
	Config
	Bucket                string `json:"bucket"`
	BlockPublicAcls       bool   `json:"block_public_acls"`
	BlockPublicPolicy     bool   `json:"block_public_policy"`
	IgnorePublicAcls      bool   `json:"ignore_public_acls"`
	RestrictPublicBuckets bool   `json:"restrict_public_buckets"`
}

// GetS3BucketConfig returns config for aws_s3_bucket
func GetS3BucketConfig(s *s3.Bucket) []AWSResourceConfig {

	resourceConfigs := make([]AWSResourceConfig, 0)

	cf := S3BucketConfig{
		Config: Config{
			Name: s.BucketName,
			Tags: s.Tags,
		},
		Bucket:        s.BucketName,
		AccessControl: strings.ToLower(s.AccessControl),
	}

	// add sse configurations
	if s.BucketEncryption != nil {
		sseRules := make([]ServerSideEncryptionRule, 0)
		for _, sseRule := range s.BucketEncryption.ServerSideEncryptionConfiguration {
			if sseRule.ServerSideEncryptionByDefault != nil {
				defaultConfig := DefaultSSEConfig{
					KMSMasterKeyID: sseRule.ServerSideEncryptionByDefault.KMSMasterKeyID,
					SSEAlgorithm:   sseRule.ServerSideEncryptionByDefault.SSEAlgorithm,
				}
				sseRules = append(sseRules, ServerSideEncryptionRule{
					ServerSideEncryptionByDefault: []DefaultSSEConfig{defaultConfig},
				})
			}
		}
		cf.BucketEncryption = []ServerSideEncryptionConfig{{
			ServerSideEncryptionConfiguration: sseRules,
		}}
	}

	// add logging configurasions
	if s.LoggingConfiguration != nil {
		cf.Logging = []LoggingConfig{{
			DestinationBucketName: s.LoggingConfiguration.DestinationBucketName,
			LogFilePrefix:         s.LoggingConfiguration.LogFilePrefix,
		}}
	}

	// add website configurations
	if s.WebsiteConfiguration != nil {
		cf.WebsiteConfiguration = []WebsiteConfig{{
			IndexDocument:         s.WebsiteConfiguration.IndexDocument,
			ErrorDocument:         s.WebsiteConfiguration.ErrorDocument,
			RedirectAllRequestsTo: s.WebsiteConfiguration.RedirectAllRequestsTo,
			RoutingRules:          s.WebsiteConfiguration.RoutingRules,
		}}
	}

	// add versioning configurations
	if s.VersioningConfiguration != nil {
		var status bool
		if s.VersioningConfiguration.Status == "Enabled" {
			status = true
		}
		cf.Versioning = []VersioningConfig{{
			Status: status,
		}}
	}

	// add aws_s3_bucket
	resourceConfigs = append(resourceConfigs, AWSResourceConfig{
		Resource: cf,
		Metadata: s.AWSCloudFormationMetadata,
	})

	// add aws_s3_bucket_public_access_block
	if s.PublicAccessBlockConfiguration != nil {
		resourceConfigs = append(resourceConfigs, AWSResourceConfig{
			Resource: S3BucketPublicAccessBlockConfig{
				Config: Config{
					Name: s.BucketName,
				},
				Bucket:                fmt.Sprintf("aws_s3_bucket.%s", s.BucketName),
				BlockPublicAcls:       s.PublicAccessBlockConfiguration.BlockPublicAcls,
				BlockPublicPolicy:     s.PublicAccessBlockConfiguration.BlockPublicPolicy,
				IgnorePublicAcls:      s.PublicAccessBlockConfiguration.IgnorePublicAcls,
				RestrictPublicBuckets: s.PublicAccessBlockConfiguration.RestrictPublicBuckets,
			},
			Metadata: s.AWSCloudFormationMetadata,
			Type:     PublicAccessBlock,
			Name:     s.BucketName,
		})
	}

	return resourceConfigs
}
