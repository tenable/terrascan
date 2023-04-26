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
func ValidateCategoryInput(categories []string) (bool, []string) {
	flag := false
	var invalidInputs []string
	for _, category := range categories {
		category = EnsureUpperCaseTrimmed(category)
		if !find(AcceptedCategories, category) {
			flag = true
			invalidInputs = append(invalidInputs, category)
		}
	}

	if flag {
		return false, invalidInputs
	}

	return true, invalidInputs
}

// CheckCategory validates if the category of policy rule is present in the list of specified categories
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
