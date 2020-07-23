package cloudprovider

import (
	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// CloudProvider defines the interface which every cloud provider needs to implement
// to claim support in terrascan
type CloudProvider interface {
	CreateNormalizedJSON(output.AllResourceConfigs) (interface{}, error)
}
