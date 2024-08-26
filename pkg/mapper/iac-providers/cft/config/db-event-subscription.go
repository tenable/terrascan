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
	"github.com/awslabs/goformation/v7/cloudformation/rds"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// DBEventSubscriptionConfig holds config for aws_db_event_subscription resource
type DBEventSubscriptionConfig struct {
	Config
	SnsTopicArn     string   `json:"sns_topic"`
	Enabled         bool     `json:"enabled,omitempty"`
	EventCategories []string `json:"event_categories,omitempty"`
	SourceIds       []string `json:"source_ids,omitempty"`
	SourceType      string   `json:"source_type,omitempty"`
}

// GetDBEventSubscriptionConfig returns config for aws_db_event_subscription resource
// aws_db_event_subscription
func GetDBEventSubscriptionConfig(d *rds.EventSubscription) []AWSResourceConfig {
	cf := DBEventSubscriptionConfig{
		Config:          Config{},
		SnsTopicArn:     d.SnsTopicArn,
		Enabled:         functions.GetVal(d.Enabled),
		EventCategories: d.EventCategories,
		SourceIds:       d.SourceIds,
		SourceType:      functions.GetVal(d.SourceType),
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: d.AWSCloudFormationMetadata,
	}}
}
