package cloudProvider

import (
	"reflect"

	awsProvider "github.com/accurics/terrascan/pkg/cloud-providers/aws"
)

// SupportedCloudType data type for supported IaC provider
type supportedCloudType string

// supported IaC providers
const (
	aws supportedCloudType = "aws"
)

// map of supported IaC providers
var supportedCloudProviders map[supportedCloudType]reflect.Type

// initializes a map of supported IaC providers
func init() {

	supportedCloudProviders = make(map[supportedCloudType]reflect.Type)

	// aws support
	supportedCloudProviders[aws] = reflect.TypeOf(awsProvider.AWSProvider{})
}
