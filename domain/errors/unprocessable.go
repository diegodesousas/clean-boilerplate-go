package errors

import "github.com/diegodesousas/clean-boilerplate-go/infra/errors"

type Unprocessable struct {
	*errors.ApplicationError
}

func NewUnprocessable(message string) *Unprocessable {
	return &Unprocessable{
		errors.NewApplicationError().WithMessage(message),
	}
}
