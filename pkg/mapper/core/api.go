package core

import (
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
)

// Mapper defines the base API that each IaC provider mapper must implement.
type Mapper interface {
	// Map transforms the provider specific template to terrascan native format.
	Map(doc *utils.IacDocument, params ...map[string]interface{}) (output.AllResourceConfigs, error)
}
