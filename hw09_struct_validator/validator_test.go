package hw09structvalidator

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID           string `json:"id" validate:"len:36"`
		Name         string
		Age          uint8    `validate:"min:18|max:50"`
		Email        string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role         UserRole `validate:"len:5|in:admin,stuff"`
		Phones       []string `validate:"len:11|regexp:^\\+\\d+$"`
		ChildrenAges []int    `validate:"max:17"`
		Sex          string   `validate:"in:male|female"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte `validate:"len:10"`
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Request struct {
		Host        string
		Path        string
		Port        uint16   `validate:"min:1|max:65535"`
		Method      string   `validate:"in:GET,POST,PUT,DELETE"`
		Headers     []string `validate:"regexp:^\\S+:.*$"`
		Body        string
		AccessToken Token `validate:"len:10"`
	}
)

func lookupForError(ve ValidationError, vErrs ValidationErrors) bool {
	var equal bool
	for _, v := range vErrs {
		equal = true
		equal = equal && v.Field == ve.Field
		equal = equal && v.Step.Op == ve.Step.Op
		equal = equal && v.Step.Param == ve.Step.Param
		equal = equal && v.Val == ve.Val
		equal = equal && errors.Is(ve.Err, v.Err)

		if equal {
			return true
		}
	}
	return false
}

func checkErrorsMatch(t *testing.T, expected, actual ValidationErrors) {
	t.Helper()

	require.Equalf(t, len(expected), len(actual), "Different number of errors")
	checkErrorsInclusion(t, expected, actual)
}

func checkErrorsInclusion(t *testing.T, vErrs1, vErrs2 ValidationErrors) {
	t.Helper()

	for _, v := range vErrs2 {
		found := lookupForError(v, vErrs1)
		require.Truef(t, found,
			"{Field: %s, Validator: %v, Error: %v} expected but wasn't found",
			v.Field, v.Step, v.Err,
		)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name:        "test working with not struct types",
			in:          math.Pi,
			expectedErr: ErrValueIsNotAStruct,
		}, {
			name: "Request: test working with unsupported types",
			in: Request{
				Host: "google.com", Port: 443, Path: "/search", Method: "GET", Body: "", AccessToken: Token{},
				Headers: []string{"Content-Type: text/html", "User-Agent: Mozilla 5.0"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "AccessToken", Err: ErrValidationUnsupportedType,
				},
			},
		}, {
			name: "Response: test working with other tags",
			in:   Response{Code: 403, Body: "<HTML></HTML>"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code", Val: int64(403), Step: ValidationStep{Op: "in", Param: "200,404,500"},
					Err: ErrValidationIn,
				},
			},
		}, {
			name:        "App:test successful validation",
			in:          App{Version: "1.2.3"},
			expectedErr: ValidationErrors{},
		}, {
			name: "Payload:check that []byte will be validated as a string",
			in:   Token{Payload: []byte("fdght4df6b9")},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Payload", Val: "fdght4df6b9", Step: ValidationStep{Op: "len", Param: "10"},
					Err: ErrValidationLen,
				},
			},
		}, {
			name: "User:validation case with error in one validation rule format",
			in: User{
				ID: "1", Name: "John Doe", Age: 30,
				Email: "joe_doe@domain.com", Role: "Developer",
				Phones:       []string{"+123456789", "23456789012"},
				ChildrenAges: []int{1, 6, 21},
				Sex:          "male",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID", Val: "1", Step: ValidationStep{Op: "len", Param: "36"},
					Err: ErrValidationLen,
				}, ValidationError{
					Field: "Role", Step: ValidationStep{Op: "len", Param: "5"}, Val: "Developer",
					Err: ErrValidationLen,
				}, ValidationError{
					Field: "Role", Step: ValidationStep{Op: "in", Param: "admin,stuff"}, Val: "Developer",
					Err: ErrValidationIn,
				}, ValidationError{
					Field: "Phones", Step: ValidationStep{Op: "len", Param: "11"}, Val: "+123456789",
					Err: ErrValidationLen,
				}, ValidationError{
					Field: "Phones", Step: ValidationStep{Op: "regexp", Param: `^\+\d+$`}, Val: "23456789012",
					Err: ErrValidationRegexp,
				}, ValidationError{
					Field: "ChildrenAges", Step: ValidationStep{Op: "max", Param: "17"}, Val: int64(21),
					Err: ErrValidationMax,
				}, ValidationError{
					Field: "Sex", Val: "male",
					Err: ErrValidationTag,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt

			actualErr := Validate(tt.in)
			var actual, expected ValidationErrors
			if errors.As(actualErr, &actual) && errors.As(tt.expectedErr, &expected) {
				checkErrorsMatch(t, expected, actual)
			} else {
				require.Truef(t, errors.Is(actualErr, tt.expectedErr),
					"Expected error %v, but got %v", tt.expectedErr, actualErr)
			}
		})
	}
}
