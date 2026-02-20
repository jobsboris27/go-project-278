package validator

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

func FormatValidationErrors(ve validator.ValidationErrors) ErrorResponse {
	errorsMap := make(map[string]string)

	for _, e := range ve {
		field := e.Field()
		snake := ToSnakeCase(field)

		switch e.Tag() {
		case "required":
			errorsMap[snake] = "field is required"
		case "url":
			errorsMap[snake] = "must be a valid URL"
		case "min":
			errorsMap[snake] = "minimum length is " + e.Param()
		case "max":
			errorsMap[snake] = "maximum length is " + e.Param()
		default:
			errorsMap[snake] = "invalid value"
		}
	}

	return ErrorResponse{Errors: errorsMap}
}

func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && !unicode.IsUpper(rune(s[i-1])) {
				result.WriteRune('_')
			} else if i > 0 && i < len(s)-1 && unicode.IsLower(rune(s[i+1])) {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
