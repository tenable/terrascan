/*
    Copyright (C) 2020 Accurics, Inc.

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
