package checks

import (
	"errors"
	"fmt"
	"strings"
)

type In struct {
	string string
	values []string
}

func NewCheckIn(string string, values []string) In {
	return In{
		string: string,
		values: values,
	}
}

func (i In) Validate() (bool, error) {
	for _, v := range i.values {
		if i.string == v {
			return true, nil
		}
	}
	return false, nil
}

func (i In) ValidationError(field string) error {
	return errors.New(fmt.Sprintf("field %q not have value in: %s", field, strings.Join(i.values, ", ")))
}
