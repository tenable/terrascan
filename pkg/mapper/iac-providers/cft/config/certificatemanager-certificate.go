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
	"github.com/awslabs/goformation/v7/cloudformation/certificatemanager"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// CertificateManagerCertificateConfig holds config for CertificateManagerCertificate
type CertificateManagerCertificateConfig struct {
	Config
	DomainName       string `json:"domain_name"`
	ValidationMethod string `json:"validation_method"`
}

// GetCertificateManagerCertificateConfig returns config for CertificateManagerCertificate
// aws_acm_certificate
func GetCertificateManagerCertificateConfig(c *certificatemanager.Certificate) []AWSResourceConfig {
	cf := CertificateManagerCertificateConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(c.Tags),
		},
		DomainName:       c.DomainName,
		ValidationMethod: functions.GetVal(c.ValidationMethod),
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: c.AWSCloudFormationMetadata,
	}}
}
