package utils

import (
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type PaginatedResponse struct {
	Data       any            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

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

func Trim(input string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(input), " "))
}

func CapitalizeFirst(input string) string {
	if len(input) == 0 {
		return input
	}
	return strings.ToUpper(input[:1]) + strings.ToLower(input[1:])
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

func ParseInt(s string) (*int, error) {
	if s == "" || s == "null" {
		return nil, nil
	}

	result, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func ParseFloat(s string) (*float64, error) {
	if s == "" || s == "null" {
		return nil, nil
	}

	result, err := strconv.ParseFloat(s, 64)
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

func ParsePagination(pageStr, limitStr string, defaultLimit int) (page, limit int) {
	page = 1
	limit = defaultLimit

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	return page, limit
}

func CalculatePagination(total, page, limit int) (totalPages, offset int) {
	totalPages = int(math.Ceil(float64(total) / float64(limit)))
	offset = (page - 1) * limit
	return totalPages, offset
}

func BuildPaginatedResponse(data any, total, page, limit int) *PaginatedResponse {
	totalPages, _ := CalculatePagination(total, page, limit)

	return &PaginatedResponse{
		Data: data,
		Pagination: PaginationMeta{
			Page:      page,
			PageSize:  limit,
			Total:     total,
			TotalPage: totalPages,
		},
	}
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

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(phone)
}

func IsValidUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(strings.ToLower(uuid))
}

func ContainsString(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func RemoveString(slice []string, item string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

func UniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(slice))

	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

func GetMapValueOrDefault(m map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}

func GetStringValueOrDefault(m map[string]interface{}, key string, defaultValue string) string {
	if value, exists := m[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}
