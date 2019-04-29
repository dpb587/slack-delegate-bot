package config

import (
	"fmt"
	"strings"
)

func KeyValueTuple(from interface{}) (string, interface{}, error) {
	fromMap, ok := from.(map[interface{}]interface{})
	if !ok {
		return "", nil, fmt.Errorf("expected map[string]interface{}: found %v", from)
	}

	var keys []string

	for k := range fromMap {
		kStr, ok := k.(string)
		if !ok {
			return "", nil, fmt.Errorf("expected string key: found %v", k)
		}

		keys = append(keys, kStr)
	}

	if len(keys) == 0 {
		return "", nil, fmt.Errorf("expected exactly one key-value tuple: found none")
	} else if len(keys) > 1 {
		return "", nil, fmt.Errorf("expected exactly one key-value tuple: found multiple (%s)", strings.Join(keys, ", "))
	}

	return keys[0], fromMap[keys[0]], nil
}
