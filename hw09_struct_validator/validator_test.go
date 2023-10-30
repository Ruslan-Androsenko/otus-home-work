package hw09structvalidator

import (
	"encoding/json"
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
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"regexp:^\\d+$|len:11"`
		meta   json.RawMessage //nolint:unused
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
				ID:     "24458918609ae919f6592976848b9be123ab",
				Name:   "John Smith",
				Age:    25,
				Email:  "admin@example.ru",
				Role:   "admin",
				Phones: []string{"27401234567"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "4d0acc63962a43f30a2cfff3e4f759839182",
				Name:   "John Doe",
				Age:    52,
				Email:  "test.ru",
				Role:   "guest",
				Phones: []string{"7900123ab67", "790012345"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: ErrValueIsGreaterThanSpecified},
				ValidationError{Field: "Email", Err: ErrValueDoesNotMatchRegExp},
				ValidationError{Field: "Role", Err: ErrValueOutsideSpecifiedRange},
				ValidationError{Field: "Phones", Err: ErrValueDoesNotMatchRegExp},
				ValidationError{Field: "Phones", Err: ErrLengthDoesNotMatchSpecified},
			},
		},
		{
			in: User{
				ID:     "1f2ffb504440fda52c3fd95ba4252ab1abcd1234",
				Name:   "Jane Doe",
				Age:    16,
				Email:  "janedoe@example.ru",
				Role:   "stuff",
				Phones: []string{"79001234567", "27401234567"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrLengthDoesNotMatchSpecified},
				ValidationError{Field: "Age", Err: ErrValueIsLessThanSpecified},
			},
		},
		{
			in: App{
				Version: "1.0.5",
			},
			expectedErr: nil,
		},
		{
			in: App{},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrLengthDoesNotMatchSpecified},
			},
		},
		{
			in: App{
				Version: "1.1.3-rc.2",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrLengthDoesNotMatchSpecified},
			},
		},
		{
			in: Token{}, expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("X-Real-IP"),
				Payload:   []byte("{'id': 34, 'type': 'example'}"),
				Signature: []byte("f6f7ff173a277e3b172e66d98cf4d2ac"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "Successfully",
			},
			expectedErr: nil,
		},
		{
			in: Response{},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrValueOutsideSpecifiedRange},
			},
		},
		{
			in: Response{
				Code: 204,
				Body: "Out of range",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrValueOutsideSpecifiedRange},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.EqualValues(t, tt.expectedErr, err)
		})
	}
}
