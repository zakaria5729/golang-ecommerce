package utils

import (
	"fmt"
	"strings"
)

func ParseCommaSeparatedString(input string) []string {
	var result []string
	if input != "" {
		result = strings.Split(input, ",")
		for i, field := range result {
			result[i] = strings.TrimSpace(field)
		}
	}
	return result
}

func ParseUint(s string) (*uint, error) {
	if s == "" || s == "null" {
		return nil, nil
	}

	var result uint
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func ParseBoolPtr(s string) *bool {
	if s == "" {
		return nil
	}

	result := s == "true"
	return &result
}

func BuildSelectFields(defaultFields []string, optionalFields []string, include []string) []string {
	selectFields := make([]string, len(defaultFields))
	copy(selectFields, defaultFields)

	optionalMap := make(map[string]bool)
	for _, field := range optionalFields {
		optionalMap[field] = true
	}

	for _, field := range include {
		if optionalMap[field] {
			selectFields = append(selectFields, field)
		}
	}

	return selectFields
}
