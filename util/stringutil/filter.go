package stringutil

// ListContains returns true if the string, s, is contained in the list, l, of strings
func ListContains(l []string, s string) bool {
	for _, x := range l {
		if x == s {
			return true
		}
	}
	return false
}

// FilterElementsFromList filters out elements from a given list
func FilterElementsFromList(l []string, without ...string) []string {
	filtered := []string{}
	for _, s := range l {
		if !ListContains(without, s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
