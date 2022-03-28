/*
    Copyright (C) 2022 Accurics, Inc.

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

package initialize

import (
	"github.com/accurics/terrascan/pkg/policy"
)

type environmentPolicy struct {
	regoTemplate     string
	metadataFileName string
	resourceType     string
	policyMetadata   policy.RegoMetadata
}

type environmentPolicyMetadata struct {
	RuleName        string      `json:"ruleName"`
	RegoName        string      `json:"ruleTemplateName"`
	RuleArgument    interface{} `json:"ruleArgument"`
	Severity        string      `json:"severity"`
	RuleDisplayName string      `json:"ruleDisplayName"`
	Category        string      `json:"category"`
	RuleReferenceID string      `json:"ruleReferenceId"`
	Version         int         `json:"version"`
	RuleTemplate    string      `json:"ruleTemplate"`
	ResourceType    string      `json:"resourceType"`
}
