package filters

import (
	"github.com/accurics/terrascan/pkg/policy"
)

// RegoMetadataPreLoadFilter is a pre load filter
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
			RerefenceIDsFilterSpecification{scanRules},
			CategoryFilterSpecification{categories: categories},
			SeverityFilterSpecification{severity: severity},
			PolicyTypeFilterSpecification{policyTypes: policyTypes},
		},
	}
}

// IsFiltered checks whether a RegoMetada should be filtered or not
func (r *RegoMetadataPreLoadFilter) IsFiltered(regoMetadata *policy.RegoMetadata) bool {
	// if length of skip rules is RegoMetada is not filtered
	if len(r.skipRules) < 1 {
		return false
	}
	refIDsSpec := RerefenceIDsFilterSpecification{r.skipRules}
	return refIDsSpec.IsSatisfied(regoMetadata)
}

// IsAllowed checks whether a RegoMetada should be allowed or not
func (r *RegoMetadataPreLoadFilter) IsAllowed(regoMetadata *policy.RegoMetadata) bool {
	andSpec := AndFilterSpecification{r.filterSpecs}
	return andSpec.IsSatisfied(regoMetadata)
}

// RegoDataFilter is a pre scan filter
type RegoDataFilter struct{}

// Filter func filters based on resource type
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
