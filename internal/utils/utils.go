package utils

import "fmt"

func GetStringValue(data map[string]any, key string) (string, error) {
	value, ok := data[key]

	if !ok {
		return "", fmt.Errorf("the key \"%v\" is missing", key)
	}

	stringValue, ok := value.(string)

	if !ok || stringValue == "" {
		return "", fmt.Errorf("the value of the key \"%v\" is incorrect", key)
	}

	return stringValue, nil
}
