package checks

import (
	"fmt"
)

type Len struct {
	string string
	length int
}

func NewCheckLen(str string, length int) Len {
	return Len{
		string: str,
		length: length,
	}
}

func (l Len) Validate() (bool, error) {
	return len(l.string) == l.length, nil
}

func (l Len) ValidationError(field string) error {
	return fmt.Errorf("field %q must have len %d", field, l.length)
}
