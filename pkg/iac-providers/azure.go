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

package iacprovider

import (
	"reflect"

	azurev1 "github.com/accurics/terrascan/pkg/iac-providers/azure/v1"
)

// terraform specific constants
const (
	azure                  supportedIacType    = "azure"
	azureV1                supportedIacVersion = "v1"
	azureDefaultIacVersion                     = azureV1
)

// register azure as an IaC provider with terrascan
func init() {
	// register iac provider
	RegisterIacProvider(azure, azureV1, azureDefaultIacVersion, reflect.TypeOf(azurev1.ARM{}))
}
