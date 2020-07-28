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
