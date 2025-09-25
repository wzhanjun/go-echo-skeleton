package utils

import (
	"encoding/json"
)

func ToJson(data interface{}) string {
	result, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(result)
}

func TruncateString(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}

func ParseJSONOrDefault[T any](jsonStr string, defaultVal T) T {
	if jsonStr == "" {
		return defaultVal
	}

	var obj T
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		return defaultVal
	}

	return obj
}
