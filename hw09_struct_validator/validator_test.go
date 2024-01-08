package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Vasya Pupkin",
				Age:    49,
				Email:  "VasyaPupkin@gmail.com",
				Role:   "admin",
				Phones: []string{"12345678901", "12345678901"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Vasya Pupkin",
				Age:    10,
				Email:  "VasyaPupkin@gmailcom",
				Role:   "worker",
				Phones: []string{"1234567890", "12345678901"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   errors.New(`field "Age" less than 18`),
				},
				ValidationError{
					Field: "Email",
					Err:   errors.New(`field "Email" format is invalid`),
				},
				ValidationError{
					Field: "Role",
					Err:   errors.New(`field "Role" not have value in: admin, stuff`),
				},
				ValidationError{
					Field: "Phones",
					Err:   errors.New(`field "Phones" must have len 11`),
				},
			},
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 404,
				Body: "{id:123}",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 0,
				Body: "{id:123}",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   errors.New(`field "Code" not have value in: 200, 404, 500`),
				},
			},
		},
	}

	for index, tt := range tests {
		index := index
		t.Run(fmt.Sprintf("case %d", index), func(t *testing.T) {
			tt := tt
			t.Parallel()
			testname := fmt.Sprintf("%d", index)
			t.Run(testname, func(t *testing.T) {
				t.Parallel()
				err := Validate(tt.in)
				if tt.expectedErr == nil {
					require.NoError(t, err)
				}
				var vError ValidationErrors
				if errors.As(err, &vError) {
					require.Equal(t, tt.expectedErr.Error(), err.Error())
				} else {
					require.ErrorIs(t, err, tt.expectedErr)
				}
			})
			_ = tt
		})
	}
}
