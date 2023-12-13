package inputgraph

func copyVisited(m map[string]struct{}) map[string]struct{} {
	m2 := map[string]struct{}{}
	for k := range m {
		m2[k] = struct{}{}
	}
	return m2
}

func uniqStrs(all []string) []string {
	m := map[string]struct{}{}
	for _, v := range all {
		m[v] = struct{}{}
	}
	ret := []string{}
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}
