package httpserver

import (
	"github.com/accurics/terrascan/pkg/config"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/utils"
)

type webhookDenyRuleMatcher struct {
}

// This class should check if one of the violations found is relevant for the specified K8s deny rules
func (g *webhookDenyRuleMatcher) match(violation results.Violation, denyRules config.K8sDenyRules) bool {
	if &denyRules == nil {
		return false
	}

	// Currently we support:
	// 1. A minimum severity level
	// 2. A category list
	// In case one of the conditions is met, we return true. (We perform an OR between the rules)
	if len(denyRules.DeniedSeverity) > 0 && utils.CheckSeverity(violation.Severity, denyRules.DeniedSeverity) {
		return true
	}

	if denyRules.Categories != nil {
		for _, category := range denyRules.Categories {
			if category == violation.Category {
				return true
			}
		}
	}

	return false
}
