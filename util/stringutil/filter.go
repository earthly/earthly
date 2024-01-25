package stringutil

import "golang.org/x/exp/slices"

// FilterElementsFromList filters out elements from a given list
func FilterElementsFromList(l []string, without ...string) []string {
	filtered := []string{}
	for _, s := range l {
		if !slices.Contains(without, s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
