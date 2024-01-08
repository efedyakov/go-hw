package checks

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

const (
	checkLen    = "len"
	checkRegexp = "regexp"
	checkIn     = "in"
	checkMin    = "min"
	checkMax    = "max"

	tagName        = "validate"
	tagSeparator   = "|"
	checkSeparator = ":"
	valueSeparator = ","
)

var (
	ErrCheckStringNotValid = errors.New(`check string not valid`)
)

type Check interface {
	Validate() (bool, error)
	ValidationError(field string) error
}

func GetCheckString(field reflect.StructField, value reflect.Value) ([]Check, error) {
	tag, ok := field.Tag.Lookup(tagName)
	if !ok {
		return nil, nil
	}
	checksStr := strings.Split(tag, tagSeparator)
	checks := make([]Check, 0, len(checksStr))
	for _, checkStr := range checksStr {
		ss := strings.Split(checkStr, checkSeparator)
		checkName, checkVal := ss[0], ss[1]
		switch checkName {
		case checkLen:
			length, err := strconv.Atoi(checkVal)
			if err != nil {
				return nil, ErrCheckStringNotValid
			}
			if value.Type().String() == "[]string" {
				for j := 0; j < value.Len(); j++ {
					checks = append(checks, NewCheckLen(value.Index(j).String(), length))
				}
			} else {
				checks = append(checks, NewCheckLen(value.String(), length))
			}
		case checkRegexp:
			if checkVal == "" {
				return nil, ErrCheckStringNotValid
			}
			if value.Type().String() == "[]string" {
				for j := 0; j < value.Len(); j++ {
					checks = append(checks, NewCheckRegexp(value.Index(j).String(), checkVal))
				}
			} else {
				checks = append(checks, NewCheckRegexp(value.String(), checkVal))
			}
		case checkIn:
			if checkVal == "" {
				return nil, ErrCheckStringNotValid
			}
			vals := strings.Split(checkVal, valueSeparator)
			if value.Type().String() == "[]string" {
				for j := 0; j < value.Len(); j++ {
					checks = append(checks, NewCheckIn(value.Index(j).String(), vals))
				}
			} else {
				checks = append(checks, NewCheckIn(value.String(), vals))
			}
		}
	}
	return checks, nil
}

func GetCheckInt(field reflect.StructField, value reflect.Value) ([]Check, error) {
	tag, ok := field.Tag.Lookup(tagName)
	if !ok {
		return nil, nil
	}
	checksStr := strings.Split(tag, tagSeparator)
	checks := make([]Check, 0, len(checksStr))
	for _, checkStr := range checksStr {
		ss := strings.Split(checkStr, checkSeparator)
		checkName, checkVal := ss[0], ss[1]
		switch checkName {
		case checkIn:
			if checkVal == "" {
				return nil, ErrCheckStringNotValid
			}
			vals := strings.Split(checkVal, valueSeparator)
			if value.Type().String() == "[]int" {
				for j := 0; j < value.Len(); j++ {
					checks = append(checks, NewCheckIn(strconv.FormatInt(value.Index(j).Int(), 10), vals))
				}
			} else {
				checks = append(checks, NewCheckIn(strconv.FormatInt(value.Int(), 10), vals))
			}
		case checkMin:
			min, err := strconv.Atoi(checkVal)
			if err != nil {
				return nil, err
			}
			if value.Type().String() == "[]int" {
				for j := 0; j < value.Len(); j++ {
					checks = append(checks, NewCheckMin(value.Index(j).Int(), int64(min)))
				}
			} else {
				checks = append(checks, NewCheckMin(value.Int(), int64(min)))
			}
		case checkMax:
			max, err := strconv.Atoi(checkVal)
			if err != nil {
				return nil, err
			}
			if value.Type().String() == "[]int" {
				for j := 0; j < value.Len(); j++ {
					checks = append(checks, NewCheckMax(value.Index(j).Int(), int64(max)))
				}
			} else {
				checks = append(checks, NewCheckMax(value.Int(), int64(max)))
			}
		}
	}
	return checks, nil
}
