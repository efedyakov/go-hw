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

	arrString = "[]string"
	arrInt    = "[]int"
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
			checks = getChecksLen(checks, value, length)
		case checkRegexp:
			if checkVal == "" {
				return nil, ErrCheckStringNotValid
			}
			checks = getChecksRegexp(checks, value, checkVal)
		case checkIn:
			if checkVal == "" {
				return nil, ErrCheckStringNotValid
			}
			checks = getChecksIn(checks, value, checkVal)
		}
	}
	return checks, nil
}

func getChecksLen(checks []Check, value reflect.Value, length int) []Check {
	if value.Type().String() == arrString {
		for j := 0; j < value.Len(); j++ {
			return append(checks, NewCheckLen(value.Index(j).String(), length))
		}
	}
	return append(checks, NewCheckLen(value.String(), length))
}

func getChecksRegexp(checks []Check, value reflect.Value, regexp string) []Check {
	if value.Type().String() == arrString {
		for j := 0; j < value.Len(); j++ {
			return append(checks, NewCheckRegexp(value.Index(j).String(), regexp))
		}
	}
	return append(checks, NewCheckRegexp(value.String(), regexp))
}

func getChecksIn(checks []Check, value reflect.Value, valsstr string) []Check {
	vals := strings.Split(valsstr, valueSeparator)
	if value.Type().String() == arrString {
		for j := 0; j < value.Len(); j++ {
			return append(checks, NewCheckIn(value.Index(j).String(), vals))
		}
	}
	return append(checks, NewCheckIn(value.String(), vals))
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
			checks = getChecksIn_int(checks, value, checkVal)
		case checkMin:
			min, err := strconv.Atoi(checkVal)
			if err != nil {
				return nil, err
			}
			checks = getChecksMin(checks, value, min)
		case checkMax:
			max, err := strconv.Atoi(checkVal)
			if err != nil {
				return nil, err
			}
			checks = getChecksMax(checks, value, max)
		}
	}
	return checks, nil
}

func getChecksIn_int(checks []Check, value reflect.Value, valsstr string) []Check {
	vals := strings.Split(valsstr, valueSeparator)
	if value.Type().String() == arrInt {
		for j := 0; j < value.Len(); j++ {
			return append(checks, NewCheckIn(strconv.FormatInt(value.Index(j).Int(), 10), vals))
		}
	}
	return append(checks, NewCheckIn(strconv.FormatInt(value.Int(), 10), vals))
}

func getChecksMax(checks []Check, value reflect.Value, max int) []Check {
	if value.Type().String() == arrInt {
		for j := 0; j < value.Len(); j++ {
			return append(checks, NewCheckMax(value.Index(j).Int(), int64(max)))
		}
	}
	return append(checks, NewCheckMax(value.Int(), int64(max)))
}

func getChecksMin(checks []Check, value reflect.Value, min int) []Check {
	if value.Type().String() == arrInt {
		for j := 0; j < value.Len(); j++ {
			return append(checks, NewCheckMin(value.Index(j).Int(), int64(min)))
		}
	}
	return append(checks, NewCheckMin(value.Int(), int64(min)))
}
