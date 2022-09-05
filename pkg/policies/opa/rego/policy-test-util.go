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

package rego

import (
	"encoding/json"
	"os"

	"github.com/tenable/terrascan/pkg/policy"
)

// LoadRegoMetadata reads rego meta data file
func LoadRegoMetadata(file string) (*policy.RegoMetadata, error) {
	metadata, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// Read metadata into struct
	regoMetadata := policy.RegoMetadata{}
	if err = json.Unmarshal(metadata, &regoMetadata); err != nil {
		return nil, err
	}
	return &regoMetadata, nil
}
