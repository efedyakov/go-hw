package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/efedyakov/go-hw/hw09_struct_validator/checks"

)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	for i, vErr := range v {
		if i > 0 {
			builder.WriteRune('\n')
		}
		builder.WriteString(vErr.Err.Error())
	}
	return builder.String()
}

var validate *validator.Validate

func Validate(v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return errors.New(" not a structure")
	}
	iv := reflect.ValueOf(v)
	vErrors := make(ValidationErrors, 0, iv.NumField())
	var err error
	for i := 0; i < iv.Type().NumField(); i++ {
		vErrors, err = validateField(iv.Type().Field(i), iv.Field(i), vErrors)
		if err != nil {
			var ve ValidationErrors
			if errors.As(err, &ve) {
				vErrors = append(vErrors, ve...)
			} else {
				return err
			}
		}
	}
	if len(vErrors) > 0 {
		return vErrors
	}
	return nil
}

func validateField(
	field reflect.StructField,
	value reflect.Value,
	vErrors ValidationErrors,
) (ValidationErrors, error) {
	if !field.IsExported() {
		return vErrors, nil
	}
	if isNested(field) {
		return vErrors, Validate(value.Interface())
	}

	checks, err := getChecks(field, value)
	if err != nil {
		return vErrors, err
	}
	for _, check := range checks {
		valid, err := check.Validate()
		if err != nil {
			return vErrors, err
		}
		if !valid {
			vErrors = append(vErrors, ValidationError{
				Field: field.Name,
				Err:   check.ValidationError(field.Name),
			})
		}
	}
	return vErrors, nil
}

func isNested(field reflect.StructField) bool {
	if field.Type.Kind() != reflect.Struct {
		return false
	}
	tag, ok := field.Tag.Lookup("validate")
	if !ok {
		return false
	}
	return tag == "nested"
}

func getChecks(field reflect.StructField, value reflect.Value) ([]checks.Check, error) {
	switch field.Type.Kind() { //nolint: exhaustive
	case reflect.String:
		checks, err := checks.GetCheckString(field, value)
		if err != nil {
			return nil, err
		}
		return checks, nil
	case reflect.Int:
		checks, err := checks.GetCheckInt(field, value)
		if err != nil {
			return nil, err
		}
		return checks, nil
	case reflect.Slice:
		switch field.Type.String() {
		case "[]string":
			checks, err := checks.GetCheckString(field, value)
			if err != nil {
				return nil, err
			}
			return checks, nil
		case "[]int":
			checks, err := checks.GetCheckInt(field, value)
			if err != nil {
				return nil, err
			}
			return checks, nil
		default:
			return nil, nil
		}
	default:
		return nil, nil
	}
}
