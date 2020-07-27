package cloudprovider

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
)

var (
	errCloudNotSupported = fmt.Errorf("cloud type not supported")
)

// NewCloudProvider returns a new CloudProvider
func NewCloudProvider(cloudType string) (cloudProvider CloudProvider, err error) {

	// get CloudProvider from supportedCloudProviders
	cloudProviderObject, supported := supportedCloudProviders[supportedCloudType(cloudType)]
	if !supported {
		zap.S().Errorf("cloud type '%s' not supported", cloudType)
		return cloudProvider, errCloudNotSupported
	}

	return reflect.New(cloudProviderObject).Interface().(CloudProvider), nil
}

// IsCloudSupported returns true/false depending on whether the cloud
// provider is supported in terrascan or not
func IsCloudSupported(cloudType string) bool {
	if _, supported := supportedCloudProviders[supportedCloudType(cloudType)]; !supported {
		return false
	}
	return true
}
