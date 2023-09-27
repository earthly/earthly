package stringutil

import "regexp"

// NamedGroupMatches returns a map with all found named regex groups as keys and an array of all matches as the value
// and an array of all the keys in the order they were set in the regex (this is since the map keys order won't be predictable).
func NamedGroupMatches(s string, re *regexp.Regexp) (map[string][]string, []string) {
	all := make(map[string][]string)
	res := re.FindAllStringSubmatch(s, -1)
	names := make([]string, 0, len(re.SubexpNames()))
	for groupIdx, groupName := range re.SubexpNames() {
		if groupName == "" {
			continue
		}
		for _, matchRes := range res {
			if matchRes[groupIdx] != "" {
				if len(all[groupName]) == 0 {
					// only add the group name once and only if we have matches for it
					names = append(names, groupName)
				}
				all[groupName] = append(all[groupName], matchRes[groupIdx])
			}
		}
	}
	return all, names
}
