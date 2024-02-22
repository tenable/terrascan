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
	"encoding/json"

	"github.com/awslabs/goformation/v7/cloudformation/s3"
)

// S3BucketPolicyConfig holds config for aws_s3_bucket_policy
type S3BucketPolicyConfig struct {
	Config
	PolicyDocument string `json:"policy"`
	Bucket         string `json:"bucket"`
}

// GetS3BucketPolicyConfig returns config for aws_s3_bucket_policy
// aws_s3_bucket_policy
func GetS3BucketPolicyConfig(p *s3.BucketPolicy) []AWSResourceConfig {
	cf := S3BucketPolicyConfig{
		Config: Config{
			Name: p.Bucket,
		},
		Bucket: p.Bucket,
	}

	policyDocument, err := json.Marshal(p.PolicyDocument)
	if err == nil {
		cf.PolicyDocument = string(policyDocument)
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: p.AWSCloudFormationMetadata,
	}}
}
