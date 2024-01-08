package checks

import (
	"errors"
	"fmt"
)

type Min struct {
	value int64
	min   int64
}

func NewCheckMin(value, min int64) Min {
	return Min{
		value: value,
		min:   min,
	}
}

func (m Min) Validate() (bool, error) {
	return m.value >= m.min, nil
}

func (m Min) ValidationError(field string) error {
	return errors.New(fmt.Sprintf("field %q less than %d", field, m.min))
}
