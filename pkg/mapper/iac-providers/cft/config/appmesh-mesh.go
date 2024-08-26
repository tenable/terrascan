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
	"github.com/awslabs/goformation/v7/cloudformation/appmesh"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// AppMeshEgressFilterBlock holds config for AppMeshEgressFilter
type AppMeshEgressFilterBlock struct {
	Type string `json:"type"`
}

// AppMeshSpecBlock holds config for AppMeshSpec
type AppMeshSpecBlock struct {
	EgressFilter []AppMeshEgressFilterBlock `json:"egress_filter"`
}

// AppMeshMeshConfig holds config for AppMeshMesh
type AppMeshMeshConfig struct {
	Config
	Name string             `json:"name"`
	Spec []AppMeshSpecBlock `json:"spec"`
}

// GetAppMeshMeshConfig returns config for AppMeshMesh
// aws_appmesh_mesh
func GetAppMeshMeshConfig(m *appmesh.Mesh) []AWSResourceConfig {
	var spec []AppMeshSpecBlock
	if m.Spec != nil {
		spec = make([]AppMeshSpecBlock, 1)

		if m.Spec.EgressFilter != nil {
			spec[0].EgressFilter = make([]AppMeshEgressFilterBlock, 1)

			spec[0].EgressFilter[0].Type = m.Spec.EgressFilter.Type
		}
	}

	cf := AppMeshMeshConfig{
		Config: Config{
			Name: functions.GetVal(m.MeshName),
			Tags: functions.PatchAWSTags(m.Tags),
		},
		Name: functions.GetVal(m.MeshName),
		Spec: spec,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: m.AWSCloudFormationMetadata,
	}}
}
