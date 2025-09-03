package utils

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
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

func ParsePagination(pageStr, pageSizeStr string) (page, pageSize int) {
	page = 1
	pageSize = constants.DefaultPageSize

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if l, err := strconv.Atoi(pageSizeStr); err == nil && l > 0 && l <= constants.MaxPageSize {
			pageSize = l
		}
	}

	return page, pageSize
}

func CalculatePagination(total, page, pageSize int) (totalPages, offset int) {
	if pageSize <= 0 {
		return 0, 0
	}

	totalPages = (total + pageSize - 1) / pageSize
	if page < 1 {
		page = 1
	}

	if totalPages > 0 && page > totalPages {
		page = totalPages
	}

	offset = (page - 1) * pageSize
	return totalPages, offset
}

func BuildPaginatedResponse(data any, total, page, pageSize int) *models.PaginatedResponse {
	totalPages, _ := CalculatePagination(total, page, pageSize)

	return &models.PaginatedResponse{
		Data: data,
		Pagination: models.PaginationMeta{
			Page:          page,
			PageSize:      pageSize,
			TotalElements: total,
			TotalPages:    totalPages,
			IsFirst:       page == 1,
			IsLast:        page >= totalPages,
		},
	}
}

func BuildSelectFields(defaultFields []string, optionalFields []string, include []string) []string {
	includeAll := false

	for _, field := range include {
		if strings.ToLower(field) == constants.All {
			includeAll = true
			break
		}
	}

	if includeAll {
		return []string{"*"}
	}

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

func BuildSortingOrder(sortBy string, sortOrder string, fields *[]string) string {
	if sortBy == "" {
		return ""
	}

	allowedFields := map[string]bool{
		constants.FieldID:        true,
		constants.FieldCreatedAt: true,
		constants.FieldUpdatedAt: true,
	}

	if fields != nil && len(*fields) > 0 {
		for _, field := range *fields {
			allowedFields[field] = true
		}
	}

	if !allowedFields[sortBy] {
		return ""
	}

	if strings.EqualFold(sortOrder, constants.SortOrderDesc) {
		return sortBy + " " + constants.SortOrderDesc
	}

	return sortBy + " " + constants.SortOrderAsc
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
