package checks

import (
	"errors"
	"fmt"
	"regexp"
)

type Regexp struct {
	string string
	expr   string
}

func NewCheckRegexp(string, expr string) Regexp {
	return Regexp{
		string: string,
		expr:   expr,
	}
}

func (r Regexp) Validate() (bool, error) {
	reg, err := regexp.Compile(r.expr)
	if err != nil {
		return false, err
	}
	return reg.Match([]byte(r.string)), nil
}

func (r Regexp) ValidationError(field string) error {
	return errors.New(fmt.Sprintf("field %q format is invalid", field))
}
