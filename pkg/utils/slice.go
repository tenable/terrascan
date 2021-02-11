package utils

// IsSliceEqual checks if two slices of string are equal or not
func IsSliceEqual(list1, list2 []string) bool {
	if len(list1) != len(list2) {
		return false
	}
	for i := range list1 {
		if list1[i] != list2[i] {
			return false
		}
	}
	return true
}
