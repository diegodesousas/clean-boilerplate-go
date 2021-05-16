package errors

import (
	"fmt"
)

type ValidationErrors struct {
	Errors map[string]ValidationError `json:"errors"`
}

func (e ValidationErrors) Error() string {
	return ""
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: map[string]ValidationError{},
	}
}

type ValidationError struct {
	Field  string                    `json:"-"`
	Errors map[ValidationType]string `json:"messages"`
}

func (e *ValidationErrors) Add(field string, errorType ValidationType, args ...interface{}) {
	validationError, ok := e.Errors[field]
	if !ok {
		defaultArg := []interface{}{field}
		validationError = ValidationError{
			Field: field,
			Errors: map[ValidationType]string{
				errorType: fmt.Sprintf(Messages[errorType], append(defaultArg, args...)...),
			},
		}

		e.Errors[field] = validationError
	}

	defaultArg := []interface{}{field}
	validationError.Errors[errorType] = fmt.Sprintf(Messages[errorType], append(defaultArg, args...)...)
}

func (e ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}
