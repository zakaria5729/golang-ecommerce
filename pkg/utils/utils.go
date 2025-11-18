package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"

	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/option"
	"github.com/easy-comerce/backend/pkg/response"
	"golang.org/x/crypto/bcrypt"
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

func UintToString(s uint) string {
	return strconv.FormatUint(uint64(s), 10)
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

func ParseStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ParsePagination(pageStr, pageSizeStr string) (page, pageSize int) {
	page = 1
	pageSize = c.DefaultPageSize

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if l, err := strconv.Atoi(pageSizeStr); err == nil && l > 0 && l <= c.MaxPageSize {
			pageSize = l
		}
		if pageSize > c.MaxPageSize {
			pageSize = c.MaxPageSize
		}
	}

	return page, pageSize
}

func CalculatePagination(total int64, page int, pageSize int) (totalPages int, offset int) {
	if pageSize <= 0 || total < 0 {
		return 0, 0
	}

	if page < 1 {
		page = 1
	}

	if total == 0 {
		return 0, 0
	}

	totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	if page > totalPages {
		page = totalPages
	}

	offset = (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	return totalPages, offset
}

func BuildPaginatedResponse(data any, total int64, page int, pageSize int) *response.PaginatedResponse {
	totalPages, _ := CalculatePagination(total, page, pageSize)

	return &response.PaginatedResponse{
		Data: data,
		Pagination: response.PaginationMeta{
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
		if strings.ToLower(field) == c.All {
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
	if sortBy == "" || sortOrder == "" {
		return ""
	}

	allowedFields := map[string]bool{
		c.FieldID:        true,
		c.FieldCreatedAt: true,
		c.FieldUpdatedAt: true,
	}

	if fields != nil && len(*fields) > 0 {
		for _, field := range *fields {
			allowedFields[field] = true
		}
	}

	if !allowedFields[sortBy] {
		return ""
	}

	if strings.EqualFold(sortOrder, c.SortOrderDesc) {
		return sortBy + " " + c.SortOrderDesc
	}

	return sortBy + " " + c.SortOrderAsc
}

func BuildSortingOrders(sortables []option.SortOption, fields *[]string) string {
	if len(sortables) == 0 {
		return ""
	}

	allowedFields := map[string]bool{
		c.FieldID:        true,
		c.FieldCreatedAt: true,
		c.FieldUpdatedAt: true,
	}

	if fields != nil && len(*fields) > 0 {
		for _, f := range *fields {
			allowedFields[f] = true
		}
	}
	var parts []string

	for _, s := range sortables {
		if s.SortBy == "" || s.SortOrder == "" {
			continue
		}

		if !allowedFields[s.SortBy] {
			continue
		}

		order := c.SortOrderAsc
		if strings.EqualFold(s.SortOrder, c.SortOrderDesc) {
			order = c.SortOrderDesc
		}

		parts = append(parts, s.SortBy+" "+order)
	}

	return strings.Join(parts, ", ")
}

func GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
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

func BuildFullImageURL(publicDomain string, pathKey *string) *string {
	if pathKey == nil || *pathKey == "" || publicDomain == "" {
		return nil
	}

	var url string
	if strings.HasPrefix(publicDomain, "http://") || strings.HasPrefix(publicDomain, "https://") {
		url = fmt.Sprintf("%s/%s", publicDomain, *pathKey)
	} else {
		url = fmt.Sprintf("https://%s/%s", publicDomain, *pathKey)
	}

	return &url
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, target any, method string) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return false
	}
	return true
}

func HashPassword(password string) (error, string) {
	if password == "" {
		return errors.New("password is empty"), ""
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err, ""
	}
	return nil, string(hashedPassword)
}

func CheckPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GetProjectRootPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..")
}

func GetLogFolderPath() string {
	return filepath.Join(GetProjectRootPath(), c.LogFolderName)
}

func ExtractNameFromEmail(email string, capitalize bool) string {
	at := strings.Index(email, "@")
	if at == -1 {
		return ""
	}
	name := email[:at]
	if capitalize {
		name = CapitalizeFirst(name)
	}
	return name
}
