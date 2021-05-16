package errors

import (
	"encoding/json"
	"fmt"
)

type ApplicationError struct {
	Fields  map[string]interface{} `json:"fields,omitempty"`
	Message string                 `json:"message,omitempty"`
}

func NewApplicationError() *ApplicationError {
	return &ApplicationError{
		Fields: map[string]interface{}{},
	}
}

func (e *ApplicationError) WithFields(key string, value interface{}) *ApplicationError {
	e.Fields[key] = value

	return e
}

func (e *ApplicationError) WithMessage(message string) *ApplicationError {
	e.Message = message

	return e
}

func (e *ApplicationError) Error() string {
	bytes, _ := json.Marshal(e.Fields)

	return fmt.Sprintf("application-error: Fields=%s", string(bytes))
}
