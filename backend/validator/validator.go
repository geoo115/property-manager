package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "validation errors"
	}

	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// Validator interface for custom validation
type Validator interface {
	Validate() error
}

// String validation functions
func ValidateRequired(value string, fieldName string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: "is required",
			Value:   value,
		}
	}
	return nil
}

func ValidateMinLength(value string, minLength int, fieldName string) *ValidationError {
	if utf8.RuneCountInString(value) < minLength {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("must be at least %d characters", minLength),
			Value:   value,
		}
	}
	return nil
}

func ValidateMaxLength(value string, maxLength int, fieldName string) *ValidationError {
	if utf8.RuneCountInString(value) > maxLength {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("must not exceed %d characters", maxLength),
			Value:   value,
		}
	}
	return nil
}

func ValidateEmail(email string, fieldName string) *ValidationError {
	if email == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: "is required",
			Value:   email,
		}
	}

	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be a valid email address",
			Value:   email,
		}
	}
	return nil
}

func ValidatePhone(phone string, fieldName string) *ValidationError {
	if phone == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: "is required",
			Value:   phone,
		}
	}

	// Remove common phone number separators
	cleanPhone := strings.ReplaceAll(phone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "(", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ")", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "+", "")

	// Basic phone validation (10-15 digits)
	phoneRegex := regexp.MustCompile(`^\d{10,15}$`)
	if !phoneRegex.MatchString(cleanPhone) {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be a valid phone number (10-15 digits)",
			Value:   phone,
		}
	}
	return nil
}

func ValidatePassword(password string, fieldName string) *ValidationError {
	if password == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: "is required",
			Value:   "",
		}
	}

	if len(password) < 8 {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be at least 8 characters long",
			Value:   "",
		}
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return &ValidationError{
			Field:   fieldName,
			Message: "must contain at least one uppercase letter, lowercase letter, digit, and special character",
			Value:   "",
		}
	}

	return nil
}

func ValidateRole(role string, fieldName string) *ValidationError {
	validRoles := []string{"admin", "tenant", "landlord", "maintenanceTeam"}
	for _, validRole := range validRoles {
		if role == validRole {
			return nil
		}
	}
	return &ValidationError{
		Field:   fieldName,
		Message: fmt.Sprintf("must be one of: %s", strings.Join(validRoles, ", ")),
		Value:   role,
	}
}

// Numeric validation functions
func ValidatePositiveFloat(value float64, fieldName string) *ValidationError {
	if value <= 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be greater than 0",
			Value:   value,
		}
	}
	return nil
}

func ValidateNonNegativeFloat(value float64, fieldName string) *ValidationError {
	if value < 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be non-negative",
			Value:   value,
		}
	}
	return nil
}

func ValidatePositiveInt(value int, fieldName string) *ValidationError {
	if value <= 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be greater than 0",
			Value:   value,
		}
	}
	return nil
}

func ValidateNonNegativeInt(value int, fieldName string) *ValidationError {
	if value < 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be non-negative",
			Value:   value,
		}
	}
	return nil
}

// Sanitization functions
func SanitizeString(input string) string {
	// Remove leading/trailing spaces
	sanitized := strings.TrimSpace(input)

	// Remove null bytes and other control characters
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")

	return sanitized
}

func SanitizeHTML(input string) string {
	// Basic HTML sanitization - remove common HTML tags
	htmlTags := []string{
		"<script>", "</script>",
		"<iframe>", "</iframe>",
		"<object>", "</object>",
		"<embed>", "</embed>",
		"<form>", "</form>",
	}

	result := input
	for _, tag := range htmlTags {
		result = strings.ReplaceAll(result, tag, "")
		result = strings.ReplaceAll(result, strings.ToUpper(tag), "")
	}

	return SanitizeString(result)
}

// Utility function to collect validation errors
func CollectValidationErrors(errors ...*ValidationError) ValidationErrors {
	var result ValidationErrors
	for _, err := range errors {
		if err != nil {
			result = append(result, *err)
		}
	}
	return result
}
