package parseutil

import (
	"errors"
	"strings"
)

// StringToMap expects to get a string in the form of key1=val1,key2=val2,...
// and returns a map with the keys and values
func StringToMap(str string) (map[string]string, error) {
	pairs := strings.Split(str, ",")
	kvp := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		if strings.TrimSpace(pair) == "" {
			continue
		}
		k, v, ok := strings.Cut(pair, "=")
		if !ok {
			return nil, errors.New("key/value must be set with =")
		}
		kvp[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return kvp, nil
}
