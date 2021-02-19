package utils

// AcceptedCategories is the list of all policy categories
var AcceptedCategories []string = []string{
	"LOGGING AND MONITORING",
	"COMPLIANCE VALIDATION",
	"RESILIENCE",
	"SECURITY BEST PRACTICES",
	"INFRASTRUCTURE SECURITY",
	"IDENTITY AND ACCESS MANAGEMENT",
	"CONFIGURATION AND VULNERABILITY ANALYSIS",
	"DATA PROTECTION",
}

// ValidateCategoryInput validates input for --category flag
func ValidateCategoryInput(categories []string) bool {
	for _, category := range categories {
		category = EnsureUpperCaseTrimmed(category)
		if !find(AcceptedCategories, category) {
			return false
		}
	}
	return true
}

// CheckCategory validates if the category of policy rule is present in the list of specificed categories
func CheckCategory(ruleCategory string, desiredCategories []string) bool {
	ruleCategory = EnsureUpperCaseTrimmed(ruleCategory)
	for i, category := range desiredCategories {
		desiredCategories[i] = EnsureUpperCaseTrimmed(category)
	}

	return find(desiredCategories, ruleCategory)
}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
