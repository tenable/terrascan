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

package initialize

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tenable/terrascan/pkg/policy"
)

type environmentPolicy struct {
	regoTemplate     string
	metadataFileName string
	resourceType     string
	policyMetadata   policy.RegoMetadata
}

func newPolicy(ruleMetadata environmentPolicyMetadata) (environmentPolicy, error) {
	var policy environmentPolicy
	var templateArgs map[string]interface{}

	policy.regoTemplate = "package accurics\n\n" + ruleMetadata.RuleTemplate
	policy.metadataFileName = ruleMetadata.RuleReferenceID + ".json"
	policy.resourceType = ruleMetadata.ResourceType

	policy.policyMetadata.Name = ruleMetadata.RuleName
	policy.policyMetadata.File = ruleMetadata.RegoName + ".rego"
	policy.policyMetadata.ResourceType = ruleMetadata.ResourceType
	policy.policyMetadata.Severity = ruleMetadata.Severity
	policy.policyMetadata.Description = ruleMetadata.RuleDisplayName
	policy.policyMetadata.ReferenceID = ruleMetadata.RuleReferenceID
	policy.policyMetadata.ID = ruleMetadata.RuleReferenceID
	policy.policyMetadata.Category = ruleMetadata.Category
	policy.policyMetadata.Version = ruleMetadata.Version

	templateString, ok := ruleMetadata.RuleArgument.(string)
	if !ok {
		return policy, fmt.Errorf("incorrect rule argument type, must be a string")
	}
	err := json.Unmarshal([]byte(templateString), &templateArgs)
	if err != nil {
		return policy, fmt.Errorf("error occurred while unmarshaling rule arguments into map[string]interface{}, error: '%w'", err)
	}
	policy.policyMetadata.TemplateArgs = templateArgs

	return policy, nil
}

func (p environmentPolicy) getType() string {
	provider := strings.ToLower(p.resourceType)

	if strings.HasPrefix(provider, "azure") {
		return "azure"
	}

	if strings.HasPrefix(provider, "google") {
		return "gcp"
	}

	if strings.HasPrefix(provider, "kubernetes") {
		return "k8s"
	}

	return strings.Split(provider, "_")[0]
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
