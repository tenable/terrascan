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
