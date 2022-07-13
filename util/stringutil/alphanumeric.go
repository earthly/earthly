package stringutil

import "regexp"

var alphanumericRegexp *regexp.Regexp

func init() {
	var err error
	alphanumericRegexp, err = regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
}

// AlphanumericOnly removes all non-alphanumeric characters from the string
func AlphanumericOnly(s string) string {
	return alphanumericRegexp.ReplaceAllString(s, "")
}
