package validator

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	return fmt.Sprintf("validation failed: %s", v[0].Message)
}

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

func (v *ValidationErrors) AddError(field, message string) {
	*v = append(*v, ValidationError{Field: field, Message: message})
}

type Validator interface {
	Validate() ValidationErrors
}

func ValidateRequired(value, fieldName string) ValidationErrors {
	var errors ValidationErrors
	if strings.TrimSpace(value) == "" {
		errors.AddError(fieldName, fmt.Sprintf("%s is required", fieldName))
	}
	return errors
}

func ValidateMinLength(value, fieldName string, minLength int) ValidationErrors {
	var errors ValidationErrors
	if len(strings.TrimSpace(value)) < minLength {
		errors.AddError(fieldName, fmt.Sprintf("%s must be at least %d characters long", fieldName, minLength))
	}
	return errors
}

func ValidateMaxLength(value, fieldName string, maxLength int) ValidationErrors {
	var errors ValidationErrors
	if len(value) > maxLength {
		errors.AddError(fieldName, fmt.Sprintf("%s must not exceed %d characters", fieldName, maxLength))
	}
	return errors
}

func ValidateURL(value, fieldName string) ValidationErrors {
	var errors ValidationErrors
	if value != "" {
		if _, err := url.ParseRequestURI(value); err != nil {
			errors.AddError(fieldName, fmt.Sprintf("%s must be a valid URL", fieldName))
		}
	}
	return errors
}

func ValidateRegex(value, fieldName, pattern, message string) ValidationErrors {
	var errors ValidationErrors
	if value != "" {
		matched, err := regexp.MatchString(pattern, value)
		if err != nil || !matched {
			errors.AddError(fieldName, message)
		}
	}
	return errors
}

func ValidatePositiveInteger(value uint, fieldName string) ValidationErrors {
	var errors ValidationErrors
	if value == 0 {
		errors.AddError(fieldName, fmt.Sprintf("%s must be a positive integer", fieldName))
	}
	return errors
}

func ValidatePassword(password, fieldName string) ValidationErrors {
	var errors ValidationErrors

	if len(password) < 8 {
		errors.AddError(fieldName, "password must be at least 8 characters long")
	}

	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasDigit {
		errors.AddError(fieldName, "password must contain at least one digit")
	}

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	if !hasLetter {
		errors.AddError(fieldName, "password must contain at least one letter")
	}

	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`).MatchString(password)
	if !hasSpecial {
		errors.AddError(fieldName, "password must contain at least one special character")
	}

	return errors
}

func MergeValidationErrors(errorsList ...ValidationErrors) ValidationErrors {
	var result ValidationErrors
	for _, errors := range errorsList {
		result = append(result, errors...)
	}
	return result
}
