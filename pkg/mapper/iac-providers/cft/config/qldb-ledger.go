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

import "github.com/awslabs/goformation/v5/cloudformation/qldb"

// QldbLedgerConfig holds config for aws_qldb_ledger resource
type QldbLedgerConfig struct {
	Config
	Name               string `json:"name,omitempty"`
	PermissionsMode    string `json:"permissions_mode"`
	DeletionProtection bool   `json:"deletion_protection"`
}

// GetQldbLedgerConfig returns config for aws_qldb_ledger resource
func GetQldbLedgerConfig(q *qldb.Ledger) []AWSResourceConfig {

	cf := QldbLedgerConfig{
		Config: Config{
			Name: q.Name,
			Tags: q.Tags,
		},
		Name:               q.Name,
		PermissionsMode:    q.PermissionsMode,
		DeletionProtection: q.DeletionProtection,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: q.AWSCloudFormationMetadata,
	}}
}
