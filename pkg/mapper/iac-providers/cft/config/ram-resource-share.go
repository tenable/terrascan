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
	"github.com/awslabs/goformation/v6/cloudformation/ram"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// RAMResourceShareConfig holds config for RAMResourceShare
type RAMResourceShareConfig struct {
	Config
	Name                    string `json:"name"`
	AllowExternalPrincipals bool   `json:"allow_external_principals"`
}

// GetRAMResourceShareConfig returns config for RAMResourceShare
func GetRAMResourceShareConfig(r *ram.ResourceShare) []AWSResourceConfig {
	cf := RAMResourceShareConfig{
		Config: Config{
			Name: r.Name,
			Tags: r.Tags,
		},
		Name:                    r.Name,
		AllowExternalPrincipals: functions.GetBool(r.AllowExternalPrincipals),
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: r.AWSCloudFormationMetadata,
	}}
}
