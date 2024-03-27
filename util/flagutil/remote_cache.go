package flagutil

import (
	"github.com/pkg/errors"
	"strings"
)

func ParseImageNameAndAttrs(s string) (string, map[string]string, error) {
	entries := strings.Split(s, ",")
	imageName := entries[0]
	attrs := make(map[string]string)
	var err error
	for _, entry := range entries[1:] {
		pair := strings.Split(strings.TrimSpace(entry), "=")
		if len(pair) != 2 {
			return "", attrs, errors.Errorf("failed to parse remote cache attribute: expected a key=value pair while parsing %q", entry)
		}
		attrs[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
	}
	return imageName, attrs, err
}
