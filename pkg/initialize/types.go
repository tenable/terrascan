package initialize

import (
	"github.com/accurics/terrascan/pkg/policy"
)

type commercialPolicy struct {
	regoTemplate     string
	metadataFileName string
	resourceType     string
	policyMetadata   policy.RegoMetadata
}

type commercialPolicyMetadata struct {
	RuleName        string      `json:"ruleName"`
	RegoName        string      `json:"rule"`
	RuleTemplateID  string      `json:"ruleTemplateId"`
	RuleArgument    interface{} `json:"ruleArgument"`
	Severity        string      `json:"severity"`
	RuleDisplayName string      `json:"ruleDisplayName"`
	Category        string      `json:"category"`
	RuleReferenceId string      `json:"ruleReferenceId"`
	Version         int         `json:"version"`
	RuleTemplate    string      `json:"ruleTemplate"`
	ResourceType    string      `json:"resourceType"`
}
