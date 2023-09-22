package stringutil

import "regexp"

func NamedGroupMatches(s string, re *regexp.Regexp) map[string][]string {
	all := make(map[string][]string)
	res := re.FindAllStringSubmatch(s, -1)
	for groupIdx, groupName := range re.SubexpNames() {
		if groupName == "" {
			continue
		}
		for _, matchRes := range res {
			if matchRes[groupIdx] != "" {
				all[groupName] = append(all[groupName], matchRes[groupIdx])
			}
		}
	}
	return all
}
