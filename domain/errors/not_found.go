package errors

import (
	"errors"
)

type EntityNotFound struct {
	error
}

func NewEntityNotFound(message string) *EntityNotFound {
	return &EntityNotFound{
		errors.New(message),
	}
}
