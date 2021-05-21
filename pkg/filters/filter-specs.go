package filters

import (
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/utils"
)

// PolicyTypeFilterSpecification is policy type based Filter Spec
type PolicyTypeFilterSpecification struct {
	policyType string
}

// IsSatisfied implementation for policy type based Filter spec
func (p PolicyTypeFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	// if policy type is not present for metadata, return true
	if len(r.PolicyType) < 1 {
		return true
	}
	return p.policyType == r.PolicyType
}

// ResourceTypeFilterSpecification is resource type based Filter Spec
type ResourceTypeFilterSpecification struct {
	resourceType string
}

// IsSatisfied implementation for resource type based Filter spec
func (rs ResourceTypeFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	// if resource type is not present for metadata, return true
	if len(r.ResourceType) < 1 {
		return true
	}
	return rs.resourceType == r.ResourceType
}

// RerefenceIDFilterSpecification is reference ID based Filter Spec
type RerefenceIDFilterSpecification struct {
	ReferenceID string
}

// IsSatisfied implementation for reference ID based Filter spec
func (rs RerefenceIDFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	return rs.ReferenceID == r.ReferenceID
}

// RerefenceIDsFilterSpecification is reference IDs based Filter Spec
type RerefenceIDsFilterSpecification struct {
	ReferenceIDs []string
}

// IsSatisfied implementation for reference IDs based Filter spec
func (rs RerefenceIDsFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	if len(rs.ReferenceIDs) < 1 {
		return true
	}
	isSatisfied := false
	for _, refID := range rs.ReferenceIDs {
		rfIDSpec := RerefenceIDFilterSpecification{refID}
		if rfIDSpec.IsSatisfied(r) {
			isSatisfied = true
			break
		}
	}
	return isSatisfied
}

// CategoryFilterSpecification is categories based Filter Spec
type CategoryFilterSpecification struct {
	categories []string
}

// IsSatisfied implementation for category based Filter spec
func (c CategoryFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	if len(c.categories) < 1 {
		return true
	}
	return utils.CheckCategory(r.Category, c.categories)
}

// SeverityFilterSpecification is severity based Filter Spec
type SeverityFilterSpecification struct {
	severity string
}

// IsSatisfied implementation for severity based Filter spec
func (s SeverityFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	if len(s.severity) < 1 {
		return true
	}
	return utils.CheckSeverity(r.Severity, s.severity)
}

// AndFilterSpecification is a logical AND Filter spec
type AndFilterSpecification struct {
	filterSpecs []policy.FilterSpecification
}

// IsSatisfied implementation for And Filter spec
func (a AndFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	if len(a.filterSpecs) < 1 {
		return false
	}
	isSatisfied := true
	for _, filterSpec := range a.filterSpecs {
		isSatisfied = isSatisfied && filterSpec.IsSatisfied(r)
		if !isSatisfied {
			return isSatisfied
		}
	}
	return isSatisfied
}
