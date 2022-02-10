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

import "github.com/awslabs/goformation/v5/cloudformation/rds"

type DBEventSubscriptionConfig struct {
	Config
	SnsTopicArn     string   `json:"sns_topic"`
	Enabled         bool     `json:"enabled,omitempty"`
	EventCategories []string `json:"event_categories,omitempty"`
	SourceIds       []string `json:"source_ids,omitempty"`
	SourceType      string   `json:"source_type,omitempty"`
}

func GetDBEventSubscriptionConfig(db *rds.EventSubscription) []AWSResourceConfig {

	cf := DBEventSubscriptionConfig{
		Config:          Config{},
		SnsTopicArn:     db.SnsTopicArn,
		Enabled:         db.Enabled,
		EventCategories: db.EventCategories,
		SourceIds:       db.SourceIds,
		SourceType:      db.SourceType,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: db.AWSCloudFormationMetadata,
	}}

}
