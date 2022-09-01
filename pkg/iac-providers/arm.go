/*
    Copyright (C) 2022 Tenable, Inc.

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

	armv1 "github.com/tenable/terrascan/pkg/iac-providers/arm/v1"
)

// terraform specific constants
const (
	arm                  supportedIacType    = "arm"
	armV1                supportedIacVersion = "v1"
	armDefaultIacVersion                     = armV1
)

// register arm as an IaC provider with terrascan
func init() {
	// register iac provider
	RegisterIacProvider(arm, armV1, armDefaultIacVersion, reflect.TypeOf(armv1.ARMV1{}))
}
