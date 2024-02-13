package stringutil

// StrOrDefault returns str or defaultStr when str is empty.
func StrOrDefault(str, defaultStr string) string {
	if str == "" {
		return defaultStr
	}
	return str
}
