package awsProvider

import (
	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// CreateNormalizedJson creates a normalized json for the given input
func (a *AWSProvider) CreateNormalizedJson(allResourcesConfig output.AllResourceConfigs) (interface{}, error) {
	return allResourcesConfig, nil
}
