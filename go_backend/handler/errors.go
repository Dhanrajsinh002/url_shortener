package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// formatValidationError converts validator's raw, technical field errors
// into a single human-readable string, e.g. "password must be at least
// 8 characters" instead of the default Go struct-path dump.
func formatValidationError(err error) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return err.Error() // not a validation error — fall back to raw message
	}

	messages := make([]string, 0, len(ve))
	for _, fe := range ve {
		field := strings.ToLower(fe.Field())
		switch fe.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", field))
		case "min":
			messages = append(messages, fmt.Sprintf("%s must be at least %s characters", field, fe.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("%s must be at most %s characters", field, fe.Param()))
		default:
			messages = append(messages, fmt.Sprintf("%s is invalid", field))
		}
	}
	return strings.Join(messages, "; ")
}