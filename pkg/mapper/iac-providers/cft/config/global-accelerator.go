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
	"github.com/awslabs/goformation/v7/cloudformation/globalaccelerator"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// GlobalAcceleratorConfig holds config for aws_globalaccelerator_accelerator resource
type GlobalAcceleratorConfig struct {
	Config
	Name          string `json:"name"`
	IPAddressType string `json:"ip_address_type"`
	Enabled       bool   `json:"enabled"`
}

// GetGlobalAcceleratorConfig returns config for aws_globalaccelerator_accelerator resource
// aws_globalaccelerator_accelerator
func GetGlobalAcceleratorConfig(g *globalaccelerator.Accelerator) []AWSResourceConfig {
	cf := GlobalAcceleratorConfig{
		Config: Config{
			Name: g.Name,
			Tags: functions.PatchAWSTags(g.Tags),
		},
		Name:          g.Name,
		Enabled:       functions.GetVal(g.Enabled),
		IPAddressType: functions.GetVal(g.IpAddressType),
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: g.AWSCloudFormationMetadata,
	}}

}
