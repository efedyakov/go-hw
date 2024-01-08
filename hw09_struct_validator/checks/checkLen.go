package checks

import (
	"errors"
	"fmt"
)

type Len struct {
	string string
	length int
}

func NewCheckLen(string string, length int) Len {
	return Len{
		string: string,
		length: length,
	}
}

func (l Len) Validate() (bool, error) {
	return len(l.string) == l.length, nil
}

func (l Len) ValidationError(field string) error {
	return errors.New(fmt.Sprintf("field %q must have len %d", field, l.length))
}
