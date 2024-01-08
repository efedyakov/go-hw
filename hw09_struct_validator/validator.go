package hw09structvalidator

import (
	"errors"
	"fmt"
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

func Validate(v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return errors.New(" not a structure")
	}
	iv := reflect.ValueOf(v)
	vErrors := make(ValidationErrors, 0, iv.NumField())
	var err error
	for i := 0; i < iv.Type().NumField(); i++ {
		vErrors, err = validateField(iv.Type().Field(i), iv.Field(i), vErrors, prefix)
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
	prefix string,
) (ValidationErrors, error) {
	if !field.IsExported() {
		return vErrors, nil
	}
	if isNested(field) {
		prefix = fmt.Sprintf("%s%s->", prefix, field.Name)
		return vErrors, Validate(value.Interface(), prefix)
	}
	ext := getRulesExtractor(field.Type.Kind(), field.Type)
	if ext == nil {
		return vErrors, nil
	}
	rules, err := ext.Extract(field, value)
	if err != nil {
		return vErrors, err
	}
	for _, r := range rules {
		valid, err := r.Validate()
		if err != nil {
			return vErrors, err
		}
		if !valid {
			vErrors = append(vErrors, ValidationError{
				Field: field.Name,
				Err:   r.ValidationError(prefix + field.Name),
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

func getChecks(field reflect.StructField, value reflect.Value /*kind reflect.Kind, fieldType reflect.Type*/) checks.Check {
	switch field.Type.Kind() {
	case reflect.String:
		checks, err := checks.GetCheckString(field, value)
		if err != nil {
			return nil
		}
		return checks
	case reflect.Int:
		/*if _, ok := extractors["int"]; !ok {
			extractors["int"] = rule.NewIntRulesExtractor("validate", "|", ":", ",")
		}
		return extractors["int"]*/
	case reflect.Slice:
		//return getSliceRulesExtractor(field.Type)
	default:
		return nil
	}
}
