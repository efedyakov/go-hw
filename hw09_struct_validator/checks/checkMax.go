package checks

import (
	"fmt"
)

type Max struct {
	value int64
	max   int64
}

func NewCheckMax(value, max int64) Max {
	return Max{
		value: value,
		max:   max,
	}
}

func (m Max) Validate() (bool, error) {
	return m.value <= m.max, nil
}

func (m Max) ValidationError(field string) error {
	return fmt.Errorf("field %q greater than %d", field, m.max)
}
