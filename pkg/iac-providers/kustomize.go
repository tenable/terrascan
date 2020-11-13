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

	kustomizev3 "github.com/accurics/terrascan/pkg/iac-providers/kustomize/v3"
)

// kustomize specific constants
const (
	kustomize                  supportedIacType    = "kustomize"
	kustomizeV3                supportedIacVersion = "v3"
	kustomizeDefaultIacVersion                     = kustomizeV3
)

// register kustomize as an IaC provider with terrascan
func init() {
	// register iac provider
	RegisterIacProvider(kustomize, kustomizeV3, kustomizeDefaultIacVersion, reflect.TypeOf(kustomizev3.KustomizeV3{}))
}
