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

package output

// ResourceConfig describes a resource present in IaC
type ResourceConfig struct {
	ID     string      `json:"id"`
	Name   string      `json:"name"`
	Source string      `json:"source"`
	Line   int         `json:"line"`
	Type   string      `json:"type"`
	Config interface{} `json:"config"`
}

// AllResourceConfigs is a list/slice of resource configs present in IaC
type AllResourceConfigs map[string][]ResourceConfig
