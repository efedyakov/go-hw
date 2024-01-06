package hw09structvalidator

import (
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

var validate *validator.Validate

func Validate(v interface{}) error {
	// Place your code here.
	return nil //validate.Struct(v)
}
