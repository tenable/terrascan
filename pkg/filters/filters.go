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

package filters

import (
	"github.com/tenable/terrascan/pkg/policy"
)

// RegoMetadataPreLoadFilter is a pre load filter
// this filter would be while the policy files are processed by policy engine
type RegoMetadataPreLoadFilter struct {
	scanRules   []string
	skipRules   []string
	categories  []string
	policyTypes []string
	severity    string
	filterSpecs []policy.FilterSpecification
}

// NewRegoMetadataPreLoadFilter is a constructor func for RegoMetadataPreLoadFilter
func NewRegoMetadataPreLoadFilter(scanRules, skipRules, categories, policyTypes []string, severity string) *RegoMetadataPreLoadFilter {
	return &RegoMetadataPreLoadFilter{
		scanRules:   scanRules,
		skipRules:   skipRules,
		categories:  categories,
		policyTypes: policyTypes,
		severity:    severity,
		// add applicable filter specs to the list
		filterSpecs: []policy.FilterSpecification{
			ReferenceIDsFilterSpecification{scanRules},
			CategoryFilterSpecification{categories: categories},
			SeverityFilterSpecification{severity: severity},
			PolicyTypesFilterSpecification{policyTypes: policyTypes},
		},
	}
}

// IsFiltered checks whether a RegoMetadata should be filtered or not
func (r *RegoMetadataPreLoadFilter) IsFiltered(regoMetadata *policy.RegoMetadata) bool {
	// if skip rules are specified, RegoMetadata is not filtered
	if len(r.skipRules) < 1 {
		return false
	}
	refIDsSpec := ReferenceIDsFilterSpecification{r.skipRules}
	return refIDsSpec.IsSatisfied(regoMetadata)
}

// IsAllowed checks whether a RegoMetadata should be allowed or not
func (r *RegoMetadataPreLoadFilter) IsAllowed(regoMetadata *policy.RegoMetadata) bool {
	andSpec := AndFilterSpecification{r.filterSpecs}
	return andSpec.IsSatisfied(regoMetadata)
}

// RegoDataFilter is a pre scan filter,
// it will be used by policy engine before the evaluation of resources start
type RegoDataFilter struct{}

// Filter func will filter based on resource type
func (r *RegoDataFilter) Filter(rmap map[string]*policy.RegoData, input policy.EngineInput) map[string]*policy.RegoData {
	// if resource config is empty, return original map
	if len(*input.InputData) < 1 {
		return rmap
	}
	tempMap := make(map[string]*policy.RegoData)
	for resType := range *input.InputData {
		for k := range rmap {
			resFilterSpec := ResourceTypeFilterSpecification{resType}
			if resFilterSpec.IsSatisfied(&rmap[k].Metadata) {
				tempMap[k] = rmap[k]
			}
		}
	}
	return tempMap
}
