package awsprovider

import (
	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// CreateNormalizedJSON creates a normalized json for the given input
func (a *AWSProvider) CreateNormalizedJSON(allResourcesConfig output.AllResourceConfigs) (interface{}, error) {
	return allResourcesConfig, nil
}
